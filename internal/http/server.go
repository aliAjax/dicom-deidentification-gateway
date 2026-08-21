package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/example/dicom-deidentification-gateway/internal/config"
	"github.com/example/dicom-deidentification-gateway/internal/deidentification/domain"
	deid "github.com/example/dicom-deidentification-gateway/internal/deidentification/infrastructure"
	dicomd "github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	dicomi "github.com/example/dicom-deidentification-gateway/internal/dicom/infrastructure"
	"github.com/example/dicom-deidentification-gateway/internal/platform"
	routingd "github.com/example/dicom-deidentification-gateway/internal/routing/domain"
	routingi "github.com/example/dicom-deidentification-gateway/internal/routing/infrastructure"
	studyd "github.com/example/dicom-deidentification-gateway/internal/study/domain"
	studyi "github.com/example/dicom-deidentification-gateway/internal/study/infrastructure"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Server struct {
	cfg       config.Config
	log       *slog.Logger
	instances *dicomi.MemoryStore
	parser    dicomd.Parser
	profiles  *deid.Service
	targets   *routingi.Memory
	studies   *studyi.Memory
	mu        sync.RWMutex
	exports   map[string]map[string]any
}

func structStudy(uid string) studyd.Study {
	return studyd.Study{ID: platform.NewID("study"), StudyUID: uid, Status: "received"}
}

func New(cfg config.Config, log *slog.Logger) *Server {
	return &Server{cfg: cfg, log: log, instances: dicomi.NewMemoryStore(), parser: dicomi.Parser{}, profiles: deid.NewService(), targets: routingi.New(), studies: studyi.New(), exports: map[string]map[string]any{}}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.health)
	mux.HandleFunc("/metrics", s.metrics)
	mux.HandleFunc("/v1/dicom/instances", s.instancesHandler)
	mux.HandleFunc("/v1/deidentification/profiles", s.profilesHandler)
	mux.HandleFunc("/v1/instances/", s.instanceHandler)
	mux.HandleFunc("/v1/studies/", s.studiesHandler)
	mux.HandleFunc("/v1/targets", s.targetsHandler)
	mux.HandleFunc("/v1/exports", s.exportsHandler)
	return mux
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "time": time.Now().UTC()})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# HELP dicom_instances_total received instances\ndicom_instances_total", lenMust(s.instances.List(context.Background(), "", 0)))
}
func lenMust(v []dicomd.Instance, _ error) int { return len(v) }
func (s *Server) instancesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		status := r.URL.Query().Get("status")
		list, err := s.instances.List(r.Context(), status, 100)
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, 200, list)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, s.cfg.MaxBodyBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErr(w, platform.Invalid("read DICOM payload", err))
		return
	}
	i, err := s.parser.Parse(r.Context(), body)
	if err != nil {
		writeErr(w, err)
		return
	}
	if old, e := s.instances.FindByHash(r.Context(), i.ContentHash); e == nil {
		writeJSON(w, 200, old)
		return
	}
	if err := dicomi.Validate(i); err != nil {
		i.Status = dicomd.Quarantined
	}
	_ = s.persistSpool(i.ID, body)
	if err := s.instances.Save(r.Context(), i); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, i)
}
func (s *Server) persistSpool(id string, b []byte) error {
	if err := os.MkdirAll(s.cfg.SpoolDir, 0750); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.cfg.SpoolDir, id+".dcm"), b, 0600)
}
func (s *Server) profilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var req struct {
		Name  string `json:"name"`
		Rules []struct {
			Tag    string `json:"tag"`
			Action string `json:"action"`
			Value  string `json:"value"`
		} `json:"rules"`
		DateShiftDays   int  `json:"date_shift_days"`
		KeepPrivateTags bool `json:"keep_private_tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, platform.Invalid("invalid profile", err))
		return
	}
	p := domain.Profile{ID: platform.NewID("profile"), Name: req.Name, DateShiftDays: req.DateShiftDays, KeepPrivateTags: req.KeepPrivateTags}
	for _, x := range req.Rules {
		p.Rules = append(p.Rules, domain.Rule{Tag: x.Tag, Action: deid.ParseAction(x.Action), Value: x.Value})
	}
	if err := s.profiles.SaveProfile(r.Context(), p); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, p)
}
func (s *Server) instanceHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeErr(w, platform.NotFound("instance"))
		return
	}
	id := parts[2]
	i, err := s.instances.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	if len(parts) == 3 {
		writeJSON(w, 200, i)
		return
	}
	if len(parts) >= 4 && parts[3] == "deidentify" && r.Method == http.MethodPost {
		var req struct {
			ProfileID string `json:"profile_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, err)
			return
		}
		p, err := s.profiles.Profile(r.Context(), req.ProfileID)
		if err != nil {
			writeErr(w, err)
			return
		}
			tags, err := s.profiles.Apply(r.Context(), p, i.Tags)
			if err != nil {
				i = rollbackInstance(i)
				_ = s.instances.Update(r.Context(), i)
				writeErr(w, err)
				return
			}
		i.Tags = tags
		i.PatientID = tags["PatientID"]
		i.Status = dicomd.Deidentified
		if err := s.instances.Update(r.Context(), i); err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, 200, i)
		return
	}
	writeErr(w, platform.NotFound("operation"))
}
func (s *Server) targetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v, e := s.targets.List(r.Context())
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, v)
		return
	}
	var t routingd.Target
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeErr(w, err)
		return
	}
	if t.ID == "" {
		t.ID = platform.NewID("target")
	}
	if err := s.targets.Save(r.Context(), t); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, t)
}
func (s *Server) studiesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeErr(w, platform.NotFound("study"))
		return
	}
	uid := parts[2]
	st, err := s.studies.ByUID(r.Context(), uid)
	if err != nil {
		st = structStudy(uid)
	}
	if len(parts) >= 4 && parts[3] == "send" {
		writeJSON(w, missingStudyStatus(), map[string]any{"study_uid": uid, "status": "queued"})
		return
	}
	writeJSON(w, 200, st)
}
func (s *Server) exportsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			StudyID     string   `json:"study_id"`
			InstanceIDs []string `json:"instance_ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErr(w, err)
			return
		}
		if !exportRequestValid(req.StudyID, len(req.InstanceIDs)) {
			writeErr(w, platform.Invalid("study and instance ids are required", nil))
			return
		}
		id := platform.NewID("export")
		h := sha256.Sum256([]byte(strings.Join(req.InstanceIDs, ",")))
		e := map[string]any{"id": id, "study_id": req.StudyID, "instance_count": len(req.InstanceIDs), "status": "completed", "hash": hex.EncodeToString(h[:]), "created_at": time.Now().UTC()}
		s.mu.Lock()
		s.exports[id] = e
		s.mu.Unlock()
		writeJSON(w, 202, e)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/exports/")
	s.mu.RLock()
	e, ok := s.exports[id]
	s.mu.RUnlock()
	if !ok {
		writeErr(w, platform.NotFound("export"))
		return
	}
	writeJSON(w, 200, e)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, err error) {
	status := 500
	var p *platform.Error
	if errors.As(err, &p) {
		switch p.Code {
		case platform.CodeInvalid:
			status = 400
		case platform.CodeNotFound:
			status = 404
		case platform.CodeConflict:
			status = 409
		case platform.CodeUnauthorized:
			status = 401
		}
	}
	writeJSON(w, status, map[string]any{"error": err.Error()})
}
func NewHTTPServer(cfg config.Config, h http.Handler) *http.Server {
	return &http.Server{Addr: cfg.HTTPAddr, Handler: h, ReadTimeout: cfg.ReadTimeout, WriteTimeout: cfg.WriteTimeout, IdleTimeout: 60 * time.Second}
}
func RequestLogger(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}

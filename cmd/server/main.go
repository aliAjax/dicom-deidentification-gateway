package main

import (
	"context"
	"github.com/example/dicom-deidentification-gateway/internal/auth"
	"github.com/example/dicom-deidentification-gateway/internal/config"
	dicomadapter "github.com/example/dicom-deidentification-gateway/internal/dicom/adapter"
	dicomdomain "github.com/example/dicom-deidentification-gateway/internal/dicom/domain"
	httpapi "github.com/example/dicom-deidentification-gateway/internal/http"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := httpapi.New(cfg, log)
	h := auth.Middleware(auth.APIKeyValidator{Key: cfg.APIKey}, httpapi.RequestLogger(log, srv.Handler()))
	server := httpapi.NewHTTPServer(cfg, h)
	go func() {
		log.Info("http server started", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server stopped", "error", err)
		}
	}()
	dimseCtx, dimseCancel := context.WithCancel(context.Background())
	defer dimseCancel()
	go func() {
		log.Info("dimse listener started", "addr", cfg.DICOMAddr)
		err := dicomadapter.ServeDIMSE(dimseCtx, cfg.DICOMAddr, func(ctx context.Context, s *dicomadapter.Session, raw []byte) error {
			pdu, err := dicomdomain.DecodePDU(raw)
			if err != nil {
				return err
			}
			if pdu.Type == dicomdomain.PDUDATATF {
				return s.WritePDU(ctx, raw)
			}
			return nil
		})
		if err != nil {
			log.Error("dimse listener stopped", "error", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dimseCancel()
	_ = server.Shutdown(shutdown)
	log.Info("server shutdown")
}

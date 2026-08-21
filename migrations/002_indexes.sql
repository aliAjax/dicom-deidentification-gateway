CREATE INDEX IF NOT EXISTS idx_transfer_jobs_status ON transfer_jobs(status);
CREATE INDEX IF NOT EXISTS idx_exports_study ON exports(study_id);

CREATE TABLE IF NOT EXISTS dicom_instances (id TEXT PRIMARY KEY, patient_id TEXT NOT NULL, study_uid TEXT NOT NULL, series_uid TEXT NOT NULL, sop_instance_uid TEXT NOT NULL, sop_class_uid TEXT, modality TEXT, source_ae TEXT, file_path TEXT NOT NULL, content_hash TEXT NOT NULL UNIQUE, status TEXT NOT NULL, size_bytes BIGINT NOT NULL, created_at TIMESTAMPTZ NOT NULL, updated_at TIMESTAMPTZ NOT NULL);
CREATE INDEX IF NOT EXISTS idx_dicom_instances_study ON dicom_instances(study_uid);
CREATE INDEX IF NOT EXISTS idx_dicom_instances_status ON dicom_instances(status);
CREATE TABLE IF NOT EXISTS deidentification_profiles (id TEXT PRIMARY KEY,name TEXT NOT NULL,rules JSONB NOT NULL,created_at TIMESTAMPTZ NOT NULL);
CREATE TABLE IF NOT EXISTS transfer_targets (id TEXT PRIMARY KEY,name TEXT NOT NULL,address TEXT NOT NULL,port INTEGER NOT NULL,enabled BOOLEAN NOT NULL);
CREATE TABLE IF NOT EXISTS transfer_jobs (id TEXT PRIMARY KEY,instance_id TEXT NOT NULL,target_id TEXT NOT NULL,status TEXT NOT NULL,attempts INTEGER NOT NULL,last_error TEXT);
CREATE TABLE IF NOT EXISTS exports (id TEXT PRIMARY KEY,study_id TEXT NOT NULL,path TEXT NOT NULL,hash TEXT NOT NULL,status TEXT NOT NULL,instance_count INTEGER NOT NULL);

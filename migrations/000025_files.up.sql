CREATE TABLE patient_files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL REFERENCES patients(id),
  uploaded_by UUID NOT NULL REFERENCES users(id),
  source TEXT NOT NULL DEFAULT 'patient' CHECK (
    source IN ('patient', 'facility')
  ),
  title TEXT NOT NULL,
  doc_type TEXT NOT NULL CHECK (
    doc_type IN (
      'lab_result', 'letter', 'scan',
      'prescription', 'other'
    )
  ),
  storage_key TEXT NOT NULL UNIQUE,
  content_type TEXT NOT NULL,
  size_bytes BIGINT NOT NULL CHECK (
    size_bytes > 0 AND size_bytes <= 10485760
  ),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX ON patient_files (patient_id, created_at);
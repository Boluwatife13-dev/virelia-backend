CREATE TABLE record_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL REFERENCES patients(id),
  facility_id UUID NOT NULL REFERENCES facilities(id),
  author_id UUID NOT NULL REFERENCES users(id),
  appointment_id UUID REFERENCES appointments(id),
  kind TEXT NOT NULL CHECK (
    kind IN ('visit', 'vitals', 'diagnosis', 'note')
  ),
  data JSONB NOT NULL,
  version INT NOT NULL DEFAULT 1 CHECK (version >= 1),
  replaces_id UUID REFERENCES record_entries(id),
  correction_reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK ((replaces_id IS NULL) = (version = 1)),
  CHECK (replaces_id IS NULL OR correction_reason IS NOT NULL)
);

CREATE UNIQUE INDEX record_entries_one_correction
ON record_entries (replaces_id)
WHERE replaces_id IS NOT NULL;

CREATE INDEX ON record_entries (patient_id, created_at);
CREATE INDEX ON record_entries (facility_id, created_at);

CREATE FUNCTION record_entries_locked() RETURNS trigger AS $$
BEGIN
  RAISE EXCEPTION 'record entries cannot be changed or deleted, add a correction instead';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER record_entries_no_change
BEFORE UPDATE OR DELETE ON record_entries
FOR EACH ROW EXECUTE FUNCTION record_entries_locked();
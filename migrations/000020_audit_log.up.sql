-- Who did what, to which patient, when and from where.
CREATE TABLE audit_log (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  -- Who did it. No foreign key, so the log stays
  -- complete even if a user is later removed.
  actor_id UUID,
  actor_role TEXT,
  facility_id UUID,

  -- What they did, e.g. 'facility.approve'.
  action TEXT NOT NULL,

  -- Whose record it touched, if any.
  patient_id UUID,

  -- The thing acted on, e.g. a staff member's id.
  target_id TEXT,

  details JSONB NOT NULL DEFAULT '{}',
  ip TEXT,
  user_agent TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  -- Each row stores a fingerprint of itself and of
  -- the row before it. Changing an old row breaks
  -- the chain, so tampering shows up.
  prev_hash TEXT NOT NULL,
  hash TEXT NOT NULL
);

CREATE INDEX ON audit_log (patient_id, created_at);
CREATE INDEX ON audit_log (actor_id, created_at);
CREATE INDEX ON audit_log (facility_id, created_at);

-- Fills prev_hash and hash on every new row.
CREATE FUNCTION audit_log_chain() RETURNS trigger AS $$
DECLARE
  last_hash TEXT;
BEGIN
  -- One insert at a time, so the chain never forks.
  PERFORM pg_advisory_xact_lock(424242);

  SELECT hash INTO last_hash
  FROM audit_log
  ORDER BY id DESC
  LIMIT 1;

  NEW.prev_hash := COALESCE(last_hash, '');

  NEW.hash := encode(sha256(convert_to(concat_ws('|',
    NEW.prev_hash, NEW.id, NEW.actor_id, NEW.actor_role,
    NEW.facility_id, NEW.action, NEW.patient_id,
    NEW.target_id, NEW.details, NEW.created_at
  ), 'UTF8')), 'hex');

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_log_chain
BEFORE INSERT ON audit_log
FOR EACH ROW EXECUTE FUNCTION audit_log_chain();

-- Refuses any change or delete. The log can only grow.
CREATE FUNCTION audit_log_locked() RETURNS trigger AS $$
BEGIN
  RAISE EXCEPTION 'audit_log cannot be changed or deleted';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_log_no_change
BEFORE UPDATE OR DELETE ON audit_log
FOR EACH ROW EXECUTE FUNCTION audit_log_locked();

CREATE TRIGGER audit_log_no_truncate
BEFORE TRUNCATE ON audit_log
FOR EACH STATEMENT EXECUTE FUNCTION audit_log_locked();
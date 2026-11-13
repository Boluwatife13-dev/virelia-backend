CREATE TABLE consents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL REFERENCES patients(id),
  subject TEXT NOT NULL CHECK (subject IN ('referral', 'transfer')),
  subject_id UUID NOT NULL,
  from_facility_id UUID NOT NULL REFERENCES facilities(id),
  to_facility_id UUID NOT NULL REFERENCES facilities(id),
  decision TEXT NOT NULL CHECK (decision IN ('granted', 'refused')),
  share_records BOOLEAN NOT NULL DEFAULT false,
  method TEXT NOT NULL CHECK (method IN ('app', 'in_person')),
  given_by UUID NOT NULL REFERENCES users(id),
  witness_name TEXT,
  form_reference TEXT,
  expires_at TIMESTAMPTZ,
  withdrawn_at TIMESTAMPTZ,
  withdrawn_by UUID REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (
    method <> 'in_person'
    OR (witness_name IS NOT NULL AND form_reference IS NOT NULL)
  ),
  CHECK (decision = 'granted' OR share_records = false)
);

CREATE UNIQUE INDEX consents_one_active
ON consents (subject, subject_id)
WHERE withdrawn_at IS NULL;

CREATE INDEX ON consents (patient_id, created_at);
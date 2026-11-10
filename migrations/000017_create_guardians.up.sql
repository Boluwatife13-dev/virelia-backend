
ALTER TABLE patients DROP CONSTRAINT patients_contact;


CREATE TABLE guardians (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  child_id UUID NOT NULL
    REFERENCES patients(id) ON DELETE CASCADE,

  guardian_id UUID NOT NULL
    REFERENCES patients(id) ON DELETE CASCADE,

  relationship TEXT NOT NULL CHECK (
    relationship IN ('mother', 'father', 'guardian', 'other')
  ),

  added_by UUID REFERENCES users(id),

  
  verified_at TIMESTAMPTZ,

 
  expires_at TIMESTAMPTZ NOT NULL,

  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE (child_id, guardian_id),
  CHECK (child_id <> guardian_id)
);

CREATE INDEX ON guardians (guardian_id);
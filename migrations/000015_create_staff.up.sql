CREATE TABLE staff (
  user_id UUID PRIMARY KEY
    REFERENCES users(id) ON DELETE CASCADE,

  facility_id UUID REFERENCES facilities(id),

  role TEXT NOT NULL CHECK (
    role IN ('platform_admin', 'hospital_admin',
             'doctor', 'front_desk')
  ),

  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  licence_number TEXT,
  specialty TEXT,

  status TEXT NOT NULL DEFAULT 'active' CHECK (
    status IN ('active', 'suspended', 'removed')
  ),

  totp_secret TEXT,
  totp_enabled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  CHECK ((role = 'platform_admin') = (facility_id IS NULL))
);

CREATE INDEX ON staff (facility_id);

CREATE TABLE invites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  facility_id UUID NOT NULL REFERENCES facilities(id),

  email TEXT NOT NULL,

  role TEXT NOT NULL CHECK (
    role IN ('hospital_admin', 'doctor', 'front_desk')
  ),

  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  licence_number TEXT,
  specialty TEXT,

  invited_by UUID REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  accepted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
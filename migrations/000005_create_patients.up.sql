CREATE TABLE patients (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(id),
  network_id TEXT NOT NULL UNIQUE,

  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  date_of_birth DATE NOT NULL,
  sex TEXT NOT NULL CHECK (sex IN ('male', 'female')),

  nin_hash TEXT UNIQUE,
  state TEXT,
  address TEXT,
  emergency_contact_name TEXT,
  emergency_contact_phone TEXT,

  terms_accepted_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
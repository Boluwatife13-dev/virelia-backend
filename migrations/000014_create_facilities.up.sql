CREATE TABLE facilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  name TEXT NOT NULL,

  facility_type TEXT NOT NULL CHECK (
    facility_type IN ('hospital', 'clinic', 'primary_health_centre')
  ),

  licence_number TEXT NOT NULL,
  licence_expires_at DATE NOT NULL,

  state TEXT NOT NULL,
  address TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,

  admin_first_name TEXT NOT NULL,
  admin_last_name TEXT NOT NULL,
  admin_email TEXT NOT NULL,

  status TEXT NOT NULL DEFAULT 'pending' CHECK (
    status IN ('pending', 'approved', 'rejected', 'suspended')
  ),

  rejection_reason TEXT,
  reviewed_by UUID REFERENCES users(id),
  reviewed_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON facilities (status);
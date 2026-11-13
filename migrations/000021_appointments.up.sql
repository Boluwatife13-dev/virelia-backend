CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE doctor_profiles (
  user_id UUID PRIMARY KEY
    REFERENCES staff(user_id) ON DELETE CASCADE,
  fee_naira INT NOT NULL DEFAULT 0 CHECK (fee_naira >= 0),
  slot_minutes INT NOT NULL DEFAULT 30
    CHECK (slot_minutes IN (15, 20, 30, 45, 60)),
  bio TEXT,
  accepting BOOLEAN NOT NULL DEFAULT true,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE doctor_hours (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  doctor_id UUID NOT NULL
    REFERENCES staff(user_id) ON DELETE CASCADE,
  weekday SMALLINT NOT NULL CHECK (weekday BETWEEN 0 AND 6),
  starts TIME NOT NULL,
  ends TIME NOT NULL,
  CHECK (starts < ends),
  UNIQUE (doctor_id, weekday, starts)
);

CREATE TABLE appointments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL
    REFERENCES patients(id) ON DELETE CASCADE,
  doctor_id UUID NOT NULL REFERENCES staff(user_id),
  facility_id UUID NOT NULL REFERENCES facilities(id),
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL DEFAULT 'booked' CHECK (
    status IN (
      'booked', 'checked_in', 'completed',
      'cancelled', 'no_show'
    )
  ),
  reason TEXT,
  booked_by UUID REFERENCES users(id),
  cancelled_by UUID REFERENCES users(id),
  cancel_reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (starts_at < ends_at),

  EXCLUDE USING gist (
    doctor_id WITH =,
    tstzrange(starts_at, ends_at) WITH &&
  ) WHERE (status <> 'cancelled'),

  EXCLUDE USING gist (
    patient_id WITH =,
    tstzrange(starts_at, ends_at) WITH &&
  ) WHERE (status <> 'cancelled')
);

CREATE INDEX ON appointments (doctor_id, starts_at);
CREATE INDEX ON appointments (patient_id, starts_at);
CREATE INDEX ON appointments (facility_id, starts_at);
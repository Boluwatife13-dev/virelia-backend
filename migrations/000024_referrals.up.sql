CREATE TABLE referrals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL REFERENCES patients(id),
  from_facility_id UUID NOT NULL REFERENCES facilities(id),
  from_doctor_id UUID NOT NULL REFERENCES staff(user_id),
  to_facility_id UUID NOT NULL REFERENCES facilities(id),
  specialty TEXT NOT NULL,
  urgency TEXT NOT NULL CHECK (urgency IN ('urgent', 'routine')),
  reason TEXT NOT NULL,
  question TEXT,
  status TEXT NOT NULL DEFAULT 'awaiting_consent' CHECK (
    status IN (
      'awaiting_consent', 'awaiting_sender', 'sent',
      'accepted', 'declined', 'refused', 'cancelled'
    )
  ),
  share_records BOOLEAN NOT NULL DEFAULT false,
  sender_id UUID REFERENCES users(id),
  sent_at TIMESTAMPTZ,
  respond_by TIMESTAMPTZ,
  reminded_at TIMESTAMPTZ,
  decided_by UUID REFERENCES users(id),
  decided_at TIMESTAMPTZ,
  decline_reason TEXT,
  appointment_id UUID REFERENCES appointments(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (from_facility_id <> to_facility_id),
  CHECK (status <> 'accepted' OR appointment_id IS NOT NULL),
  CHECK (status <> 'declined' OR decline_reason IS NOT NULL)
);

CREATE INDEX ON referrals (patient_id, created_at);
CREATE INDEX ON referrals (from_facility_id, status);
CREATE INDEX ON referrals (to_facility_id, status);

CREATE TABLE referral_records (
  referral_id UUID NOT NULL REFERENCES referrals(id) ON DELETE CASCADE,
  entry_id UUID NOT NULL REFERENCES record_entries(id),
  PRIMARY KEY (referral_id, entry_id)
);

CREATE TABLE transfers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  patient_id UUID NOT NULL REFERENCES patients(id),
  requesting_facility_id UUID NOT NULL REFERENCES facilities(id),
  holding_facility_id UUID NOT NULL REFERENCES facilities(id),
  requested_by UUID NOT NULL REFERENCES users(id),
  reason TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'awaiting_consent' CHECK (
    status IN (
      'awaiting_consent', 'awaiting_holder', 'sent',
      'declined', 'refused', 'cancelled'
    )
  ),
  respond_by TIMESTAMPTZ,
  reminded_at TIMESTAMPTZ,
  decided_by UUID REFERENCES users(id),
  decided_at TIMESTAMPTZ,
  decline_reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (requesting_facility_id <> holding_facility_id),
  CHECK (status <> 'declined' OR decline_reason IS NOT NULL)
);

CREATE INDEX ON transfers (patient_id, created_at);
CREATE INDEX ON transfers (holding_facility_id, status);
CREATE INDEX ON transfers (requesting_facility_id, status);

CREATE TABLE transfer_records (
  transfer_id UUID NOT NULL REFERENCES transfers(id) ON DELETE CASCADE,
  entry_id UUID NOT NULL REFERENCES record_entries(id),
  PRIMARY KEY (transfer_id, entry_id)
);

CREATE TABLE consultation_records (
  appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
  entry_id UUID NOT NULL REFERENCES record_entries(id),
  attached_by UUID NOT NULL REFERENCES users(id),
  attached_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (appointment_id, entry_id)
);
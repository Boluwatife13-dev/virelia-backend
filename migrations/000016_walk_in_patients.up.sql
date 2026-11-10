ALTER TABLE patients
  ALTER COLUMN user_id DROP NOT NULL;


ALTER TABLE patients
  ALTER COLUMN terms_accepted_at DROP NOT NULL;


ALTER TABLE patients
  ADD COLUMN phone TEXT;


ALTER TABLE patients
  ADD COLUMN registered_by UUID REFERENCES users(id);

ALTER TABLE patients
  ADD COLUMN registered_facility_id UUID
  REFERENCES facilities(id);


ALTER TABLE patients
  ADD CONSTRAINT patients_contact
  CHECK (user_id IS NOT NULL OR phone IS NOT NULL);

CREATE INDEX patients_phone_idx ON patients (phone);
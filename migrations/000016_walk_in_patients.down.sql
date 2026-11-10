
DELETE FROM patients WHERE user_id IS NULL;

UPDATE patients
SET terms_accepted_at = created_at
WHERE terms_accepted_at IS NULL;

ALTER TABLE patients DROP CONSTRAINT patients_contact;

ALTER TABLE patients DROP COLUMN registered_facility_id;
ALTER TABLE patients DROP COLUMN registered_by;
ALTER TABLE patients DROP COLUMN phone;

ALTER TABLE patients ALTER COLUMN terms_accepted_at SET NOT NULL;
ALTER TABLE patients ALTER COLUMN user_id SET NOT NULL;
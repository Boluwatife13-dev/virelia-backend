DROP TABLE guardians;

DELETE FROM patients
WHERE user_id IS NULL AND phone IS NULL;

ALTER TABLE patients
  ADD CONSTRAINT patients_contact
  CHECK (user_id IS NOT NULL OR phone IS NOT NULL);
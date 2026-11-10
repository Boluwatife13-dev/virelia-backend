ALTER TABLE facilities
  DROP CONSTRAINT facilities_two_person;

UPDATE facilities
SET status = 'pending'
WHERE status = 'checked';

ALTER TABLE facilities DROP COLUMN check_note;
ALTER TABLE facilities DROP COLUMN checked_at;
ALTER TABLE facilities DROP COLUMN checked_by;

ALTER TABLE facilities
  DROP CONSTRAINT facilities_status_check;

ALTER TABLE facilities
  ADD CONSTRAINT facilities_status_check CHECK (
    status IN ('pending', 'approved', 'rejected', 'suspended')
  );

ALTER TABLE facilities DROP CONSTRAINT facilities_status_check;
ALTER TABLE facilities ADD CONSTRAINT facilities_status_check CHECK (status IN ('pending', 'checked', 'approved', 'rejected', 'suspended'));
ALTER TABLE facilities ADD COLUMN checked_by UUID REFERENCES users(id);
ALTER TABLE facilities ADD COLUMN checked_at TIMESTAMPTZ;
ALTER TABLE facilities ADD COLUMN check_note TEXT;
ALTER TABLE facilities ADD CONSTRAINT facilities_two_person CHECK (status <> 'approved' OR (checked_by IS NOT NULL AND reviewed_by IS NOT NULL AND checked_by <> reviewed_by)) NOT VALID;

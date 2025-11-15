-- Restore verification_status and rejection_reason to organizers
ALTER TABLE organizers ADD COLUMN verification_status TEXT DEFAULT 'pending';
ALTER TABLE organizers ADD COLUMN rejection_reason TEXT;

-- Rename about back to contacts in organizers table
ALTER TABLE organizers RENAME COLUMN about TO contacts;

-- Rename about back to cv in volunteers table
ALTER TABLE volunteers RENAME COLUMN about TO cv;

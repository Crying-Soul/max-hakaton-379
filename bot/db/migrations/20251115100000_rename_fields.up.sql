-- Rename cv to about in volunteers table
ALTER TABLE volunteers RENAME COLUMN cv TO about;

-- Rename contacts to about in organizers table
ALTER TABLE organizers RENAME COLUMN contacts TO about;

-- Drop verification_status and rejection_reason from organizers
-- (since they're now tracked in organizer_verification_requests table)
ALTER TABLE organizers DROP COLUMN IF EXISTS verification_status;
ALTER TABLE organizers DROP COLUMN IF EXISTS rejection_reason;

-- Organizer verification history keeps every submission/change
CREATE TABLE organizer_verification_requests (
    id SERIAL PRIMARY KEY,
    organizer_id BIGINT NOT NULL REFERENCES organizers(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    organizer_comment TEXT,
    admin_comment TEXT,
    reviewed_by BIGINT REFERENCES admins(id),
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP
);

CREATE INDEX idx_organizer_verification_requests_organizer ON organizer_verification_requests(organizer_id);
CREATE INDEX idx_organizer_verification_requests_status ON organizer_verification_requests(status);

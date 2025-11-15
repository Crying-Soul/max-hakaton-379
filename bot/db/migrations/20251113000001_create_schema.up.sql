-- Enable useful extensions for spatial indexing
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Users represent all bot participants
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username TEXT,
    name TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('volunteer', 'organizer', 'admin')),
    state TEXT NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    location_lat DECIMAL(10, 8),
    location_lon DECIMAL(11, 8)
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_state ON users(state);
CREATE INDEX idx_users_is_blocked ON users(is_blocked);
CREATE INDEX idx_users_location_gist ON users USING gist (location_lon, location_lat);

-- Admins table links administrators back to users
CREATE TABLE admins (
    id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Volunteers store volunteer-specific information
CREATE TABLE volunteers (
    id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    cv TEXT,
    search_radius INT DEFAULT 10,
    category_ids INT[]
);

CREATE INDEX idx_volunteers_category_ids ON volunteers USING gin (category_ids);

-- Organizers with verification details
CREATE TABLE organizers (
    id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    organization_name TEXT NOT NULL,
    verification_status TEXT DEFAULT 'pending',
    rejection_reason TEXT,
    contacts TEXT,
    verified_at TIMESTAMP,
    verified_by BIGINT REFERENCES admins(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_organizers_verification_status ON organizers(verification_status);

-- Categories for events
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_categories_is_active ON categories(is_active);

-- Events organized by organizers
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    chat BIGINT,
    date TIMESTAMP NOT NULL,
    duration_hours INT,
    location TEXT NOT NULL,
    location_lat DECIMAL(10, 8) NOT NULL,
    location_lon DECIMAL(11, 8) NOT NULL,
    category_id INT REFERENCES categories(id),
    organizer_id BIGINT REFERENCES organizers(id),
    contacts TEXT,
    max_volunteers INT NOT NULL,
    current_volunteers INT DEFAULT 0,
    status TEXT DEFAULT 'open',
    cancelled_reason TEXT,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_organizer_id ON events(organizer_id);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_date ON events(date);
CREATE INDEX idx_events_category_id ON events(category_id);
CREATE INDEX idx_events_location_gist ON events USING gist (location_lon, location_lat);
CREATE INDEX idx_events_status_date ON events(status, date);

-- Applications that volunteers submit to events
CREATE TABLE volunteer_applications (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    volunteer_id BIGINT REFERENCES volunteers(id) ON DELETE CASCADE,
    status TEXT DEFAULT 'pending',
    rejection_reason TEXT,
    reviewed_by BIGINT REFERENCES organizers(id),
    applied_at TIMESTAMP DEFAULT NOW(),
    reviewed_at TIMESTAMP
);

CREATE INDEX idx_volunteer_applications_event_id ON volunteer_applications(event_id);
CREATE INDEX idx_volunteer_applications_volunteer_id ON volunteer_applications(volunteer_id);
CREATE INDEX idx_volunteer_applications_status ON volunteer_applications(status);
CREATE INDEX idx_volunteer_applications_volunteer_status ON volunteer_applications(volunteer_id, status);
CREATE INDEX idx_volunteer_applications_event_status ON volunteer_applications(event_id, status);
CREATE UNIQUE INDEX idx_volunteer_applications_event_volunteer_unique ON volunteer_applications(event_id, volunteer_id);

-- Participants that joined the event after approval
CREATE TABLE event_participants (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    volunteer_id BIGINT REFERENCES volunteers(id) ON DELETE CASCADE,
    application_id INT REFERENCES volunteer_applications(id) ON DELETE SET NULL,
    joined_chat_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_event_participants_event_id ON event_participants(event_id);
CREATE INDEX idx_event_participants_volunteer_id ON event_participants(volunteer_id);
CREATE UNIQUE INDEX idx_event_participants_unique ON event_participants(event_id, volunteer_id);

-- Media uploaded for events
CREATE TABLE event_media (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES events(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT NOW(),
    uploaded_by BIGINT REFERENCES users(id)
);

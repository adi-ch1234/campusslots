CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workshops (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    location TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    capacity INT NOT NULL CHECK (capacity > 0),
    available_seats INT NOT NULL CHECK (available_seats >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_workshops_available_lte_capacity CHECK (available_seats <= capacity)
);

CREATE TABLE registrations (
    id BIGSERIAL PRIMARY KEY,
    workshop_id BIGINT NOT NULL REFERENCES workshops(id),
    student_id BIGINT NOT NULL REFERENCES students(id),
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'CANCELLED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ,
    CONSTRAINT chk_registrations_cancelled_at_consistency CHECK (
        (status = 'ACTIVE' AND cancelled_at IS NULL)
        OR (status = 'CANCELLED' AND cancelled_at IS NOT NULL)
    )
);

-- Guarantees a student can have at most ONE active registration per workshop,
-- while still keeping cancelled history and allowing re-registration.
CREATE UNIQUE INDEX uq_registrations_active_student_workshop
    ON registrations (student_id, workshop_id)
    WHERE status = 'ACTIVE';

-- Speeds up lookups for "list a student's registrations" and workshop audits.
CREATE INDEX idx_registrations_student_id ON registrations (student_id);
CREATE INDEX idx_registrations_workshop_id ON registrations (workshop_id);
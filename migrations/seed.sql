INSERT INTO students (id, display_name) VALUES
    (1, 'James Bond'),
    (2, 'Harry Potter'),
    (3, 'Padme Amidala'),
    (4, 'Christopher Nolan'),
    (5, 'Golang Go')
ON CONFLICT (id) DO NOTHING;

-- Advance the auto-increment sequence so future inserts don't collide with seeded IDs
SELECT setval('students_id_seq', COALESCE((SELECT MAX(id) FROM students), 1));

INSERT INTO workshops (id, title, location, starts_at, status, capacity, available_seats) VALUES
    (1, 'Intro to Go', 'Lecture Hall A', NOW() + INTERVAL '7 days', 'OPEN', 10, 10),
    (2, 'PostgreSQL Transactions Deep Dive', 'Lab 204', NOW() + INTERVAL '10 days', 'OPEN', 25, 25)
ON CONFLICT (id) DO NOTHING;

SELECT setval('workshops_id_seq', COALESCE((SELECT MAX(id) FROM workshops), 1));
-- ################## Create tables ##################

CREATE TABLE snippets (
	id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created TIMESTAMPTZ NOT NULL,
	expires TIMESTAMPTZ NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMPTZ NOT NULL
);

-- ################## Insert new rows ##################

INSERT INTO snippets (title, content, created, expires) VALUES (
	'An old slient pond',
	'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence
again.\n\n– Matsuo Bashō',
	NOW(),
	NOW() + INTERVAL '10 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
	'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to
blow.\n\n– Natsume Soseki',
	NOW(),
	NOW() + INTERVAL '20 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
	'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s
face.\n\n– Murakami Kijo',
	NOW(),
	NOW() + INTERVAL '15 days'
);

-- ################## Create a new user and grant permissions ##################

CREATE USER web WITH PASSWORD '3c523592-852d-42be-915c-d5931792e39e';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO web;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO web;

package database

const schema = `
CREATE TABLE IF NOT EXISTS accounts (
	id TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE COLLATE NOCASE,
	name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	role TEXT NOT NULL CHECK (role IN ('member', 'staff')),
	validation_status TEXT NOT NULL CHECK (validation_status IN ('unvalidated', 'validated')),
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	csrf_token TEXT NOT NULL,
	created_at TEXT NOT NULL,
	expires_at TEXT NOT NULL,
	revoked_at TEXT
);
CREATE INDEX IF NOT EXISTS sessions_account_id ON sessions(account_id);

CREATE TABLE IF NOT EXISTS posts (
	id TEXT PRIMARY KEY,
	author_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	body TEXT NOT NULL,
	created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS posts_created_at ON posts(created_at DESC);

CREATE TABLE IF NOT EXISTS likes (
	post_id TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
	account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	created_at TEXT NOT NULL,
	PRIMARY KEY (post_id, account_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
	id TEXT PRIMARY KEY,
	actor_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	target_account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	action TEXT NOT NULL,
	from_status TEXT,
	to_status TEXT,
	created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS audit_logs_created_at ON audit_logs(created_at DESC);
`

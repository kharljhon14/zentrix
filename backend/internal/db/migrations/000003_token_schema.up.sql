CREATE TABLE IF NOT EXISTS tokens(
    hash bytea PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users on DELETE CASCADE,
    expiry TIMESTAMPTZ NOT NULL,
    scope text NOT NULL
);
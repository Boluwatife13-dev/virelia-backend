-- Hardware keys (or Windows Hello) for platform admins.
CREATE TABLE security_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL
    REFERENCES users(id) ON DELETE CASCADE,

  credential_id BYTEA NOT NULL UNIQUE,

  -- The key's public half and counters, as JSON.
  credential JSONB NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_used_at TIMESTAMPTZ
);

CREATE INDEX ON security_keys (user_id);

-- A short-lived challenge between the two steps
-- of registering or using a key.
CREATE TABLE key_challenges (
  user_id UUID NOT NULL
    REFERENCES users(id) ON DELETE CASCADE,

  purpose TEXT NOT NULL
    CHECK (purpose IN ('register', 'login')),

  data JSONB NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,

  PRIMARY KEY (user_id, purpose)
);

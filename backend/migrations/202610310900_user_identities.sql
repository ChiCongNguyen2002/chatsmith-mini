CREATE TABLE user_identities (
    provider       text        NOT NULL,
    subject        text        NOT NULL,
    user_id        uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    email          text,
    email_verified boolean     NOT NULL DEFAULT false,
    created_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (provider, subject)
);

CREATE INDEX user_identities_verified_email ON user_identities (email) WHERE email_verified;

INSERT INTO user_identities (provider, subject, user_id, email, email_verified)
SELECT 'google', google_sub, id, email, true
FROM users
WHERE google_sub IS NOT NULL;

ALTER TABLE users DROP COLUMN google_sub;

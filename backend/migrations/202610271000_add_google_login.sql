ALTER TABLE users
    ADD COLUMN google_sub text UNIQUE,
    ADD COLUMN email      text,
    ADD COLUMN name       text,
    ADD COLUMN picture    text;

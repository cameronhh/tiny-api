CREATE TABLE IF NOT EXISTS endpoints
(
    id SERIAL,
    url TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT 0.00,
    CONSTRAINT endpoints_pkey PRIMARY KEY (id)
)
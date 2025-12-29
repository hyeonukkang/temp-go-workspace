CREATE TABLE IF NOT EXISTS file_info (
    id SERIAL PRIMARY KEY,
    file_name TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL
);

INSERT INTO file_info (file_name, uploaded_at) VALUES ('sample.txt', NOW());

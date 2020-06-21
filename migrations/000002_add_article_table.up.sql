CREATE TABLE article (
    id serial PRIMARY KEY,
    title varchar(127) NOT NULL DEFAULT '',
    content text NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

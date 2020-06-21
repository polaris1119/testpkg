CREATE TABLE gopher (
    id serial PRIMARY KEY,
    username varchar(31) NOT NULL DEFAULT '',
    email varchar(63) NOT NULL DEFAULT '',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

comment on table gopher is 'gopher用户表';
comment on column gopher.username is '用户名';

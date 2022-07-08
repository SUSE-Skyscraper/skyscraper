-- +goose Up
create table policies
(
    id    uuid         not null primary key,
    ptype varchar(255) not null,
    v0    varchar(255) not null,
    v1    varchar(255) not null,
    v2    varchar(255) null default null,
    v3    varchar(255) null default null,
    v4    varchar(255) null default null,
    v5    varchar(255) null default null,
    unique (ptype, v0, v1, v2, v3, v4, v5)
);

insert into policies (id, ptype, v0, v1, v2)
values (uuid_generate_v5('6ba7b812-9dad-11d1-80b4-00c04fd430c8', concat('p', 'readonly', '/*', 'GET')), 'p', 'readonly',
        '/*', 'GET')
on conflict do nothing;

insert into policies (id, ptype, v0, v1, v2)
values (uuid_generate_v5('6ba7b812-9dad-11d1-80b4-00c04fd430c8', concat('p', 'admin', '/*', '(.*)')), 'p', 'admin',
        '/*', '(.*)')
on conflict do nothing;

insert into policies (id, ptype, v0, v1, v2)
values (uuid_generate_v5('6ba7b812-9dad-11d1-80b4-00c04fd430c8', concat('p', 'contributor', '/*', '(.*)')), 'p',
        'contributor', '/*', '(.*)')
on conflict do nothing;

-- +goose Down
drop table policies;

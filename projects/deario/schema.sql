-- auto-generated definition
create table diarys
(
    content    TEXT default ''                                 not null,
    created    TEXT default ''                                 not null,
    date       TEXT default ''                                 not null,
    id         TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    uid        TEXT default ''                                 not null,
    updated    TEXT default ''                                 not null,
    aiFeedback TEXT default ''                                 not null
);

-- auto-generated definition
create table push_keys
(
    created TEXT default ''                                 not null,
    id      TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    token   TEXT default ''                                 not null,
    uid     TEXT default ''                                 not null,
    updated TEXT default ''                                 not null
);

-- auto-generated definition
create table diary_settings
(
    id                TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    uid               TEXT default ''                                 not null
        unique,
    random_range_days INTEGER default 30                              not null,
    created           TEXT default ''                                 not null,
    updated           TEXT default ''                                 not null
);

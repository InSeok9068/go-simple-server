-- auto-generated definition
create table diarys (
    id TEXT default (
        'r' || lower(hex(randomblob(7)))
    ) not null primary key,
    uid TEXT default '' not null,
    date TEXT default '' not null,
    content TEXT default '' not null,
    aiFeedback TEXT default '' not null,
    aiImage TEXT default '' not null,
    created TEXT default '' not null,
    updated TEXT default '' not null
);

-- auto-generated definition
create table push_keys (
    uid TEXT default '' not null,
    id TEXT default (
        'r' || lower(hex(randomblob(7)))
    ) not null primary key,
    token TEXT default '' not null,
    created TEXT default '' not null,
    updated TEXT default '' not null
);
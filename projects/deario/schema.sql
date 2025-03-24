-- auto-generated definition
create table diarys
(
    content TEXT default ''                                 not null,
    created TEXT default ''                                 not null,
    date    TEXT default ''                                 not null,
    id      TEXT default ('r' || lower(hex(randomblob(7)))) not null
        primary key,
    uid     TEXT default ''                                 not null,
    updated TEXT default ''                                 not null
);
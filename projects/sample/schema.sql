CREATE TABLE authors
(
    bio     TEXT             DEFAULT ''                                 NOT NULL,
    created TEXT             DEFAULT ''                                 NOT NULL,
    id      TEXT PRIMARY KEY DEFAULT ('r' || lower(hex(randomblob(7)))) NOT NULL,
    name    TEXT             DEFAULT ''                                 NOT NULL,
    updated TEXT             DEFAULT ''                                 NOT NULL
);
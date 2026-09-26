CREATE TABLE dns_records(
    domain     VARCHAR(255) PRIMARY KEY,
    ip         VARCHAR(45) NOT NULL,
    created_at VARCHAR(32) NOT NULL DEFAULT '',
    updated_at VARCHAR(32) NOT NULL DEFAULT ''
);
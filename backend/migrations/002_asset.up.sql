CREATE TABLE gm_catalog (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    parent_id   VARCHAR(32),
    name        VARCHAR(128) NOT NULL,
    sort_code   INT DEFAULT 0,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL
);

CREATE TABLE gm_catalog_config (
    id           VARCHAR(32) PRIMARY KEY,
    tenant_id    VARCHAR(32) NOT NULL DEFAULT 'default',
    catalog_id   VARCHAR(32) NOT NULL,
    config_key   VARCHAR(128) NOT NULL,
    config_value TEXT,
    created_at   DATETIME(3) NOT NULL,
    updated_at   DATETIME(3) NOT NULL,
    deleted_at   DATETIME(3) NULL,
    UNIQUE KEY uk_catalog_config_key (catalog_id, config_key)
);

CREATE TABLE gm_asset (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    title       VARCHAR(256) NOT NULL,
    type        VARCHAR(32) NOT NULL,
    status      TINYINT NOT NULL DEFAULT 0,
    catalog_id  VARCHAR(32),
    description TEXT,
    created_by  VARCHAR(32) NOT NULL,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL,
    KEY idx_asset_tenant_catalog (tenant_id, catalog_id),
    KEY idx_asset_tenant_status (tenant_id, status)
);

CREATE TABLE gm_asset_file (
    id            VARCHAR(32) PRIMARY KEY,
    tenant_id     VARCHAR(32) NOT NULL DEFAULT 'default',
    asset_id      VARCHAR(32) NOT NULL,
    version       INT NOT NULL DEFAULT 1,
    storage_path  VARCHAR(512) NOT NULL,
    mime_type     VARCHAR(128),
    file_size     BIGINT,
    checksum      VARCHAR(64),
    is_master     TINYINT NOT NULL DEFAULT 1,
    created_at    DATETIME(3) NOT NULL,
    updated_at    DATETIME(3) NOT NULL,
    deleted_at    DATETIME(3) NULL,
    KEY idx_asset_file_asset (asset_id)
);

CREATE TABLE gm_asset_metadata (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    asset_id    VARCHAR(32) NOT NULL,
    meta_key    VARCHAR(128) NOT NULL,
    meta_value  TEXT,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    UNIQUE KEY uk_asset_meta_key (asset_id, meta_key)
);

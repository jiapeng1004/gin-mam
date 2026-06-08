CREATE TABLE gm_transcode_group (
    id              VARCHAR(32) PRIMARY KEY,
    tenant_id       VARCHAR(32) NOT NULL DEFAULT 'default',
    name            VARCHAR(128) NOT NULL,
    message         TEXT,
    group_type      INT NOT NULL DEFAULT 0,
    param           TEXT,
    strategy_type   INT NOT NULL DEFAULT 0,
    default_flag    TINYINT NOT NULL DEFAULT 0,
    available_flag  TINYINT NOT NULL DEFAULT 1,
    created_by      VARCHAR(32) NOT NULL,
    created_at      DATETIME(3) NOT NULL,
    updated_at      DATETIME(3) NOT NULL,
    deleted_at      DATETIME(3) NULL
);

CREATE TABLE gm_transcode_profile (
    id               VARCHAR(32) PRIMARY KEY,
    tenant_id        VARCHAR(32) NOT NULL DEFAULT 'default',
    name             VARCHAR(128) NOT NULL,
    alias            VARCHAR(128),
    transcode_type   INT NOT NULL DEFAULT 0,
    definition_type  INT NOT NULL DEFAULT 0,
    param            TEXT,
    special_param    TEXT,
    created_at       DATETIME(3) NOT NULL,
    updated_at       DATETIME(3) NOT NULL,
    deleted_at       DATETIME(3) NULL
);

CREATE TABLE gm_transcode_task (
    id                  VARCHAR(32) PRIMARY KEY,
    tenant_id           VARCHAR(32) NOT NULL DEFAULT 'default',
    asset_id            VARCHAR(32) NOT NULL,
    asset_file_id       VARCHAR(32),
    transcode_group_id  VARCHAR(32) NOT NULL,
    profile_id          VARCHAR(32),
    status              TINYINT NOT NULL DEFAULT 0,
    external_job_id     VARCHAR(128),
    output_path         VARCHAR(512),
    error_msg           TEXT,
    callback_payload    TEXT,
    created_by          VARCHAR(32) NOT NULL,
    created_at          DATETIME(3) NOT NULL,
    updated_at          DATETIME(3) NOT NULL,
    KEY idx_transcode_task_asset (tenant_id, asset_id),
    KEY idx_transcode_task_status (tenant_id, status)
);
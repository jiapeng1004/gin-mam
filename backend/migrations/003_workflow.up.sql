CREATE TABLE gm_workflow_def (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    name        VARCHAR(128) NOT NULL,
    audit_level INT NOT NULL,
    status      TINYINT NOT NULL DEFAULT 1,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL
);

CREATE TABLE gm_workflow_level_user (
    id               VARCHAR(32) PRIMARY KEY,
    tenant_id        VARCHAR(32) NOT NULL DEFAULT 'default',
    workflow_def_id  VARCHAR(32) NOT NULL,
    level            INT NOT NULL,
    user_id          VARCHAR(32) NOT NULL,
    UNIQUE KEY uk_workflow_level_user (workflow_def_id, level, user_id)
);

CREATE TABLE gm_workflow_instance (
    id               VARCHAR(32) PRIMARY KEY,
    tenant_id        VARCHAR(32) NOT NULL DEFAULT 'default',
    workflow_def_id  VARCHAR(32) NOT NULL,
    asset_id         VARCHAR(32) NOT NULL,
    level            INT NOT NULL DEFAULT 0,
    audit_level      INT NOT NULL,
    audit_status     TINYINT NOT NULL DEFAULT 0,
    locker_id        VARCHAR(32) NULL,
    created_by       VARCHAR(32) NOT NULL,
    created_at       DATETIME(3) NOT NULL,
    updated_at       DATETIME(3) NOT NULL,
    deleted_at       DATETIME(3) NULL,
    KEY idx_workflow_instance_asset (tenant_id, asset_id),
    KEY idx_workflow_instance_status (tenant_id, audit_status)
);

CREATE TABLE gm_workflow_instance_level_user (
    id           VARCHAR(32) PRIMARY KEY,
    tenant_id    VARCHAR(32) NOT NULL DEFAULT 'default',
    instance_id  VARCHAR(32) NOT NULL,
    level        INT NOT NULL,
    user_id      VARCHAR(32) NOT NULL,
    audit_status TINYINT NOT NULL DEFAULT 0,
    UNIQUE KEY uk_instance_level_user (instance_id, level, user_id)
);

CREATE TABLE gm_workflow_operate (
    id           VARCHAR(32) PRIMARY KEY,
    tenant_id    VARCHAR(32) NOT NULL DEFAULT 'default',
    instance_id  VARCHAR(32) NOT NULL,
    operator_id  VARCHAR(32) NOT NULL,
    action       TINYINT NOT NULL,
    remark       TEXT,
    created_at   DATETIME(3) NOT NULL,
    KEY idx_workflow_operate_instance (instance_id)
);
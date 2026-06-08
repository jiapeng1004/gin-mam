CREATE TABLE gm_user (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    username    VARCHAR(64) NOT NULL,
    password    VARCHAR(128) NOT NULL,
    nickname    VARCHAR(64),
    org_id      VARCHAR(32),
    status      TINYINT NOT NULL DEFAULT 1,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL,
    UNIQUE KEY uk_tenant_username (tenant_id, username)
);

CREATE TABLE gm_role (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    name        VARCHAR(64) NOT NULL,
    code        VARCHAR(64) NOT NULL,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL,
    UNIQUE KEY uk_tenant_code (tenant_id, code)
);

CREATE TABLE gm_org (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    parent_id   VARCHAR(32),
    name        VARCHAR(128) NOT NULL,
    sort_code   INT DEFAULT 0,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL
);

CREATE TABLE gm_menu (
    id          VARCHAR(32) PRIMARY KEY,
    tenant_id   VARCHAR(32) NOT NULL DEFAULT 'default',
    parent_id   VARCHAR(32),
    title       VARCHAR(128) NOT NULL,
    path        VARCHAR(256),
    component   VARCHAR(256),
    permission  VARCHAR(128),
    type        TINYINT NOT NULL,
    sort_code   INT DEFAULT 0,
    created_at  DATETIME(3) NOT NULL,
    updated_at  DATETIME(3) NOT NULL,
    deleted_at  DATETIME(3) NULL
);

CREATE TABLE gm_sys_config (
    id           VARCHAR(32) PRIMARY KEY,
    tenant_id    VARCHAR(32) NOT NULL DEFAULT 'default',
    config_key   VARCHAR(128) NOT NULL,
    config_value TEXT,
    category     VARCHAR(64),
    remark       VARCHAR(256),
    sort_code    INT DEFAULT 0,
    created_at   DATETIME(3) NOT NULL,
    updated_at   DATETIME(3) NOT NULL,
    deleted_at   DATETIME(3) NULL,
    UNIQUE KEY uk_tenant_key (tenant_id, config_key)
);

CREATE TABLE gm_user_role (
    user_id VARCHAR(32) NOT NULL,
    role_id VARCHAR(32) NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE gm_role_menu (
    role_id VARCHAR(32) NOT NULL,
    menu_id VARCHAR(32) NOT NULL,
    PRIMARY KEY (role_id, menu_id)
);
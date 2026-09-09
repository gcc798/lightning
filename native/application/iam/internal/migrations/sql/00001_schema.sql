-- +goose Up
CREATE TABLE s_api_permission (
  id bigint PRIMARY KEY, parent_id bigint DEFAULT 0, module varchar(64) NOT NULL,
  code varchar(128) NOT NULL, name varchar(64) NOT NULL, node_type smallint DEFAULT 2 NOT NULL,
  action varchar(32) DEFAULT '*' NOT NULL, method varchar(16), path varchar(255), sort bigint DEFAULT 0,
  status smallint DEFAULT 0, remark varchar(500), create_by bigint, update_by bigint,
  created_time timestamptz, updated_time timestamptz
);
CREATE UNIQUE INDEX idx_s_api_permission_code ON s_api_permission(code);
CREATE INDEX idx_s_api_permission_module ON s_api_permission(module);
CREATE INDEX idx_s_api_permission_parent_id ON s_api_permission(parent_id);

CREATE TABLE s_auth_client (
  client_id varchar(64) PRIMARY KEY, grant_type varchar(255), device_type varchar(32), status smallint DEFAULT 0,
  timeout bigint DEFAULT 604800, active_timeout bigint DEFAULT 1800, remark varchar(500), create_by bigint,
  created_time timestamptz, update_by bigint, updated_time timestamptz
);
CREATE TABLE s_menu (
  id bigint PRIMARY KEY, menu_name varchar(64) NOT NULL, parent_id bigint DEFAULT 0, sort bigint DEFAULT 0,
  path varchar(255), component varchar(255), query varchar(255), is_frame smallint DEFAULT 0,
  is_cache smallint DEFAULT 0, menu_type smallint NOT NULL, visible smallint DEFAULT 0, status smallint DEFAULT 0,
  perms varchar(255), icon varchar(64), remark varchar(500), create_by bigint, update_by bigint,
  created_time timestamptz, updated_time timestamptz
);
CREATE INDEX idx_s_menu_parent_id ON s_menu(parent_id);
CREATE TABLE s_org (
  id bigint PRIMARY KEY, parent_id bigint DEFAULT 0, ancestors varchar(512), org_name varchar(128) NOT NULL,
  org_code varchar(64), org_type varchar(32) DEFAULT 'company', leader varchar(64), phone varchar(32),
  email varchar(128), status smallint DEFAULT 0, sort bigint DEFAULT 0, remark varchar(500), create_by bigint,
  update_by bigint, created_time timestamptz, updated_time timestamptz
);
CREATE UNIQUE INDEX idx_s_org_org_code ON s_org(org_code);
CREATE INDEX idx_s_org_parent_id ON s_org(parent_id);
CREATE TABLE s_role (
  id bigint PRIMARY KEY, role_key varchar(64) NOT NULL, role_name varchar(64) NOT NULL, sort bigint DEFAULT 0,
  status smallint DEFAULT 0, data_scope smallint DEFAULT 1, is_system boolean DEFAULT false,
  remark varchar(500), create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz
);
CREATE UNIQUE INDEX idx_s_role_role_key ON s_role(role_key);
CREATE TABLE s_user (
  id bigint PRIMARY KEY, user_name varchar(64) NOT NULL, nick_name varchar(64), user_type smallint DEFAULT 0,
  org_id bigint DEFAULT 0, email varchar(128), phonenumber varchar(32), sex smallint DEFAULT 2,
  avatar varchar(512), password varchar(255), status smallint DEFAULT 0, sort bigint DEFAULT 0,
  login_ip varchar(64), login_date bigint, open_id varchar(128), union_id varchar(128), remark varchar(500),
  create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz
);
CREATE UNIQUE INDEX idx_s_user_user_name ON s_user(user_name);
CREATE UNIQUE INDEX idx_s_user_email ON s_user(email) WHERE email <> '';
CREATE UNIQUE INDEX idx_s_user_phonenumber ON s_user(phonenumber) WHERE phonenumber <> '';
CREATE UNIQUE INDEX idx_s_user_open_id ON s_user(open_id) WHERE open_id <> '';
CREATE UNIQUE INDEX idx_s_user_union_id ON s_user(union_id) WHERE union_id <> '';
CREATE INDEX idx_s_user_org_id ON s_user(org_id);

CREATE TABLE m_role_api_permission (id bigint PRIMARY KEY, role_id bigint NOT NULL, permission_id bigint NOT NULL, source smallint NOT NULL DEFAULT 0, create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz);
CREATE UNIQUE INDEX idx_role_api_permission ON m_role_api_permission(role_id, permission_id, source);
CREATE TABLE m_user_api_permission (id bigint PRIMARY KEY, user_id bigint NOT NULL, permission_id bigint NOT NULL, source smallint NOT NULL DEFAULT 0, create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz);
CREATE UNIQUE INDEX idx_user_api_permission ON m_user_api_permission(user_id, permission_id, source);
CREATE TABLE m_menu_api_permission (id bigint PRIMARY KEY, menu_id bigint NOT NULL, permission_id bigint NOT NULL, create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz);
CREATE UNIQUE INDEX idx_menu_api_permission ON m_menu_api_permission(menu_id, permission_id);
CREATE TABLE m_role_menu (id bigint PRIMARY KEY, role_id bigint NOT NULL, menu_id bigint NOT NULL, create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz);
CREATE INDEX idx_role_menu ON m_role_menu(role_id, menu_id);
CREATE TABLE m_user_role (id bigint PRIMARY KEY, user_id bigint NOT NULL, role_id bigint NOT NULL, create_by bigint, update_by bigint, created_time timestamptz, updated_time timestamptz);
CREATE UNIQUE INDEX idx_user_role ON m_user_role(user_id, role_id);

-- +goose Down
DROP TABLE IF EXISTS m_user_role, m_role_menu, m_menu_api_permission, m_user_api_permission, m_role_api_permission, s_user, s_role, s_org, s_menu, s_auth_client, s_api_permission CASCADE;

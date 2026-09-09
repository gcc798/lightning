-- +goose Up
CREATE TABLE s_config (
  id bigint PRIMARY KEY, name varchar(128) NOT NULL, code varchar(128) NOT NULL, data jsonb NOT NULL,
  remark varchar(500), create_by bigint, created_time timestamptz, update_by bigint, updated_time timestamptz
);
CREATE UNIQUE INDEX idx_s_config_code ON s_config(code);
CREATE TABLE s_dict_data (
  id bigint PRIMARY KEY, parent_id bigint DEFAULT 0, sort bigint DEFAULT 0, dict_label varchar(128),
  dict_value varchar(128), dict_type varchar(128), is_default boolean DEFAULT false, status smallint DEFAULT 0,
  remark varchar(500), create_by bigint, created_time timestamptz, update_by bigint, updated_time timestamptz
);
CREATE INDEX idx_s_dict_data_dict_type ON s_dict_data(dict_type);
CREATE INDEX idx_s_dict_data_parent_id ON s_dict_data(parent_id);
CREATE TABLE s_login_log (
  id bigint PRIMARY KEY, user_name varchar(64), ipaddr varchar(64), login_location varchar(128),
  browser varchar(64), os varchar(64), status smallint DEFAULT 0, msg varchar(500), login_time timestamptz, client_id varchar(64)
);
CREATE INDEX idx_s_login_log_login_time ON s_login_log(login_time);
CREATE INDEX idx_s_login_log_user_name ON s_login_log(user_name);
CREATE TABLE s_oper_log (
  id bigint PRIMARY KEY, title varchar(128), business_type varchar(32), method varchar(255), request_method varchar(16),
  device_type varchar(32), oper_name varchar(64), oper_url varchar(1024), oper_ip varchar(64), oper_location varchar(128),
  oper_param text, json_result text, status char(1), error_msg text, oper_time timestamptz, cost_time bigint, user_agent varchar(512)
);
CREATE INDEX idx_s_oper_log_oper_time ON s_oper_log(oper_time);

INSERT INTO s_config (id, name, code, data, remark, create_by, created_time, update_by, updated_time) VALUES
(1880159541355580001, '微信集成配置', 'integration.wechat', '{"enabled":false,"appId":"","secret":"","templateId":""}', '微信小程序登录与消息能力运行期配置', 0, NOW(), 0, NOW()),
(1880159541355580002, '短信集成配置', 'integration.sms', '{"enabled":false,"accessKeyId":"","accessKeySecret":"","signName":"","templateCode":""}', '短信服务运行期配置', 0, NOW(), 0, NOW()),
(1880159541355580003, '邮件集成配置', 'integration.email', '{"enabled":false,"host":"","port":0,"username":"","password":"","from":""}', '邮件服务运行期配置', 0, NOW(), 0, NOW()),
(1880159541355580004, '验证码配置', 'auth.captcha', '{"image":{"enabled":false,"length":4,"width":120,"height":40,"expire":300},"sms":{"enabled":false,"length":6,"expire":300,"template":"SMS_CODE_TEMPLATE","provider":"aliyun"},"email":{"enabled":false,"length":6,"expire":300,"template":"验证码：%s"}}', '图形、短信和邮件验证码运行期配置', 0, NOW(), 0, NOW()),
(1880159541355580005, '调度器配置', 'scheduler', '{"enabled":true,"refreshIntervalSeconds":5,"jobs":{"data-cleanup":{"enabled":true,"cron":"0 0 2 * * *"}}}', '仅允许调度代码注册的后台任务', 0, NOW(), 0, NOW());

-- +goose Down
DROP TABLE IF EXISTS s_oper_log, s_login_log, s_dict_data, s_config CASCADE;

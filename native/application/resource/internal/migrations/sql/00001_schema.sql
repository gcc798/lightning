-- +goose Up
CREATE TABLE biz_attachment (
  id bigint PRIMARY KEY, file_name varchar(255) NOT NULL, file_key varchar(512) NOT NULL, file_size bigint NOT NULL,
  file_type varchar(128), file_ext varchar(32), business_type varchar(64), business_id varchar(64),
  business_field varchar(64), is_public boolean DEFAULT false, access_url varchar(1024), metadata jsonb,
  status smallint DEFAULT 0, expire_time timestamptz, create_by bigint, create_time timestamptz, update_time timestamptz
);
CREATE INDEX idx_biz_attachment_business_id ON biz_attachment(business_id);
CREATE INDEX idx_biz_attachment_business_type ON biz_attachment(business_type);
CREATE INDEX idx_biz_attachment_file_key ON biz_attachment(file_key);

-- +goose Down
DROP TABLE IF EXISTS biz_attachment CASCADE;

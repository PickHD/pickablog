CREATE SEQUENCE IF NOT EXISTS tag_seq;
CREATE TABLE IF NOT EXISTS tag (
      id INTEGER DEFAULT nextval('tag_seq') PRIMARY KEY,
      name VARCHAR UNIQUE NOT NULL,
      created_at TIMESTAMPTZ,
      created_by VARCHAR,
      updated_at TIMESTAMPTZ,
      updated_by VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_on_tag_id on tag (id);
CREATE TRIGGER log_insert_tag BEFORE INSERT ON tag FOR EACH ROW EXECUTE PROCEDURE log_insert();
CREATE TRIGGER log_update_tag BEFORE UPDATE ON tag FOR EACH ROW EXECUTE PROCEDURE log_update();

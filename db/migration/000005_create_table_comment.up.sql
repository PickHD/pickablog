CREATE SEQUENCE IF NOT EXISTS comment_seq;
CREATE TABLE IF NOT EXISTS comments (
      id INTEGER DEFAULT nextval('comment_seq') PRIMARY KEY,
      comment TEXT NULL,
      user_id INTEGER NOT NULL,
      article_id INTEGER NOT NULL,
      created_at TIMESTAMPTZ,
      created_by VARCHAR,
      updated_at TIMESTAMPTZ,
      updated_by VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_on_comment_id ON comments (id);
CREATE TRIGGER log_insert_comment BEFORE INSERT ON comments FOR EACH ROW EXECUTE PROCEDURE log_insert();
CREATE TRIGGER log_update_comment BEFORE UPDATE ON comments FOR EACH ROW EXECUTE PROCEDURE log_update();

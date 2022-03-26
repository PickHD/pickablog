CREATE SEQUENCE IF NOT EXISTS like_seq;
CREATE TABLE IF NOT EXISTS likes (
      id INTEGER DEFAULT nextval('like_seq') PRIMARY KEY,
      like_count INTEGER NOT NULL,
      user_id INTEGER NOT NULL,
      article_id INTEGER NOT NULL,
      created_at TIMESTAMPTZ,
      created_by VARCHAR,
      updated_at TIMESTAMPTZ,
      updated_by VARCHAR
);

CREATE TRIGGER log_insert_like BEFORE INSERT ON likes FOR EACH ROW EXECUTE PROCEDURE log_insert();
CREATE TRIGGER log_update_like BEFORE UPDATE ON likes FOR EACH ROW EXECUTE PROCEDURE log_update();

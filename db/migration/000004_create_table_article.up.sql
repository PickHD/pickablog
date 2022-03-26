CREATE SEQUENCE IF NOT EXISTS article_seq;
CREATE TABLE IF NOT EXISTS article (
      id INTEGER DEFAULT nextval('article_seq') PRIMARY KEY,
      title TEXT NOT NULL,
      slug TEXT NOT NULL,
      body TEXT NULL,
      footer TEXT NULL,
      user_id INTEGER NOT NULL,
      comments INTEGER[] NULL,
      likes INTEGER[] NULL,
      tags INTEGER[] NULL,
      created_at TIMESTAMPTZ,
      created_by VARCHAR,
      updated_at TIMESTAMPTZ,
      updated_by VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_on_article_id ON article (id);
CREATE INDEX IF NOT EXISTS idx_on_article_slug ON article (slug);
CREATE TRIGGER log_insert_article BEFORE INSERT ON article FOR EACH ROW EXECUTE PROCEDURE log_insert();
CREATE TRIGGER log_update_article BEFORE UPDATE ON article FOR EACH ROW EXECUTE PROCEDURE log_update();

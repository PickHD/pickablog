DROP INDEX IF EXISTS idx_on_article_id;
DROP INDEX IF EXISTS idx_on_article_slug;
DROP TRIGGER IF EXISTS log_insert_article ON article CASCADE;
DROP TRIGGER IF EXISTS log_update_article ON article CASCADE;
DROP TABLE IF EXISTS article CASCADE;
DROP SEQUENCE IF EXISTS article_seq CASCADE;
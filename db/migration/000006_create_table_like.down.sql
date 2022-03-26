-- schema role section
DROP TRIGGER IF EXISTS log_insert_like ON likes CASCADE;
DROP TRIGGER IF EXISTS log_update_like ON likes CASCADE;
DROP TABLE IF EXISTS likes CASCADE;
DROP SEQUENCE IF EXISTS like_seq CASCADE;
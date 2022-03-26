-- schema user section
DROP INDEX IF EXISTS idx_on_user_id;
DROP INDEX IF EXISTS idx_on_user_email;
DROP TRIGGER IF EXISTS log_insert_user ON "user" CASCADE;
DROP TRIGGER IF EXISTS log_update_user ON "user" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP SEQUENCE IF EXISTS user_seq CASCADE;

-- schema role section
DROP TRIGGER IF EXISTS log_insert_role ON role CASCADE;
DROP TRIGGER IF EXISTS log_update_role ON role CASCADE;
DROP TABLE IF EXISTS role CASCADE;
DROP SEQUENCE IF EXISTS role_seq CASCADE;
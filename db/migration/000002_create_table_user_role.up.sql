-- schema role section
CREATE SEQUENCE IF NOT EXISTS role_seq;
CREATE TABLE IF NOT EXISTS role (
      id INTEGER DEFAULT nextval('role_seq') PRIMARY KEY,
      name VARCHAR UNIQUE NOT NULL,
      created_at TIMESTAMPTZ,
      created_by VARCHAR,
      updated_at TIMESTAMPTZ,
      updated_by VARCHAR
);

CREATE TRIGGER log_insert_role BEFORE INSERT ON role FOR EACH ROW EXECUTE PROCEDURE log_insert();
CREATE TRIGGER log_update_role BEFORE UPDATE ON role FOR EACH ROW EXECUTE PROCEDURE log_update();

-- schema user section
CREATE SEQUENCE IF NOT EXISTS user_seq;
CREATE TABLE IF NOT EXISTS "user" (
     id INTEGER DEFAULT nextval('user_seq') PRIMARY KEY,
     full_name VARCHAR NOT NULL,
     email VARCHAR UNIQUE NOT NULL,
     password VARCHAR NOT NULL,
     role_id INTEGER NOT NULL,
     created_at TIMESTAMPTZ,
     created_by VARCHAR,
     updated_at TIMESTAMPTZ,
     updated_by VARCHAR
);

CREATE INDEX IF NOT EXISTS idx_on_user_id on "user" (id);
CREATE INDEX IF NOT EXISTS idx_on_user_email on "user" (email);
CREATE TRIGGER log_insert_user BEFORE INSERT ON "user" FOR EACH ROW EXECUTE PROCEDURE log_insert();
CREATE TRIGGER log_update_user BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE PROCEDURE log_update();
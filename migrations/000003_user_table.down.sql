ALTER TABLE url_entity
DROP CONSTRAINT fk_url_entity_user;

ALTER TABLE url_entity
DROP COLUMN user_id;

DROP TABLE IF EXISTS users;
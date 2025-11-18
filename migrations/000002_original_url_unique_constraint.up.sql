DELETE FROM url_entity t1
USING url_entity t2
WHERE t1.original_url = t2.original_url
  AND t1.id > t2.id;

ALTER TABLE url_entity
ADD CONSTRAINT unique_original_url UNIQUE (original_url);
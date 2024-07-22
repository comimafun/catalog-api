-- alter table circle add column "rating"
ALTER TABLE "circle"
ADD COLUMN rating varchar(10) CHECK (rating IN ('GA', 'PG', 'M'));

CREATE INDEX circle_rating_idx ON "circle" ("rating");
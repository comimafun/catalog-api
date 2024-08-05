alter table "circle"
add column "batch" integer;

create index idx_circle_batch on "circle" ("batch");
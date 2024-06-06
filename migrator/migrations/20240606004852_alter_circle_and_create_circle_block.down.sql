alter table "circle"
drop column if exists "circle_block_id";

drop table if exists "circle_block";

drop index if exists "idx_circle_circle_block";
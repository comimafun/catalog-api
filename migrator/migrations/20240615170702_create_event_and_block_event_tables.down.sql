drop table if exists block_event;

drop table if exists "event";

alter table "circle"
drop column if exists "event_id";

drop index if exists "idx_circle_event_id";
drop index if exists "idx_circle_event_id";

drop index if exists "idx_block_event_event_id";

drop index if exists "idx_block_event_circle_id";

drop index if exists "idx_block_event_prefix_postfix_event_id";

alter table "circle"
drop column if exists "event_id";

drop table if exists block_event;

drop table if exists "event";
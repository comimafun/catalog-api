drop index if exists idx_circle_batch;

alter table circle
drop column if exists "batch"
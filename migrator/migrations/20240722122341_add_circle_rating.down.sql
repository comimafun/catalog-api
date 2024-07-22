-- drop index circle_rating_idx
drop if exists index circle_rating_idx;

alter table "circle"
drop column "rating";
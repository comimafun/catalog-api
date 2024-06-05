drop index if exists "idx_user_circle_id";

alter table "user"
drop column if exists "circle_id";

drop index if exists "idx_user_upvote_user_id";

drop index if exists "idx_user_upvote_circle_id";

drop table if exists "user_upvote";

drop index if exists "idx_user_bookmark_user_id";

drop index if exists "idx_user_bookmark_circle_id";

drop table if exists "user_bookmark";
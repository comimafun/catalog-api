alter table "user"
add column "circle_id" integer references "circle" ("id") on delete set null;

create index "idx_user_circle_id" on "user" ("circle_id");

create table
    "user_upvote" (
        "user_id" integer not null,
        "circle_id" integer not null,
        "created_at" timestamp not null default current_timestamp,
        primary key ("user_id", "circle_id"),
        foreign key ("user_id") references "user" ("id") on delete cascade,
        foreign key ("circle_id") references "circle" ("id") on delete cascade
    );

create index "idx_user_upvote_user_id" on "user_upvote" ("user_id");

create index "idx_user_upvote_circle_id" on "user_upvote" ("circle_id");

create table
    "user_bookmark" (
        "user_id" integer not null,
        "circle_id" integer not null,
        "created_at" timestamp not null default current_timestamp,
        primary key ("user_id", "circle_id"),
        foreign key ("user_id") references "user" ("id") on delete cascade,
        foreign key ("circle_id") references "circle" ("id") on delete cascade
    );

create index "idx_user_bookmark_user_id" on "user_bookmark" ("user_id");

create index "idx_user_bookmark_circle_id" on "user_bookmark" ("circle_id");
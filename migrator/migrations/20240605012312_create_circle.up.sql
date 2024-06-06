create type day as enum ('first', 'second', 'both');

create table
    "circle" (
        "id" serial primary key,
        "name" varchar(255) not null,
        "slug" varchar(255) not null unique,
        "picture_url" varchar(255),
        "url" varchar(255),
        "facebook_url" varchar(255),
        "twitter_url" varchar(255),
        "instagram_url" varchar(255),
        "description" text,
        "batch" integer,
        "verified" boolean default false,
        "publish" boolean default false,
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        "deleted_at" timestamp,
        "user_id" integer,
        "day" day,
        foreign key ("user_id") references "user" ("id") on delete cascade
    );

create index "idx_circle_user_id" on "circle" ("user_id");

create index "idx_circle_verified" on "circle" ("verified");

create index "idx_circle_batch" on "circle" ("batch");

create index "idx_circle_day" on "circle" ("day");
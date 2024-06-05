create table
    "refresh_token" (
        "id" serial primary key,
        "token" varchar(255) not null unique,
        "access_token" varchar(255) not null unique,
        "user_id" integer not null,
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        "expires_at" timestamp not null,
        foreign key ("user_id") references "user" ("id") on delete cascade
    );

create index "idx_refresh_token_user_id" on "refresh_token" ("user_id");
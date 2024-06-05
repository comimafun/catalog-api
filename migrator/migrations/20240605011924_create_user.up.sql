create table
    "user" (
        "id" serial primary key,
        "email" varchar(255) not null unique,
        "name" varchar(255) not null,
        "hash" varchar(255) not null,
        "profile_picture_url" varchar(255),
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        "deleted_at" timestamp
    )
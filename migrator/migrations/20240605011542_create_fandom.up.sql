create table
    "fandom" (
        "id" serial primary key,
        "name" varchar(255) not null,
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        "deleted_at" timestamp
    )
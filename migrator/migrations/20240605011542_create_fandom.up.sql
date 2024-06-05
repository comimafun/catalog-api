create table
    "fandom" (
        "id" serial primary key,
        "name" varchar(255) not null,
        "visible" boolean not null default true,
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        "deleted_at" timestamp
    );

create index "idx_fandom_visible" on "fandom" ("visible");
create table
    "circle_fandom" (
        "circle_id" int not null,
        "fandom_id" int not null,
        "created_at" timestamp not null default current_timestamp,
        primary key ("circle_id", "fandom_id"),
        foreign key ("circle_id") references "circle" ("id") on delete cascade,
        foreign key ("fandom_id") references "fandom" ("id") on delete cascade
    );
create table
    "circle_fandom" (
        "circle_id" int not null,
        "fandom_id" int not null,
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        primary key ("circle_id", "fandom_id")
    );
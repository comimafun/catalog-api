create table
    "event" (
        id serial primary key,
        "name" varchar(255) not null,
        slug varchar(255) not null unique,
        "description" text,
        started_at timestamp not null,
        ended_at timestamp not null,
        created_at timestamp not null default current_timestamp,
        updated_at timestamp not null default current_timestamp,
        deleted_at timestamp
    );

create table
    "block_event" (
        id serial primary key,
        event_id integer not null,
        circle_id integer not null unique,
        prefix varchar(10) not null,
        postfix varchar(10) not null,
        created_at timestamp not null default current_timestamp,
        updated_at timestamp not null default current_timestamp,
        deleted_at timestamp,
        unique (prefix, postfix, event_id),
        foreign key (circle_id) references "circle" (id) on delete set null,
        foreign key (event_id) references "event" (id) on delete cascade
    );

alter table "circle"
add column "event_id" integer;

alter table "circle" add foreign key ("event_id") references "event" (id) on delete set null;

-- add index by event id
create index "idx_circle_event_id" on "circle" ("event_id");
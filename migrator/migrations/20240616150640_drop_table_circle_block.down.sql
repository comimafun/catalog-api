create table
    "circle_block" (
        id serial primary key,
        circle_id int,
        prefix varchar(2) not null,
        postfix varchar(8) not null,
        unique (prefix, postfix),
        created_at timestamp not null default current_timestamp,
        updated_at timestamp not null default current_timestamp,
        deleted_at timestamp,
        foreign key (circle_id) references "circle" (id) on delete set null
    );
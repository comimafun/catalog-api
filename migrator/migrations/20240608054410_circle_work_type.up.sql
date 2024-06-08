create table
    "circle_work_type" (
        circle_id integer not null,
        work_type_id integer not null,
        created_at timestamp not null default current_timestamp,
        primary key (circle_id, work_type_id),
        foreign key (circle_id) references "circle" (id) on delete cascade,
        foreign key (work_type_id) references "work_type" (id) on delete cascade
    )
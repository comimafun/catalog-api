create table
  "report" (
    id serial primary key,
    user_id integer not null,
    circle_id integer not null,
    reason varchar(255),
    created_at timestamp not null default current_timestamp,
    foreign key (circle_id) references "circle" (id) on delete cascade,
    foreign key (user_id) references "user" (id) on delete cascade
  );
  
-- add index
create index "idx_report_user_id" on "report" ("user_id");

create index "idx_report_circle_id" on "report" ("circle_id");
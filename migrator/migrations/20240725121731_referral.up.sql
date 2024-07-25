create table
    "referral" (
        "id" serial primary key,
        "referral_code" varchar(255) not null unique,
        "circle_id" integer not null unique,
        "created_at" timestamp not null,
        "updated_at" timestamp not null
    );

alter table "referral" add foreign key ("circle_id") references "circle" (id) on delete cascade;

create index "idx_referral_circle_id" on "referral" ("circle_id");

alter table "circle"
add column "used_referral_code_id" integer;

alter table "circle" add foreign key ("used_referral_code_id") references "referral" (id) on delete set null;
create table
    "product" (
        "id" serial primary key,
        "name" varchar(255) not null,
        "image_url" varchar(255) not null,
        "circle_id" integer not null,
        "created_at" timestamp not null default current_timestamp,
        "updated_at" timestamp not null default current_timestamp,
        "deleted_at" timestamp,
        foreign key ("circle_id") references "circle" ("id") on delete cascade
    );

create index "idx_product_circle_id" on "product" ("circle_id");
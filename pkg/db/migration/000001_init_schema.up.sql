CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "picture" varchar,
  "locale" varchar NOT NULL,
  "board_limit" int NOT NULL,
  "created_at" timestamptz NOT NULL,
  "verified_email" boolean
);

CREATE TABLE "boards" (
  "id" varchar PRIMARY KEY,
  "project_name" varchar UNIQUE NOT NULL,
  "picture" varchar,
  "owner_id" varchar NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE TABLE "tasks" (
  "id" varchar PRIMARY KEY,
  "order" int NOT NULL,
  "name" varchar NOT NULL,
  "board_id" varchar NOT NULL
);

CREATE TABLE "members" (
  "user_id" varchar NOT NULL,
  "board_id" varchar NOT NULL
);

CREATE TABLE "task_items" (
  "id" varchar PRIMARY KEY,
  "description" varchar NOT NULL,
  "order" int NOT NULL,
  "tag" varchar,
  "assignee" varchar,
  "status_id" varchar NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "comments" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "task_id" varchar NOT NULL,
  "comment" varchar NOT NULL,
  "created_at" timestamptz NOT NULL
);

ALTER TABLE "boards" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "tasks" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id") ON DELETE CASCADE;

ALTER TABLE "members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "members" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id") ON DELETE CASCADE;

ALTER TABLE "task_items" ADD FOREIGN KEY ("assignee") REFERENCES "users" ("id");

ALTER TABLE "task_items" ADD FOREIGN KEY ("status_id") REFERENCES "tasks" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("task_id") REFERENCES "task_items" ("id") ON DELETE CASCADE;

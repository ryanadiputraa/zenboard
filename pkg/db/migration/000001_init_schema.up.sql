CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "picture" varchar,
  "locale" varchar NOT NULL,
  "board_limit" int NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE TABLE "boards" (
  "id" varchar PRIMARY KEY,
  "project_name" varchar UNIQUE NOT NULL,
  "picture" varchar,
  "owner_id" varchar NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE TABLE "task_status" (
  "id" varchar PRIMARY KEY,
  "order" int NOT NULL,
  "name" varchar NOT NULL,
  "board_id" varchar NOT NULL
);

CREATE TABLE "members" (
  "user_id" varchar NOT NULL,
  "board_id" varchar NOT NULL
);

CREATE TABLE "tasks" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "order" int NOT NULL,
  "tag" varchar,
  "assignee" varchar NOT NULL,
  "board_id" varchar NOT NULL,
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

ALTER TABLE "boards" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "task_status" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id");

ALTER TABLE "members" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "members" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("assignee") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("board_id") REFERENCES "boards" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("status_id") REFERENCES "task_status" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

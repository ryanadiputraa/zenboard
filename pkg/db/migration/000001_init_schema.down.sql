ALTER TABLE "comments" DROP CONSTRAINT IF EXISTS "comments_task_id_fkey";
ALTER TABLE "comments" DROP CONSTRAINT IF EXISTS "comments_user_id_fkey";
ALTER TABLE "task_items" DROP CONSTRAINT IF EXISTS "task_items_status_id_fkey";
ALTER TABLE "task_items" DROP CONSTRAINT IF EXISTS "task_items_assignee_fkey";
ALTER TABLE "members" DROP CONSTRAINT IF EXISTS "members_board_id_fkey";
ALTER TABLE "members" DROP CONSTRAINT IF EXISTS "members_user_id_fkey";
ALTER TABLE "tasks" DROP CONSTRAINT IF EXISTS "tasks_board_id_fkey";
ALTER TABLE "boards" DROP CONSTRAINT IF EXISTS "boards_owner_id_fkey";

DROP TABLE IF EXISTS "task_items";
DROP TABLE IF EXISTS "comments";
DROP TABLE IF EXISTS "members";
DROP TABLE IF EXISTS "tasks";
DROP TABLE IF EXISTS "boards";
DROP TABLE IF EXISTS "users";
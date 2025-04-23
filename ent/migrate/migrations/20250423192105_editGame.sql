-- Rename a column from "user_player_first_id" to "user_white_id"
ALTER TABLE "chesses" RENAME COLUMN "user_player_first_id" TO "user_white_id";
-- Rename a column from "user_player_second_id" to "user_black_id"
ALTER TABLE "chesses" RENAME COLUMN "user_player_second_id" TO "user_black_id";
-- Modify "chesses" table
ALTER TABLE "chesses" DROP CONSTRAINT "chesses_users_player_first_id", DROP CONSTRAINT "chesses_users_player_second_id", ALTER COLUMN "status" SET DEFAULT 'waiting', ADD COLUMN "result" character varying NOT NULL DEFAULT '0-0', ADD CONSTRAINT "chesses_users_black_id" FOREIGN KEY ("user_black_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "chesses_users_white_id" FOREIGN KEY ("user_white_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

-- Create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "email" character varying NOT NULL, "name" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "password" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create "chesses" table
CREATE TABLE "chesses" ("id" uuid NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "status" character varying NOT NULL DEFAULT 'in_progress', "user_player_first_id" uuid NOT NULL, "user_player_second_id" uuid NOT NULL, "user_winner_id" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "chesses_users_player_first_id" FOREIGN KEY ("user_player_first_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "chesses_users_player_second_id" FOREIGN KEY ("user_player_second_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "chesses_users_winner_id" FOREIGN KEY ("user_winner_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);

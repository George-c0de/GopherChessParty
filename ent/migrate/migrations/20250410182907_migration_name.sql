-- Create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "email" character varying NOT NULL, "name" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "password" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create "chesses" table
CREATE TABLE "chesses" ("id" uuid NOT NULL, "status" smallint NOT NULL DEFAULT 0, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "user_chesses_as_first" uuid NULL, "user_chesses_as_second" uuid NULL, "user_chesses_won" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "chesses_users_chesses_as_first" FOREIGN KEY ("user_chesses_as_first") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "chesses_users_chesses_as_second" FOREIGN KEY ("user_chesses_as_second") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT "chesses_users_chesses_won" FOREIGN KEY ("user_chesses_won") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);

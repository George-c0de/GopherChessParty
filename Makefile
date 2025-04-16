add_schema:
	ent new

migrate:
	atlas migrate apply --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:postgres@localhost:5432/gopher_chess?search_path=public&sslmode=disable" \
  --revisions-schema public

generate_schema:
	go run entgo.io/ent/cmd/ent generate ./ent/schema


revision:
	atlas migrate diff migration_name \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "postgres://postgres:postgres@localhost:5432/gopher_chess?search_path=public&sslmode=disable"

lint:
	golangci-lint run
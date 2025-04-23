add_schema:
	ent new

migrate:
	atlas migrate apply --dir "file://ent/migrate/migrations" \
  --url "postgres://postgres:postgres@localhost:5432/gopher_chess?search_path=public&sslmode=disable" \
  --revisions-schema public

generate_schema:
	go run entgo.io/ent/cmd/ent generate ./ent/schema

checkpoint:
	atlas migrate checkpoint \
		--dir "file://ent/migrate/migrations" \
		--dev-url "docker://postgres/15/dev"

revision:
	atlas migrate diff DeleteWinner \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/dev"


lint:
	golangci-lint run

format:
	golangci-lint fmt
.PHONY: postgres adminer migrate

postgres:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=secret postgres

adminer:
	docker run --rm -ti -p 8080:8080 adminer

migrate:
	C:\Users\user\Downloads\migrate.windows-386\migrate.exe -path ./migrations -database postgres://postgres:secret@localhost/postgres?sslmode=disable up

migrate-down:
	C:\Users\user\Downloads\migrate.windows-386\migrate.exe -path ./migrations -database postgres://postgres:secret@localhost/postgres?sslmode=disable down
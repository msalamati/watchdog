docker run --rm -d --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mypass -e POSTGRES_DB=postgres postgres:13
POSTGRES_URL=postgresql://postgres:mypass@localhost:5432/postgres?sslmode=disable go run main.go

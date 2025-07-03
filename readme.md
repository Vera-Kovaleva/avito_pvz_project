# Полезные команды

Запустить Docker Desktop

- docker compose down && docker system prune --volumes --force && docker compose up -d
- pgcli --host 127.0.0.1 --port 5432 --username postgres

- go test -count=1 ./...

# Запуск
cp .env.sample .env
go run .

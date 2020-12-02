# Mattrax

[![Go Report Card](https://goreportcard.com/badge/github.com/mattrax/Mattrax)](https://goreportcard.com/report/github.com/mattrax/Mattrax)

Open Source MDM (Mobile Device Management) System. Supporting Windows, Linux and macOS. There are plans to implement IOS, Android and ChromeOS in the future.

## Project Status

Mattrax is under super **heavy development**. The Mattrax Cloud (SaaS version) will launch early 2021 and the self-hosted version will be available not long after that. Currently, the project is not ready for prime time so **do not expose it to the internet** and **expect bugs**.

## Developing

Mattrax is built using [Go](https://golang.org) and uses the [PostgreSQL](https://www.postgresql.org) database. To setup a development environment use the instructions below.

```bash
# Terminal 1: The Postgres database
docker-compose -f ./docker-compose-dev.yml up
# Terminal 2: Go Backend
go run ./cmd/mattrax --db "postgres://mattrax_db:password@localhost/mattrax?sslmode=disable" --domain mdm.example.com --tlscert "./certs/tls.crt" --tlskey "./certs/tls.key" --debug
# Terminal 3: Vue Frontend
cd dashboard
npm i
API_BASE_URL="https://mdm.example.com/api" npm run dev
# Open http://localhost:3000 in your browser and use the default credentials logged to the console.
```

This project uses [sqlc](https://github.com/kyleconroy/sqlc) so the command `sqlc generate` is used to generate the `internal/db` package from `sql/queries.sql`.

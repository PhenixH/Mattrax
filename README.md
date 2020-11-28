# Mattrax

[![Go Report Card](https://goreportcard.com/badge/github.com/mattrax/Mattrax)](https://goreportcard.com/report/github.com/mattrax/Mattrax)

Open Source MDM (Mobile Device Management) System. Supporting Windows, Linux and macOS. There are plans to implement IOS, Android and ChromeOS in the future.

## Project Status

Mattrax is under super **heavy development**. The Mattrax Cloud (SaaS version) will launch early 2021 and the self-hosted version will be available not long after that. Currently, the project is not ready for prime time so **do not expose it to the internet** and **expect bugs**.

## Running

**Instructions are out of date and will be updated in near future**

This project requires an external [PostgreSQL](https://www.postgresql.org/) database. The `sql/schema.sql` and `sql/deploy.sql` should be run on a blank database to configure it. You should change the `deploy.sql` to fit your deployment settings. Then start the Go binary (with arguments `--db "postgres://localhost/Mattrax" --domain mdm.example.com`) and your MDM server will be working.

## Developing

This project uses [sqlc](https://github.com/kyleconroy/sqlc) so the command `sqlc generate` is used to generate the `internal/db` package from `sql/queries.sql`.

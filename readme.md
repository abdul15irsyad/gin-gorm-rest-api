# Gin Gorm REST API

project for learning golang with gin framework and gorm as ORM

## Tech Stack

- Golang go1.22.2
- Gin Framework
- Gorm
- PostgreSQL

## Installation

1. pull repository

   ```bash
   git clone git@github.com:abdul15irsyad/gin-gorm-rest-api.git
   cd gin-gorm-rest-api
   ```

2. install dependencies
   ```bash
   go mod download
   ```
3. configure environment, copy from `.env.example` to `.env` and adjust to your setup like database, port, jwt secret, etc.
   ```bash
   cp .env.example .env
   ```
4. run the project

   ```bash
   go run main.go
   ```

   or with [air](https://github.com/cosmtrek/air) (live reload)

   ```bash
   air
   ```

5. run the seeder for dummy datas

   ```bash
   go run cmd/seed/main.go
   ```

## Postman

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/6292564-55c171e1-8b56-4b55-a1f8-bd97378281c1?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D6292564-55c171e1-8b56-4b55-a1f8-bd97378281c1%26entityType%3Dcollection%26workspaceId%3De14a18ea-da74-4b90-b978-d57d03cd3ded#?env%5Bgin%20gorm%20local%5D=W3sia2V5IjoiYmFzZV91cmwiLCJ2YWx1ZSI6Imh0dHA6Ly9sb2NhbGhvc3Q6MzAwMSIsImVuYWJsZWQiOnRydWUsInR5cGUiOiJkZWZhdWx0Iiwic2Vzc2lvblZhbHVlIjoiaHR0cDovL2xvY2FsaG9zdDozMDAxIiwic2Vzc2lvbkluZGV4IjowfV0=)

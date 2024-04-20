# Learning Gin Gorm

project for learning golang with gin framework and gorm as ORM

## Tech Stack

- Golang go1.22.2
- Gin Framework
- Gorm
- PostgreSQL

## Installation

1. pull repository

   ```bash
   git clone git@github.com:abdul15irsyad/belajar-gin.git
   cd belajar-gin
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

   or with air (live reload)

   ```bash
   air
   ```

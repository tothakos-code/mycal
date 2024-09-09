<p align="center">
<h1 align="center">Vue Go PostgreSQL Auth Template</h1>
</p>

Minimalistic boilerplate for Vue, Go, PostgreSQL.

## Features

- 🚀 Client and Server Protected Routes
- 🚀 User Authentication
- 🚀 Go stdlib `net/http`
- 🚀 PostgreSQL no orm
- 🚀 Clean Architecture
- 🚀 Docker

---

## Setup

### Frontend

#### 1. Run Vite

```bash
npm install
npm run dev
```

### Backend

#### 1. Set your Environment Variables

```py
# Here is an example .env file
PORT=8080
JWT_TOKEN_DURATION_HOURS=24
JWT_SECRET=meow

DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=gpa # database name
DB_USERNAME=root
DB_PASSWORD=password
DB_SCHEMA=public
```

#### 2. Change makefile to match your database name

```makefile
db:
	docker compose up -d;
  # default is "gpa" as specified in .env
	psql -h localhost -p 5432 -U root -d gpa;

```

#### 3. Set up database schema

connect to db

```bash
make db
```

execute this sql query`./internal/database/schema.sql` in your database

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(254) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

```

#### 4. Run Server

```bash
# run db
make db

# run server
go run cmd/api/main.go

# or use air for hot reload
air
```

---

## What you need to know

### Frontend-Backend Communication

#### Success:

- Backend: `sendJson(w, response, http.StatusOK)`
  - `sendJson` is just util function that just sends a JSON response to the client
- Frontend: `response.json()`

#### Error:

- Backend: `http.Error(w, "Invalid Request", http.StatusBadRequest)`
- Frontend: `response.text()`

### Frontend Stuff

- `/router/index.tsx` is where you declare and manage your routes
- `/context/AuthContext.tsx` is where you can find the `useAuth` hook
  - This hook gives you access to the `user` object from the backend
  - Manages the expiration of the JWT token
- `/Providers.tsx` is where you can add more `providers` or `wrappers`

---

### Inspired By:

[go-blueprint](https://github.com/Melkeydev/go-blueprint) <br>
[How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/#an-opportunity-to-hide-the-requestresponse-types-away)

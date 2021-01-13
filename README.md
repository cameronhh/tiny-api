# [Tiny API](https://tiny-api.dev) - A Tiny API To Make Tiny APIs

![Website Preview](https://github.com/cameronhh/tiny-api-client/blob/master/.github/repo-image.png)

- _Quickly prototyping some front-end code and need a couple of endpoints to return some static JSON?_
- _Perhaps you're integrating a front-end with a back-end that doesn't have a test server?_
- **_Make an endpoint that returns static JSON in seconds with [Tiny API](https://tiny-api.dev)_**

This repo is the server behind [https://tiny-api.dev](https://tiny-api.dev). The Tiny API Server is a REST API written in Go using [Gin](https://github.com/gin-gonic/gin).
The code for the front end can be found [here](https://github.com/cameronhh/tiny-api-client).

## Setting up the development environment

### 1. Start the test database

Tiny API is built on top of Postgres (I know, SQLite would have been tinier).
You can quickly run a test database with docker using the following command.

```
docker run -d -p 5432:5432 \
--name tiny-api-postgres \
-e POSTGRES_USER=dev \
-e POSTGRES_PASSWORD=mysecretpassword \
-e POSTGRES_DB=tiny-api-dev \
postgres
```

### 2. Set Environment Variables

Running the server requires the following environment variables to be set:

```
DB_CONNECTION=postgresql://dev:mysecretpassword@localhost:5432/tiny-api-dev?sslmode=disable
PORT=8080
GIN_ENV=development
CLIENT_URL=http://localhost:3000
```

To simplify things, create a `.env` and then run:

```
export $(cat .env | xargs)
```

### 3. Run the tests

Running the tests will populate the test database with tables if they don't already exist.
The tests can be run with:

```
go test
```

### 4. Run the server

To run the server, run:

```
go run .
```

## Things To Do:

- ~~Allow only URL safe characters to be used for temp endpoints~~
- Add an endpoint to check if a new temp endpoint already exists or not
- Remove (or at least comment out) endpoints that are unused
- Improve the directory structure (tests in particular)
- ~~Write a database script for bootstrapping the tables~~
- Add an expiry for temp endpoints, and a helpful response message once those endpoints are hit after expiring

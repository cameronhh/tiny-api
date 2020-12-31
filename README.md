# [Tiny API](https://tiny-api.dev) - A Tiny API To Make Tiny APIs
![Website Preview](https://github.com/cameronhh/tiny-api/blob/master/.github/repo-image.png)

- *Quickly prototyping some front-end code and need a couple of endpoints to return some static JSON?*
- *Perhaps you're integrating a front-end with a back-end that doesn't have a test server?*
- ***Make an endpoint that returns static JSON in seconds with [Tiny API](https://tiny-api.dev)***


This repo is the server behind [tiny-api.dev](https://tiny-api.dev). The Tiny API Server is a REST API written in Go using [Gin](https://github.com/gin-gonic/gin).

## Setting up the development environment

### 1. Start the test database

Tiny API is built on top of Postgres (I know, SQLite would have been tinier).
You can quickly run a test database with docker using the following command.

```
docker run -d -p 5432:5432 \
--name tiny-api-dbms \
-e POSTGRES_USER=postgres \
-e POSTGRES_PASSWORD=mysecretpassword \
-e POSTGRES_DB=tiny-api-db \
postgres
```

### 2. Set Environment Variables

Running the server requires the following environment variables to be set:

```
DB_HOSTNAME=localhost
DB_USERNAME=postgres
DB_PASSWORD=mysecretpassword
DB_NAME=tiny-api-db
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

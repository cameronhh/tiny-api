# Tiny API - A Tiny API For Making Tiny APIs

## Usage

### Run a test database

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

### Setting Environment Variables

Running the server requires the following environment variables to be set:

```
DB_HOSTNAME=localhost
DB_USERNAME=postgres
DB_PASSWORD=mysecretpassword
DB_NAME=tiny-api-db
GIN_ENV=development
CLIENT_URL=http://localhost:3000
```

To simplify things, create a file called `.env` and then run:

```
export $(cat .env | xargs)
```

### Run the server

To run the server, run:

```
go run .
```

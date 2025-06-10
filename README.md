# Backend for SVAR FilterBuilder

=======


### How to use

- configure DB connection
- import DB dump ( dump.sql)
- run the app

```bash
go build
./filter-backend-go
```

### REST API

#### Get all data from the tablesave

```
POST /api/data/{table}
```

Body can contain a filtering query

#### Get unique field values

```
GET /api/data/{table}/{field}/suggest
```

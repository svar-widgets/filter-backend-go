SVAR Query - Backend Demo
=======


### How to use

- configure DB connection
- import DB dump ( dump.sql)
- run the app

```bash
./query
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

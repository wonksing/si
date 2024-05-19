# si(storage interface)

`si` is a package designed to help developers (mostly myself) read and write data to various destinations by wrapping standard and non-standard packages. Its main function is to convert between structs and primary types (bytes, strings, etc.), allowing you to skip the encoding and decoding routines.

- file 
- tcp
- sql
- http
- websocket ([gorilla](https://github.com/gorilla/websocket))
- kafka ([sarama](https://github.com/Shopify/sarama))
- elasticsearch ([go-elasticsearch](https://github.com/elastic/go-elasticsearch))
- ftp ([jlaffaye](https://github.com/jlaffaye/ftp))

## Installation

```bash
go get -u github.com/wonksing/si/v2
```

## Quick Start

1. sql
```go
connStr := "host=testpghost port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
driver := "postgres"
db, _ := sql.Open(driver, connStr)
sqldb := sisql.NewSqlDB(db).WithTagKey("si")

type BoolTest struct {
    Nil      string `json:"nil" si:"nil"`
    True_1   bool   `json:"true_1" si:"true_1"`
    True_2   bool   `json:"true_2" si:"true_2"`
    False_1  bool   `json:"false_1" si:"false_1"`
    False_2  bool   `json:"false_2" si:"false_2"`
    Ignore_3 string `si:"-"`
}
query := `
    select null as nil,
        null as true_1, '1' as true_2, 
        0 as false_1, '0' as false_2
    union all
    select null as nil,
        1 as true_1, '1' as true_2,
        0 as false_1, '0' as false_2
`

m := []BoolTest{}
_, err := sqldb.QueryStructs(query, &m)

```
## Test

Test flags include the following.

- `ONLINE_TEST`
  - Runs tests that actually connect to the storages such as PostgreSQL database.
- `LONG_TEST`
  - Runs tests that takes long time to complete.

Examples of running the tests.

```bash
ONLINE_TEST=1 LONG_TEST=1 go test ./...
ONLINE_TEST=1 go test -run SKIP -bench . -benchtime 100x -benchmem
```

## Versions

### v2.2.1

- Moved from [github.com/go-wonk/si](github.com/go-wonk/si) 


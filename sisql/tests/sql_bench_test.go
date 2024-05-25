package sisql_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sisql"
	"github.com/wonksing/si/v2/tests/testmodels"
)

func BenchmarkSqlDB_Exec(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	sqldb := sisql.NewSqlDB(db)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		query := `insert into student(email_address, name, borrowed)
		values($1, $2, $3)`

		sqldb.Exec(query, i, i, 0)
		// byt, _ := json.Marshal(m)
		// json.Unmarshal(byt, &o)
	}
}

func BenchmarkSqlDB_QueryIntoMap(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		sqldb := sisql.NewSqlDB(db) // sio.SqlColumn{"decimal_", sio.SqlColTypeFloat64},
		// sio.SqlColumn{"numeric_", sio.SqlColTypeFloat64},
		// sio.SqlColumn{"char_arr_", sio.SqlColTypeUints8},

		query := `
			select null as nil,
				'123'::varchar(255) as str,
				123::integer as int2_,
				123::decimal(24,4) as decimal_,
				123::numeric(24,4) as numeric_,
				123::bigint as bigint_,
				'{"abcde", "lunch"}'::char(5)[] as char_arr_,
				'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
				'0123'::bytea as bytea_,
				to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

		m := make([]map[string]interface{}, 0)
		sqldb.QueryMaps(query, &m)
		// byt, _ := json.Marshal(m)
		// json.Unmarshal(byt, &o)
	}

}

/*
Benchmark on json and mapstructure

goos: windows
goarch: amd64
pkg: github.com/wonksing/si/v2/sisql/tests/sql_test
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkSqlDB_QueryIntoAny_Struct-8                 100           1683737 ns/op            7236 B/op        151 allocs/op
BenchmarkSqlDB_QueryIntoAny_Struct2-8                100           2903963 ns/op           11857 B/op        252 allocs/op
PASS
*/
func BenchmarkSqlDB_QueryIntoAny_Struct(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	for i := 0; i < b.N; i++ {
		sqldb := sisql.NewSqlDB(db) // sio.SqlColumn{"decimal_", sio.SqlColTypeFloat64},
		// sio.SqlColumn{"numeric_", sio.SqlColTypeFloat64},
		// sio.SqlColumn{"char_arr_", sio.SqlColTypeUints8},

		query := `
			select null as nil,
				'123'::varchar(255) as str,
				123::integer as int2_,
				123::decimal(24,4) as decimal_,
				123::numeric(24,4) as numeric_,
				123::bigint as bigint_,
				'{"abcde", "lunch"}'::char(5)[] as char_arr_,
				'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
				'0123'::bytea as bytea_,
				to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

		tl := testmodels.TableList{}
		_, err := sqldb.QueryStructs(query, &tl)
		require.Nil(b, err)

		expected := `[{"nil":"","int2_":123,"decimal_":123,"numeric_":123,"bigint_":123,"char_arr_":"e2FiY2RlLGx1bmNofQ==","varchar_arr_":"e2FiY2RlLGx1bmNofQ==","bytea_":"0123","time_":"2022-01-01T12:12:12Z"}]`
		assert.Equal(b, expected, tl.String())
	}
}

func BenchmarkSqlDB_QueryMaps(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	sqldb := sisql.NewSqlDB(db)
	for i := 0; i < b.N; i++ {

		query := `
			select null as nil,
				'123'::varchar(255) as str,
				123::integer as int2_,
				123::decimal(24,4) as decimal_,
				123::numeric(24,4) as numeric_,
				123::bigint as bigint_,
				'{"abcde", "lunch"}'::char(5)[] as char_arr_,
				'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
				'0123'::bytea as bytea_,
				to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

		// tl := Table{}
		tl := make([]map[string]interface{}, 0, 16)
		sqldb.QueryMaps(query, &tl)
	}
}

func BenchmarkSqlDB_QueryStructs(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	sqldb := sisql.NewSqlDB(db)
	query := `
		select null as nil,
			'123'::varchar(255) as str,
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_,
			123::bigint as bigint_,
			'{"abcde", "lunch"}'::char(5)[] as char_arr_,
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123'::bytea as bytea_,
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
			union all
	`
	query = strings.Repeat(query, 50)
	query += `
		select null as nil,
			'123'::varchar(255) as str,
			123::integer as int2_,
			123::decimal(24,4) as decimal_,
			123::numeric(24,4) as numeric_,
			123::bigint as bigint_,
			'{"abcde", "lunch"}'::char(5)[] as char_arr_,
			'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
			'0123'::bytea as bytea_,
			to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tl := testmodels.TableList{}
		sqldb.QueryStructs(query, &tl)
	}
}

func BenchmarkSqlDB_QueryStructsWithColumn(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	sqldb := sisql.NewSqlDB(db) // sio.SqlColumn{Name: "decimal_", Type: sio.SqlColTypeFloat64},
	// sio.SqlColumn{Name: "numeric_", Type: sio.SqlColTypeFloat64},
	// sio.SqlColumn{Name: "char_arr_", Type: sio.SqlColTypeUints8},

	for i := 0; i < b.N; i++ {

		query := `
			select null as nil,
				'123'::varchar(255) as str,
				123::integer as int2_,
				123::decimal(24,4) as decimal_,
				123::numeric(24,4) as numeric_,
				123::bigint as bigint_,
				'{"abcde", "lunch"}'::char(5)[] as char_arr_,
				'{"abcde", "lunch"}'::varchar(50)[] as varchar_arr_,
				'0123'::bytea as bytea_,
				to_timestamp('20220101121212', 'YYYYMMDDHH24MISS') as time_
		`

		tl := testmodels.TableList{}
		sqldb.QueryStructs(query, &tl)
	}
}

type Student struct {
	ID           int    `json:"id"`
	EmailAddress string `json:"email_address"`
	Name         string `json:"name"`
	Borrowed     bool   `json:"borrowed"`
}

/*
goos: darwin
goarch: arm64
pkg: github.com/wonksing/si/v2/sisql/tests/sql_test
BenchmarkSqlDB_QueryStructsStudent-8   	    1902	    647721 ns/op	    2642 B/op	      63 allocs/op
PASS

goos: windows
goarch: amd64
pkg: github.com/wonksing/si/v2/sisql/tests/sql_test
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkSqlDB_QueryStructsStudent-8         100           8626656 ns/op          992037 B/op      31466 allocs/op
PASS

goos: windows
goarch: amd64
pkg: github.com/wonksing/si/v2/sisql/tests/sql_test
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
BenchmarkSqlDB_QueryStructsStudent-8         100           7502967 ns/op          747775 B/op      25015 allocs/op
PASS
*/
func BenchmarkSqlDB_QueryStructsStudent(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	require.NotNil(b, db)

	sqldb := sisql.NewSqlDB(db, sisql.WithTagKey("json"))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		query := `select id, email_address, name, borrowed from student`

		var tl []Student
		sqldb.QueryStructs(query, &tl)
		// fmt.Println(tl)
	}
}

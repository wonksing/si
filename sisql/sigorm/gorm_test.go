package sigorm

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_OpenPostgres(t *testing.T) {
	db, _, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	_, err = OpenPostgres(db)
	require.Nil(t, err)
}

func Test_OpenPostgresWithConfig(t *testing.T) {
	db, _, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	gdb, err := OpenPostgresWithConfig(db, &gorm.Config{DryRun: true})
	require.Nil(t, err)

	require.EqualValues(t, true, gdb.DryRun)
}

func Test_NewPostgresDialector(t *testing.T) {
	d := NewPostgresDialector(postgres.Config{
		DriverName: "pgx",
	})

	require.NotNil(t, d)
}

func Test_OpenMysql(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"version"})
	rows.AddRow("5.7")
	mock.ExpectQuery("SELECT VERSION()").WithoutArgs().WillReturnRows(rows)

	_, err = OpenMysql(db)
	require.Nil(t, err)
}

func Test_OpenMysqlWithConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"version"})
	rows.AddRow("5.7")
	mock.ExpectQuery("SELECT VERSION()").WithoutArgs().WillReturnRows(rows)
	gdb, err := OpenMysqlWithConfig(db, &gorm.Config{DryRun: true})
	require.Nil(t, err)

	require.EqualValues(t, true, gdb.DryRun)
}

func Test_NewMysqlDialector(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"version"})
	rows.AddRow("5.7")
	mock.ExpectQuery("SELECT VERSION()").WithoutArgs().WillReturnRows(rows)
	d := NewMysqlDialector(mysql.Config{Conn: db})

	require.NotNil(t, d)
}

func Test_Open(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Nil(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"version"})
	rows.AddRow("5.7")
	mock.ExpectQuery("SELECT VERSION()").WithoutArgs().WillReturnRows(rows)

	d := NewMysqlDialector(mysql.Config{Conn: db})

	gdb, err := Open(d, &gorm.Config{})
	require.Nil(t, err)
	require.NotNil(t, gdb)
}

package errbadconn

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertWhenErrBadConnIsReturned(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO mytable").WillReturnError(driver.ErrBadConn)

	// when
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a transaction", err)
	}

	_, err = tx.Exec("INSERT INTO mytable(a, b) VALUES (?, ?)", "A", "B")

	// then
	assert.Equal(t, driver.ErrBadConn, err)
}
package errbadconn

import (
	"context"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertWhenErrBadConnIsReturned(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected while opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	ep := mock.ExpectPrepare("INSERT INTO mytable")
	ep.ExpectExec().WithArgs("A", "B").WillReturnError(driver.ErrBadConn)

	// when
	conn, err := db.Conn(context.Background())
	if err != nil {
		t.Errorf("an error '%s' was not expected while opening a connection", err)
	}

	_, err = conn.BeginTx(context.Background(), nil) // tx will remain open
	if err != nil {
		t.Errorf("an error '%s' was not expected while opening a transaction", err)
	}

	stmt, err := conn.PrepareContext(context.Background(), "INSERT INTO mytable(a, b) VALUES (?, ?)") // on purpose not in tx
	if err != nil {
                t.Errorf("an error '%s' was not expected while preparing statement", err)
        }

	stmt.ExecContext(context.Background(), "A", "B") // hangs...

	// then
	assert.Equal(t, driver.ErrBadConn, err)
}

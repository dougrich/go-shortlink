package shortlink

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	assert := assert.New(t)
	// Create a new data base
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: MemDBSchema,
	})
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := shortlinkTxn{db.Txn(true)}

	s, err := txn.getRedirect("/test")
	assert.NoError(err)
	assert.Nil(s)
	err = txn.setRedirect("/test", "https://test.com", http.StatusFound)
	assert.NoError(err)
	s, err = txn.getRedirect("/test")
	assert.NoError(err)
	assert.NotNil(s)
	assert.Equal("/test", s.Path)
	assert.Equal("https://test.com", s.Redirect)
	assert.Equal(http.StatusFound, s.StatusCode)
}

func TestHandlerFound(t *testing.T) {
	assert := assert.New(t)
	// Create a new data base
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: MemDBSchema,
	})
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := shortlinkTxn{db.Txn(true)}

	err = txn.setRedirect("/test", "https://test.com", http.StatusFound)
	assert.NoError(err)
	txn.Commit()

	handler := GetRedirectHandler(db, "")

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusFound, rr.Code)
	location := rr.Header().Get("Location")
	assert.Equal("https://test.com", location)
}

func TestHandlerNotFound(t *testing.T) {
	assert := assert.New(t)
	// Create a new data base
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: MemDBSchema,
	})
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := shortlinkTxn{db.Txn(true)}

	err = txn.setRedirect("/test", "https://test.com", http.StatusFound)
	assert.NoError(err)
	txn.Commit()

	handler := GetRedirectHandler(db, "")

	req, err := http.NewRequest("GET", "/test34", nil)
	assert.NoError(err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusNotFound, rr.Code)
}

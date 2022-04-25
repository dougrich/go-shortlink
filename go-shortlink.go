package shortlink

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-memdb"
)

const (
	TablenameShortlink = "shortlinks"
)

var (
	ErrGuildIDEmpty = fmt.Errorf("")
	MemDBSchema     = map[string]*memdb.TableSchema{
		"shortlinks": &memdb.TableSchema{
			Name: TablenameShortlink,
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Path"},
				},
			},
		},
	}
)

type shortlinkTxn struct {
	*memdb.Txn
}

type Shortlink struct {
	Path       string
	Redirect   string
	StatusCode int
}

func (s shortlinkTxn) getRedirect(path string) (*Shortlink, error) {
	l, err := s.First(TablenameShortlink, "id", path)
	if err != nil && err != memdb.ErrNotFound {
		return nil, err
	} else if l == nil {
		return nil, nil
	} else {
		return l.(*Shortlink), nil
	}
}

func (s shortlinkTxn) setRedirect(path string, redirect string, statusCode int) error {
	l, err := s.First(TablenameShortlink, "id", path)
	if err != nil && err != memdb.ErrNotFound {
		return err
	}
	if l != nil {
		err = s.Delete(TablenameShortlink, l)
		if err != nil && err != memdb.ErrNotFound {
			return err
		}
	}
	return s.Insert(TablenameShortlink, &Shortlink{
		Path:       path,
		Redirect:   redirect,
		StatusCode: statusCode,
	})
}

func respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

func GetRedirectHandler(db *memdb.MemDB, basepath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path[:len(basepath)] == basepath {
			path = path[len(basepath):]
		}
		txn := shortlinkTxn{db.Txn(false)}
		sl, err := txn.getRedirect(path)
		if err != nil {
			txn.Abort()
			respond(w, http.StatusInternalServerError)
		} else if sl == nil {
			txn.Commit()
			respond(w, http.StatusNotFound)
		} else {
			txn.Commit()
			http.Redirect(w, r, sl.Redirect, sl.StatusCode)
		}
	})
}

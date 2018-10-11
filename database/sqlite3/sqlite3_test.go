package sqlite3

import (
	"database/sql"
	"fmt"
	"github.com/ges-sh/migrate"
	dt "github.com/ges-sh/migrate/database/testing"
	_ "github.com/ges-sh/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	dir, err := ioutil.TempDir("", "sqlite3-driver-test")
	if err != nil {
		return
	}
	defer func() {
		os.RemoveAll(dir)
	}()
	fmt.Printf("DB path : %s\n", filepath.Join(dir, "sqlite3.db"))
	p := &Sqlite{}
	addr := fmt.Sprintf("sqlite3://%s", filepath.Join(dir, "sqlite3.db"))
	d, err := p.Open(addr)
	if err != nil {
		t.Fatalf("%v", err)
	}

	db, err := sql.Open("sqlite3", filepath.Join(dir, "sqlite3.db"))
	if err != nil {
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			return
		}
	}()
	dt.Test(t, d, []byte("CREATE TABLE t (Qty int, Name string);"))
	driver, err := WithInstance(db, &Config{})
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := d.Drop(); err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migration",
		"ql", driver)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("UP")
	err = m.Up()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

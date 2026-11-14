package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCreatesDurableSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "maulid.db")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	tables := []string{"accounts", "sessions", "posts", "likes", "audit_logs"}
	for _, table := range tables {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&count)
		if err != nil || count != 1 {
			t.Fatalf("table %s missing (err=%v count=%d)", table, err, count)
		}
	}

	db.Close()
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer reopened.Close()
	var count int
	if err := reopened.QueryRow(`SELECT COUNT(*) FROM accounts`).Scan(&count); err != nil {
		t.Fatalf("query reopened database: %v", err)
	}
}

func TestOpenRejectsEmptyPath(t *testing.T) {
	if _, err := Open(" "); err == nil {
		t.Fatal("expected empty path to fail")
	}
}

func TestMigrateRejectsNil(t *testing.T) {
	if err := Migrate(nil); err == nil {
		t.Fatal("expected nil database to fail")
	}
}

func TestMigrateReportsClosedDatabase(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	db.Close()
	if err := Migrate(db); err == nil {
		t.Fatal("expected closed database to fail")
	}
}

func TestOpenUsesPrivatePermissions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "private")
	path := filepath.Join(dir, "maulid.db")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat database: %v", err)
	}
	if fileInfo.Mode().Perm() != 0o600 {
		t.Fatalf("database permissions = %o", fileInfo.Mode().Perm())
	}
	dirInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat directory: %v", err)
	}
	if dirInfo.Mode().Perm() != 0o700 {
		t.Fatalf("database directory permissions = %o", dirInfo.Mode().Perm())
	}
}

func TestOpenReportsFilesystemFailures(t *testing.T) {
	blockingFile := filepath.Join(t.TempDir(), "blocked")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	if _, err := Open(filepath.Join(blockingFile, "maulid.db")); err == nil {
		t.Fatal("expected directory creation failure")
	}

	readOnly := filepath.Join(t.TempDir(), "readonly")
	if err := os.Mkdir(readOnly, 0o500); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer os.Chmod(readOnly, 0o700)
	if _, err := Open(filepath.Join(readOnly, "maulid.db")); err == nil {
		t.Fatal("expected database file creation failure")
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open memory database: %v", err)
	}
	defer db.Close()
	for i := 0; i < 2; i++ {
		if err := Migrate(db); err != nil {
			t.Fatalf("migrate %d: %v", i, err)
		}
	}
}

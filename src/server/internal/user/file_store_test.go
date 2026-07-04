package user

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStoreLifecycle(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	created, err := store.Create("Ada", "ADA@example.com", "secret")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != 1 || created.Email != "ada@example.com" {
		t.Fatalf("Create() = %+v", created)
	}
	if _, err := store.Create("Duplicate", "ada@example.com", "secret"); !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("duplicate Create() error = %v", err)
	}
	if _, err := store.Authenticate("ada@example.com", "secret"); err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	contents, err := os.ReadFile(filepath.Join(filepath.Dir(store.path), "users.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if bytes.Contains(contents, []byte(`"password": "secret"`)) || len(contents) == 0 {
		t.Fatal("stored password was not hashed")
	}

	users, err := store.List()
	if err != nil || len(users) != 1 {
		t.Fatalf("List() = %+v, %v", users, err)
	}
	if err := store.Delete(created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	users, err = store.List()
	if err != nil || len(users) != 0 {
		t.Fatalf("List() after delete = %+v, %v", users, err)
	}
}

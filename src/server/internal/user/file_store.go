package user

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type fileData struct {
	Users []User `json:"user_array"`
}

type FileStore struct {
	path string
	mu   sync.RWMutex
}

func NewFileStore(path string) (*FileStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}
	store := &FileStore{path: path}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := store.write(fileData{Users: []User{}}); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("inspect data file: %w", err)
	}
	return store, nil
}

func (s *FileStore) Create(name, email, password string) (PublicUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return PublicUser{}, err
	}
	email = strings.TrimSpace(strings.ToLower(email))
	for _, existing := range data.Users {
		if strings.EqualFold(existing.Email, email) {
			return PublicUser{}, ErrEmailAlreadyExists
		}
	}

	nextID := 1
	for _, existing := range data.Users {
		if existing.ID >= nextID {
			nextID = existing.ID + 1
		}
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return PublicUser{}, fmt.Errorf("hash password: %w", err)
	}
	created := User{ID: nextID, Name: strings.TrimSpace(name), Email: email, Password: string(passwordHash)}
	data.Users = append(data.Users, created)
	if err := s.write(data); err != nil {
		return PublicUser{}, err
	}
	return created.Public(), nil
}

func (s *FileStore) Authenticate(email, password string) (PublicUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := s.read()
	if err != nil {
		return PublicUser{}, err
	}
	for _, candidate := range data.Users {
		if strings.EqualFold(candidate.Email, strings.TrimSpace(email)) {
			if bcrypt.CompareHashAndPassword([]byte(candidate.Password), []byte(password)) == nil || candidate.Password == password {
				return candidate.Public(), nil
			}
		}
	}
	return PublicUser{}, ErrInvalidCredentials
}

func (s *FileStore) List() ([]PublicUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}
	users := make([]PublicUser, 0, len(data.Users))
	for _, stored := range data.Users {
		users = append(users, stored.Public())
	}
	return users, nil
}

func (s *FileStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return err
	}
	for index, stored := range data.Users {
		if stored.ID == id {
			data.Users = append(data.Users[:index], data.Users[index+1:]...)
			return s.write(data)
		}
	}
	return ErrNotFound
}

func (s *FileStore) read() (fileData, error) {
	contents, err := os.ReadFile(s.path)
	if err != nil {
		return fileData{}, fmt.Errorf("read users: %w", err)
	}
	if len(contents) == 0 {
		return fileData{Users: []User{}}, nil
	}
	var data fileData
	if err := json.Unmarshal(contents, &data); err != nil {
		return fileData{}, fmt.Errorf("decode users: %w", err)
	}
	if data.Users == nil {
		data.Users = []User{}
	}
	return data, nil
}

func (s *FileStore) write(data fileData) error {
	contents, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encode users: %w", err)
	}
	temporary := s.path + ".tmp"
	if err := os.WriteFile(temporary, contents, 0o600); err != nil {
		return fmt.Errorf("write users: %w", err)
	}
	if err := os.Rename(temporary, s.path); err != nil {
		return fmt.Errorf("replace users: %w", err)
	}
	return nil
}

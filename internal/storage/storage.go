package storage

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Book struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type Store struct {
	root string
}

func NewDefault() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get user home dir: %w", err)
	}
	return New(filepath.Join(home, ".sendtokindle"))
}

func New(root string) (*Store, error) {
	if strings.TrimSpace(root) == "" {
		return nil, errors.New("storage root is empty")
	}

	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("abs storage root: %w", err)
	}

	if err := os.MkdirAll(abs, 0o700); err != nil {
		return nil, fmt.Errorf("mkdir storage root: %w", err)
	}

	return &Store{root: abs}, nil
}

func (s *Store) Root() string {
	return s.root
}

func (s *Store) List() ([]Book, error) {
	entries, err := os.ReadDir(s.root)
	if err != nil {
		return nil, fmt.Errorf("read storage dir: %w", err)
	}

	books := make([]Book, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("stat entry %q: %w", entry.Name(), err)
		}
		if !info.Mode().IsRegular() {
			continue
		}
		books = append(books, Book{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(books, func(i, j int) bool {
		if books[i].ModTime.Equal(books[j].ModTime) {
			return books[i].Name < books[j].Name
		}
		return books[i].ModTime.After(books[j].ModTime)
	})
	return books, nil
}

func (s *Store) Open(name string) (*os.File, fs.FileInfo, error) {
	fullPath, err := s.safeJoin(name)
	if err != nil {
		return nil, nil, err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, err
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}
	if !info.Mode().IsRegular() {
		_ = f.Close()
		return nil, nil, fmt.Errorf("not a regular file: %q", name)
	}
	return f, info, nil
}

func (s *Store) Save(originalName string, r io.Reader) (string, error) {
	filename := SanitizeFilename(originalName)
	if filename == "" {
		return "", errors.New("empty filename after sanitization")
	}

	fullPath, err := s.safeJoin(filename)
	if err != nil {
		return "", err
	}

	tmpPath := fullPath + ".uploading"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer func() { _ = tmpFile.Close() }()

	if _, err := io.Copy(tmpFile, r); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, fullPath); err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("rename temp file: %w", err)
	}

	return filename, nil
}

func (s *Store) Delete(name string) error {
	fullPath, err := s.safeJoin(name)
	if err != nil {
		return err
	}
	if err := os.Remove(fullPath); err != nil {
		return err
	}
	return nil
}

func (s *Store) safeJoin(name string) (string, error) {
	filename := SanitizeFilename(name)
	if filename == "" {
		return "", errors.New("invalid filename")
	}

	full := filepath.Join(s.root, filename)
	abs, err := filepath.Abs(full)
	if err != nil {
		return "", fmt.Errorf("abs file path: %w", err)
	}

	rel, err := filepath.Rel(s.root, abs)
	if err != nil {
		return "", fmt.Errorf("rel file path: %w", err)
	}
	if rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", errors.New("path escapes storage root")
	}
	return abs, nil
}

func SanitizeFilename(input string) string {
	name := strings.TrimSpace(input)
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "\x00", "")

	name = strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':':
			return '_'
		default:
			if r < 0x20 || r == 0x7f {
				return -1
			}
			return r
		}
	}, name)

	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")
	if name == "" || name == "." || name == ".." {
		return ""
	}

	const maxLen = 200
	if len(name) > maxLen {
		name = name[:maxLen]
	}
	return name
}


package fs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

var (
	// ErrNoCurrentDir returned on creating the content file before the current
	// directory was created.
	ErrNoCurrentDir = errors.New("current dir was not created")

	// ErrInvalidDest returned when provided save path is not a valid directory.
	ErrInvalidDest = errors.New("invalid destination path")
)

type PathResolver struct {
	// Dest contains a base directory for all folders that will be created.
	// It should be a valid absolute or relative path and must already exist.
	Dest string

	absDest     string
	currentDest string
}

// NewPathResolver returns new *PathResolver along with resolving provided dest
// value. If an incorrect dest value is provided an error will be returned.
func NewPathResolver(dest string) (*PathResolver, error) {
	absDest, err := filepath.Abs(dest)
	if err != nil {
		return nil, ErrInvalidDest
	}

	stat, err := os.Stat(absDest)
	if err != nil || !stat.IsDir() {
		return nil, ErrInvalidDest
	}

	return &PathResolver{
		Dest:    dest,
		absDest: absDest,
	}, nil
}

// CreateFolder creates a new dir with the provided name within PathResolver.Dest dir.
// Created dir will be marked as current dir. If dir already exists the creation
// step will be skipped.
func (p *PathResolver) CreateFolder(name string) error {
	err := os.Chdir(p.absDest)
	if err != nil {
		return fmt.Errorf("cannot resolve path for folder %s: %w", name, err)
	}

	p.currentDest = path.Join(p.absDest, name)

	err = os.Mkdir(name, fs.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("cannot create folder %s: %w", p.currentDest, err)
	}

	return nil
}

// CreateFile creates a new file with the corresponding name and returns its
// io.WriteCloser. If FSResolver.currentDest wasn't created before, then an error
// will be returned.
func (p *PathResolver) CreateFile(name string) (io.WriteCloser, error) {
	if p.currentDest == "" {
		return nil, ErrNoCurrentDir
	}

	filePath := path.Join(p.currentDest, name)

	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create file %s: %w", filePath, err)
	}

	return f, nil
}

// Remove removes file or dir by its name. If the provided path wasn't found the
// deletion step will be skipped.
func (p *PathResolver) Remove(name string) {
	_ = os.Remove(name)
}

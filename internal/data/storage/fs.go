package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type RootPath = string

type File struct {
	Path           string
	ModifyDateTime time.Time
	IsDir          bool
}

func (f File) Name() string {
	return path.Base(f.Path)
}

type cacheEntry struct {
	hit   time.Time
	files []File
}

// A [Storage] implementation that uses the host file system as the provider.
// Output are files stored a in directory.
// Input is a relative directory path to where files are stored. Relative path starts from the directory
// that is associated on the storage initialization [Root].
type FileSystemStorage struct {
	Root RootPath

	cache map[string]cacheEntry
}

// Initializes a [FileSystemStorage] structure based on a [root] path.
// If [root] path does not exist, one will be created under the same path.
func NewFileSystemStorage(root RootPath) (*FileSystemStorage, error) {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	_, err := os.Stat(root)

	if os.IsExist(err) {
		return &FileSystemStorage{Root: root, cache: map[string]cacheEntry{}}, nil
	}

	err = os.MkdirAll(root, os.ModePerm)

	if err != nil {
		return nil, err
	}

	return &FileSystemStorage{Root: root, cache: map[string]cacheEntry{}}, nil
}

func (s *FileSystemStorage) Lookup(i string) ([]File, error) {
	d, p, err := s.dirInfo(i)

	if err != nil {
		return nil, err
	}

	if c, ok := s.cache[i]; ok && d.ModTime().After(c.hit) {
		return c.files, nil
	}

	dfs, err := os.ReadDir(p)

	if err != nil {
		return nil, err
	}

	fs := make([]File, len(dfs))
	for i, df := range dfs {
		f, err := df.Info()

		if err != nil {
			return nil, err
		}

		fs[i] = s.fileInfoToFile(p, f)
	}

	return fs, err
}

func (s *FileSystemStorage) Store(i string, fs []File) ([]File, error) {
	_, p, err := s.dirInfo(i)

	if err != nil {
		return nil, err
	}

	nfs := make([]File, len(fs))
	for i, f := range fs {
		np := path.Join(p, f.Name())
		err = Move(f.Path, path.Join(p, f.Name()))

		if err != nil {
			return nil, err
		}

		f.Path = np
		nfs[i] = f
	}

	// refresh stat to lookup new mod time
	d, _, err := s.dirInfo(i)

	if err == nil {
		s.updateCache(i, nfs, d.ModTime())
	}

	return nfs, err
}

func (s *FileSystemStorage) Delete(i string) error {
	_, p, err := s.dirInfo(i)

	if err != nil {
		return err
	}

	err = os.RemoveAll(p)

	if err == nil {
		delete(s.cache, i)
	}

	return err
}

func (s *FileSystemStorage) dirInfo(i string) (os.FileInfo, string, error) {
	if !strings.HasSuffix(i, "/") {
		i = i + "/"
	}

	p := fmt.Sprintf("%s%s", s.Root, i)
	d, err := os.Stat(p)

	if os.IsNotExist(err) {
		err = os.Mkdir(p, os.ModePerm)
	}

	return d, p, err
}

func (s *FileSystemStorage) updateCache(i string, fs []File, t time.Time) {
	e, ok := s.cache[i]
	nfs := fs

	if ok {
		nfs = append(nfs, e.files...)
	}

	s.cache[i] = cacheEntry{hit: t, files: nfs}
}

func (s FileSystemStorage) fileInfoToFile(base string, file os.FileInfo) File {
	return File{
		Path:           path.Join(base, file.Name()),
		ModifyDateTime: file.ModTime(),
		IsDir:          file.IsDir(),
	}
}

func Move(source, destination string) error {
	err := os.Rename(source, destination)
	if err != nil && strings.Contains(err.Error(), "invalid cross-device link") {
		return moveCrossDevice(source, destination)
	}
	return err
}

func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return err
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return err
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return err
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return err
	}
	os.Remove(source)
	return nil
}

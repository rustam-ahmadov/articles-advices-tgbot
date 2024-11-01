package file

import (
	"encoding/gob"
	"errors"
	"fmt"
	"m/lib/e"
	"m/lib/storage"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	defaultPerm = 0744
)

type FStorage struct {
	basePath string
}

var ErrNoSavedPages = errors.New("no saved page")

func New(basePath string) FStorage {
	return FStorage{basePath: basePath}
}

func (s FStorage) Save(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("cannot save page", err) }()

	//path to directory where file will be stored
	//ex: /documents + /jack-nicholson => /documents/jack-nicholson
	fPath := filepath.Join(s.basePath, p.UserName)

	//here we create /jack-nicholson directory
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	//hash of username and url
	fName, err := fileName(p)
	if err != nil {
		return err
	}

	//documents/jack-nicholson/23hsdlfk
	fPath = filepath.Join(fPath, fName)

	//create this file 23hsdlfk
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	//save page in file
	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}
	return nil
}

func (s FStorage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("cannot pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s FStorage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s FStorage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check is file %s exists", path)

		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func (s FStorage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)

	if err != nil {
		msg := fmt.Sprintf("can't open the file: %s", filepath)
		return nil, e.Wrap(msg, err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		msg := fmt.Sprintf("can't decode the file: %s", filepath)
		return nil, e.Wrap(msg, err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

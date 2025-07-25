package pswd

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func isFile(p string) bool {
	s, err := os.Stat(p)
	if err != nil {
		return false
	}
	if s.IsDir() {
		return false
	}
	return true
}
func isDir(p string) bool {
	s, err := os.Stat(p)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (p *Pswd) keysDir(dirs ...string) string {
	return path.Join(path.Join(dirs...), ".keys")
}

func (p *Pswd) privateKey(dirs ...string) string {
	return path.Join(p.keysDir(dirs...), "private.asc")
}
func (p *Pswd) publicKey(dirs ...string) string {
	return path.Join(p.keysDir(dirs...), "public.asc")
}

func (p *Pswd) Path(elem ...string) string {
	return path.Join(p.storagePath, path.Join(elem...))
}
func (p *Pswd) Passfile(name string) string {
	return p.Path(name) + ".asc"
}
func (p *Pswd) passfileToName(pf string) string {
	rootPath := strings.TrimPrefix(pf, p.storagePath)
	relativePath := strings.TrimPrefix(rootPath, "/")
	name := strings.TrimSuffix(relativePath, ".asc")
	return name
}

type passwordGetter = func() (string, error)

func walk(dir string, filter func(path string, d fs.DirEntry) bool) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !filter(path, d) {
			if d.IsDir() {
				return filepath.SkipDir // пропустить всю скрытую директорию
			}
			return nil // пропустить скрытый файл
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil

}

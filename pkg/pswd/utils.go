package pswd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

func (p *Pswd) getKeyIdPath(passname string) (string, error) {
	parent := filepath.Dir(p.storagePath)
	for dir := p.Path(filepath.Dir(passname)); dir != parent; dir = filepath.Dir(dir) {
		if p.isInit(p.passfileToName(dir)) {
			return keyFile(dir), nil
		}
	}
	return "", fmt.Errorf("key not found")
}
func (p *Pswd) isInit(name string) bool {
	kf := keyFile(p.storagePath, name)
	s, err := os.Stat(kf)
	if err != nil {
		return false
	}
	return !s.IsDir()
}
func (p *Pswd) getKeyIdWithPath(passname string) (string, string, error) {
	kf, err := p.getKeyIdPath(passname)
	if err != nil {
		return "", "", err
	}
	content, err := os.ReadFile(kf)
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(string(content)), kf, nil
}
func (p *Pswd) getKeyId(passname string) (string, error) {
	id, _, err := p.getKeyIdWithPath(passname)
	return id, err
}

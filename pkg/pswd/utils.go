package pswd

import (
	"os"
	"path"
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

func (p *Pswd) keysDir(dir string) string {
	return path.Join(dir, ".keys")
}

func (p *Pswd) privateKey(dir string) string {
	return path.Join(p.keysDir(dir), "private.asc")
}
func (p *Pswd) publicKey(dir string) string {
	return path.Join(p.keysDir(dir), "public.asc")
}

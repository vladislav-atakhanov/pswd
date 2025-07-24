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

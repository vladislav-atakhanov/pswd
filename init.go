package pswd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

func keyFile(dir ...string) string {
	return path.Join(append(dir, ".key-id")...)
}

func (p *Pswd) Init(name string, id string, master func(key string) (string, error), names chan string) (string, bool, error) {
	if names != nil {
		defer close(names)
	}
	dir := p.Path(name)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", false, fmt.Errorf("create passdir: %w", err)
	}
	files, err := walk(dir, func(fp string, de fs.DirEntry) bool {
		if fp == p.storagePath {
			return true
		}
		if strings.HasPrefix(de.Name(), ".") {
			return false
		}
		if de.IsDir() {
			return true
		}
		return strings.HasSuffix(de.Name(), ".asc")
	})
	if err != nil {
		return "", false, err
	}
	if len(files) == 0 {
		if err := writeKeyId(keyFile(dir), id); err != nil {
			return "", true, fmt.Errorf("write key id: %w", err)
		}
		return dir, false, nil
	}

	passwords := map[string]string{}
	oldKeys := map[string]struct{}{}
	for _, f := range files {
		name := p.passfileToName(f)
		oldId, oldPath, err := p.getKeyIdWithPath(name)
		if err != nil {
			return "", false, fmt.Errorf("key for %s not found", name)
		}
		if oldId == id {
			oldKeys[oldPath] = struct{}{}
			continue
		}
		password, err := p.ShowLazy(name, func(key string) (string, error) {
			if keypass, ok := passwords[key]; ok {
				return keypass, nil
			}
			keypass, err := master(key)
			if err != nil {
				return "", err
			}
			passwords[key] = keypass
			return keypass, nil
		})
		if err != nil {
			return "", true, err
		}
		if _, err := p.InsertWithKey(id, name, password); err != nil {
			return "", true, err
		}
		oldKeys[oldPath] = struct{}{}
		if names != nil {
			go func(name string) {
				names <- name
			}(name)
		}
	}
	dst := keyFile(dir)
	for p := range oldKeys {
		if p == dst {
			continue
		}
		if s, _ := isSubPath(path.Dir(p), dir); s {
			continue
		}
		if err := os.Remove(p); err != nil {
			return "", true, fmt.Errorf("delete key: %w", err)
		}
	}
	if err := writeKeyId(dst, id); err != nil {
		return "", true, fmt.Errorf("write key id: %w", err)
	}
	return dir, true, nil
}

func writeKeyId(dst string, id string) error {
	return os.WriteFile(dst, []byte(id), 0644)
}

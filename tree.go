package pswd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Passtype int

const (
	PassUnknown Passtype = iota
	PassFile
	PassDir
)

func (p *Pswd) Type(name string) (Passtype, error) {
	if name == "" {
		return PassDir, nil
	}
	s, err := os.Stat(p.Path(name))
	if err == nil {
		if s.IsDir() {
			return PassDir, nil
		}
		return PassUnknown, nil
	}

	s, err = os.Stat(p.Passfile(name))
	if err != nil {
		return PassUnknown, fmt.Errorf("%s is not in the %s", name, filepath.Base(p.storagePath))
	}
	return PassFile, nil
}

type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
}

func (p *Pswd) Tree(name string) (*TreeNode, error) {
	switch t, _ := p.Type(name); t {
	case PassUnknown:
		return nil, fmt.Errorf("%s is unknown", name)
	case PassFile:
		return buildTree(p.Passfile(name), true)
	case PassDir:
		return buildTree(p.Path(name), true)
	}
	return nil, fmt.Errorf("%s is unknown", name)
}

func buildTree(path string, force bool) (*TreeNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	n := info.Name()
	d := info.IsDir()

	if !force {
		if strings.HasPrefix(n, ".") {
			return nil, nil
		}
		if !d && !strings.HasSuffix(n, ".asc") {
			return nil, nil
		}
	}

	node := &TreeNode{
		Name:  n,
		IsDir: d,
	}

	if !d {
		node.Name = strings.TrimSuffix(n, ".asc")
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childNode, err := buildTree(childPath, false)
		if err != nil {
			return nil, err
		}
		if childNode != nil {
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}

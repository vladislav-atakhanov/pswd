package vault

import "fmt"

func (v *Vault) Rename(id contentKey, label string) error {
	item, ok := v.Content[id]
	if !ok {
		return fmt.Errorf("password %s not found", id.String())
	}
	item.Label = label
	v.Content[id] = item
	return nil
}

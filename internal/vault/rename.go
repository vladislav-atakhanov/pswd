package vault

import "fmt"

func (v *Vault) Rename(id contentKey, label string) error {
	item, ok := v.content[id]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, id.String())
	}
	item.Label = label
	v.content[id] = item
	return nil
}

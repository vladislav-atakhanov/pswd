package vault

import "fmt"

func (v *Vault) Rename(id contentKey, label string) error {
	if err := validateLabel(label); err != nil {
		return err
	}
	item, ok := v.content[id]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, id.String())
	}
	item.Label = label
	v.content[id] = item
	return nil
}

package vault

import (
	"io"
	"time"

	"github.com/vladislav-atakhanov/pswd/internal/uuid"
)

func (v *Vault) Add(plain io.Reader, label string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	v.Content[id] = Item{
		content:    plain,
		Label:      label,
		LastUpdate: uint64(time.Now().Unix()),
	}
	return nil
}

package vault

import "github.com/vladislav-atakhanov/pswd/internal/uuid"

type Entry struct {
	ID         uuid.V4
	Label      string
	LastUpdate uint64
}

func (v *Vault) List() []Entry {
	entries := make([]Entry, 0, len(v.content))
	for id, item := range v.content {
		entries = append(entries, Entry{
			ID:         id,
			Label:      item.Label,
			LastUpdate: item.LastUpdate,
		})
	}
	return entries
}

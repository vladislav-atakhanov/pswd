package vault

import (
	"bytes"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func (v *Vault) Remove(id contentKey) error {
	item, ok := v.content[id]
	if !ok {
		return nil
	}

	if item.length > 0 {
		v.orphanedSpans = append(v.orphanedSpans, Span{Start: item.start, Length: item.length})
	}

	if buf, ok := item.content.(*bytes.Buffer); ok {
		mem.ZeroBytes(buf.Bytes())
	}

	delete(v.content, id)
	return nil
}

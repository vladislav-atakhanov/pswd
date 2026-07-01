package vault

import (
	"bytes"
	"fmt"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func (v *Vault) Remove(id contentKey) error {
	item, ok := v.content[id]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, id.String())
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

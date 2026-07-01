package vault

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

func (v *Vault) Update(id contentKey, plain io.Reader) error {
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

	item.content = plain
	item.LastUpdate = uint64(time.Now().Unix())
	v.content[id] = item
	return nil
}

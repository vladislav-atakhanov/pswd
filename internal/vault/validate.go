package vault

import "fmt"

const MaxLabelLength = 1 << 16

func validateLabel(label string) error {
	if len(label) == 0 {
		return fmt.Errorf("%w: label is empty", ErrInvalidLabel)
	}
	if len(label) > MaxLabelLength {
		return fmt.Errorf("%w: label length %d exceeds maximum %d", ErrInvalidLabel, len(label), MaxLabelLength)
	}
	return nil
}

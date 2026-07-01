package main

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

const clipboardTimeout = 10 * time.Second

func clipboardWrite(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return err
	}
	fmt.Printf("Copied to clipboard (will clear in %v)\n", clipboardTimeout)
	go func() {
		time.Sleep(clipboardTimeout)
		clipboard.WriteAll("")
	}()
	return nil
}

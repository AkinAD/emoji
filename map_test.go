package emoji

import (
	"testing"
)

func TestEmojiMap(t *testing.T) {
	for name, code := range emojiMap {
		got := Replace(name)
		if got != code {
			t.Fatalf("test case %q fail: got: %v, expected: %v", name, got, code)
		}
	}
}

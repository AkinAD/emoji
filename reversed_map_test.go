package emoji

import (
	"fmt"
	"testing"
)

func TestEmojiExists(t *testing.T) {
	x := reverseEmojiMap["❤️"]
	if ":red_heart:" != x {
		t.Fatal("Emoji not found")
	}
	y := reverseEmojiMap["1️⃣"]
	if ":one:" != y {
		t.Fatal("Emoji not found")
	}
	z := reverseEmojiMap["❤️‍🔥"]
	if ":heart_on_fire:" != z {
		t.Fatal("Emoji not found")
	}
	testCases := []struct {
		name  string
		emoji string
		want  string
	}{
		{
			name:  "heart❤️",
			emoji: "❤️",
			want:  ":red_heart:",
		},
		{
			name:  "Number 1",
			emoji: "1️⃣",
			want:  ":one:",
		},
		{
			name:  "Complex heart on fire 1",
			emoji: "❤️‍🔥",
			want:  ":heart_on_fire:",
		},
		// {
		// 	name:  "SkinTone Bias : 👎🏿",
		// 	emoji: "👎🏿",
		// 	want:  ":one:",
		// },
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := reverseEmojiMap[tt.emoji]
			if tt.want != got {
				t.Fatalf("test case %v fail: got: %v, expected: %v", tt.name, got, tt.want)
			}
		})
	}

	msg := "country flag alias 🇺🇸"
	m := Deparse(msg)
	fmt.Printf("m: %s\n", m)
	msg = "country flag alias 🇬🇧"
	m = Deparse(msg)
	fmt.Printf("m: %s\n", m)
	msg = Parse(":flag-gb:")
	fmt.Printf("msg: %s\n", msg)

}

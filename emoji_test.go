package emoji

import (
	"testing"
)

func TestEmoji(t *testing.T) {
	tt := []struct {
		input    Emoji
		expected string
	}{
		{input: GrinningFace, expected: "\U0001F600"},
		{input: EyeInSpeechBubble, expected: "\U0001F441\uFE0F\u200D\U0001F5E8\uFE0F"},
		{input: ManGenie, expected: "\U0001F9DE\u200D\u2642\uFE0F"},
		{input: Badger, expected: "\U0001F9A1"},
		{input: FlagForTurkey, expected: "\U0001F1F9\U0001F1F7"},
	}

	for i, tc := range tt {
		got := tc.input.String()
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestEmojiWithTone(t *testing.T) {
	tt := []struct {
		input    EmojiWithTone
		tone     Tone
		expected string
	}{
		{input: WavingHand, tone: Tone(""), expected: "\U0001F44B"},
		{input: WavingHand, tone: Default, expected: "\U0001F44B"},
		{input: WavingHand, tone: Light, expected: "\U0001F44B\U0001F3FB"},
		{input: WavingHand, tone: MediumLight, expected: "\U0001F44B\U0001F3FC"},
		{input: WavingHand, tone: Medium, expected: "\U0001F44B\U0001F3FD"},
		{input: WavingHand, tone: MediumDark, expected: "\U0001F44B\U0001F3FE"},
		{input: WavingHand, tone: Dark, expected: "\U0001F44B\U0001F3FF"},
	}

	for i, tc := range tt {
		got := tc.input.Tone(tc.tone)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestEmojiWithTones(t *testing.T) {
	tt := []struct {
		input    EmojiWithTone
		tones    []Tone
		expected string
	}{
		{input: WomanAndManHoldingHands, tones: []Tone{}, expected: "\U0001f46b"},
		{input: WomanAndManHoldingHands, tones: []Tone{MediumLight}, expected: "\U0001f46b\U0001F3FC"},
		{input: WomanAndManHoldingHands, tones: []Tone{Medium, Dark}, expected: "\U0001f469\U0001F3FD\u200d\U0001f91d\u200d\U0001f468\U0001F3FF"},
	}

	for i, tc := range tt {
		got := tc.input.Tone(tc.tones...)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestCountryFlag(t *testing.T) {
	tt := []struct {
		input    string
		expected Emoji
	}{
		{input: "tr", expected: FlagForTurkey},
		{input: "TR", expected: FlagForTurkey},
		{input: "us", expected: FlagForUnitedStates},
		{input: "gb", expected: FlagForUnitedKingdom},
	}

	for i, tc := range tt {
		got, err := CountryFlag(tc.input)
		if err != nil {
			t.Fatalf("test case %v fail: %v", i+1, err)
		}
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestCountryFlagError(t *testing.T) {
	tt := []struct {
		input string
		fail  bool
	}{
		{input: "tr", fail: false},
		{input: "a", fail: true},
		{input: "tur", fail: true},
	}

	for i, tc := range tt {
		_, err := CountryFlag(tc.input)
		if (err != nil) != tc.fail {
			t.Fatalf("test case %v fail: %v", i+1, err)
		}
	}
}

func TestNewEmojiTone(t *testing.T) {
	tt := []struct {
		input    []string
		expected EmojiWithTone
	}{
		{input: nil, expected: EmojiWithTone{}},
		{input: []string{}, expected: EmojiWithTone{}},
		{input: []string{"\U0001f64b@"}, expected: PersonRaisingHand},
		{
			input:    []string{"\U0001f46b@", "\U0001f469@\u200d\U0001f91d\u200d\U0001f468@"},
			expected: WomanAndManHoldingHands,
		},
	}

	for i, tc := range tt {
		got := newEmojiWithTone(tc.input...)
		if got != tc.expected {
			t.Fatalf("test case %v fail: got: %v, expected: %v", i+1, got, tc.expected)
		}
	}
}

func TestContainsEmoji(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     bool
	}{
		{
			name:     "empty input string",
			inputStr: "",
			want:     false,
		},
		{
			name:     "string without emoji",
			inputStr: "hello! This is a simple string without any emoji",
			want:     false,
		},
		{
			name:     "numbers in string",
			inputStr: "qwerty1",
			want:     false,
		},
		{
			name:     "emoji number before string",
			inputStr: "1️⃣qwerty",
			want:     true,
		},
		{
			name:     "emoji number in string",
			inputStr: "qwerty 1️⃣",
			want:     true,
		},
		{
			name:     "several emojis and number in string",
			inputStr: "1️⃣hello world 2️⃣4️⃣ clock 8 🕓 7️⃣",
			want:     true,
		},
		{
			name:     "only emoji in string",
			inputStr: `🥰`,
			want:     true,
		},
		{
			name:     "emoji in the middle of a string",
			inputStr: `hi 😀 how r u?`,
			want:     true,
		},
		{
			name:     "emoji in the end of a string",
			inputStr: `hi! how r u doing?🤔`,
			want:     true,
		},
		{
			name:     "heart emoji in string",
			inputStr: "I ❤️ you",
			want:     true,
		},
		{
			name:     "Skin tone emoji 1",
			inputStr: "for you 👍🏿",
			want:     true,
		},
		{
			name:     "Skin tone complex emoji ",
			inputStr: "for you 👩🏾‍❤️‍👨🏿",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsEmoji(tt.inputStr); got != tt.want {
				t.Errorf("ContainsEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveEmojis(t *testing.T) {
	tests := []struct {
		name     string
		inputStr string
		want     string
	}{
		{
			name:     "string without emoji",
			inputStr: "string without emoji",
			want:     "string without emoji",
		},
		{
			name:     "string with numbers",
			inputStr: "1qwerty2",
			want:     "1qwerty2",
		},
		{
			name:     "string with emoji numbers",
			inputStr: "1️⃣qwerty2",
			want:     "qwerty2",
		},
		{
			name:     "string with emojis",
			inputStr: "❤️🛶😂",
			want:     "",
		},
		{
			name:     "string with unicode 14 emoji",
			inputStr: "te\U0001FAB7st",
			want:     "test",
		},
		{
			name:     "remove rare emojis",
			inputStr: "🧖 hello 🦋world",
			want:     "hello world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveEmojis(tt.inputStr); got != tt.want {
				t.Errorf("RemoveEmojis() = [%v], want [%v]", got, tt.want)
			}
		})
	}
}
func BenchmarkEmoji(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = WavingHand.String()
	}
}

func BenchmarkEmojiWithTone(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = WavingHand.Tone(Medium)
	}
}

func BenchmarkEmojiWithToneTwo(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = WomanAndManHoldingHands.Tone(Medium, Dark)
	}
}

func BenchmarkCountryFlag(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = CountryFlag("tr")
	}
}

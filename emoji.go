package emoji

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidTone = errors.New("tone is not a known or valid skintone")
)

// Base attributes
const (
	TonePlaceholder = "@"
	flagBaseIndex   = '\U0001F1E6' - 'a'
)

// Skin tone colors
const (
	Default     Tone = ""
	Light       Tone = "\U0001F3FB"
	MediumLight Tone = "\U0001F3FC"
	Medium      Tone = "\U0001F3FD"
	MediumDark  Tone = "\U0001F3FE"
	Dark        Tone = "\U0001F3FF"
)

// Emoji defines an emoji object with no skin variations.
type Emoji string

// String returns string representation of the simple emoji.
func (e Emoji) String() string {
	return string(e)
}

// EmojiWithTone defines an emoji object that has skin tone options.
type EmojiWithTone struct {
	oneTonedCode string
	twoTonedCode string
	defaultTone  Tone
}

// newEmojiWithTone constructs a new emoji object that has skin tone options.
func newEmojiWithTone(codes ...string) EmojiWithTone {
	if len(codes) == 0 {
		return EmojiWithTone{}
	}

	one := codes[0]
	two := codes[0]

	if len(codes) > 1 {
		two = codes[1]
	}

	return EmojiWithTone{
		oneTonedCode: one,
		twoTonedCode: two,
	}
}

// withDefaultTone sets default tone for an emoji and returns it.
func (e EmojiWithTone) withDefaultTone(tone string) EmojiWithTone {
	e.defaultTone = Tone(tone)

	return e
}

// String returns string representation of the emoji with default skin tone.
func (e EmojiWithTone) String() string {
	return strings.ReplaceAll(e.oneTonedCode, TonePlaceholder, e.defaultTone.String())
}

// Tone returns string representation of the emoji with given skin tone.
func (e EmojiWithTone) Tone(tones ...Tone) string {
	// if no tone given, return with default skin tone
	if len(tones) == 0 {
		return e.String()
	}

	str := e.twoTonedCode
	replaceCount := 1

	// if one tone given or emoji doesn't have twoTonedCode, use oneTonedCode
	// Also, replace all with one tone
	if len(tones) == 1 {
		str = e.oneTonedCode
		replaceCount = -1
	}

	// replace tone one by one
	for _, t := range tones {
		// use emoji's default tone
		if t == Default {
			t = e.defaultTone
		}

		str = strings.Replace(str, TonePlaceholder, t.String(), replaceCount)
	}

	return str
}

// Tone defines skin tone options for emojis.
type Tone string

// String returns string representation of the skin tone.
func (t Tone) String() string {
	return string(t)
}

// CountryFlag returns a country flag emoji from given country code.
// Full list of country codes: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
func CountryFlag(code string) (Emoji, error) {
	if len(code) != 2 {
		return "", fmt.Errorf("not valid country code: %q", code)
	}

	code = strings.ToLower(code)
	flag := countryCodeLetter(code[0]) + countryCodeLetter(code[1])

	return Emoji(flag), nil
}

// countryCodeLetter shifts given letter byte as flagBaseIndex.
func countryCodeLetter(l byte) string {
	return string(rune(l) + flagBaseIndex)
}

// ContainsEmoji checks whether a given string contains any emojis.
func ContainsEmoji(s string) bool {
	in := s
	var cRunes []rune

	for len(in) > 0 {
		r, size := utf8.DecodeRuneInString(in)
		cRunes = append(cRunes, r)
		c := string(cRunes)
		_, ok1 := reverseEmojiMap[c]
		_, ok2 := reverseEmojiMap[in]
		if ok1 || ok2 {
			// Found alias
			if NumberMap[c] && !numRegex.MatchString(in) {
				in = in[size:]
				continue // false alarm, found regular digit
			}
			cRunes = nil
			return true
		}

		if numRegex.MatchString(in) {
			return true
		}

		if s := RunesToHexKey([]rune{r}); len(s) >= 4 {
			in = in[size:]
			continue
		}
		// Flush cRunes if any
		if len(cRunes) > 0 {
			cRunes = nil
		}
		in = in[size:]
	}

	return false
}

// RemoveEmojis removes all emojis from the s string and returns a new string.
func RemoveEmojis(in string) string {
	var cRunes []rune
	var output strings.Builder
	// var potentialNumEmoji []rune

	for len(in) > 0 {
		in = numRegex.ReplaceAllString(in, "")
		if numRegex.MatchString(in) {
			continue
		}
		r, size := utf8.DecodeRuneInString(in)
		cRunes = append(cRunes, r)
		c := string(cRunes)
		_, ok1 := reverseEmojiMap[c]
		_, ok2 := reverseEmojiMap[in]
		if (ok1 || ok2) && !NumberMap[c] {
			// Found alias
			output.WriteString("")
			cRunes = nil
		}

		if s := RunesToHexKey([]rune{r}); len(s) >= 4 {
			in = in[size:]
			continue
		}
		// Flush cRunes if any
		if len(cRunes) > 0 {
			output.WriteString(string(cRunes))
			cRunes = nil
		}
		in = in[size:]
	}
	return strings.TrimSpace(output.String())
}

// FindAll finds all emojis in given string and return as an array of strings. If there are no emojis it returns a nil-slice.
func FindAll(in string) []string {
	var emojis []string = make([]string, 0)
	var cRunes []rune

	potentialNums := numRegex.FindAllString(in, -1)
	placeholder := "-"
	prevWasJoiner := false
	in = numRegex.ReplaceAllString(in, placeholder)

	for len(in) > 0 {
		r, size := utf8.DecodeRuneInString(in)
		cRunes = append(cRunes, r)
		c := string(cRunes)
		_, ok1 := reverseEmojiMap[c]
		_, ok2 := reverseEmojiMap[in]
		if (ok1 || ok2) && !NumberMap[c] && !isSkinTone(c) && !prevWasJoiner {
			// Found alias
			emojis = append(emojis, c)
			cRunes = nil
		}
		if c == placeholder {
			numEmoji := potentialNums[0] // get the first match
			emojis = append(emojis, numEmoji)
			potentialNums = potentialNums[1:]
			continue
		}
		if isSkinTone(c) {
			emojis[len(emojis)-1] += c
			cRunes = nil
		}
		if isJoiner(c) || prevWasJoiner {
			emojis[len(emojis)-1] += c
			cRunes = nil
			prevWasJoiner = !prevWasJoiner
			in = in[size:]
			continue
		}
		if s := RunesToHexKey([]rune{r}); len(s) >= 4 {
			in = in[size:]
			continue
		}
		// Flush cRunes if any
		if len(cRunes) > 0 {
			cRunes = nil
		}
		in = in[size:]
	}
	return emojis
}

// Checks whether an emoji in a given string has a skin tone applied
func HasTone(in string) bool {
	return isSkinTone(in)
}

// Returns one skin tone of the first emoji that has a tone applied, if none are found, an empty string is returned
func GetTone(in string) Tone {
	first := toneRegex.FindString(in)
	t, _ := matchToneToInternal(first) // error is disregarded since regex match is either valid or empty
	return t
}

// Returns all valid skin tones in a given string having emojis with tones applied
func GetAllTones(in string) []Tone {
	tones := make([]Tone, 0)
	matches := toneRegex.FindAllString(in, -1)
	for _, v := range matches {
		t, _ := matchToneToInternal(v) // error is disregarded since regex match is either valid or empty
		tones = append(tones, t)
	}
	return tones
}

// Return random emoji
func Random() string {
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// return emojiMap[r.Intn(len(e))]
	return ""
}

func isSkinTone(in string) bool {
	_, err := matchToneToInternal(in)
	ok := toneRegex.MatchString(in)
	return err == nil || ok
}

func matchToneToInternal(in string) (Tone, error) {
	switch in {
	case Light.String():
		return Light, nil
	case MediumLight.String():
		return MediumLight, nil
	case Medium.String():
		return Medium, nil
	case MediumDark.String():
		return MediumDark, nil
	case Dark.String():
		return Dark, nil
	default:
		return "", ErrInvalidTone
	}
}

func isJoiner(in string) bool {
	return joinerRegex.MatchString(in)
}

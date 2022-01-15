package emoji

import (
	"fmt"
	"strings"
	"unicode/utf8"
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
	msg := s
	var cRunes []rune

	for len(msg) > 0 {
		r, size := utf8.DecodeRuneInString(msg)
		cRunes = append(cRunes, r)
		c := string(cRunes)
		_, ok1 := reverseEmojiMap[c]
		_, ok2 := reverseEmojiMap[msg]
		if ok1 || ok2 {
			// Found alias
			if NumberMap[c] && !numRegex.MatchString(msg) {
				msg = msg[size:]
				continue // false alarm, found regular digit
			}
			cRunes = nil
			return true
		}

		if numRegex.MatchString(msg) {
			return true
		}

		if s := RunesToHexKey([]rune{r}); len(s) >= 4 {
			msg = msg[size:]
			continue
		}
		// Flush cRunes if any
		if len(cRunes) > 0 {
			cRunes = nil
		}
		msg = msg[size:]
	}

	return false
}

// RemoveEmojis removes all emojis from the s string and returns a new string.
func RemoveEmojis(msg string) string {
	var cRunes []rune
	var output strings.Builder
	// var potentialNumEmoji []rune

	for len(msg) > 0 {
		msg = numRegex.ReplaceAllString(msg, "")
		if numRegex.MatchString(msg) {
			continue
		}
		r, size := utf8.DecodeRuneInString(msg)
		cRunes = append(cRunes, r)
		c := fmt.Sprintf("%s", string(cRunes))
		_, ok1 := reverseEmojiMap[c]
		_, ok2 := reverseEmojiMap[msg]
		if (ok1 || ok2) && !NumberMap[c] {
			// Found alias
			output.WriteString("")
			cRunes = nil
		}

		if s := RunesToHexKey([]rune{r}); len(s) >= 4 {
			msg = msg[size:]
			continue
		}
		// Flush cRunes if any
		if len(cRunes) > 0 {
			output.WriteString(string(cRunes))
			cRunes = nil
		}
		msg = msg[size:]
	}
	return strings.TrimSpace(output.String())
}

// FindAll finds all emojis in given string and return as an array of strings. If there are no emojis it returns a nil-slice.
func FindAll(in string) []string {
	var emojis []string = make([]string, 0)

	var cRunes []rune
	var output strings.Builder

	// var numLookup = make(map[string]string)

	// var potentialNumEmoji []rune
	// potentialNumPos := numRegex.FindAllStringIndex(in, -1)
	potentialNums := numRegex.FindAllString(in, -1)
	// for i := 0; i < len(potentialNumPos); i++ {
	// 	numLoc[fmt.Sprintf("%d", potentialNumPos[i][0])] = potentialNums[i]
	// }
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
			output.WriteString(string(cRunes))
			cRunes = nil
		}
		in = in[size:]
	}
	return emojis
}

func isSkinTone(in string) bool {
	switch in {
	case Light.String():
		return true
	case MediumLight.String():
		return true
	case Medium.String():
		return true
	case MediumDark.String():
		return true
	case Dark.String():
		return true
	default:
		return false
	}
}
func isJoiner(in string) bool {
	return joinerRegex.MatchString(in)
}

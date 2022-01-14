package emoji

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

var (
	flagRegex = regexp.MustCompile(`^:flag-([a-zA-Z]{2}):$`)
	// numRegex  = regexp.MustCompile(`\x{FE0F}|\x{20E3}`) // match digits based on unique sequence
	// numRegex = regexp.MustCompile(`\x{FE0F}|\x{20E3}|(?i)20E3|(?i)FE0F`) // case insensitive
	// numRegex = regexp.MustCompile(`[0-9-]\x{FE0F}|\x{20E3}|(?i)20E3|(?i)FE0F`)
	// numRegex = regexp.MustCompile(`(?P<digit>\d)(\x{FE0F}|\x{20E3}|(?i)20E3|(?i)FE0F<other>)`) //
	// numRegex = regexp.MustCompile(`(?P<digit>\d)(\x{FE0F}\x{20E3}|(?i)20E3|(?i)FE0F<other>)`)        // named :match any digit emoji
	numRegex = regexp.MustCompile(`(?P<digit>\*|\#|\d)(\x{FE0F}\x{20E3}|(?i)20E3|(?i)FE0F<other>)`) // named: match any digit emoji and #️⃣*️⃣
)

type Replacer struct {
	matched bytes.Buffer
}

func NewReplacer() *Replacer {
	return &Replacer{matched: bytes.Buffer{}}
}

// Replace replaces emoji aliases (:pizza:) with unicode representation.
func (p *Replacer) Replace(input string) string {
	p.matched.Reset()
	return replaceInternal(input, &p.matched)
}

// Replace replaces emoji aliases (:pizza:) with unicode representation.
func Replace(input string) string {
	return replaceInternal(input, &bytes.Buffer{})
}

// Parse is an alias for Replace
func Parse(input string) string {
	return Replace(input)
}

// replaceInternal replaces emoji aliases (:pizza:) with unicode representation.
func replaceInternal(input string, matched *bytes.Buffer) string {
	var output strings.Builder
	output.Grow(len(input))

	for _, r := range input {
		// when it's not `:`, it might be inner or outer of the emoji alias
		if r != ':' {
			// if matched is empty, it's the outer of the emoji alias
			if matched.Len() == 0 {
				output.WriteRune(r)
				continue
			}

			matched.WriteRune(r)

			// if it's space, the alias's not valid.
			// reset matched for breaking the emoji alias
			if unicode.IsSpace(r) {
				output.WriteString(unsafeString(matched))
				matched.Reset()
			}
			continue
		}

		// r is `:` now
		// if matched is empty, it's the beginning of the emoji alias
		if matched.Len() == 0 {
			matched.WriteByte(':')
			continue
		}

		// it's the end of the emoji alias
		match := unsafeString(matched)
		matched.WriteByte(':')
		alias := unsafeString(matched)

		// check for emoji alias
		if code, ok := Find(alias); ok {
			output.WriteString(code)
			matched.Reset()
			continue
		}

		// not found any emoji
		output.WriteString(match)
		// it might be the beginning of the another emoji alias
		matched.Reset()
		matched.WriteByte(':')
	}

	// if matched not empty, add it to output
	if matched.Len() != 0 {
		output.WriteString(unsafeString(matched))
		matched.Reset()
	}

	return output.String()
}

// Map returns the emojis map.
// Key is the alias of the emoji.
// Value is the code of the emoji.
func Map() map[string]string {
	return emojiMap
}

// AppendAlias adds new emoji pair to the emojis map.
func AppendAlias(alias, code string) error {
	if c, ok := emojiMap[alias]; ok {
		return fmt.Errorf("emoji already exist: %q => %+q", alias, c)
	}

	for _, r := range alias {
		if unicode.IsSpace(r) {
			return fmt.Errorf("emoji alias is not valid: %q", alias)
		}
	}

	emojiMap[alias] = code

	return nil
}

// Exist checks existence of the emoji by alias.
func Exist(alias string) bool {
	_, ok := Find(alias)

	return ok
}

// Find returns the emoji code by alias.
func Find(alias string) (string, bool) {
	if code, ok := emojiMap[alias]; ok {
		return code, true
	}

	if flag := checkFlag(alias); len(flag) > 0 {
		return flag, true
	}

	return "", false
}

// checkFlag finds flag emoji for `flag-[CODE]` pattern
func checkFlag(alias string) string {
	if matches := flagRegex.FindStringSubmatch(alias); len(matches) == 2 {
		flag, _ := CountryFlag(matches[1])

		return flag.String()
	}

	return ""
}

// checkNumber finds a number emoji
// func checkNumber(emoji string) (string, bool) {

// 	return ""
// }

func unsafeString(matched *bytes.Buffer) string {
	buf := matched.Bytes()
	return *(*string)(unsafe.Pointer(&buf))
}
func Deparse(msg string) string {
	var cRunes []rune
	var output strings.Builder

	for len(msg) > 0 {
		// dealing with a number emoji is convoluted
		result := ReplaceAllStringSubmatchFunc(numRegex, msg, func(groups []string) string {
			//groups[1] is the digit portion of an emoji
			// E.g 7️⃣ is represented as "7\ufe0f\u20e3", the regex seperates 7 out and discards the rest (groups[2])
			// groups[0] is the original emoji. in this exapmle 7️⃣
			return reverseEmojiMap[groups[0]]
		})
		msg = result

		r, size := utf8.DecodeRuneInString(msg)
		cRunes = append(cRunes, r)
		c := fmt.Sprintf("%s", string(cRunes))

		if alias, ok := reverseEmojiMap[c]; ok && !NumberMap[c] {
			// Found alias
			normalizedStr := normalizedString(msg)
			lge, s := longestEmoji(normalizedStr)
			if lge != "" {
				output.WriteString(lge)
				size = s
			} else {
				output.WriteString(alias)
			}
			// Reset current rune
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
	return output.String()
}

// ReversedMap returns the reversed emoji map of aliases
// Key is the code of the emoji
// Value is the alias
func ReversedMap() map[string]string {
	return reverseEmojiMap
}

func FindReverse(unicode string) (string, bool) {
	if alias, ok := reverseEmojiMap[unicode]; ok {
		return alias, true
	}
	return "", false
}

// RunesToHexKey - Convert a slice of runes to hex string representation of their Unicode Code Point value
func RunesToHexKey(runes []rune) (output string) {
	// Build a slice of hex representations of each rune
	hexParts := []string{}
	for _, rune := range runes {
		hexParts = append(hexParts, fmt.Sprintf("%X", rune))
	}

	// Join the hex strings with a hypen - this is the key used in the emojis map
	output = strings.Join(hexParts, "-")
	return
}

func normalizedString(input string) string {
	runes := []rune(input)
	for idx, r := range runes {
		if hk := RunesToHexKey([]rune{r}); len(hk) < 4 {
			return string(runes[:idx+1])
		}
	}
	return input
}

func longestEmoji(normalizedStr string) (string, int) {
	runes := []rune(normalizedStr)
	size := 0
	for len(runes) > 0 {
		emoji := fmt.Sprintf("%s", string(runes))
		if alias, ok := reverseEmojiMap[emoji]; ok {
			for _, r := range runes {
				size += utf8.RuneLen(r)
			}
			return alias, size
		}
		runes = runes[:len(runes)-1]
	}
	return "", 0
}

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

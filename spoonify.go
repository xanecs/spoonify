package main

import (
	"errors"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

const (
	stateStart = iota
	stateV1
	stateV1L
	stateV1E
	stateV2
	stateV3
	stateDiph
)

func main() {
	js.Global.Set("spoon", map[string]interface{}{
		"to":   toSpoon,
		"from": fromSpoon,
	})
}

func isVowel(input rune) bool {
	return strings.ContainsRune("AEIOUÄÖÜaeiouäöü", input)
}

func isValidDiphtong(input string) bool {
	for _, diph := range []string{"ei", "au", "eu", "äu", "ie"} {
		if strings.ToLower(input) == diph {
			return true
		}
	}
	return false
}

func spoonifyVowel(input string) string {
	return string([]rune(input)[0]) + "lew" + strings.ToLower(input)
}

func toSpoon(input string) (output string) {
	output = ""
	stor := ""
	for _, c := range input {
		if !isVowel(c) {
			if len(stor) > 0 {
				output += spoonifyVowel(stor)
				stor = ""
			}
			output += string(c)
		} else {
			if len(stor) == 0 {
				stor = string(c)
			} else {
				if isValidDiphtong(stor + string(c)) {
					output += spoonifyVowel(stor + string(c))
				} else {
					output += spoonifyVowel(stor) + spoonifyVowel(string(c))
				}
				stor = ""
			}
		}
	}
	if len(stor) > 0 {
		output += spoonifyVowel(stor)
	}
	return
}

func fromSpoon(input string) (output string, err error) {
	output = ""
	err = nil
	stor := ""
	state := stateStart
	upper := false
	for _, c := range input {
		switch state {
		case stateStart:
			if isVowel(c) {
				state = stateV1
				if strings.ToLower(string(c)) != string(c) {
					upper = true
				}
			} else {
				output += string(c)
			}
		case stateV1:
			if c == rune('l') {
				state = stateV1L
			} else {
				err = errors.New("Vowel not followed by 'l'")
				return
			}
		case stateV1L:
			if c == rune('e') {
				state = stateV1E
			} else {
				err = errors.New("Vowel not followed by 'e'")
				return
			}
		case stateV1E:
			if c == rune('w') {
				state = stateV2
			} else {
				err = errors.New("Vowel not followed by 'w'")
				return
			}
		case stateV2:
			if isVowel(c) {
				stor = string(c)
				state = stateV3
			} else {
				err = errors.New("'lew' not followed by vowel")
				return
			}
		case stateV3:
			st := stor
			if upper {
				st = strings.ToUpper(st)
			}
			if isValidDiphtong(st + string(c)) {
				output += st + string(c)
				state = stateStart
			} else {
				output += st
				if isVowel(c) {
					state = stateV1
				} else {
					output += string(c)
					state = stateStart
				}
			}

			stor = ""
			upper = false

		}
	}
	if len(stor) > 0 {
		output += stor
	}
	return
}

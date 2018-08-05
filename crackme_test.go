package crackme

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

// TestVector is for challenges with Expected values
type TestVector struct {
	Challenge
	Expected string
}

var set5 = TestVector{
	Challenge: Challenge{
		Rounds:  4096,
		KeyLen:  40,
		Salt:    []byte("saltSALTsaltSALTsaltSALTsaltSALTsalt"),
		PRF:     "HMAC-SHA256",
		Pwd:     "passwordPASSWORDpassword",
		Hint:    "Set 5 https://github.com/ircmaxell/quality-checker/blob/master/tmp/gh_18/PHP-PasswordLib-master/test/Data/Vectors/pbkdf2-draft-josefsson-sha256.test-vectors",
		prfHash: sha256.New,
		Dk:      nil,
	},
	// 348c89dbcbd32b2f32d814b8116e84cf2b17347ebc1800181c4e2a1fb8dd53e1c635518c7dac47e9
	// From http://stackoverflow.com/a/5136918/1304076
	Expected: "348c89dbcbd32b2f32d814b8116e84cf2b17347ebc1800181c4e2a1fb8dd53e1c635518c7dac47e9",
}

// Pass checks of Dk is Expected
func TestPBKDF2(t *testing.T) {
	if set5.Dk == nil {
		set5.DeriveKey()
	}
	e, _ := hex.DecodeString(set5.Expected)
	if bytes.Compare(set5.Dk, e) != 0 {
		t.Error("Didn't derive expected key")
	}
}

// String for test vector challenge
func (tvec TestVector) String() string {
	if tvec.Dk == nil {
		tvec.DeriveKey()
	}
	r := fmt.Sprintf("Passwd:\t\"%s\"\n", tvec.Pwd)
	r += tvec.Challenge.String()
	r += fmt.Sprintf("Expect:\t%s\n", tvec.Expected)

	return r
}

func TestBitHint(t *testing.T) {
	type vec struct {
		in       string
		bits     int
		expected string
	}

	/* For tests I need the first byte of the sha256 hash of some strings

	Using

		for p in one two three four "governor washout beak" "glassy ubiquity absence" "splendor excel rarefy"; do
			h=$(echo -n $p | shasum -a256 | cut -b1-2)
			echo "$p:  $h"
		done

	I got

		one:  76
		two:  3f
		three:  8b
		four:  04
	*/

	vecs := []vec{
		{"one", 1, "0b0"},
		{"two", 1, "0b0"},
		{"three", 1, "0b1"},
		{"four", 1, "0b0"},

		{"one", 2, "0b01"},
		{"two", 2, "0b00"},
		{"three", 2, "0b10"},
		{"four", 2, "0b00"},

		{"one", 3, "0b011"},
		{"two", 3, "0b001"},
		{"three", 3, "0b100"},
		{"four", 3, "0b000"},

		{"governor washout beak", 1, "0b0"},
		{"glassy ubiquity absence", 1, "0b1"},
		{"splendor excel rarefy", 1, "0b0"},

		{"governor washout beak", 2, "0b01"},
		{"glassy ubiquity absence", 2, "0b11"},
		{"splendor excel rarefy", 2, "0b01"},

		{"governor washout beak", 3, "0b011"},
		{"glassy ubiquity absence", 3, "0b111"},
		{"splendor excel rarefy", 3, "0b010"},
	}

	for _, v := range vecs {
		result := MakeBitHint(v.in, v.bits)
		if result != v.expected {
			t.Errorf("For s = %q expected %q but got %q", v.in, v.expected, result)
		}
	}

}

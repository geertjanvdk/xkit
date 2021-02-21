// Copyright (c) 2021, Geert JM Vanderkelen

package xid

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/bits"
	"strings"
)

var urlSafeAlphabet = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")
var nanoDefaultSize = 21
var nanoMinSize = 4
var nanoMaxSize = nanoDefaultSize

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
const base56Alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"

type nanoID struct {
	size     int
	alphabet []rune
}

// NanoID is implements Nano ID found at https://github.com/ai/nanoid.
//
// To configure the size, use chaining:
//
//     id := NanoID().SetSize(12).Get().
//
// To set a different alphabet method chaining:
//
//     fmt.Println(NanoID().SetAlphabet("base56").SetSize(13)).
//
// When anything goes wrong generating the Nano ID, the empty string is returned.
func NanoID() *nanoID {
	return &nanoID{
		size:     nanoDefaultSize,
		alphabet: urlSafeAlphabet,
	}
}

// SetSize sets the size of the Nano ID.
// Panics when size is lower than 4 and higher than 21.
func (n *nanoID) SetSize(size int) *nanoID {
	if size < nanoMinSize || size > nanoDefaultSize {
		panic(fmt.Sprintf("(nanoid) size must be %d<=size<=%d", nanoMinSize, nanoMaxSize))
	}
	n.size = size
	return n
}

// SetAlphabet sets the alphabet according to string s. Some named
// alphabets are available when setting s to either "base62", "base58",
// or "base56". If s is empty, the default Nano ID alphabet is used (which
// is URL safe).
func (n *nanoID) SetAlphabet(s string) *nanoID {
	switch strings.ToLower(s) {
	case "base62":
		n.alphabet = []rune(base62Alphabet)
	case "base58":
		n.alphabet = []rune(base58Alphabet)
	case "base56":
		n.alphabet = []rune(base56Alphabet)
	case "":
		n.alphabet = urlSafeAlphabet
	default:
		n.alphabet = []rune(s)
	}

	return n
}

func (n *nanoID) generate() string {
	alphaLen := len(n.alphabet)

	mask := (2 << (31 - bits.LeadingZeros32(uint32(alphaLen-1|1)))) - 1
	step := int(math.Ceil((1.6 * float64(mask*n.size)) / float64(alphaLen)))
	id := make([]rune, n.size)
	b := make([]byte, step)

	var j int
	for {
		if _, err := rand.Read(b); err != nil {
			return ""
		}
		for i := step - 1; i > 0; i-- {
			pos := int(b[i] & byte(mask))
			if pos < alphaLen {
				id[j] = n.alphabet[pos]
				j++
				if j == n.size {
					return string(id)
				}
			}
		}
	}
}

// String returns the Nano ID as string. Note that executing
// this function will always re-generate a new Nano ID.
func (n nanoID) String() string {
	return n.generate()
}

// Get returns the Nano ID as string. Note that executing
// this function will always re-generate a new Nano ID.
func (n nanoID) Get() string {
	return n.generate()
}

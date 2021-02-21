// Copyright (c) 2020, 2021, Geert JM Vanderkelen

package xid

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

const (
	uuidv1 = iota + 1
	_
	_
	uuidv4
	uuidv5
)

const (
	uuidVariantRFC4122 = iota + 1
)

var reVersion = map[int]*regexp.Regexp{
	// UUID v1
	uuidv1: regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}`),
	// UUID v4
	uuidv4: regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`),
	// UUID v5
	uuidv5: regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}`),
}

// NilUUID is a UUID with all bits set to 0.
// As a string, it is represented as "00000000-0000-0000-0000-000000000000".
var NilUUID = UUID{}

// UUIDNamespace defines a UUID v5 namespace.
type UUIDNamespace struct {
	Name string
	UUID UUID
}

var (
	UUIDURLNamespace = UUIDNamespace{
		Name: "URL",
		UUID: UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
	}
)

// UUID defines a type big enough to hold a a RFC 4122 Universally Unique IDentifier.
type UUID [16]byte

// CodeAsHex returns the 36 character long string string representation of the UUID.
// It looks like `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
func (u UUID) String() string {
	// using `fmt.Sprintf` is not as efficient when creating lots of UUID
	// represented as string.
	us := make([]byte, 36)

	hex.Encode(us[0:8], u[0:4]) // time-low, 4 bytes
	us[8] = '-'
	hex.Encode(us[9:13], u[4:6]) // time-mid, 2 bytes
	us[13] = '-'
	hex.Encode(us[14:18], u[6:8]) // time-high-and-version, 2 bytes
	us[18] = '-'
	hex.Encode(us[19:23], u[8:10]) // clock-seq-and-reserved + clock-seq-low, 2 bytes
	us[23] = '-'
	hex.Encode(us[24:], u[10:]) // node

	return string(us)
}

// Version returns the version number of the UUID as defined by RFC 4122.
func (u UUID) Version() int {
	return int(u[6] >> 4)
}

// IsNil returns true if all bits are zero.
func (u UUID) IsNil() bool {
	for _, o := range u {
		if o != 0 {
			return false
		}
	}
	return true
}

// MarshalJSON returns the JSON encoding of u.
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.String() + "\""), nil
}

// UnmarshalJSON parses the JSON-encoded data and stores the result in u.
func (u *UUID) UnmarshalJSON(data []byte) error {
	n := UUIDFromString(strings.Trim(string(data), "\""))

	if n.IsNil() {
		return fmt.Errorf("invalid UUID string representation")
	}

	*u = n
	return nil
}

// UUIDFromString creates a UUID instance from a UUID represented by the string s.
// If s is not a UUID, the nil UUID (all bits 0) is returned.
// The string can be with or without hyphens.
func UUIDFromString(s string) UUID {
	// decode the string of hexadecimal into bytes
	us := strings.Replace(s, "-", "", -1)
	octets, err := hex.DecodeString(us)
	if err != nil || len(octets) != 16 {
		return NilUUID
	}

	u := UUID{}
	for i, o := range octets {
		u[i] = o
	}

	return u
}

// UUIDv4 returns a UUID version 4 (random).
// See https://tools.ietf.org/html/rfc4122.
func UUIDv4() UUID {
	u := UUID{}

	_, _ = rand.Read(u[:]) // what could fail...

	encodeUUIDVersion(&u, 4)
	encodeUUIDVariant(&u, uuidVariantRFC4122)

	return u
}

// UUIDv5 returns a UUID using namespace and name.
// See https://tools.ietf.org/html/rfc4122.
func UUIDv5(namespace UUIDNamespace, name string) UUID {
	u := UUID{}

	cs := sha1.New()
	cs.Write(namespace.UUID[:])
	cs.Write([]byte(name))
	copy(u[:], cs.Sum(nil))

	encodeUUIDVersion(&u, 5)
	encodeUUIDVariant(&u, uuidVariantRFC4122)

	return u
}

// UUIDIsValid returns whether the string represents a UUID v1, or v4. Note that value is not
// trimmed of spaces.
func UUIDIsValid(value string) bool {
	if len(value) != 36 {
		return false
	}

	for _, ver := range reVersion {
		if ver.Match([]byte(value)) {
			return true
		}
	}

	return false
}

// IsUUIDV1 returns true if value is UUID v1.
func IsUUIDV1(value string) bool {
	if len(value) != 36 {
		return false
	}

	return reVersion[uuidv1].Match([]byte(value))
}

// IsUUIDV4 returns true if value is UUID v4.
func IsUUIDV4(value string) bool {
	if len(value) != 36 {
		return false
	}

	return reVersion[uuidv4].Match([]byte(value))
}

// IsUUIDV5 returns true if value is UUID v4.
func IsUUIDV5(value string) bool {
	if len(value) != 36 {
		return false
	}

	return reVersion[uuidv5].Match([]byte(value))
}

func encodeUUIDVersion(u *UUID, version int) {
	// Version number is in the most significant 4 bits of the time
	// stamp (bits 4 through 7 of the time_hi_and_version field)
	u[6] = byte((int(u[6]) & 0x0f) | version<<4)
}

func encodeUUIDVariant(u *UUID, variant int) {
	switch variant {
	case uuidVariantRFC4122:
		// RFC 4122 variant: Set the two most significant bits (bits 6 and 7) of the
		// clock_seq_hi_and_reserved to zero and one, respectively.
		fallthrough
	default:
		u[8] = u[8]&(0xff>>2) | 0x80
	}
}

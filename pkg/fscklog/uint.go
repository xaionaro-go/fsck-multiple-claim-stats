package fscklog

import (
	"math"
	"strconv"
)

type Uint64Null struct {
	Uint64 uint64
	Valid  bool
}

var maxUint8Len uint
var maxUint8FirstChar uint8
var maxUint64Len uint
var maxUint64FirstChar uint8

func init() {
	maxUint8String := strconv.FormatUint(math.MaxUint8, 10)
	maxUint8Len = uint(len(maxUint8String))
	maxUint8FirstChar = maxUint8String[0]
	maxUint64String := strconv.FormatUint(math.MaxUint64, 10)
	maxUint64Len = uint(len(maxUint64String))
	maxUint64FirstChar = maxUint64String[0]
}

func isUint8(s string) bool {
	if uint(len(s)) > maxUint8Len {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	if uint(len(s)) == maxUint8Len && s[0] == maxUint8FirstChar {
		_, err := strconv.ParseUint(s, 10, 8)
		return err == nil
	}
	return true
}

func isUint64(s string) bool {
	if uint(len(s)) > maxUint64Len {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	if uint(len(s)) == maxUint64Len && s[0] == maxUint64FirstChar {
		_, err := strconv.ParseUint(s, 10, 64)
		return err == nil
	}
	return true
}

func mustParseUint8(s string) uint8 {
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		panic(err)
	}
	return uint8(v)
}

func mustParseUint64(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

package vkconv

import (
	"fmt"
	"strings"
)

func VarcharAsString(bytes [256]byte) string {
	r := strings.Builder{}

	for _, b := range bytes {
		if b == byte(0) { // skip null bytes
			continue
		}

		r.WriteByte(b)
	}

	return strings.TrimSpace(r.String())
}

func StringAsVarchar(src string) [256]byte {
	if len(src) > 256 {
		panic(fmt.Errorf("failed convert string '%s' to vulkan varchar: len is greater than 256", src))
	}

	label := [256]byte{}
	for i, b := range src {
		label[i] = byte(b)
	}

	return label
}

// NormalizeString convert string(varchar) -> [256]byte -> string(256)
func NormalizeString(src string) string {
	return VarcharAsString(StringAsVarchar(src))
}

func NormalizeStringList(list []string) []string {
	normalized := make([]string, 0, len(list))

	for _, src := range list {
		normalized = append(normalized, NormalizeString(src))
	}

	return normalized
}

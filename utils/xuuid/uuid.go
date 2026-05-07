package xuuid

import (
	"github.com/google/uuid"
	"strings"
)

func ShortUuid() string {
	return Uuid()[0:16]
}
func Uuid() string {
	u, _ := uuid.New().MarshalText()
	uuidStr := strings.Replace(string(u), "-", "", -1)
	return string(uuidStr)
}

func UuidToUpper() string {
	return strings.ToLower(Uuid())
}

func ShortUuidToUpper() string {
	return UuidToUpper()[0:16]
}

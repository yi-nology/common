package xmapping

import (
	"github.com/yi-nology/common/tools/encoding"
	"io"
)

// UnmarshalTomlBytes unmarshals TOML bytes into the given v.
func UnmarshalTomlBytes(content []byte, v interface{}, opts ...UnmarshalOption) error {
	b, err := encoding.TomlToJson(content)
	if err != nil {
		return err
	}

	return UnmarshalJsonBytes(b, v, opts...)
}

// UnmarshalTomlReader unmarshals TOML from the given io.Reader into the given v.
func UnmarshalTomlReader(r io.Reader, v interface{}, opts ...UnmarshalOption) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return UnmarshalTomlBytes(b, v, opts...)
}

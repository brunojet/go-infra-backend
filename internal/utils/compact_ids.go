package utils

import "encoding/base64"

// compositeMaxPartSize defined in consts.go

func EncodeCompositeKey[T AnyInt](srcs ...T) (string, error) {
	if len(srcs) == 0 {
		return "", errCompositeKeyValuesRequired
	}
	var err error
	buf := make([]byte, 0, len(srcs)*compositeMaxPartSize)
	for _, src := range srcs {
		buf, err = appendVarintToBuffer(buf, src)
		if err != nil {
			return "", err
		}
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func DecodeCompositeKey[T AnyInt](encoded string, dsts ...*T) error {
	if len(dsts) == 0 {
		return errCompositeKeyValuesRequired
	}
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return errInvalidCompositeKeyEncoding(err)
	}
	idx := 0
	for _, d := range dsts {
		nextIdx, err := readVarintFromBufferAt(payload, idx, d)
		if err != nil {
			return err
		}
		idx = nextIdx
	}
	return nil
}

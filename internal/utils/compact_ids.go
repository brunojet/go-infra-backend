package utils

import "encoding/base64"

// compositeMaxPartSize defined in consts.go

func EncodeCompositeKey[T AnyInt](srcs ...T) (string, error) {
	if len(srcs) == 0 {
		return "", errCompositeKeyValuesRequired
	}
	buf := make([]byte, 0, len(srcs)*compositeMaxPartSize)
	for _, src := range srcs {
		if isReadyForCompressedEncoding(src) {
			if err := appendVarintCompressed(&buf, src); err != nil {
				return "", err
			}
		} else {
			if err := appendVarintUncompressed(&buf, src); err != nil {
				return "", err
			}
		}
	}
	used := len(buf)
	return base64.RawURLEncoding.EncodeToString(buf[:used]), nil
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
		nextIdx, err := restoreVarIntFromBuffer(payload, idx, d)
		if err != nil {
			return err
		}
		idx = nextIdx
	}
	return nil
}

package utils

import "encoding/binary"

// convertInt64ToTPtr converts an int64 to the specified integer type T and stores it in dst.
// It returns an error if dst is nil, src is negative, or if the value cannot be represented in type T.
func convertInt64ToTPtr[T AnyInt](src int64, dst *T) error {
	if dst == nil {
		return errDstNil
	}
	if src < 0 {
		return errInvalidCompactInt
	}
	v, err := convertInt64ToT[T](src)
	if err != nil {
		return err
	}
	*dst = v
	return nil
}

// reserveAndWrite checks available capacity, reslices dst to make room for
// `needed` bytes, and then calls writeFn(pos) to perform the actual write at
// the starting position. It returns an error when capacity is insufficient.
func reserveAndWrite(dst *[]byte, needed int, writeFn func(pos int)) error {
	if dst == nil {
		return errDstNil
	}
	free := cap(*dst) - len(*dst)
	if free < needed {
		return errInsufficientBuffer(needed)
	}
	old := len(*dst)
	*dst = (*dst)[:old+needed]
	writeFn(old)
	return nil
}

func isReadyForCompressedEncoding[T AnyInt](src T) bool {
	switch any(src).(type) {
	case int16:
		return false
	case int32:
		vv := int64(src)
		return byte(uint64(vv)>>shiftInt16ToTopByte) == zeroByte
	default:
		vv := int64(src)
		return byte(uint64(vv)>>shiftInt64ToTopByte) == zeroByte
	}
}

// appendVarintCompressed writes a zig-zag varint for value into *dst using available
// capacity; it does not allocate. If there isn't enough free capacity it
// returns an error. The caller must ensure dst is non-nil.
func appendVarintCompressed[T AnyInt](dst *[]byte, src T) error {
	vv := int64(src)
	if vv <= 0 {
		return errInvalidCompactInt
	}
	// zig-zag encode and emit as 7-bit LEB128 groups with 1-byte prefix
	ux := uint64((vv << 1) ^ (vv >> zigZagSignShift))
	needed := 1
	tmp := ux
	for tmp >= uint64(compactVarintMsb) {
		needed++
		tmp >>= 7
	}
	if needed > compactVarintMaxBytes {
		return errCompactVarintTooLarge(needed)
	}
	// inline reservation to avoid closure allocation in hot path
	dataBytes := needed
	total := dataBytes + 1
	if dst == nil {
		return errDstNil
	}
	free := cap(*dst) - len(*dst)
	if free < total {
		return errInsufficientBuffer(total)
	}
	old := len(*dst)
	*dst = (*dst)[:old+total]
	pos := old
	(*dst)[pos] = byte(compactVarintPrefixBase + byte(dataBytes))
	pos++
	dataBytes-- // adjust to number of data bytes following prefix
	for i := 0; i < dataBytes; i++ {
		(*dst)[pos+i] = byte(ux) | compactVarintMsb
		ux >>= 7
	}
	(*dst)[pos+dataBytes] = byte(ux)
	return nil
}

func appendVarintUncompressed[T AnyInt](dst *[]byte, src T) error {
	vv := int64(src)
	if vv < 0 {
		return errInvalidCompactInt
	}
	if dst == nil {
		return errDstNil
	}
	needed := getIntTypeLen[T]()
	// inline reservation to avoid closure allocation in hot path
	free := cap(*dst) - len(*dst)
	if free < needed {
		return errInsufficientBuffer(needed)
	}
	old := len(*dst)
	*dst = (*dst)[:old+needed]
	pos := old
	// write big-endian using encoding/binary for common sizes
	switch needed {
	case int16Size:
		binary.BigEndian.PutUint16((*dst)[pos:pos+int16Size], uint16(vv))
	case int32Size:
		binary.BigEndian.PutUint32((*dst)[pos:pos+int32Size], uint32(vv))
	default:
		binary.BigEndian.PutUint64((*dst)[pos:pos+int64Size], uint64(vv))
	}
	return nil
}

func restoreVarIntFromBuffer[T AnyInt](src []byte, startOffset int, dst *T) (int, error) {
	if startOffset < 0 || startOffset >= len(src) {
		return 0, errOffsetOutOfRange
	}
	prefix := src[startOffset]
	if prefix == compactVarintMsb || prefix > compactVarintPrefixMax {
		return 0, errInvalidCompactIntPrefix(prefix)
	}
	var headerLen int
	var dataLen int
	if prefix > compactVarintPrefixBase {
		headerLen = 1
		dataLen = int(prefix - compactVarintPrefixBase)
	} else {
		dataLen = getIntTypeLen[T]()
	}
	dataStart := startOffset + headerLen
	endOffset := dataStart + dataLen
	if endOffset > len(src) {
		return 0, errInvalidCompactIntPayload
	}
	var value int64
	if prefix > compactVarintPrefixBase {
		// compressed: data are LEB128 little-endian 7-bit groups
		var ux uint64
		shift := 0
		for i := dataStart; i < endOffset; i++ {
			b := src[i]
			ux |= uint64(b&0x7F) << uint(shift)
			shift += 7
		}
		// zig-zag decode
		// original: ux = uint64((vv << 1) ^ (vv >> 63))
		// decode back to signed int64
		vv := int64((ux >> 1) ^ uint64((-(ux & 1))))
		value = vv
	} else {
		// uncompressed: big-endian integer bytes
		for i := dataStart; i < endOffset; i++ {
			value = (value << byteSizeShift) | int64(src[i])
		}
	}
	if err := convertInt64ToTPtr(value, dst); err != nil {
		return 0, err
	}
	return endOffset, nil
}

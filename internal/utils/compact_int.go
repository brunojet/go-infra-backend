package utils

// Pacote utils provê utilitários para manipulação eficiente de inteiros compactados em buffers binários.
//
// Este arquivo implementa funções para serializar e desserializar inteiros (int16, int32, int64) em formato compacto,
// otimizando o uso de espaço em buffers binários, especialmente para valores pequenos.
//
// Funções principais:
//   - appendVarintToBuffer: Serializa um inteiro para um buffer de bytes em formato compacto.
//   - readVarintFromBufferAt: Lê um inteiro compactado de um buffer de bytes a partir de um offset.
//
// Uso típico: serialização/deserialização de identificadores, índices ou outros valores inteiros em protocolos binários customizados.

// appendVarintToBuffer serializa um valor inteiro v do tipo T para o slice de bytes dst,
// utilizando um formato compacto:
//   - Para valores zero, grava apenas o header compactIntHeaderMask (0x80).
//   - Para valores pequenos, grava apenas os bytes significativos e um header compactado (MSB 1, 7 bits = n de bytes).
//   - Para valores grandes, grava todos os bytes do tipo.
//
// Parâmetros:
//   - dst: slice de bytes de destino (deve ter capacidade suficiente).
//   - v: valor inteiro a ser serializado (int16, int32 ou int64).
//
// Retorna:
//   - Slice de bytes resultante com o valor serializado.
//   - Erro caso o buffer seja insuficiente ou o valor seja inválido.
func appendVarintToBuffer[T AnyInt](dst []byte, v T) ([]byte, error) {
	needed := getIntTypeLen[T]()
	free := cap(dst) - len(dst)

	if free < needed {
		return dst, errInsufficientBuffer(needed)
	}

	vv := int64(v)
	if vv < 0 {
		return dst, errInvalidCompactInt
	}

	startOffset := len(dst)

	dst = dst[:startOffset+needed]

	isCompressReady := needed > 2 && byte(vv>>((needed*8)-8)) == 0x00
	pos := startOffset
	started := true

	if isCompressReady {
		pos += 1
		started = false
	}

	for i := needed - 1; i >= 0; i-- {
		b := byte(vv >> (i * 8))
		if !started {
			if b == 0 {
				continue
			}
			started = true
		}
		dst[pos] = b
		pos++
	}

	if !started {
		dst = dst[:startOffset+1]
		dst[startOffset] = compactIntHeaderMask
		return dst, nil
	} else if isCompressReady {
		n := pos - startOffset - 1
		dst[startOffset] = compactIntHeaderMask | byte(n)
		dst = dst[:pos]
	}

	return dst, nil
}

// readVarintFromBufferAt desserializa um valor inteiro do tipo T a partir do slice de bytes src,
// começando no offset informado, armazenando o resultado em out.
//
// Interpretação:
//   - Se o header for compactIntHeaderMask (0x80), retorna zero.
//   - Se o header tiver MSB 1, os 7 bits menos significativos indicam quantos bytes seguem.
//   - Caso contrário, lê todos os bytes do tipo.
//
// Parâmetros:
//   - src: slice de bytes de origem.
//   - offset: posição inicial para leitura.
//   - out: ponteiro para variável onde o valor lido será armazenado.
//
// Retorna:
//   - Quantidade de bytes consumidos na leitura.
//   - Erro caso o buffer seja insuficiente ou o formato seja inválido.
func readVarintFromBufferAt[T AnyInt](src []byte, offset int, out *T) (int, error) {
	needed := getIntTypeLen[T]()
	bufLen := len(src)

	if offset > bufLen {
		return 0, errOffsetOutOfRange
	}

	first := src[offset]

	// Header zero: valor zero
	if first == compactIntHeaderMask {
		*out = T(0)
		return offset + 1, nil // consome só o header
	} else if first > compactIntHeaderMask { // Header compactado: MSB 1
		realBytes := int(first & compactIntRealBytesMask)
		if realBytes+1 > needed {
			return 0, errVarintCompactadoInvalido
		}
		offset++ // pula o header
		needed = realBytes
	}
	if offset+needed > bufLen {
		return 0, errBufferInsuficiente
	}
	var v int64
	for i := 0; i < needed; i++ {
		v = (v << 8) | int64(src[offset+i])
	}
	*out = T(v)
	return offset + needed, nil
}

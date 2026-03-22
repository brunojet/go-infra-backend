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
package utils

// appendVarintToBuffer serializa um valor inteiro v do tipo T para o slice de bytes dst,
// utilizando um formato compacto que economiza espaço para valores pequenos.
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
		dst = dst[:1]
		dst[len(dst)-1] = 0x80
		return dst, nil
	} else if isCompressReady {
		n := pos - startOffset - 1
		dst[startOffset] = 0x80 | byte(n)
		// Ajusta o comprimento do slice para refletir apenas os bytes realmente usados
		dst = dst[:pos]
	}

	return dst, nil
}

// readVarintFromBufferAt desserializa um valor inteiro do tipo T a partir do slice de bytes src,
// começando no offset informado, armazenando o resultado em out.
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
	pos := offset

	// Regra customizada: header compactIntZeroHeader significa valor zero
	if first == compactIntZeroHeader {
		*out = T(0)
		return 1, nil // consome só o header
	} else if first > compactIntZeroHeader { // compactado: header tem MSB 0, os 7 bits restantes indicam quantos bytes seguem
		realBytes := int(first & 0x7F)
		if realBytes+1 > needed {
			return 0, errVarintCompactadoInvalido
		}
		pos++ // pula o header
		needed = realBytes
	}

	if pos+needed > bufLen {
		return 0, errBufferInsuficiente
	}
	var v int64
	for i := 0; i < needed; i++ {
		v = (v << 8) | int64(src[pos+i])
	}
	*out = T(v)
	consumed := pos + needed - offset
	return consumed, nil
}

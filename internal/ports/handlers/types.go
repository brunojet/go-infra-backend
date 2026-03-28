package handlers

import "regexp"

// Regex para validação de IDs em handlers. Ex: para validar IDs do tipo int64, use Int64GtZero, para int32 use Int32GtZero, etc.
// As regex abaixo validam apenas o formato e o número máximo de dígitos para cada tipo,
// mas não garantem que o valor está dentro do range exato do tipo. Para garantir,
// faça a conversão para int64/int32/int16 e compare o valor.
var (
	// int64: até 19 dígitos, mas precisa validar valor <= 9223372036854775807
	Int64GtZero = regexp.MustCompile(`^[1-9][0-9]{0,18}$`)
	// int32: até 10 dígitos, mas precisa validar valor <= 2147483647
	Int32GtZero = regexp.MustCompile(`^[1-9][0-9]{0,9}$`)
	// int16: até 5 dígitos, mas precisa validar valor <= 32767
	Int16GteZero = regexp.MustCompile(`^(0|[1-9][0-9]{0,4})$`)
	// base64 url safe (sem padding)
	Base64UrlSafe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

type HandlerParameters struct {
	HandlerPath      string
	IDValidationRule *regexp.Regexp
}

type PaginationResponse struct {
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	TotalItems int64  `json:"totalItems"`
	OrderBy    string `json:"orderBy,omitempty"`
	Order      string `json:"order,omitempty"`
}

type ListResponses[R any] struct {
	Data       []R                `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

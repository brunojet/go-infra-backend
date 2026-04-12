package streams

import (
	"errors"
	"fmt"
)

var (
	errMarshalRequestBody  = errors.New("bffclient/streams: marshal request body")
	errDecodeResponse      = errors.New("bffclient/streams: decode response")
	errDownloadFile        = errors.New("bffclient/streams: download file")
	errUnexpectedNoContent = errors.New("bffclient/streams: JsonResponseStream used with 204 NoContent — use NoBodyResponseStream instead")
)

func wrapErr(sentinel, cause error) error {
	return fmt.Errorf("%w: %w", sentinel, cause)
}

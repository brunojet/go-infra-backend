package streams

import (
	"io"
)

// ---------------------------------------------------------------------------
// FileDownloadStream
// ---------------------------------------------------------------------------

// FileDownloadStream implements BffHttpResponseStream for binary/file downloads.
//
// Decode copies the upstream response body directly to the provided io.Writer
// (e.g. an http.ResponseWriter or an os.File) without buffering the entire
// body in memory — suitable for large file responses.
//
// Use this instead of JsonResponseStream when the upstream returns raw bytes
// (images, PDFs, archives, etc.) rather than a JSON-encoded payload.
type FileDownloadStream struct {
	httpResponse
	w io.Writer
}

// NewFileDownloadStream creates a FileDownloadStream that pipes the response
// body into w. w must not be nil.
func NewFileDownloadStream(w io.Writer) *FileDownloadStream {
	return &FileDownloadStream{httpResponse: NewHttpResponse(0), w: w}
}

// Decode copies the response body to the underlying writer.
func (s *FileDownloadStream) Decode(r io.Reader) error {
	if _, err := io.Copy(s.w, r); err != nil {
		return wrapErr(errDownloadFile, err)
	}
	return nil
}

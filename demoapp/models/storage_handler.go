package models

import "context"

// StorageHandler é a interface mínima que serviços de armazenamento devem implementar.
// Implementações possíveis: DefaultStorageService (local filesystem), S3StorageService, etc.
type StorageHandler interface {
	DeleteFile(ctx context.Context, path string) error
	FileExists(ctx context.Context, path string) bool
	UploadFile(ctx context.Context, path string, data []byte) error
}

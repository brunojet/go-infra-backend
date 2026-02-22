package adapters

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestLocalS3EventQueue_Start(t *testing.T) {
	keyPath := "unsigned"
	bucket := "test-bucket"
	fullPath := filepath.Join(os.TempDir(), bucket, keyPath)
	os.RemoveAll(fullPath)
	os.MkdirAll(fullPath, 0755)
	defer os.RemoveAll(fullPath)

	queue := NewLocalS3EventQueue(context.Background(), bucket, keyPath)

	received := make(chan events.S3Event, 1)
	stop := queue.Start(func(event any) {
		s3evt, ok := event.(events.S3Event)
		if ok {
			received <- s3evt
		}
	})
	defer stop()

	// Cria arquivo para simular upload
	key := "file1.txt"
	filePath := filepath.Join(fullPath, key)
	os.WriteFile(filePath, []byte("conteudo de teste"), 0644)

	select {
	case evt := <-received:
		assert.Equal(t, 1, len(evt.Records))
		rec := evt.Records[0]
		assert.Equal(t, "ObjectCreated:Put", rec.EventName)
		assert.Equal(t, bucket, rec.S3.Bucket.Name)
		assert.Equal(t, filepath.Join(keyPath, key), rec.S3.Object.Key)
		assert.NotZero(t, rec.S3.Object.Size)
		assert.NotEmpty(t, rec.S3.Object.ETag)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout esperando evento S3")
	}
}

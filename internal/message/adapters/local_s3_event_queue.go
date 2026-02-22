package adapters

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/fsnotify/fsnotify"
)

// LocalS3EventQueue simula recebimento de eventos SQS contendo eventos S3 (PutObject)
type LocalS3EventQueue struct {
	watcherPath string
	storagePath string
	filePath    string
}

// NewLocalS3EventQueue cria um novo watcher para o diretório local
func NewLocalS3EventQueue(ctx context.Context, storagePath, filePath string) *LocalS3EventQueue {
	watcherPath := filepath.Join(storagePath, filePath)
	// Se storagePath não for absoluto, assume diretório base em os.TempDir()
	if !filepath.IsAbs(storagePath) {
		watcherPath = filepath.Join(os.TempDir(), storagePath, filePath)
	}
	if err := os.MkdirAll(watcherPath, 0755); err != nil {
		panic(fmt.Sprintf("Erro ao criar diretório watcherPath: %v", err))
	}
	return &LocalS3EventQueue{
		watcherPath: watcherPath,
		storagePath: storagePath,
		filePath:    filePath,
	}
}

// Start inicia o watcher e chama o callback para cada evento S3 PutObject simulado
func (l *LocalS3EventQueue) Start(onMessage func(event any)) (stop func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Erro ao criar watcher: %v", err)
	}
	if err := watcher.Add(l.watcherPath); err != nil {
		log.Fatalf("Erro ao monitorar diretório %s: %v", l.filePath, err)
	}
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					key := filepath.Join(l.filePath, filepath.Base(event.Name))
					filePath := event.Name
					fileInfo, err := os.Stat(filePath)
					var size int64
					var etag string
					if err == nil {
						size = fileInfo.Size()
						etag = calcMD5(filePath)
					}
					s3Event := events.S3Event{
						Records: []events.S3EventRecord{{
							EventVersion:      "2.1",
							EventSource:       "aws:s3",
							AWSRegion:         "us-east-1",
							EventTime:         time.Now(),
							EventName:         "ObjectCreated:Put",
							PrincipalID:       events.S3UserIdentity{PrincipalID: "LOCAL"},
							RequestParameters: events.S3RequestParameters{SourceIPAddress: "127.0.0.1"},
							ResponseElements:  map[string]string{"x-amz-request-id": "local", "x-amz-id-2": "local"},
							S3: events.S3Entity{
								SchemaVersion:   "1.0",
								ConfigurationID: "local-config",
								Bucket: events.S3Bucket{
									Name:          l.storagePath,
									OwnerIdentity: events.S3UserIdentity{PrincipalID: "LOCAL"},
									Arn:           "arn:aws:s3:::" + l.storagePath,
								},
								Object: events.S3Object{
									Key:       key,
									Size:      size,
									ETag:      etag,
									VersionID: "1",
									Sequencer: "local-seq",
								},
							},
						}},
					}
					onMessage(s3Event)
				}
			case err := <-watcher.Errors:
				log.Printf("Erro no watcher: %v", err)
			case <-quit:
				watcher.Close()
				return
			}
		}
	}()
	return func() { close(quit) }
}

// calcMD5 calcula o hash MD5 do arquivo para simular o ETag do S3
func calcMD5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

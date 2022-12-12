package files

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/go-co-op/gocron"
	"github.com/kurin/blazer/b2"
)

var maxLimitFile int = 10

type LocalToB2 struct {
	Files                  *db.FileModel
	BucketName             string
	AccountId              string
	AccountKey             string
	B2Prefix               string
	MappingUploadLocalLink string
	Bucket                 *b2.Bucket
	Context                context.Context
}

func NewB2Scheduler(
	accountId string,
	accountKey string,
	bucketName string,
	b2Prefix string,
	mappingUploadLocalLink string,
	files *db.FileModel,
) (*LocalToB2, error) {
	ctx := context.Background()
	b2, err := b2.NewClient(ctx, accountId, accountKey)

	if err != nil {
		return nil, err
	}

	bucket, err := b2.Bucket(ctx, bucketName)

	if err != nil {
		return nil, err
	}

	lb2 := &LocalToB2{
		files,
		bucketName,
		accountId,
		accountKey,
		b2Prefix,
		mappingUploadLocalLink,
		bucket,
		ctx,
	}

	return lb2, nil
}

func (l *LocalToB2) StartScheduler(limit int, sleep int, idle int) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(idle).Second().Do(l.UploadLocalToB2Task, limit, sleep)

	s.StartBlocking()
}

func (l *LocalToB2) UploadLocalToB2(localPath string) (*string, error) {
	suffixPath := strings.TrimPrefix(localPath, l.MappingUploadLocalLink)
	newLocalPath := filepath.Join(RootUploadPath, suffixPath)

	f, err := os.Open(newLocalPath)
	if err != nil {
		log.Println("UploadLocalToB2Task:", err)
		return nil, err
	}
	defer f.Close()

	obj := l.Bucket.Object(suffixPath)
	w := obj.NewWriter(l.Context)

	defer w.Close()

	if _, err := io.Copy(w, f); err != nil {
		log.Println("UploadLocalToB2Task:", err)
		return nil, err
	}

	url := obj.URL()

	return &url, nil
}

func (l *LocalToB2) UploadLocalToB2Task(limit int, sleep int) error {
	log.Println("UploadLocalToB2Task: RUNING")

	fileLimit := maxLimitFile
	if limit > 0 && limit < maxLimitFile {
		fileLimit = limit
	}

loop:
	files, err := l.Files.NotUpload(fileLimit)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Upload %s", file.LocalPath)
		cloudLink, err := l.UploadLocalToB2(file.LocalPath)
		if err == nil {
			l.Files.UploadCloudLink(file.LocalPath, *cloudLink)
		}

		time.Sleep(time.Duration(sleep))
	}

	if len(files) >= fileLimit {
		goto loop
	}

	return nil
}

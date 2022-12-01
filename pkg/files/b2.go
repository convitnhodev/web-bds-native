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

type LocalToB2 struct {
	Files                  *db.FileModel
	BucketName             string
	AccountId              string
	AccountKey             string
	UploadToB2At           string
	B2Prefix               string
	MappingUploadLocalLink string
	Bucket                 *b2.Bucket
	Context                context.Context
}

func NewB2Scheduler(
	accountId string,
	accountKey string,
	bucketName string,
	uploadToB2At string,
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
		uploadToB2At,
		b2Prefix,
		mappingUploadLocalLink,
		bucket,
		ctx,
	}

	return lb2, nil
}

func (l *LocalToB2) StartScheduler() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Day().At(l.UploadToB2At).Do(l.UploadLocalToB2Task)

	s.StartAsync()
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

func (l *LocalToB2) UploadLocalToB2Task() error {
	log.Println("UploadLocalToB2Task: RUNING")

	files, err := l.Files.NotUpload()

	if err != nil {
		log.Println("UploadLocalToB2Task:", err)
		return err
	}

	for _, file := range files {
		cloudLink, err := l.UploadLocalToB2(file.LocalPath)
		if err == nil {
			l.Files.UploadCloudLink(file.LocalPath, *cloudLink)
		}
	}

	return nil
}

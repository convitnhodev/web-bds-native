package files

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/pkg/errors"
)

type LocalFile struct {
	Root                    string
	PrefixUploadingRootLink string
	Files                   *db.FileModel
}

func (l *LocalFile) GenNamefile(filename string) string {
	h := sha1.New()
	ext := filepath.Ext(filename)
	filenameWithoutExt := strings.TrimSuffix(filename, ext)

	h.Write([]byte(filenameWithoutExt + "-" + time.Now().Format("20060102150405")))
	sha1_hash := hex.EncodeToString(h.Sum(nil)) + ext

	return sha1_hash
}

func (l *LocalFile) UploadFile(prefix_path string, file io.Reader, fileHeader *multipart.FileHeader) (*string, error) {
	header := fileHeader.Header
	size := fileHeader.Size
	Mb := int64(1024 * 1024)
	contentType := header.Get("Content-Type")
	isValidFile := false

	if strings.HasPrefix(contentType, "image/") {
		if size >= 5*Mb {
			return nil, errors.New("Image size greater than 5Mb")
		}
		isValidFile = true
	}

	if contentType == "application/pdf" {
		if size >= 20*Mb {
			return nil, errors.New("Document size greater than 20Mb")
		}
		isValidFile = true
	}

	if strings.HasPrefix(contentType, "video/") {
		if size >= 100*Mb {
			return nil, errors.New("Video size greater than 20Mb")
		}
		isValidFile = true
	}

	if !isValidFile {
		return nil, errors.New("File type no support")
	}

	root := filepath.Join(l.Root, prefix_path)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		err := os.MkdirAll(root, os.ModeDir)
		if err != nil {
			return nil, err
		}
	}

	dstFilename := filepath.Join(root, l.GenNamefile(fileHeader.Filename))
	dst, err := os.Create(dstFilename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	// Create new row for local file
	if _, err := l.Files.Create(filepath.Join(l.PrefixUploadingRootLink, dstFilename)); err != nil {
		return nil, err
	}

	return &dstFilename, nil
}

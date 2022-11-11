package file

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type LocalFile struct {
	Root string
}

func (l *LocalFile) UploadFile(file *os.File, fileHeader *multipart.FileHeader) (*string, error) {
	defer file.Close()

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

	dstFilename := filepath.Join(l.Root, fileHeader.Filename)
	dst, err := os.Create(dstFilename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	return &dstFilename, nil
}

package files

import (
	"crypto/sha1"
	"encoding/hex"
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

	h := sha1.New()
	ext := filepath.Ext(fileHeader.Filename)
	fileName := strings.TrimSuffix(fileHeader.Filename, ext)
	h.Write([]byte(fileName))
	sha1_hash := hex.EncodeToString(h.Sum(nil)) + ext

	dstFilename := filepath.Join(root, sha1_hash)
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

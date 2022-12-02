package files

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/deeincom/deeincom/pkg/models/db"
	"github.com/pkg/errors"
)

var RootUploadPath string
var FileExists = errors.New("File is exists")

type LocalFile struct {
	Root                   string
	MappingUploadLocalLink string
	Files                  *db.FileModel
}

func (l *LocalFile) GenNamefile(filename string, bytes []byte) string {
	h := sha1.New()
	ext := filepath.Ext(filename)

	h.Write(bytes)
	sha1_hash := hex.EncodeToString(h.Sum(nil)) + ext

	return sha1_hash
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func (l *LocalFile) UploadFile(prefixPath string, file io.Reader, fileHeader *multipart.FileHeader) (*string, error) {
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

	root := filepath.Join(RootUploadPath, prefixPath)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		err := os.MkdirAll(root, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	hashFilename := l.GenNamefile(fileHeader.Filename, bytes)
	dstFilename := filepath.Join(root, hashFilename)
	localLink := JoinURL(l.MappingUploadLocalLink, strings.TrimPrefix(dstFilename, RootUploadPath))

	if _, err := os.Stat(dstFilename); err == nil {
		return &localLink, FileExists
	}

	dst, err := os.Create(dstFilename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := dst.Write(bytes); err != nil {
		return nil, err
	}

	// Create new row for local file
	if _, err := l.Files.Create(localLink); err != nil {
		return nil, err
	}

	return &localLink, nil
}

func (l *LocalFile) RemoveLocalFile(filename string) error {
	suffixPath := strings.TrimPrefix(filename, l.MappingUploadLocalLink)
	newLocalPath := filepath.Join(RootUploadPath, suffixPath)

	if err := l.Files.Remove(filename); err != nil {
		return err
	}

	return os.Remove(newLocalPath)
}

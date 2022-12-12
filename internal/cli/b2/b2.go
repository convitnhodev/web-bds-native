package b2

import (
	"log"
	"os"
	"path/filepath"

	"github.com/deeincom/deeincom/internal/cli/root"
	"github.com/deeincom/deeincom/pkg/files"
	"github.com/pkg/errors"
)

var limit int
var sleep int
var idle int
var cfg string
var uploaddir string

func init() {
	cmd := root.New("b2")

	cmd.IntVar(&limit, "limit", 10, "Limit là số file nó sẽ upload")
	cmd.IntVar(&sleep, "sleep", 5, "Sleep sau khi upload 1 file.")
	cmd.IntVar(&idle, "idle", 30, "Thời gian idle sau khi hết task.")
	cmd.StringVar(&uploaddir, "uploaddir", "./upload", "Thư mục upload files")

	fullUploaddir, err := filepath.Abs(uploaddir)
	if err != nil {
		panic(err)
	}
	files.RootUploadPath = fullUploaddir

	cmd.Action(func() error {
		return run(cmd)
	})
}

func run(c *root.Cmd) error {
	if c.App.Config.B2AccountId != "" {
		B2Scheduler, err := files.NewB2Scheduler(
			c.App.Config.B2AccountId,
			c.App.Config.B2AccountKey,
			c.App.Config.B2BucketName,
			c.App.Config.B2Prefix,
			c.App.Config.MappingUploadLocalLink,
			c.App.Files,
		)

		if err != nil {
			panic(err)
		}

		log.Println("Run b2 job")
		B2Scheduler.StartScheduler(
			limit,
			sleep,
			idle,
		)
	}

	return nil
}

func isMissingConfig(err error) bool {
	var pathError *os.PathError
	return errors.As(err, &pathError)
}

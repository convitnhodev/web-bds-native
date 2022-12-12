package version

import (
	"fmt"

	"github.com/deeincom/deeincom/internal/cli/root"
)

func init() {
	cmd := root.New("version")
	cmd.Action(func() error {
		return run(cmd)
	})
}

func run(c *root.Cmd) error {
	fmt.Println(`1.0.0`)
	return nil
}

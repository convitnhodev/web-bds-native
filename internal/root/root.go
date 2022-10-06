package root

import (
	"context"
	"flag"

	app "github.com/deeincom/deeincom"
)

// CmdSet is a slice of Cmd pointer
var CmdSet = []*Cmd{}

// Cmd is `flag.FlagSet` with exec fn
type Cmd struct {
	Context context.Context
	*flag.FlagSet
	App *app.App
	fn  func() error
}

// Action register fn
func (c *Cmd) Action(fn func() error) {
	c.fn = fn
}

// Exec execute fn
func (c *Cmd) Exec() error {
	if c.fn != nil {
		return c.fn()
	}
	return nil
}

// New return a new cmd
func New(name string) *Cmd {
	c := &Cmd{
		Context: context.Background(),
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
		fn:      func() error { return nil },
	}
	CmdSet = append(CmdSet, c)
	return c
}

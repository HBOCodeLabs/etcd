package command

import (
	"errors"

	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
)

// NewRemoveCommand returns the CLI command for "rm".
func NewRemoveCommand() cli.Command {
	return cli.Command{
		Name:  "rm",
		Usage: "remove a key",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "dir", Usage: "removes the key if it is an empty directory or a key-value pair"},
			cli.BoolFlag{Name: "recursive", Usage: "removes the key and all child keys(if it is a directory)"},
			cli.StringFlag{Name: "with-value", Value: "", Usage: "previous value"},
			cli.IntFlag{Name: "with-index", Value: 0, Usage: "previous index"},
		},
		Action: func(c *cli.Context) {
			handleAll(c, removeCommandFunc)
		},
	}
}

// removeCommandFunc executes the "rm" command.
func removeCommandFunc(c *cli.Context, client *etcd.Client) (*etcd.Response, error) {
	if len(c.Args()) == 0 {
		return nil, errors.New("Key required")
	}
	key := c.Args()[0]
	recursive := c.Bool("recursive")
	dir := c.Bool("dir")

	// TODO: distinguish with flag is not set and empty flag
	// the cli pkg need to provide this feature
	prevValue := c.String("with-value")
	prevIndex := uint64(c.Int("with-index"))

	if prevValue != "" || prevIndex != 0 {
		return client.CompareAndDelete(key, prevValue, prevIndex)
	}

	if recursive || !dir {
		return client.Delete(key, recursive)
	}

	return client.DeleteDir(key)
}

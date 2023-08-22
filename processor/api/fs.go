package api

import (
	"embed"
	"io/fs"
	"path"
)

type ProcessorFs struct {
	Content embed.FS
}

func (c *ProcessorFs) Open(name string) (fs.File, error) {
	return c.Content.Open(path.Join("webclient", "build", "static", name))
}

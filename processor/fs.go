package processor

import (
	"embed"
	"io/fs"
	"path"
)

type processorFs struct {
	content embed.FS
}

func (c *processorFs) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("webclient", "build", "static", name))
}

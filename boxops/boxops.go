package boxops

import (
	"github.com/gobuffalo/packr"
)

type BoxOps struct {
	Box packr.Box
}

// New creates new FileOps box
func New(path string) (f *BoxOps) {
	box := packr.NewBox(path)

	return &BoxOps{Box: box}
}

func (b *BoxOps) Read(filePath string) (str string, err error) {
	str, err = b.Box.FindString(filePath)
	if err != nil {
		return "", err
	}

	return str, err
}

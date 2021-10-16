package box

import "embed"

//go:embed configs/* lang/*
var FileSystem embed.FS

type Box struct{}

func (b *Box) Read(filePath string) (str string, err error) {
	var bytes []byte
	bytes, err = FileSystem.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

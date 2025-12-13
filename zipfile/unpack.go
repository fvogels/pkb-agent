package zipfile

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"os"
	"pkb-agent/util/pathlib"
)

type ZippedFile struct {
	Filename pathlib.Path
	Contents []byte
}

func Unpack(buffer []byte) ([]ZippedFile, error) {
	size := int64(len(buffer))
	bufferReader := bytes.NewReader(buffer)

	zipFileReader, err := zip.NewReader(bufferReader, size)
	if err != nil {
		return nil, err
	}

	files := []ZippedFile{}

	for _, entry := range zipFileReader.File {
		reader, err := entry.Open()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		var buffer bytes.Buffer
		if _, err := io.Copy(bufio.NewWriter(&buffer), reader); err != nil {
			return nil, err
		}

		file := ZippedFile{
			Filename: pathlib.New(entry.Name),
			Contents: buffer.Bytes(),
		}

		files = append(files, file)
	}

	return files, nil
}

func (zippedFile *ZippedFile) SaveToDirectory(directory pathlib.Path) (pathlib.Path, error) {
	destinationPath := directory.Join(zippedFile.Filename)
	reader := bytes.NewReader(zippedFile.Contents)
	writer, err := os.Create(destinationPath.String())
	if err != nil {
		return destinationPath, err
	}

	if _, err := io.Copy(writer, reader); err != nil {
		return destinationPath, err
	}

	return destinationPath, nil
}

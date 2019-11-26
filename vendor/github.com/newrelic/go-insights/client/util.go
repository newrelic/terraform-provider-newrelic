package client

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
)

func gZipBuffer(body []byte) (io.Reader, error) {
	var err error

	readBuffer := bufio.NewReader(bytes.NewReader(body))
	buffer := bytes.NewBuffer([]byte{})
	writer := gzip.NewWriter(buffer)

	_, err = readBuffer.WriteTo(writer)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

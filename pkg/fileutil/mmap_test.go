package fileutil

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	filename := path.Join(t.TempDir(), "testdata")
	file, err := os.Create(filename)
	assert.NoError(t, err)
	content := "abc123"
	_, err = file.WriteString(content)
	assert.NoError(t, err)
	assert.NoError(t, file.Close())

	file, err = os.Open(filename)
	assert.NoError(t, err)
	bys, err := Map(file)
	assert.NoError(t, err)
	assert.Equal(t, []byte(content), bys)
	assert.NoError(t, Unmap(file, bys))
}

func TestRWMap(t *testing.T) {
	var content = []byte("12345")
	var size = 1024
	filename := path.Join(t.TempDir(), "testdata")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	mapBytes, err := RWMap(f, size)
	assert.NoError(t, err)
	if Unmap(f, nil) != nil {
		t.Error("unmap nil returns not nil")
	}

	buffer := bytes.NewBuffer(mapBytes[:0])

	_, err = buffer.Write(content)
	assert.NoError(t, err)

	err = Sync(mapBytes)
	assert.NoError(t, err)

	if Unmap(f, mapBytes) != nil {
		t.Errorf("unmap mapBytes with error: %v", err)
	}

	fileContent, err := os.ReadFile(filename)
	assert.NoError(t, err)

	assert.Len(t, fileContent, size)

	assert.Equal(t, content, fileContent[0:len(content)])
}

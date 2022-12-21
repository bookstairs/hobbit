package toml

import (
	"bufio"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// EncodeToml encodes data into file using the toml format,
// encode data to tmp file, if success then rename tmp => target file
func EncodeToml(fileName string, v interface{}) error {
	tmp := fmt.Sprintf("%s.tmp", fileName)
	f, _ := os.Create(tmp)
	w := bufio.NewWriter(f)
	// write info using toml format
	if err := toml.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, fileName); err != nil {
		return fmt.Errorf("rename tmp file[%s] name error:%s", tmp, err)
	}
	return nil
}

func WriteConfig(fileName, content string) error {
	tmp := fmt.Sprintf("%s.tmp", fileName)
	f, _ := os.Create(tmp)
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(content); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	// need close file before rename file name
	if err := f.Close(); err != nil {
		return err
	}
	// if file isn't closed, will throw the file be processed the other process in windows.
	if err := os.Rename(tmp, fileName); err != nil {
		return fmt.Errorf("rename tmp file[%s] name error:%s", tmp, err)
	}
	return nil
}

// DecodeToml decodes data from file using toml format
func DecodeToml(fileName string, v interface{}) error {
	if _, err := toml.DecodeFile(fileName, v); err != nil {
		return err
	}
	return nil
}

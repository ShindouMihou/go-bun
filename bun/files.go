package bun

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type BFile struct {
	path  string
	file  *os.File
	error error
}

func File(path string) *BFile {
	return &BFile{
		path:  path,
		file:  nil,
		error: nil,
	}
}

func (file *BFile) openRead() {
	f, err := os.Open(file.path)
	file.file = f
	file.error = err
}

func (file *BFile) openWrite(trunc bool) {
	var f *os.File
	var err error

	if trunc {
		f, err = os.Create(file.path)
	} else {
		f, err = os.OpenFile(file.path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	}

	file.file = f
	file.error = err
}

func (file *BFile) clear() {
	if err := file.file.Truncate(0); err != nil {
		panic(err)
	}
	if _, err := file.file.Seek(0, 0); err != nil {
		panic(err)
	}
}

// MkdirParent creates the parent folders of the path.
func (file *BFile) mkparent() {
	if strings.Contains(file.path, "\\") || strings.Contains(file.path, "/") {
		if err := os.MkdirAll(filepath.Dir(file.path), os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func (file *BFile) read(fn func() any) any {
	file.openRead()
	defer file.close()

	if file.error != nil {
		panic(file.error)
	}

	return fn()
}

func (file *BFile) write(trunc bool, fn func() any) any {
	file.mkparent()
	file.openWrite(trunc)
	defer file.close()

	if file.error != nil {
		panic(file.error)
	}

	if trunc {
		file.clear()
	}

	return fn()
}

func (file *BFile) close() {
	// ignore the error, it's likely that it just already called
	_ = file.file.Close()
}

func (file *BFile) Text() string {
	return string(file.Bytes())
}

func (file *BFile) Bytes() []byte {
	return file.read(func() any {
		bytes, err := io.ReadAll(file.file)
		if err != nil {
			panic(err)
		}
		return bytes
	}).([]byte)
}

func (file *BFile) Json(t interface{}) {
	file.read(func() any {
		bytes, err := io.ReadAll(file.file)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(bytes, &t); err != nil {
			panic(err)
		}
		return nil
	})
}

func (file *BFile) wrt(trunc bool, bytes []byte) {
	file.write(trunc, func() any {
		if _, err := file.file.Write(bytes); err != nil {
			panic(err)
		}
		return nil
	})
}

func (file *BFile) wrtjson(trunc bool, t interface{}) {
	bytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	file.wrt(trunc, bytes)
}

func (file *BFile) wrtany(trunc bool, t any) {
	switch t.(type) {
	case string:
		file.wrt(trunc, []byte(t.(string)))
	case []byte:
		file.wrt(trunc, t.([]byte))
	default:
		file.wrtjson(true, t)
	}
}

func (file *BFile) Write(t any) {
	file.wrtany(false, t)
}

func (file *BFile) Overwrite(t any) {
	file.wrtany(true, t)
}

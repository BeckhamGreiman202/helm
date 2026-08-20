package getter

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

// FileGetter is the default getter for file:// URLs
type FileGetter struct {
	opts options
}

// NewFileGetter constructs a new FileGetter
func NewFileGetter(options ...Option) (Getter, error) {
	return &FileGetter{
		opts: newOptions(options),
	}, nil
}

// Get file
func (g *FileGetter) Get(href string, options ...Option) (*bytes.Buffer, error) {
	g.opts.apply(options)
	path, err := filepath.Abs(href)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, 4096)
	var result bytes.Buffer
	for {
		n, err := f.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return &result, nil
}

package fs

import (
	"io"
	"path"
	"reactor-crw"
)

// PathResolver defines a simple interface to resolve path for content
// before saving it.
type pathResolver interface {
	// CreateFolder creates a new folder in the "current" dir. The current
	// directory is specified in Config.SavePath provided by a user. If dir
	// with such name already exists it'll skip creation step.
	CreateFolder(name string) error

	// CreateFile creates a new file for downloaded content in the directory
	// created by the CreateFolder function and returns the corresponding
	// io.WriteCloser. If a file with such a name already exists then it'll
	// be overwritten.
	CreateFile(name string) (io.WriteCloser, error)

	// Remove removes a file or folder by its name. If provided path does not
	// exist it'll be skipped.
	Remove(name string)
}

// FileSaver defines ContentHandler implementation that will download content
// and save it to the host's file system. It resolves the corresponding file
// path with PathResolver.
type FileSaver struct {
	pr pathResolver
	t  reactor_crw.Transport
}

// NewFileSaver creates a new FileSaver instance along with a new folder that
// will act as a context for a new instance. All files will be processed within
// this new folder.
func NewFileSaver(pr pathResolver, t reactor_crw.Transport, baseFolder string) (*FileSaver, error) {
	err := pr.CreateFolder(baseFolder)
	if err != nil {
		return nil, err
	}

	return &FileSaver{
		pr: pr,
		t:  t,
	}, nil
}

// Process downloads content by the corresponding URL by making an HTTP request
// and saves the result to the file system. In case of error during saving the
// content the corresponding file will be deleted from the file system.
func (f *FileSaver) Process(url string, progress chan<- int, e chan<- error) {
	defer func() {
		progress <- 1
	}()

	data, err := f.t.FetchData(url)
	if err != nil {
		e <- err
		return
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(data)

	file, err := f.pr.CreateFile(path.Base(url))
	if err != nil {
		e <- err
		return
	}

	defer func(f io.WriteCloser) {
		_ = f.Close()
	}(file)

	_, err = io.Copy(file, data)
	if err != nil {
		f.pr.Remove(path.Base(url))
		e <- err
	}
}

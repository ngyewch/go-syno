package fs

import (
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/filestation"
	"io"
	"io/fs"
	"path/filepath"
	"time"
)

type FS struct {
	dir            string
	client         *api.Client
	fileStationApi *filestation.Api
}

func NewFS(client *api.Client, dir string) (*FS, error) {
	if dir == "" {
		return nil, fs.ErrInvalid
	}
	fileStationApi := filestation.New(client)
	getInfoResponse, err := fileStationApi.GetInfo(filestation.GetInfoRequest{
		Path:       []string{dir},
		Additional: []string{"size", "time"},
	})
	if err != nil {
		return nil, err
	}

	if len(getInfoResponse.Data.Files) == 0 {
		return nil, fs.ErrNotExist
	}
	entry := getInfoResponse.Data.Files[0]
	if entry.Additional == nil {
		return nil, fs.ErrNotExist
	}
	if !entry.IsDir {
		return nil, fs.ErrInvalid
	}

	return &FS{
		dir:            dir,
		client:         client,
		fileStationApi: fileStationApi,
	}, nil
}

func (f *FS) resolvePath(name string) string {
	resolvedPath := f.dir
	if name != "" {
		resolvedPath += "/" + name
	}
	return resolvedPath
}

func (f *FS) Open(name string) (fs.File, error) {
	r, err := f.fileStationApi.Download(filestation.DownloadRequest{
		Path: []string{f.resolvePath(name)},
		Mode: "download",
	})
	if err != nil {
		return nil, err
	}
	return &file{
		fs:   f,
		path: name,
		r:    r,
	}, nil
}

func (f *FS) Stat(name string) (fs.FileInfo, error) {
	getInfoResponse, err := f.fileStationApi.GetInfo(filestation.GetInfoRequest{
		Path:       []string{f.resolvePath(name)},
		Additional: []string{"size", "time"},
	})
	if err != nil {
		return nil, err
	}
	if getInfoResponse.Data == nil || len(getInfoResponse.Data.Files) == 0 || getInfoResponse.Data.Files[0].Additional == nil {
		return nil, fs.ErrNotExist
	}
	return &dirEntry{
		file: getInfoResponse.Data.Files[0],
	}, nil
}

func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	listResponse, err := f.fileStationApi.List(filestation.ListRequest{
		FolderPath: f.resolvePath(name),
		Additional: []string{"size", "time"},
	})
	if err != nil {
		return nil, err
	}
	var dirEntries []fs.DirEntry
	for _, file := range listResponse.Data.Files {
		dirEntries = append(dirEntries, &dirEntry{
			file: file,
		})
	}
	return dirEntries, nil
}

func (f *FS) Sub(dir string) (fs.FS, error) {
	return NewFS(f.client, f.resolvePath(dir))
}

// ----

type file struct {
	fs   *FS
	path string
	r    io.ReadCloser
}

func (f *file) Close() error {
	return f.r.Close()
}

func (f *file) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f *file) Stat() (fs.FileInfo, error) {
	return f.fs.Stat(f.path)
}

// ----

type dirEntry struct {
	file filestation.File
}

func (d *dirEntry) Name() string {
	return filepath.Base(d.file.Path)
}

func (d *dirEntry) IsDir() bool {
	return d.file.IsDir
}

func (d *dirEntry) Type() fs.FileMode {
	if d.file.IsDir {
		return fs.ModeDir
	} else {
		return 0
	}
}

func (d *dirEntry) Info() (fs.FileInfo, error) {
	return d, nil
}

func (d *dirEntry) Size() int64 {
	if d.file.Additional == nil {
		return -1
	}
	return d.file.Additional.Size
}

func (d *dirEntry) Mode() fs.FileMode {
	return d.Type()
}

func (d *dirEntry) ModTime() time.Time {
	if d.file.Additional == nil || d.file.Additional.Time == nil {
		return time.Time{}
	}
	return time.Unix(d.file.Additional.Time.Mtime, 0)
}

func (d *dirEntry) Sys() any {
	return d.file
}

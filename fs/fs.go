package fs

import (
	"fmt"
	"github.com/ngyewch/go-syno/api"
	"github.com/ngyewch/go-syno/api/filestation"
	"io"
	"io/fs"
	"path/filepath"
	"time"
)

type FS struct {
	fileStationApi *filestation.Api
}

func NewFS(client *api.Client) *FS {
	return &FS{
		fileStationApi: filestation.New(client),
	}
}

func (f *FS) Open(name string) (fs.File, error) {
	r, err := f.fileStationApi.Download(filestation.DownloadRequest{
		Path: []string{name},
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
		Path:       []string{name},
		Additional: []string{"size", "time"},
	})
	if err != nil {
		return nil, err
	}
	if !getInfoResponse.Success {
		return nil, fmt.Errorf("synology error code: %d", getInfoResponse.Error.Code)
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
		FolderPath: name,
		Additional: []string{"size", "time"},
	})
	if err != nil {
		return nil, err
	}
	if !listResponse.Success {
		return nil, fmt.Errorf("synology error code: %d", listResponse.Error.Code)
	}
	var dirEntries []fs.DirEntry
	for _, file := range listResponse.Data.Files {
		dirEntries = append(dirEntries, &dirEntry{
			file: file,
		})
	}
	return dirEntries, nil
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

package fileHandler

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/interp"
)

type fakeFileInfo struct {
	path string
}

func (f fakeFileInfo) Name() string {
	return f.path
}

func (f fakeFileInfo) Size() int64 {
	return 0
}

func (f fakeFileInfo) Mode() os.FileMode {
	return 0
}

func (f fakeFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fakeFileInfo) IsDir() bool {
	return false
}

func (f fakeFileInfo) Sys() any {
	return nil
}

func GetStatHandler(files []scenarios.File) interp.StatHandlerFunc {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	return func(ctx context.Context, path string, followSymlinks bool) (os.FileInfo, error) {
		for _, file := range files {
			if filepath.Join(cwd, file.Path) == path {
				if !file.Exists {
					return nil, os.ErrNotExist
				}
				return fakeFileInfo{
					path: file.Path,
				}, nil
			}
		}
		return interp.DefaultStatHandler()(ctx, path, followSymlinks)
	}
}

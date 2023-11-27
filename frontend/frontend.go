package frontend

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
)

//go:generate yarn install
//go:generate yarn build

//go:embed all:build/*
var content embed.FS

func FileServerWith404(root http.FileSystem) http.Handler {
	filesrv := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := root.Open(r.URL.Path)
		if err != nil && errors.Is(err, fs.ErrNotExist) {
			r.URL.Path = "/200.html"
		}
		if err == nil {
			f.Close()
		}
		filesrv.ServeHTTP(w, r)
	})
}

func Handler() (http.Handler, error) {
	subFs, err := fs.Sub(content, "build")
	if err != nil {
		return nil, err
	}
	return FileServerWith404(http.FS(subFs)), nil
}

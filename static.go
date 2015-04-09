package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _esc_localFS struct{}

var _esc_local _esc_localFS

type _esc_staticFS struct{}

var _esc_static _esc_staticFS

type _esc_file struct {
	compressed string
	size       int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_esc_localFS) Open(name string) (http.File, error) {
	f, present := _esc_data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_esc_staticFS) Open(name string) (http.File, error) {
	f, present := _esc_data[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		gr, err = gzip.NewReader(bytes.NewBufferString(f.compressed))
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (f *_esc_file) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_esc_file
	}
	return &httpFile{
		Reader:    bytes.NewReader(f.data),
		_esc_file: f,
	}, nil
}

func (f *_esc_file) Close() error {
	return nil
}

func (f *_esc_file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_esc_file) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_esc_file) Name() string {
	return f.name
}

func (f *_esc_file) Size() int64 {
	return f.size
}

func (f *_esc_file) Mode() os.FileMode {
	return 0
}

func (f *_esc_file) ModTime() time.Time {
	return time.Time{}
}

func (f *_esc_file) IsDir() bool {
	return f.isDir
}

func (f *_esc_file) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _esc_local
	}
	return _esc_static
}

var _esc_data = map[string]*_esc_file{

	"/static/dist/style.css": {
		local: "static/dist/style.css",
		size:  11,
		compressed: "" +
			"\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xd2\xd7R(\xc9O\xc9W\xd0\xd2\xe7\x02\x04\x00\x00\xff\xff\x8ah\xaf\x1e\v\x00\x00\x00",
	},

	"/static/index.html": {
		local: "static/index.html",
		size:  1176,
		compressed: "" +
			"\x1f\x8b\b\x00\x00\tn\x88\x00\xff\x9c\x94\xcfo\xd30\x14\xc7\xef\xf9+,\xc3!Մ\xcd\x10'\x9a\xe4\x02B\x1c\xb8u\x12B\xd3\x0e\x8e\xed&\xee\\;\xe49\xcd*\xb4\xff\x9d\xe7\xa4m\x9a\xb2J0\x9f\xe2\xf7\xe3\xfb\xc9\xfb\x91du\xd8\xda\"\xc9j-T\x91\x10<Y0\xc1\xea\xe2\xab\xed\xb4\v\x8a\xac\x82\b\x1dd|\xb4\x8e\x11\xd6" +
			"\xb8GR\xb7z\x9dS\x0e\xe87\x92+\x03\x01\x9f\xf7V3\t@I\xd87:\xa7A?\x05>\xdc92\xf8\b\xc9J\xaf\xf6\a\xa5\xfa\xf6\x00\xfar\x02\xa1)\x19\x9d\xca\xec\x88Q9\xb5\xbe\xa2E\x86\x88\xdd\xd1\x03\xb25M8\x87l\xc4N\x8cVJ\xa0\x959\xadCh>q.6\xe2\x89U\xdeWV\x8b\xc6\x00\x93~;ظ5" +
			"%\xf0ͯN\xb7{~\xcb>\xb2\x0f\x87\v\xdb\x1a\xc76\x10\x81\xa3^\xf1\x0f\xc81&\x9e\xb7\xe9\xbas2\x18\xef\xd2\x05\xf9\x9d\x9c\xec;\xd1\x12\xe9\x9d[\xce,X\x19\xc91\x87\xbe\x895.\x96S\xfcQ\x85\x88\xa6\xd1N}\xf7U\xba\x85*J\x92\xb3\x135\x14*`\xf6\xfd\xfb\x87\xbf]~%[o-F(\x06\xc3\xe3\x9doH" +
			">]\xbfiSՁ\xbcC\x83\xb4\x06\xe70\x1a\x963)\x04\xb3\xf15\xee|\x8a\xa8\xc5\xcck\xd6$=\x82.\xdf/\x9e\x19\xf9\xff\xc0\xcf\xc9\xf4\x94\x9c\xf3z\xe3\x94\xef\xef\xe9\x0f]\xae\xbc|ԁ>\\\xa2c\xb3\x91\xe7tONQ)\xed\x01w\x82\xdeX/E\xec.\xab=\x84\x1b\xca{\x88ͿLg\xdeI\xebA\xa3\xcci" +
			"\xa8z\x17^*r\x9a\x12\x0e3.n\x91\x95\xc5g\xd4\xd0\xe3\x14\a\x1d\xc52^\x1e\x16\x99.\x16WJ=\xa3o5\x80\xa8^\xc5\xe7H`qMc\x06S\"\x88\xab\xc0g\xa2-\x169\x17}\xb1\xa0\x9f\xbekI\xd9\xfa\x1et\\.\r\xc4\xf9@\xa0k\x1a߆\xa9\xcdp\xa5\xce3\xe4\xa1\xdb\xd37\x86\x19\xc3_\x01\xbf\xff\xe1\x87" +
			"\xf4'\x00\x00\xff\xff\xe3uץ\x98\x04\x00\x00",
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/static": {
		isDir: true,
		local: "/static",
	},

	"/static/dist": {
		isDir: true,
		local: "/static/dist",
	},
}

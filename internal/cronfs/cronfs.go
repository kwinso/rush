package cronfs

import (
	"github.com/winfsp/cgofuse/fuse"
)

type cronFS struct {
	fuse.FileSystemBase
	watcher *cronWatcher
}

func newCronFS(watcher *cronWatcher) *cronFS {
	return &cronFS{
		watcher: watcher,
	}
}

func (self *cronFS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64,
) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)

	for filename := range self.watcher.GetEntries() {
		fill(filename, nil, 0)
	}

	return 0
}

func (self *cronFS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	default:
		file, ok := self.watcher.GetEntries()[path[1:]] // remove leading slash
		if !ok {
			return -1 * fuse.ENOENT
		}

		// TODO: Maybe do a writing support
		stat.Mode = fuse.S_IFREG | 0444
		stat.Size = int64(len(file))
		return 0
	}
}

func (self *cronFS) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	endofst := ofst + int64(len(buff))
	contents, ok := self.watcher.GetEntries()[path[1:]] // remove leading slash

	if !ok {
		return 0
	}

	if endofst > int64(len(contents)) {
		endofst = int64(len(contents))
	}
	if endofst < ofst {
		return 0
	}
	n = copy(buff, contents[ofst:endofst-1])
	return
}

func (self *cronFS) Open(path string, flags int) (errc int, fh uint64) {
	idx := 0
	path = path[1:]
	for filename := range self.watcher.GetEntries() {
		if filename == path {
			return 0, uint64(idx)
		}
		idx++
	}

	return 0, 0
}

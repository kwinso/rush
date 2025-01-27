package cronfs

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/winfsp/cgofuse/fuse"
)

var activeManager *Manager = nil

type Manager struct {
	cronFS     *cronFS
	host       *fuse.FileSystemHost
	mountPoint string
	watcher    *cronWatcher
}

func GetCronFSManager() (*Manager, error) {
	if activeManager != nil {
		return activeManager, nil
	}

	watcher := NewCronWatcher()
	activeManager = &Manager{
		cronFS:  newCronFS(watcher),
		watcher: watcher,
	}

	return activeManager, nil
}

func (self *Manager) Mount(mountPoint string) error {
	if self.IsMounted() {
		log.Debug("Found active cronfs mount, will not mount again")
		return nil
	}

	self.host = fuse.NewFileSystemHost(self.cronFS)
	mountPoint, err := filepath.Abs(mountPoint)
	if err != nil {
		return err
	}
	self.mountPoint = mountPoint

	if _, err := os.Stat(mountPoint); os.IsNotExist(err) {
		log.Debug("Mount point does not exist, creating...")
		// Create the directory
		err = os.MkdirAll(mountPoint, os.ModePerm)
		if err != nil {
			return err
		}
	}

	self.watcher.Watch()
	go self.host.Mount(mountPoint, []string{})
	return nil
}

func (self *Manager) IsMounted() bool {
	return self.host != nil
}

func (self *Manager) Unmount() error {
	if self.IsMounted() {
		log.Debug("Found active cronfs mount, unmounting...")

		self.host.Unmount()

		if _, err := os.Stat(self.mountPoint); err == nil {
			log.Debug("Removing cronfs mount point directory")
			// Create the directory
			err = os.RemoveAll(self.mountPoint)
			if err != nil {
				return err
			}
		}

		self.watcher.StopWatching()
		self.host = nil
	}
	return nil
}

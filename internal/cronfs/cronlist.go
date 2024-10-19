package cronfs

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type cronWatcher struct {
	entries  map[string]string
	quitChan chan struct{}
}

func NewCronWatcher() *cronWatcher {
	return &cronWatcher{
		entries: make(map[string]string),
	}
}

func (cw *cronWatcher) Watch() {
	log.Debug("Starting cronwatcher")
	ticker := time.NewTicker(2 * time.Second)
	cw.quitChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				cw.updateCronList()
			case <-cw.quitChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (cw *cronWatcher) StopWatching() {
	log.Debug("Stopping cronwatcher")
	close(cw.quitChan)
}

func (cw *cronWatcher) GetEntries() map[string]string {
	return cw.entries
}

func (cw *cronWatcher) updateCronList() {
	list, err := getCronList()
	if err != nil {
		log.Fatal("Unable to get cronlist")
	}
	cw.entries = parseCron(list)
}

func parseCron(entries []string) map[string]string {
	cronMap := make(map[string]string)
	for _, entry := range entries {
		if strings.HasPrefix(entry, "#") {
			continue
		}
		fields := strings.Fields(entry)
		time := strings.Join(fields[:5], " ")
		cmd := strings.Join(fields[5:], " ")
		cronMap[cmd] = strings.TrimSpace(time)
	}
	return cronMap
}

func getCronList() ([]string, error) {
	// exec crontab -l
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	output = bytes.TrimSpace(output)
	return strings.Split(string(output), "\n"), nil
}

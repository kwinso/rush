package history

import (
	"fmt"
	"os"
	"path"

	"github.com/fiorix/go-readline"
)

var histPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to get user homedir", err)
		os.Exit(1)
	}
	histPath = path.Join(home, ".rush_hists")
}

func InitHistory() {
	readline.ParseAndBind("TAB: menu-complete")

	if _, err := os.Stat(histPath); os.IsNotExist(err) {
		_, err := os.Create(histPath)
		if err != nil {
			fmt.Println("Unable to create .rush_hist file in user home: ", err)
			os.Exit(1)
		}
	}
	err := readline.ReadHistoryFile(histPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SaveCmdToHistory(cmd string) {
	readline.AddHistory(cmd)
}

func SaveHistoryFile() error {
	return readline.WriteHistoryFile(histPath)
}

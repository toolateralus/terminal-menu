package state

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type State struct {
	Items    []string
	Selected int
	Running  bool
	IsDir    bool
	Path     string
	ShowHelp bool
	Previous *State
}

func (s State) IsRunning() bool {
	return s.Running
}

func NewState(previous *State, path string) State {
	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if stats.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		items := []string{}
		
		for _, file := range files {
			// try hide the hidden files.
			if !strings.HasPrefix(file.Name(), ".") {
				fileInfo, err := os.Stat(filepath.Join(path, file.Name()))
				if err != nil {
					log.Fatal(err)
				}
				// Check if the file has read permission
				if fileInfo.Mode().Perm()&0400 != 0 {
					items = append(items, file.Name())
				}
			}
		}
		
		return State{
			Items:    items,
			Selected: 0,
			Running:  true,
			Path:     path,
			Previous: previous,
			IsDir:    true,
			ShowHelp: true,
		}
	} else {
		contents, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		
		lines := strings.Split(string(contents), "\n")
		
		return State{
			Items:    lines,
			Selected: 0,
			Running:  true,
			Path:     path,
			Previous: previous,
			IsDir:    false,
			ShowHelp: false,
		}
	
	}

}
package gie

import (
	"io/ioutil"
	"strings"
)

var supportedProjections = []string{
	"etmerc", "utm",
	"aea", "leac",
	"merc",
}

// Gie manages the GIE reading and executing processes
type Gie struct {
	dir      string
	files    []string
	Commands []*Command
}

// NewGie returns a new Gie object
func NewGie(dir string) (*Gie, error) {

	g := &Gie{
		dir:      dir,
		files:    []string{},
		Commands: []*Command{},
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		file := info.Name()
		if strings.HasSuffix(file, ".gie") {
			g.files = append(g.files, dir+"/"+file)
		}
	}
	return g, nil
}

// Parse reads the .gie file and creates the commands
func (g *Gie) Parse() error {
	for _, file := range g.files {
		p, err := NewParser(file)
		if err != nil {
			return err
		}
		g.Commands = append(g.Commands, p.Commands...)
	}

	return nil
}

// IsSupported returns true iff we've ported the projection in question
func (g *Gie) IsSupported(cmd *Command) bool {

	proj := cmd.ProjectionName()

	for _, sp := range supportedProjections {
		if proj == sp {
			return true
		}
	}

	return false
}

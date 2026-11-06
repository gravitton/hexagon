package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const directory = "../images"

type figure struct {
	name string
	draw func(theme) []*panel
}

func main() {
	for _, t := range themes() {
		for _, f := range figures {
			if err := f.write(t); err != nil {
				fmt.Fprintf(os.Stderr, "diagrams: %s-%s: %v\n", f.name, t.name, err)
				os.Exit(1)
			}
		}
	}
}

func (f figure) write(t theme) error {
	svg, err := render(t, f.draw(t))
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(directory, f.name+"-"+t.name+".svg"), []byte(svg), 0o644)
}

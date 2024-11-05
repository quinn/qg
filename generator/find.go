package generator

import (
	"fmt"
)

func Find(generators []Generator, cmd string) (*Generator, error) {
	for _, gen := range generators {
		if gen.Cmd == cmd {
			return &gen, nil
		}
	}
	return nil, fmt.Errorf("generator not found: %s", cmd)
}

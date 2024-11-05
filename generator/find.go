package generator

import (
	"fmt"
)

func Find(generators []Generator, name string) (*Generator, error) {
	for _, gen := range generators {
		if gen.Cfg.Name == name {
			return &gen, nil
		}
	}
	return nil, fmt.Errorf("generator not found: %s", name)
}

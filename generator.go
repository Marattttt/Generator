package generator

import (
	"sync"

	"github.com/marattttt/generator/command"
	"github.com/marattttt/generator/drawing"
)

// Applying commmands is multithreaded
// All commands should have the same drawing parameter
type Generator struct {
	Target   *drawing.Drawing
	Commands []command.Command
}

func (g Generator) ApplyCommands() (cycles int, err error) {
	toExecute, left := command.FilterRelatedCommands(g.Commands)

	cycles = 0
	for len(toExecute) > 0 {
		var wg sync.WaitGroup
		for _, comm := range toExecute {
			wg.Add(1)
			go func(comm command.Command) {
				defer wg.Done()
				comm.Execute(g.Target)
			}(comm)
		}

		wg.Wait()
		toExecute, left = command.FilterRelatedCommands(left)
		cycles++
	}

	return cycles, nil
}

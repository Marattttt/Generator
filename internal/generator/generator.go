package generator

import (
	"sync"

	"github.com/marattttt/paperwork/generator/drawing"
)

// Applying commmands is multithreaded
// All commands should have the same drawing parameter
type Generator struct {
	Target   *drawing.Drawing
	Commands []Command
}

func (g Generator) ApplyCommands() (cycles int, err error) {
	toExecute, left := filterRelatedCommands(g.Commands)

	cycles = 0
	for len(toExecute) > 0 {
		var wg sync.WaitGroup
		for _, command := range toExecute {
			wg.Add(1)
			go func(command Command) {
				defer wg.Done()
				command.Execute(g.Target)
			}(command)
		}

		wg.Wait()
		toExecute, left = filterRelatedCommands(left)
		cycles++
	}

	return cycles, nil
}

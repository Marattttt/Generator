package generator

import (
	"image"

	"github.com/marattttt/paperwork/generator/color"
	"github.com/marattttt/paperwork/generator/drawing"
)

type Command interface {
	GetAffectedArea() image.Rectangle
	Execute(*drawing.Drawing) error
}

// Target is defined in the generator
type DrawLineCommand struct {
	Line drawing.Line
	Grad color.Gradient
}

func (command DrawLineCommand) Execute(target *drawing.Drawing) error {
	drawing.DrawLine(target, command.Line, command.Grad)
	return nil
}

func (c DrawLineCommand) GetAffectedArea() image.Rectangle {
	return c.Line.GetAffectedArea()
}

func filterRelatedCommands(unfiltered []Command) (filtered, left []Command) {
	filtered = make([]Command, 0)
	left = make([]Command, 0)

	var area1, area2 image.Rectangle
	var isAddable bool

	for _, newCom := range unfiltered {
		isAddable = true
		for _, filteredCom := range filtered {
			if !isAddable {
				break
			}

			area1 = newCom.GetAffectedArea()
			area2 = filteredCom.GetAffectedArea()

			isAddable = isAddable && area1.Intersect(area2) == image.Rectangle{}
		}

		if isAddable {
			filtered = append(filtered, newCom)
		} else {
			left = append(left, newCom)
		}
	}

	return filtered, left
}

package pattern

// Gradient is the core color elemetn used across the package
type Gradienter interface {
	GetGradient() *Gradient
}

// Holds an always sorted slice of gradient marks
// Cocurrent read operations are safe, while writing uses a mutex specific to an insttance of the struct
type Gradient struct {
	Marks []GradientMark
}

type GradientMark struct {
	Col Color
	Pos float32
}

func (g *Gradient) GetGradient() *Gradient {
	return g
}

// Changes an existing one or inserts a new mark to the gradient in ascending order
func (g *Gradient) Mark(mark GradientMark) {
	if len(g.Marks) == 0 {
		g.Marks = []GradientMark{mark}
		return
	}

	for i, m := range g.Marks {
		if m.Pos < mark.Pos {
			continue
		}

		if m.Pos == mark.Pos {
			g.Marks[i] = mark
			return
		}

		g.Marks = append(g.Marks[:i+1], g.Marks[i:]...)
		g.Marks[i] = mark
		return
	}

	g.Marks = append(g.Marks, mark)
}

// Assumes the gradient has at least 2 marks
func (g *Gradient) GetMark(start, end, pos int) GradientMark {
	if pos <= start {
		return g.Marks[0]
	}
	if pos >= end {
		return g.Marks[len(g.Marks)-1]
	}

	if len(g.Marks) == 2 && g.Marks[0] == g.Marks[1] {
		return g.Marks[0]
	}

	progress := float32(pos-start) / float32(end-start)

	for i := 1; i < len(g.Marks); i++ {
		if g.Marks[i-1].Pos <= progress && g.Marks[i].Pos >= progress {
			left := g.Marks[i-1].Col
			right := g.Marks[i].Col

			resR := left.R + uint16(progress-g.Marks[i-1].Pos*float32(right.R))
			resG := left.G + uint16(progress-g.Marks[i-1].Pos*float32(right.G))
			resB := left.B + uint16(progress-g.Marks[i-1].Pos*float32(right.B))
			resA := left.A + uint16(progress-g.Marks[i-1].Pos*float32(right.A))

			return GradientMark{
				Pos: progress,
				Col: Color{resR, resG, resB, resA},
			}
		}
	}

	return g.Marks[len(g.Marks)-1]
}

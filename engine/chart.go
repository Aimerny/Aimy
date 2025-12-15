package engine

import "github.com/hajimehoshi/ebiten/v2"

type Note struct {
	Lane   int   `json:"lane"`
	TimeMs int64 `json:"time_ms"`
	Hit    bool  `json:"hit"`
	Miss   bool  `json:"miss"`
}

type Chart struct {
	Song     string  `json:"song"`
	BPM      int     `json:"bpm"`
	OffsetMs int     `json:"offset_ms"`
	Notes    []*Note `json:"notes"`
}

func DemoChart() Chart {
	return Chart{
		Song:     "Demo",
		BPM:      120,
		OffsetMs: 0,
		Notes: []*Note{
			{0, 1000, false, false},
			{1, 1500, false, false},
			{2, 2000, false, false},
			{3, 2500, false, false},
			{0, 3000, false, false},
			{1, 3500, false, false},
			{2, 4000, false, false},
			{3, 4500, false, false},
			{0, 5000, false, false},
			{1, 5500, false, false},
			{2, 6000, false, false},
			{3, 6500, false, false},
			{0, 7000, false, false},
			{1, 7500, false, false},
			{2, 8000, false, false},
			{3, 8500, false, false},
		},
	}
}

func (n *Note) Draw(screen *ebiten.Image) {

}

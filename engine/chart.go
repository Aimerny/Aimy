package engine

type Note struct {
	Lane   int `json:"lane"`
	TimeMs int `json:"time_ms"`
}

type Chart struct {
	Song     string `json:"song"`
	BPM      int    `json:"bpm"`
	OffsetMs int    `json:"offset_ms"`
	Notes    []Note `json:"notes"`
}

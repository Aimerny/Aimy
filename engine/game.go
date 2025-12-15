package engine

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var laneKeys = [4]ebiten.Key{
	ebiten.KeyD,
	ebiten.KeyF,
	ebiten.KeyJ,
	ebiten.KeyK,
}

const (
	laneCount       int     = 4
	laneGap                 = 2
	laneWidth               = 60
	windowsLength           = 640
	windowsHeight           = 480
	judgeLineY              = 400
	judgeLineHeight         = 2
	velocityUnit            = 1
	noteHeight              = 20
	TPS                     = 60
	FULL_SCORE      float64 = 1_000_000.0

	PerfectWindow = 0.05
	GreatWindow   = 0.08
	GoodWindow    = 0.10
	MissWindow    = 0.12
)

type GameMode int

const (
	ModeTitle GameMode = iota
	ModePlaying
)

type Game struct {
	mode GameMode

	startTime     time.Time
	songStarted   bool
	songTimeMs    int64
	combo         int64
	velocity      float64
	laneAreaX     float64
	laneAreaWidth int
	judge         string
	score         float64
	scorePerNote  float64
	chartEndMs    int64

	chart       *Chart
	startButton Button

	laneImage        *ebiten.Image
	judgeLineImage   *ebiten.Image
	noteImage        *ebiten.Image
	startButtonImage *ebiten.Image
}

func NewGame(weight, height int) *Game {
	g := &Game{}
	// init properties
	g.combo = 0
	g.velocity = 0.2
	g.laneAreaWidth = laneCount*laneWidth + (laneCount-1)*laneGap
	g.laneAreaX = float64(windowsLength-g.laneAreaWidth) / 2
	g.startButton = Button{
		X:    300,
		Y:    200,
		W:    200,
		H:    80,
		Text: "Start",
	}

	g.laneImage = ebiten.NewImage(laneWidth, height)
	g.laneImage.Fill(color.RGBA{50, 50, 50, 255})

	g.judgeLineImage = ebiten.NewImage(g.laneAreaWidth, judgeLineHeight)
	g.judgeLineImage.Fill(color.RGBA{0xbb, 0xff, 0xff, 0xff})

	g.noteImage = ebiten.NewImage(laneWidth-2, noteHeight)
	g.noteImage.Fill(color.RGBA{0, 255, 255, 255})

	g.startButtonImage = ebiten.NewImage(int(g.startButton.W), int(g.startButton.H))
	g.startButtonImage.Fill(color.RGBA{0, 255, 255, 255})
	return g
}

func (g *Game) Start(chart *Chart) {
	g.startTime = time.Now()
	g.songStarted = true
	g.chart = chart
	g.score = 0.0
	g.scorePerNote = FULL_SCORE / float64(len(g.chart.Notes))
	log.Printf("chart started at %v", g.startTime)
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		return g.updateTitle()
	case ModePlaying:
		return g.updatePlaying()
	default:
		return nil
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case ModeTitle:
		g.drawTitle(screen)
	case ModePlaying:
		g.drawPlaying(screen)
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	// 背景
	screen.Fill(color.RGBA{0, 0, 0, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.startButton.X, g.startButton.Y)
	screen.DrawImage(g.startButtonImage, op)
	ebitenutil.DebugPrintAt(screen, "Click to Start", int(g.startButton.X+20), int(g.startButton.Y+8))
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	width, height := screen.Size()
	totalWidth := laneCount*laneWidth + (laneCount-1)*laneGap
	startX := (width - totalWidth) / 2

	for i := 0; i < laneCount; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(startX+i*(laneWidth+laneGap)), 0)
		screen.DrawImage(g.laneImage, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(startX), float64(judgeLineY))
	screen.DrawImage(g.judgeLineImage, op)
	_ = height
	g.drawNotes(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Time: %vms, Combo: %v, Judge: %s, Score: %d", g.songTimeMs, g.combo, g.judge, int(g.score)))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowsLength, windowsHeight
}

func (g *Game) processInputAndJudgement() {
	for lane, key := range laneKeys {
		if inpututil.IsKeyJustPressed(key) {
			g.judgeLane(lane)
		}
	}
}

func (g *Game) noteY(n *Note) float64 {
	dt := n.TimeMs - g.songTimeMs
	return judgeLineY - velocityUnit*g.velocity*float64(dt)
}

func (g *Game) judgeLane(lane int) {
	var target *Note
	for _, note := range g.chart.Notes {
		if note.Lane != lane || note.Miss || note.Hit {
			continue
		}
		target = note
		break
	}
	if target == nil {
		return
	}
	// get target note judge
	delta := g.songTimeMs - target.TimeMs
	ad := math.Abs(float64(delta) / 1000)
	heavy := 1.0

	switch {
	case ad <= PerfectWindow:
		//log.Printf("Perfect-%f\n", ad)
		g.judge = "Perfect"
		g.combo++
	case ad <= GreatWindow:
		log.Printf("Great-%f\n", ad)
		g.judge = "Great"
		g.combo++
		heavy = 0.8
	case ad <= GoodWindow:
		log.Printf("Good-%f\n", ad)
		g.judge = "Good"
		g.combo++
		heavy = 0.5
	case ad <= MissWindow:
		log.Printf("Miss-%f\n", ad)
		g.judge = "Miss"
		g.combo = 0
		heavy = 0
	default:
		//log.Printf("Null Press\n")
		return
	}
	target.Hit = true
	g.score += heavy * g.scorePerNote
}

func (g *Game) drawNotes(screen *ebiten.Image) {
	for _, note := range g.chart.Notes {
		if note.Hit || note.Miss {
			continue
		}
		y := g.noteY(note)
		// if too early, skip
		if y < -100 || y > windowsHeight {
			continue
		}

		x := g.laneX(note.Lane)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(g.noteImage, op)
	}
}

func (g *Game) laneX(lane int) float64 {
	if lane > 3 || lane < 0 {
		log.Fatal("lane out of range")
	}
	return g.laneAreaX + float64((laneWidth+laneGap)*lane)
}

func (g *Game) updateAutoMiss() {
	for _, note := range g.chart.Notes {
		if note.Hit || note.Miss {
			continue
		}
		if g.songTimeMs-note.TimeMs > MissWindow*1000 {
			log.Printf("%v Miss\n", note)
			note.Miss = true
			g.judge = "Miss"
			g.combo = 0
		}
	}
}

func (g *Game) updateTitle() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.startButton.Contains(float64(x), float64(y)) {
			g.startNewChart("chart1.json")
		}
	}
	return nil
}

func (g *Game) updatePlaying() error {
	g.songTimeMs = time.Since(g.startTime).Milliseconds()
	g.processInputAndJudgement()
	g.updateAutoMiss()
	if g.songTimeMs > g.chartEndMs {
		g.mode = ModeTitle
	}
	return nil
}

func (g *Game) startNewChart(chartPath string) {
	chart := DemoChart()
	g.chart = &chart
	lastNote := chart.Notes[len(chart.Notes)-1]
	g.chartEndMs = lastNote.TimeMs + 3000
	g.Start(&chart)
	g.mode = ModePlaying
}

type Button struct {
	X, Y float64
	W, H float64
	Text string
}

func (b Button) Contains(x, y float64) bool {
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

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
	laneCount       int = 4
	laneGap             = 2
	laneWidth           = 60
	windowsLength       = 640
	windowsHeight       = 480
	judgeLineY          = 400
	judgeLineHeight     = 2
	velocityUnit        = 1
	noteHeight          = 20
	TPS                 = 60

	PerfectWindow = 0.05
	GreatWindow   = 0.08
	GoodWindow    = 0.10
	MissWindow    = 0.12
)

type Game struct {
	startTime     time.Time
	songStarted   bool
	songTimeMs    int64
	combo         int64
	velocity      float64
	laneAreaX     float64
	laneAreaWidth int

	chart *Chart

	laneImage      *ebiten.Image
	judgeLineImage *ebiten.Image
	noteImage      *ebiten.Image
}

func NewGame(weight, height int) *Game {
	g := &Game{}
	// init properties
	g.combo = 0
	g.velocity = 0.2
	g.laneAreaWidth = laneCount*laneWidth + (laneCount-1)*laneGap
	g.laneAreaX = float64(windowsLength-g.laneAreaWidth) / 2

	g.laneImage = ebiten.NewImage(laneWidth, height)
	g.laneImage.Fill(color.RGBA{50, 50, 50, 255})

	g.judgeLineImage = ebiten.NewImage(g.laneAreaWidth, judgeLineHeight)
	g.judgeLineImage.Fill(color.RGBA{0xbb, 0xff, 0xff, 0xff})

	g.noteImage = ebiten.NewImage(laneWidth-2, noteHeight)
	g.noteImage.Fill(color.RGBA{0, 255, 255, 255})
	return g
}

func (g *Game) Start(chart *Chart) {
	g.startTime = time.Now()
	g.songStarted = true
	g.chart = chart
	log.Printf("chart started at %v", g.startTime)
}

func (g *Game) Update() error {
	if g.songStarted {
		g.songTimeMs = time.Since(g.startTime).Milliseconds()
	} else {
		chart := DemoChart()
		g.Start(&chart)
	}

	g.processInputAndJudgement()
	g.updateAutoMiss()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	width, height := screen.Size()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Time: %vms, Combo: %v", g.songTimeMs, g.combo))
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
	for _, note := range *(g.chart.Notes) {
		if note.Lane != lane || note.Miss || note.Hit {
			continue
		}
		target = &note
		break
	}
	if target == nil {
		return
	}
	// get target note judge
	delta := g.songTimeMs - target.TimeMs
	ad := math.Abs(float64(delta) / 1000)

	switch {
	case ad <= PerfectWindow:
		log.Printf("Perfect-%f\n", ad)
	case ad <= GreatWindow:
		log.Printf("Great-%f\n", ad)
	case ad <= GoodWindow:
		log.Printf("Good-%f\n", ad)
	case ad <= MissWindow:
		log.Printf("Miss-%f\n", ad)
	default:
		log.Printf("Null Press\n")
		return
	}
	target.Hit = true
}

func (g *Game) drawNotes(screen *ebiten.Image) {
	for _, note := range *(g.chart.Notes) {
		if note.Hit || note.Miss {
			continue
		}
		y := g.noteY(&note)
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
	for i := range *g.chart.Notes {
		note := &(*g.chart.Notes)[i]
		if note.Hit || note.Miss {
			continue
		}
		if g.songTimeMs-note.TimeMs > MissWindow*1000 {
			log.Printf("%v Miss\n", note)
			note.Miss = true
		}
	}
}

package engine

import (
	"fmt"
	"image/color"
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
)

type Game struct {
	startTime   time.Time
	songStarted bool
	songTimeMs  int64
	combo       int64

	laneImage *ebiten.Image
	judgeLine *ebiten.Image
}

func NewGame(weight, height int) *Game {
	g := &Game{}
	g.laneImage = ebiten.NewImage(laneWidth, height)
	g.laneImage.Fill(color.RGBA{50, 50, 50, 255})
	totalWidth := laneCount*laneWidth + (laneCount-1)*laneGap
	g.judgeLine = ebiten.NewImage(totalWidth, judgeLineHeight)
	g.judgeLine.Fill(color.RGBA{0xbb, 0xff, 0xff, 0xff})

	return g
}

func (g *Game) Start() {
	g.startTime = time.Now()
	g.songStarted = true
}

func (g *Game) Update() error {
	if g.songStarted {
		g.songTimeMs = time.Since(g.startTime).Milliseconds()
	} else {
		g.Start()
	}
	for lane, key := range laneKeys {
		if inpututil.IsKeyJustPressed(key) {
			g.handleHit(lane, g.songTimeMs)
		}
	}
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
	screen.DrawImage(g.judgeLine, op)
	_ = height
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowsLength, windowsHeight
}

func (g *Game) handleHit(lane int, songTimeMs int64) {
	fmt.Printf("lane: %d, songTimeMs: %d\n", lane, songTimeMs)
	g.combo++
}

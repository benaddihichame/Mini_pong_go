package main

import (
	"log"

	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenwidth  = 640
	screenheight = 480
	ballspeed    = 3
	paddleSpeed  = 6
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	dxdt int
	dydt int
}

type Gaming struct {
	paddle    Paddle
	ball      Ball
	score     int
	bestScore int
}

func main() {
	ebiten.SetWindowTitle("Pong hichome Go")
	ebiten.SetWindowSize(screenwidth, screenheight)
	paddle := Paddle{
		Object: Object{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}
	ball := Ball{
		Object: Object{
			X: screenwidth / 2,
			Y: screenheight / 2,
			W: 15,
			H: 15,
		},
		dxdt: ballspeed,
		dydt: ballspeed,
	}
	g := &Gaming{
		paddle: paddle,
		ball:   ball,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Gaming) Layout(outSideWidth, outSideHeight int) (int, int) {
	return screenwidth, screenheight
}

func (g *Gaming) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(g.paddle.X), float32(g.paddle.Y),
		float32(g.paddle.W), float32(g.paddle.H),
		color.White, false)
	vector.DrawFilledRect(screen,
		float32(g.ball.X), float32(g.ball.Y),
		float32(g.ball.W), float32(g.ball.H),
		color.White, false)
	scoreStr := "Score : " + fmt.Sprint((g.score))
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 10, color.White)

	bestscoreStr := "best Score : " + fmt.Sprint((g.bestScore))
	text.Draw(screen, bestscoreStr, basicfont.Face7x13, 10, 30, color.White)
}

func (g *Gaming) Update() error {
	g.paddle.MoveOneKeyPress()
	g.ball.Move()
	g.CollisionPaddel()
	g.CollisionWall()
	return nil
}

func (p *Paddle) MoveOneKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		p.Object.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Object.Y -= paddleSpeed
	}
}

func (b *Ball) Move() {
	b.Object.X += b.dxdt
	b.Object.Y += b.dydt
}

func (g *Gaming) Reset() {
	g.ball.Object.Y = screenheight / 2
	g.ball.Object.X = screenwidth / 2
	g.ball.dxdt = ballspeed
	g.ball.dydt = ballspeed
	g.score = 0
}

func (g *Gaming) CollisionWall() {
	if g.ball.Object.X <= 0 {
		g.ball.dxdt = ballspeed
	} else if g.ball.Object.X+g.ball.Object.W >= screenwidth {
		g.Reset()
	} else if g.ball.Object.Y <= 0 {
		g.ball.dydt = ballspeed
	} else if g.ball.Object.Y+g.ball.Object.H >= screenheight {
		g.ball.dydt = -ballspeed
	}
}

func (g *Gaming) CollisionPaddel() {
	if g.ball.Object.X <= g.paddle.Object.X+g.paddle.Object.W && g.ball.Object.X+g.ball.Object.W >= g.paddle.Object.X && g.ball.Object.Y >= g.paddle.Object.Y && g.ball.Object.Y <= g.paddle.Object.Y+g.paddle.Object.H {
		g.ball.dxdt = -g.ball.dxdt
		g.score++
		if g.score > g.bestScore {
			g.bestScore = g.score
		}
	}
}

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
	paddle1   Paddle
	paddle2   Paddle
	ball      Ball
	score     int
	bestScore int
}

func main() {
	ebiten.SetWindowTitle("Pong hichome Go")
	ebiten.SetWindowSize(screenwidth, screenheight)
	paddle1 := Paddle{
		Object: Object{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}
	paddle2 := Paddle{
		Object: Object{
			X: 10,
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
		paddle1: paddle1,
		paddle2: paddle2,
		ball:    ball,
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
		float32(g.paddle1.X), float32(g.paddle1.Y),
		float32(g.paddle1.W), float32(g.paddle1.H),
		color.White, false)
	vector.DrawFilledRect(screen,
		float32(g.paddle2.X), float32(g.paddle2.Y),
		float32(g.paddle2.W), float32(g.paddle2.H),
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
	g.movePaddles()
	g.ball.Move()
	g.CollisionPaddle()
	g.CollisionWall()
	return nil
}

func (g *Gaming) movePaddles() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.paddle1.Object.Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.paddle1.Object.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.paddle2.Object.Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyN) {
		g.paddle2.Object.Y += paddleSpeed
	}

	// Limiter le mouvement des paddles à l'écran
	g.paddle1.Object.Y = clamp(g.paddle1.Object.Y, 0, screenheight-g.paddle1.Object.H)
	g.paddle2.Object.Y = clamp(g.paddle2.Object.Y, 0, screenheight-g.paddle2.Object.H)
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
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
		g.Reset()
	} else if g.ball.Object.X+g.ball.Object.W >= screenwidth {
		g.Reset()
	} else if g.ball.Object.Y <= 0 {
		g.ball.dydt = ballspeed
	} else if g.ball.Object.Y+g.ball.Object.H >= screenheight {
		g.ball.dydt = -ballspeed
	}
}

func (g *Gaming) CollisionPaddle() {
	if g.ball.Object.X <= g.paddle2.Object.X+g.paddle2.Object.W && g.ball.Object.X+g.ball.Object.W >= g.paddle2.Object.X && g.ball.Object.Y+g.ball.Object.H >= g.paddle2.Object.Y && g.ball.Object.Y <= g.paddle2.Object.Y+g.paddle2.Object.H {
		g.ball.dxdt = ballspeed
		g.score++
		if g.score > g.bestScore {
			g.bestScore = g.score
		}
	}
	if g.ball.Object.X+g.ball.Object.W >= g.paddle1.Object.X && g.ball.Object.X <= g.paddle1.Object.X+g.paddle1.Object.W && g.ball.Object.Y+g.ball.Object.H >= g.paddle1.Object.Y && g.ball.Object.Y <= g.paddle1.Object.Y+g.paddle1.Object.H {
		g.ball.dxdt = -ballspeed
		g.score++
		if g.score > g.bestScore {
			g.bestScore = g.score
		}
	}
}

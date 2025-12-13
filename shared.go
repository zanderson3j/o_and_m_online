package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)


// GameInterface defines the interface all games must implement
type GameInterface interface {
	Update(gr *GameRoom) error
	Draw(screen *ebiten.Image, gr *GameRoom)
	Reset()
}

// Button represents a clickable button
type Button struct {
	x, y, width, height float64
	text                string
	enabled             bool
	hovered             bool
}

func (b *Button) Contains(x, y int) bool {
	fx, fy := float64(x), float64(y)
	return fx >= b.x && fx <= b.x+b.width && fy >= b.y && fy <= b.y+b.height
}

// Check if O&M logo was clicked
func IsLogoClicked() bool {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}

	x, y := ebiten.CursorPosition()
	// Logo is at position (15, 15) with scale 3, so roughly 120x60 pixels
	return x >= 15 && x <= 135 && y >= 15 && y <= 75
}

// Draw forest background - shared across all screens
func DrawForestBackground(screen *ebiten.Image) {
	
	// Sky gradient
	skyTop := color.RGBA{20, 40, 80, 255}
	skyBottom := color.RGBA{40, 70, 110, 255}
	for y := 0; y < screenHeight/2; y++ {
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, skyTop, false)
	}
	for y := screenHeight / 2; y < screenHeight; y++ {
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, skyBottom, false)
	}

	// Back layer - distant dark mountains/trees
	darkForest := color.RGBA{15, 35, 60, 255}
	vector.DrawFilledRect(screen, 0, 300, screenWidth, 200, darkForest, false)

	// Draw distant tree silhouettes
	for i := 0; i < 20; i++ {
		x := float32(i * 60)
		y := float32(250 + (i%3)*20)
		// Triangle tree tops
		for h := 0; h < 80; h++ {
			width := float32(h / 2)
			vector.DrawFilledRect(screen, x+20-width, y+float32(h), width*2, 2, darkForest, false)
		}
		// Trunk
		vector.DrawFilledRect(screen, x+15, y+80, 10, 40, darkForest, false)
	}

	// Middle layer - teal/green forest
	midForest := color.RGBA{20, 80, 80, 255}
	for i := 0; i < 15; i++ {
		x := float32(i*70 + 30)
		y := float32(350 + (i%2)*30)
		// Rounder tree tops
		for h := 0; h < 100; h++ {
			width := float32(h / 3)
			if h > 50 {
				width = float32((100 - h) / 3)
			}
			vector.DrawFilledRect(screen, x+25-width, y+float32(h), width*2, 2, midForest, false)
		}
		// Trunk
		vector.DrawFilledRect(screen, x+20, y+100, 10, 50, color.RGBA{15, 50, 50, 255}, false)
	}

	// Front layer - brightest teal trees
	frontForest := color.RGBA{30, 120, 100, 255}
	for i := 0; i < 10; i++ {
		x := float32(i*110 + 20)
		y := float32(450 + (i%3)*20)
		// Large foreground trees
		for h := 0; h < 150; h++ {
			width := float32(h / 2)
			if h > 75 {
				width = float32((150 - h) / 2)
			}
			vector.DrawFilledRect(screen, x+35-width, y+float32(h), width*2, 3, frontForest, false)
		}
		// Trunk
		vector.DrawFilledRect(screen, x+25, y+150, 20, 100, color.RGBA{20, 70, 60, 255}, false)
	}
}

// Draw kodama spirits - shared across all screens
func DrawKodamaSpirits(screen *ebiten.Image) {
	
	kodamaColor := color.RGBA{200, 255, 220, 200}
	kodamaGlow := color.RGBA{150, 255, 200, 100}

	kodamaPositions := [][2]float32{
		{100, 400}, {250, 350}, {450, 380}, {650, 420},
		{850, 360}, {150, 500}, {550, 480}, {900, 520},
	}

	for _, pos := range kodamaPositions {
		x, y := pos[0], pos[1]

		// Glow effect
		vector.DrawFilledCircle(screen, x, y, 12, kodamaGlow, false)

		// Simple kodama head (round with eyes)
		vector.DrawFilledCircle(screen, x, y, 8, kodamaColor, false)

		// Eyes (dark dots)
		vector.DrawFilledRect(screen, x-3, y-2, 2, 3, color.RGBA{40, 40, 40, 255}, false)
		vector.DrawFilledRect(screen, x+1, y-2, 2, 3, color.RGBA{40, 40, 40, 255}, false)

		// Simple body
		vector.DrawFilledRect(screen, x-4, y+8, 8, 10, kodamaColor, false)
	}
}

// Draw O&M logo - shared across all screens
func DrawOMLogo(screen *ebiten.Image) {
	x, y := float32(15), float32(15)
	scale := float32(3)
	p := scale

	oColor := color.RGBA{255, 220, 150, 255}
	sparkleColor := color.RGBA{255, 255, 200, 255}

	// O outer ring
	vector.DrawFilledRect(screen, x+4*p, y+3*p, 2*p, 14*p, oColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+3*p, 2*p, 14*p, oColor, false)
	vector.DrawFilledRect(screen, x+6*p, y+3*p, 6*p, 2*p, oColor, false)
	vector.DrawFilledRect(screen, x+6*p, y+15*p, 6*p, 2*p, oColor, false)
	vector.DrawFilledRect(screen, x+5*p, y+5*p, 1*p, 1*p, sparkleColor, false)

	// Ampersand "&"
	ampColor := color.RGBA{180, 230, 255, 255}
	vector.DrawFilledRect(screen, x+17*p, y+4*p, 2*p, 2*p, ampColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+6*p, 2*p, 3*p, ampColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+9*p, 4*p, 2*p, ampColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+11*p, 2*p, 3*p, ampColor, false)
	vector.DrawFilledRect(screen, x+17*p, y+14*p, 2*p, 2*p, ampColor, false)
	vector.DrawFilledRect(screen, x+20*p, y+13*p, 2*p, 3*p, ampColor, false)
	vector.DrawFilledRect(screen, x+19*p, y+8*p, 1*p, 1*p, sparkleColor, false)

	// Letter "M"
	mColor := color.RGBA{255, 180, 200, 255}
	vector.DrawFilledRect(screen, x+24*p, y+4*p, 2*p, 13*p, mColor, false)
	vector.DrawFilledRect(screen, x+26*p, y+6*p, 2*p, 4*p, mColor, false)
	vector.DrawFilledRect(screen, x+28*p, y+8*p, 2*p, 2*p, mColor, false)
	vector.DrawFilledRect(screen, x+30*p, y+6*p, 2*p, 4*p, mColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+4*p, 2*p, 13*p, mColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+5*p, 1*p, 1*p, sparkleColor, false)

	// Floating sparkles
	vector.DrawFilledRect(screen, x+2*p, y+9*p, 1*p, 1*p, color.RGBA{200, 255, 220, 200}, false)
	vector.DrawFilledRect(screen, x+38*p, y+7*p, 1*p, 1*p, color.RGBA{255, 220, 200, 200}, false)
	vector.DrawFilledRect(screen, x+36*p, y+16*p, 1*p, 1*p, color.RGBA{220, 200, 255, 200}, false)
}

// Draw Player 1 avatar (human) - shared across games
func DrawPlayer1Avatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	shirtColor := color.RGBA{45, 62, 80, 255}
	vector.DrawFilledRect(screen, x, y+42*p, 50*p, 8*p, shirtColor, false)

	neckColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledRect(screen, x+20*p, y+36*p, 10*p, 8*p, neckColor, false)

	skinColor := color.RGBA{255, 228, 200, 255}
	vector.DrawFilledRect(screen, x+12*p, y+14*p, 26*p, 24*p, skinColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+12*p, 22*p, 4*p, skinColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+36*p, 18*p, 2*p, skinColor, false)

	hairDark := color.RGBA{50, 38, 30, 255}
	hairMid := color.RGBA{75, 60, 48, 255}
	hairLight := color.RGBA{95, 78, 62, 255}

	vector.DrawFilledRect(screen, x+10*p, y+2*p, 30*p, 12*p, hairDark, false)
	vector.DrawFilledRect(screen, x+8*p, y+4*p, 4*p, 8*p, hairMid, false)
	vector.DrawFilledRect(screen, x+38*p, y+4*p, 4*p, 8*p, hairMid, false)
	vector.DrawFilledRect(screen, x+18*p, y+4*p, 8*p, 4*p, hairLight, false)
	vector.DrawFilledRect(screen, x+24*p, y+6*p, 6*p, 3*p, hairLight, false)

	vector.DrawFilledRect(screen, x+8*p, y+10*p, 6*p, 14*p, hairDark, false)
	vector.DrawFilledRect(screen, x+36*p, y+10*p, 6*p, 14*p, hairDark, false)
	vector.DrawFilledRect(screen, x+10*p, y+16*p, 3*p, 6*p, hairMid, false)

	browColor := color.RGBA{60, 48, 38, 255}
	vector.DrawFilledRect(screen, x+16*p, y+18*p, 7*p, 2*p, browColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+18*p, 7*p, 2*p, browColor, false)

	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledRect(screen, x+16*p, y+22*p, 7*p, 6*p, eyeWhite, false)
	vector.DrawFilledRect(screen, x+27*p, y+22*p, 7*p, 6*p, eyeWhite, false)

	irisColor := color.RGBA{100, 145, 120, 255}
	vector.DrawFilledRect(screen, x+17*p, y+23*p, 5*p, 5*p, irisColor, false)
	vector.DrawFilledRect(screen, x+28*p, y+23*p, 5*p, 5*p, irisColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+19*p, y+24*p, 2*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+30*p, y+24*p, 2*p, 3*p, pupilColor, false)

	vector.DrawFilledRect(screen, x+20*p, y+24*p, 2*p, 2*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+31*p, y+24*p, 2*p, 2*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+18*p, y+26*p, 1*p, 1*p, color.RGBA{200, 220, 255, 180}, false)
	vector.DrawFilledRect(screen, x+29*p, y+26*p, 1*p, 1*p, color.RGBA{200, 220, 255, 180}, false)

	noseColor := color.RGBA{240, 200, 170, 255}
	vector.DrawFilledRect(screen, x+24*p, y+28*p, 2*p, 3*p, noseColor, false)

	mouthColor := color.RGBA{200, 120, 110, 255}
	vector.DrawFilledRect(screen, x+20*p, y+33*p, 10*p, 2*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+19*p, y+32*p, 2*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+32*p, 2*p, 1*p, mouthColor, false)

	blushColor := color.RGBA{255, 180, 160, 100}
	vector.DrawFilledRect(screen, x+14*p, y+28*p, 4*p, 3*p, blushColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+28*p, 4*p, 3*p, blushColor, false)

	shadowColor := color.RGBA{235, 200, 170, 80}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 2*p, 12*p, shadowColor, false)

	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{88, 166, 255, 255}, false)
}

// Draw Player 2 avatar (dog) - shared across games
func DrawPlayer2Avatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	chestColor := color.RGBA{210, 170, 120, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, chestColor, false)

	neckColor := color.RGBA{195, 160, 110, 255}
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, neckColor, false)

	headColor := color.RGBA{200, 155, 100, 255}

	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 20*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+14*p, 18*p, 6*p, headColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+10*p, 14*p, 6*p, headColor, false)

	vector.DrawFilledRect(screen, x+16*p, y+16*p, 4*p, 3*p, color.RGBA{90, 70, 50, 80}, false)
	vector.DrawFilledRect(screen, x+30*p, y+16*p, 4*p, 3*p, color.RGBA{90, 70, 50, 80}, false)

	muzzleColor := color.RGBA{220, 185, 140, 255}
	vector.DrawFilledRect(screen, x+16*p, y+26*p, 18*p, 12*p, muzzleColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+34*p, 14*p, 4*p, muzzleColor, false)

	earDark := color.RGBA{65, 50, 38, 255}
	earInside := color.RGBA{230, 180, 140, 255}

	vector.DrawFilledRect(screen, x+10*p, y+8*p, 8*p, 12*p, earDark, false)
	vector.DrawFilledRect(screen, x+12*p, y+10*p, 4*p, 8*p, earInside, false)

	vector.DrawFilledRect(screen, x+32*p, y+8*p, 8*p, 12*p, earDark, false)
	vector.DrawFilledRect(screen, x+34*p, y+10*p, 4*p, 8*p, earInside, false)

	eyeWhite := color.RGBA{255, 255, 255, 255}
	eyeColor := color.RGBA{120, 85, 50, 255}

	vector.DrawFilledRect(screen, x+18*p, y+22*p, 6*p, 5*p, eyeWhite, false)
	vector.DrawFilledRect(screen, x+26*p, y+22*p, 6*p, 5*p, eyeWhite, false)

	vector.DrawFilledRect(screen, x+19*p, y+23*p, 4*p, 4*p, eyeColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+23*p, 4*p, 4*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+20*p, y+24*p, 2*p, 2*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+28*p, y+24*p, 2*p, 2*p, pupilColor, false)

	vector.DrawFilledRect(screen, x+21*p, y+24*p, 2*p, 2*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+29*p, y+24*p, 2*p, 2*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+19*p, y+25*p, 1*p, 1*p, color.RGBA{220, 240, 255, 180}, false)
	vector.DrawFilledRect(screen, x+27*p, y+25*p, 1*p, 1*p, color.RGBA{220, 240, 255, 180}, false)

	noseColor := color.RGBA{30, 30, 30, 255}
	vector.DrawFilledRect(screen, x+23*p, y+32*p, 4*p, 3*p, noseColor, false)
	vector.DrawFilledRect(screen, x+24*p, y+32*p, 1*p, 1*p, color.RGBA{120, 120, 120, 255}, false)

	mouthColor := color.RGBA{80, 60, 50, 255}
	vector.DrawFilledRect(screen, x+21*p, y+35*p, 8*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+20*p, y+34*p, 2*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+28*p, y+34*p, 2*p, 1*p, mouthColor, false)

	tongueColor := color.RGBA{255, 140, 150, 255}
	vector.DrawFilledRect(screen, x+24*p, y+36*p, 2*p, 2*p, tongueColor, false)

	markingColor := color.RGBA{100, 75, 50, 80}
	vector.DrawFilledRect(screen, x+14*p, y+24*p, 3*p, 10*p, markingColor, false)
	vector.DrawFilledRect(screen, x+33*p, y+24*p, 3*p, 10*p, markingColor, false)

	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 193, 7, 255}, false)
}

// DrawButton draws a styled button
func DrawButton(screen *ebiten.Image, btn *Button) {
	btnColor := color.RGBA{100, 150, 220, 255}
	borderColor := color.RGBA{70, 120, 190, 255}
	if !btn.enabled {
		btnColor = color.RGBA{140, 160, 180, 255}
		borderColor = color.RGBA{100, 120, 140, 255}
	} else if btn.hovered {
		btnColor = color.RGBA{130, 180, 240, 255}
		borderColor = color.RGBA{100, 150, 220, 255}
	}

	// Soft shadow
	vector.DrawFilledRect(screen, float32(btn.x+2), float32(btn.y+2), float32(btn.width), float32(btn.height), color.RGBA{0, 0, 0, 40}, false)

	// Button background
	vector.DrawFilledRect(screen, float32(btn.x), float32(btn.y), float32(btn.width), float32(btn.height), btnColor, false)
	vector.StrokeRect(screen, float32(btn.x), float32(btn.y), float32(btn.width), float32(btn.height), 2, borderColor, false)

	textX := int(btn.x + btn.width/2 - float64(len(btn.text))*3)
	textY := int(btn.y + btn.height/2 - 5)
	ebitenutil.DebugPrintAt(screen, btn.text, textX, textY)
	ebitenutil.DebugPrintAt(screen, btn.text, textX+1, textY) // Bold
}

package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Animation phases for the intro
type IntroPhase int

const (
	PhaseOwliveFlyingDown IntroPhase = iota
	PhasePickingUp
	PhaseFlyingTogether
	PhaseWelcomeText
	PhaseTitleText
	PhaseFlyingAway
	PhaseAvatarMarch
	PhaseComplete
)

type IntroScreen struct {
	phase        IntroPhase
	frameCount   int
	phaseFrame   int // frames within current phase

	// Owlive position and animation
	owliveX      float64
	owliveY      float64
	targetX      float64 // for smooth transitions
	targetY      float64
	wingFrame    int // for wing flapping animation

	// Millipede position
	millipedeX   float64
	millipedeY   float64
	millipedePickedUp bool

	// Text animation
	welcomeAlpha float64
	titleAlpha   float64

	// Avatar march
	marchX       float64 // X position of the marching avatars

	// Skip functionality
	skipButton   *Button

	// Audio
	audioStarted bool
}

func NewIntroScreen() *IntroScreen {
	is := &IntroScreen{
		phase:      PhaseOwliveFlyingDown,
		owliveX:    float64(screenWidth) / 2,
		owliveY:    -100, // Start above screen
		millipedeX: float64(screenWidth)/2 - 25,
		millipedeY: float64(screenHeight) - 150, // On the ground
		skipButton: &Button{
			x:       float64(screenWidth) - 120,
			y:       float64(screenHeight) - 50,
			width:   100,
			height:  35,
			text:    "SKIP",
			enabled: true,
		},
	}
	return is
}

func (is *IntroScreen) Update(gr *GameRoom) error {
	is.frameCount++
	is.phaseFrame++
	is.wingFrame++

	// Start audio on first update
	if !is.audioStarted {
		PlayIntroTheme()
		is.audioStarted = true
	}

	// Handle skip button
	x, y := ebiten.CursorPosition()
	is.skipButton.hovered = is.skipButton.Contains(x, y)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && is.skipButton.hovered {
		StartIntroFadeOut()
		is.phase = PhaseComplete
	}

	// Animation state machine
	switch is.phase {
	case PhaseOwliveFlyingDown:
		// Owlive flies down from top
		is.owliveY += 4
		// Add slight horizontal wobble
		is.owliveX = float64(screenWidth)/2 + math.Sin(float64(is.frameCount)*0.1)*30

		if is.owliveY >= float64(screenHeight)-200 {
			is.phase = PhasePickingUp
			is.phaseFrame = 0
		}

	case PhasePickingUp:
		// Owlive hovers and grabs Millipede
		is.owliveX = is.millipedeX + 25 // Center over millipede
		is.owliveY = is.millipedeY - 60

		if is.phaseFrame > 30 {
			is.millipedePickedUp = true
		}
		if is.phaseFrame > 60 {
			is.phase = PhaseFlyingTogether
			is.phaseFrame = 0
		}

	case PhaseFlyingTogether:
		// They fly in a loop/arc together
		t := float64(is.phaseFrame) * 0.02
		centerX := float64(screenWidth) / 2
		centerY := float64(screenHeight) / 2

		// Circular/figure-8 path
		is.owliveX = centerX + math.Sin(t*2)*200
		is.owliveY = centerY + math.Sin(t)*100 - 50

		if is.phaseFrame > 180 {
			is.phase = PhaseWelcomeText
			is.phaseFrame = 0
			// Set target position for smooth transition (high up, above text)
			is.targetX = float64(screenWidth)/2 - 100
			is.targetY = 50
		}

	case PhaseWelcomeText:
		// Smoothly move owl to target position
		is.owliveX += (is.targetX - is.owliveX) * 0.05
		is.owliveY += (is.targetY - is.owliveY) * 0.05

		// Gentle hover once close to target
		if math.Abs(is.owliveY-is.targetY) < 5 {
			is.owliveY = is.targetY + math.Sin(float64(is.frameCount)*0.05)*10
		}

		// Fade in "Olive and Millie Welcome you to..." after settling
		if is.phaseFrame > 30 && is.welcomeAlpha < 255 {
			is.welcomeAlpha += 3
		}

		if is.phaseFrame > 180 {
			is.phase = PhaseTitleText
			is.phaseFrame = 0
			// Keep owl at same height for title phase
			is.targetY = 50
		}

	case PhaseTitleText:
		// Fade out welcome, fade in title
		if is.welcomeAlpha > 0 {
			is.welcomeAlpha -= 5
		}
		if is.welcomeAlpha <= 0 && is.titleAlpha < 255 {
			is.titleAlpha += 3
		}
		// Smoothly move owl higher and hover
		is.owliveX += (is.targetX - is.owliveX) * 0.05
		is.owliveY += (is.targetY - is.owliveY) * 0.05
		if math.Abs(is.owliveY-is.targetY) < 5 {
			is.owliveY = is.targetY + math.Sin(float64(is.frameCount)*0.05)*10
		}

		if is.phaseFrame > 200 {
			is.phase = PhaseFlyingAway
			is.phaseFrame = 0
		}

	case PhaseFlyingAway:
		// They fly off to the right
		is.owliveX += 8
		is.owliveY -= 2

		// Fade out title
		if is.titleAlpha > 0 {
			is.titleAlpha -= 3
		}

		if is.owliveX > float64(screenWidth)+200 {
			is.phase = PhaseAvatarMarch
			is.phaseFrame = 0
			is.marchX = -float64(int(AvatarNumTypes)) * 120 // Start off-screen left (all avatars with 120 spacing)
		}

	case PhaseAvatarMarch:
		// All avatars march across the screen
		is.marchX += 2

		// Calculate total width of marching avatars
		numAvatars := int(AvatarNumTypes)
		totalWidth := float64(numAvatars) * 120

		// Complete when all avatars have marched off screen
		if is.marchX > float64(screenWidth)+150 {
			is.phase = PhaseComplete
		}
		// Also allow early transition once they've all appeared and had time on screen
		if is.marchX > totalWidth && is.phaseFrame > 240 {
			is.phase = PhaseComplete
		}

	case PhaseComplete:
		// Fade out intro music when transitioning to next screen
		StartIntroFadeOut()
		// Signal to GameRoom that intro is done
		gr.introComplete = true
		// Open avatar selection if lobby is ready, otherwise flag for later
		if gr.lobbyScreen != nil {
			gr.lobbyScreen.ShowAvatarSelection()
		} else {
			gr.needsAvatarSelectShow = true
		}
	}

	return nil
}

func (is *IntroScreen) Draw(screen *ebiten.Image) {
	// Draw forest background
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)

	// Draw Millipede (on ground or carried)
	if !is.millipedePickedUp {
		is.drawGroundMillipede(screen)
	}

	// Draw Owlive flying
	is.drawFlyingOwlive(screen)

	// Draw carried Millipede if picked up
	if is.millipedePickedUp {
		is.drawCarriedMillipede(screen)
	}

	// Draw welcome text
	if is.welcomeAlpha > 0 {
		is.drawWelcomeText(screen)
	}

	// Draw title text
	if is.titleAlpha > 0 {
		is.drawTitleText(screen)
	}

	// Draw marching avatars during avatar march phase
	if is.phase == PhaseAvatarMarch {
		is.drawMarchingAvatars(screen)
	}

	// Draw skip button
	DrawButton(screen, is.skipButton)
}

func (is *IntroScreen) drawFlyingOwlive(screen *ebiten.Image) {
	x := float32(is.owliveX)
	y := float32(is.owliveY)
	scale := float32(3.0)
	p := scale

	// Wing flap cycle (0-3)
	wingState := (is.wingFrame / 8) % 4

	// Body
	bodyColor := color.RGBA{139, 90, 60, 255}
	vector.DrawFilledRect(screen, x+12*p, y+38*p, 26*p, 12*p, bodyColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+36*p, 22*p, 4*p, bodyColor, false)

	// Belly pattern
	bellyColor := color.RGBA{200, 180, 160, 255}
	vector.DrawFilledRect(screen, x+18*p, y+40*p, 14*p, 10*p, bellyColor, false)

	// Wings - animated based on wingState
	wingColor := color.RGBA{100, 65, 40, 255}
	wingTip := color.RGBA{80, 50, 30, 255}

	switch wingState {
	case 0: // Wings up
		// Left wing
		vector.DrawFilledRect(screen, x-10*p, y+20*p, 20*p, 8*p, wingColor, false)
		vector.DrawFilledRect(screen, x-15*p, y+15*p, 15*p, 10*p, wingColor, false)
		vector.DrawFilledRect(screen, x-20*p, y+10*p, 10*p, 8*p, wingTip, false)
		// Right wing
		vector.DrawFilledRect(screen, x+40*p, y+20*p, 20*p, 8*p, wingColor, false)
		vector.DrawFilledRect(screen, x+50*p, y+15*p, 15*p, 10*p, wingColor, false)
		vector.DrawFilledRect(screen, x+60*p, y+10*p, 10*p, 8*p, wingTip, false)
	case 1: // Wings mid-up
		// Left wing
		vector.DrawFilledRect(screen, x-5*p, y+25*p, 18*p, 8*p, wingColor, false)
		vector.DrawFilledRect(screen, x-12*p, y+22*p, 12*p, 8*p, wingTip, false)
		// Right wing
		vector.DrawFilledRect(screen, x+37*p, y+25*p, 18*p, 8*p, wingColor, false)
		vector.DrawFilledRect(screen, x+50*p, y+22*p, 12*p, 8*p, wingTip, false)
	case 2: // Wings level
		// Left wing
		vector.DrawFilledRect(screen, x-8*p, y+35*p, 20*p, 6*p, wingColor, false)
		vector.DrawFilledRect(screen, x-18*p, y+34*p, 12*p, 6*p, wingTip, false)
		// Right wing
		vector.DrawFilledRect(screen, x+38*p, y+35*p, 20*p, 6*p, wingColor, false)
		vector.DrawFilledRect(screen, x+56*p, y+34*p, 12*p, 6*p, wingTip, false)
	case 3: // Wings down
		// Left wing
		vector.DrawFilledRect(screen, x-5*p, y+42*p, 18*p, 8*p, wingColor, false)
		vector.DrawFilledRect(screen, x-12*p, y+48*p, 12*p, 8*p, wingTip, false)
		// Right wing
		vector.DrawFilledRect(screen, x+37*p, y+42*p, 18*p, 8*p, wingColor, false)
		vector.DrawFilledRect(screen, x+50*p, y+48*p, 12*p, 8*p, wingTip, false)
	}

	// Head (round and wide)
	headColor := color.RGBA{120, 80, 50, 255}
	vector.DrawFilledRect(screen, x+8*p, y+18*p, 34*p, 20*p, headColor, false)
	vector.DrawFilledRect(screen, x+10*p, y+16*p, 30*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+14*p, 26*p, 4*p, headColor, false)

	// Ear tufts
	for i := 0; i < 6; i++ {
		h := float32(6 - i)
		vector.DrawFilledRect(screen, x+10*p+float32(i)*p, y+10*p+float32(i)*p, 3*p, h*p, headColor, false)
	}
	for i := 0; i < 6; i++ {
		h := float32(6 - i)
		vector.DrawFilledRect(screen, x+37*p-float32(i)*p, y+10*p+float32(i)*p, 3*p, h*p, headColor, false)
	}

	// Face disc
	faceColor := color.RGBA{180, 160, 140, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 16*p, faceColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, faceColor, false)

	// Large eyes
	eyeBgColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 7*p, eyeBgColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 7*p, eyeBgColor, false)

	eyeColor := color.RGBA{255, 200, 80, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 5*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 5*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 2*p, pupilColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 2*p, pupilColor, false)

	// Beak
	beakColor := color.RGBA{255, 180, 100, 255}
	for i := 0; i < 5; i++ {
		w := float32(5 - i)
		vector.DrawFilledRect(screen, x+25*p-w*p/2, y+30*p+float32(i)*p, w*p, 1*p, beakColor, false)
	}

	// Talons (visible when flying)
	talonColor := color.RGBA{80, 60, 40, 255}
	// Left talon
	vector.DrawFilledRect(screen, x+15*p, y+50*p, 3*p, 8*p, talonColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+56*p, 3*p, 4*p, talonColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+56*p, 3*p, 4*p, talonColor, false)
	// Right talon
	vector.DrawFilledRect(screen, x+32*p, y+50*p, 3*p, 8*p, talonColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+56*p, 3*p, 4*p, talonColor, false)
	vector.DrawFilledRect(screen, x+35*p, y+56*p, 3*p, 4*p, talonColor, false)
}

func (is *IntroScreen) drawCarriedMillipede(screen *ebiten.Image) {
	// Draw millipede below Owlive's talons - same cute inchworm style as ground
	x := float32(is.owliveX) - 40
	y := float32(is.owliveY) + 120
	p := float32(2.0)

	// Colors - cute sage green with yellow stripes
	bodyColor := color.RGBA{158, 194, 145, 255}
	stripeColor := color.RGBA{220, 200, 130, 255}
	outlineColor := color.RGBA{70, 90, 60, 255}
	legColor := color.RGBA{130, 170, 120, 255}

	// Segment positions - same arch shape as ground/avatar
	type segment struct {
		sx, sy, r float32
	}
	segments := []segment{
		{sx: 5, sy: 38, r: 4},
		{sx: 9, sy: 38, r: 3.5},
		{sx: 13, sy: 38, r: 3.5},
		{sx: 16, sy: 35, r: 3.5},
		{sx: 19, sy: 28, r: 3.5},
		{sx: 21, sy: 21, r: 3.5},
		{sx: 23, sy: 15, r: 3.5},
		{sx: 26, sy: 12, r: 3.5},
		{sx: 29, sy: 15, r: 3.5},
		{sx: 31, sy: 21, r: 3.5},
		{sx: 33, sy: 28, r: 3.5},
		{sx: 36, sy: 35, r: 3.5},
		{sx: 39, sy: 38, r: 3.5},
		{sx: 43, sy: 38, r: 3.5},
		{sx: 47, sy: 38, r: 3},
	}

	// Draw segments back to front
	for i := len(segments) - 1; i >= 0; i-- {
		seg := segments[i]
		sx := x + seg.sx*p
		sy := y + seg.sy*p
		sr := seg.r * p

		vector.DrawFilledCircle(screen, sx, sy, sr, bodyColor, false)
		if i > 0 && i < len(segments)-1 {
			vector.DrawFilledRect(screen, sx-sr*0.7, sy-sr*0.4, sr*1.4, 1.5*p, stripeColor, false)
			vector.DrawFilledRect(screen, sx-sr*0.7, sy+sr*0.2, sr*1.4, 1.5*p, stripeColor, false)
		}
		vector.StrokeCircle(screen, sx, sy, sr, 1.2*p, outlineColor, false)
	}

	// Head details
	headX := x + segments[0].sx*p
	headY := y + segments[0].sy*p
	headR := segments[0].r * p
	vector.DrawFilledCircle(screen, headX, headY, headR, bodyColor, false)
	vector.StrokeCircle(screen, headX, headY, headR, 1*p, outlineColor, false)

	// Big cute eye
	eyeX := headX - 0.5*p
	eyeY := headY - 2*p
	vector.DrawFilledCircle(screen, eyeX, eyeY, 3*p, color.RGBA{255, 255, 255, 255}, false)
	vector.StrokeCircle(screen, eyeX, eyeY, 3*p, 0.8*p, outlineColor, false)
	vector.DrawFilledCircle(screen, eyeX+0.3*p, eyeY+0.3*p, 1.5*p, color.RGBA{20, 20, 20, 255}, false)
	vector.DrawFilledCircle(screen, eyeX+0.8*p, eyeY-0.8*p, 0.6*p, color.RGBA{255, 255, 255, 255}, false)

	// Cute small smile
	vector.DrawFilledRect(screen, headX-0.5*p, headY+3*p, 2*p, 1*p, outlineColor, false)

	// Stubby legs - front pair
	vector.DrawFilledRect(screen, x+7*p, y+41*p, 2*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+11*p, y+41*p, 2*p, 4*p, legColor, false)

	// Stubby legs - back pair
	vector.DrawFilledRect(screen, x+40*p, y+41*p, 2*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+41*p, 2*p, 4*p, legColor, false)
}

func (is *IntroScreen) drawWelcomeText(screen *ebiten.Image) {
	alpha := uint8(math.Min(is.welcomeAlpha, 255))

	// Background box for text
	boxWidth := float32(600)
	boxHeight := float32(80)
	boxX := float32(screenWidth)/2 - boxWidth/2
	boxY := float32(screenHeight)/2 + 100

	boxColor := color.RGBA{30, 50, 80, alpha}
	borderColor := color.RGBA{100, 150, 220, alpha}

	vector.DrawFilledRect(screen, boxX, boxY, boxWidth, boxHeight, boxColor, false)
	vector.StrokeRect(screen, boxX, boxY, boxWidth, boxHeight, 3, borderColor, false)

	// Text
	text1 := "Olive and Millie"
	text2 := "Welcome you to..."

	text1X := int(boxX + boxWidth/2 - float32(len(text1)*6)/2)
	text2X := int(boxX + boxWidth/2 - float32(len(text2)*6)/2)

	ebitenutil.DebugPrintAt(screen, text1, text1X, int(boxY+20))
	ebitenutil.DebugPrintAt(screen, text1, text1X+1, int(boxY+20))
	ebitenutil.DebugPrintAt(screen, text2, text2X, int(boxY+45))
	ebitenutil.DebugPrintAt(screen, text2, text2X+1, int(boxY+45))
}

func (is *IntroScreen) drawTitleText(screen *ebiten.Image) {
	alpha := uint8(math.Min(is.titleAlpha, 255))

	// Larger background box for title
	boxWidth := float32(650)
	boxHeight := float32(120)
	boxX := float32(screenWidth)/2 - boxWidth/2
	boxY := float32(screenHeight)/2 - 60

	boxColor := color.RGBA{30, 50, 80, alpha}
	borderColor := color.RGBA{255, 200, 100, alpha}

	vector.DrawFilledRect(screen, boxX, boxY, boxWidth, boxHeight, boxColor, false)
	vector.StrokeRect(screen, boxX, boxY, boxWidth, boxHeight, 4, borderColor, false)

	// Title text - larger effect with multiple prints
	title1 := "OLIVE & MILLIE'S"
	title2 := "GAME ROOM"

	title1X := int(boxX + boxWidth/2 - float32(len(title1)*6)/2)
	title2X := int(boxX + boxWidth/2 - float32(len(title2)*6)/2)

	// Draw title with "glow" effect
	ebitenutil.DebugPrintAt(screen, title1, title1X, int(boxY+30))
	ebitenutil.DebugPrintAt(screen, title1, title1X+1, int(boxY+30))
	ebitenutil.DebugPrintAt(screen, title1, title1X, int(boxY+31))
	ebitenutil.DebugPrintAt(screen, title2, title2X, int(boxY+65))
	ebitenutil.DebugPrintAt(screen, title2, title2X+1, int(boxY+65))
	ebitenutil.DebugPrintAt(screen, title2, title2X, int(boxY+66))
	ebitenutil.DebugPrintAt(screen, title2, title2X+1, int(boxY+66))
}

// Draw millipede on ground without frame - cute inchworm style
func (is *IntroScreen) drawGroundMillipede(screen *ebiten.Image) {
	x := float32(is.millipedeX)
	y := float32(is.millipedeY)
	p := float32(2.5)

	// Colors - cute sage green with yellow stripes
	bodyColor := color.RGBA{158, 194, 145, 255}
	stripeColor := color.RGBA{220, 200, 130, 255}
	outlineColor := color.RGBA{70, 90, 60, 255}
	legColor := color.RGBA{130, 170, 120, 255}

	// Segment positions - same arch shape as avatar
	type segment struct {
		sx, sy, r float32
	}
	segments := []segment{
		{sx: 5, sy: 38, r: 4},
		{sx: 9, sy: 38, r: 3.5},
		{sx: 13, sy: 38, r: 3.5},
		{sx: 16, sy: 35, r: 3.5},
		{sx: 19, sy: 28, r: 3.5},
		{sx: 21, sy: 21, r: 3.5},
		{sx: 23, sy: 15, r: 3.5},
		{sx: 26, sy: 12, r: 3.5},
		{sx: 29, sy: 15, r: 3.5},
		{sx: 31, sy: 21, r: 3.5},
		{sx: 33, sy: 28, r: 3.5},
		{sx: 36, sy: 35, r: 3.5},
		{sx: 39, sy: 38, r: 3.5},
		{sx: 43, sy: 38, r: 3.5},
		{sx: 47, sy: 38, r: 3},
	}

	// Draw segments back to front
	for i := len(segments) - 1; i >= 0; i-- {
		seg := segments[i]
		sx := x + seg.sx*p
		sy := y + seg.sy*p
		sr := seg.r * p

		vector.DrawFilledCircle(screen, sx, sy, sr, bodyColor, false)
		if i > 0 && i < len(segments)-1 {
			vector.DrawFilledRect(screen, sx-sr*0.7, sy-sr*0.4, sr*1.4, 1.5*p, stripeColor, false)
			vector.DrawFilledRect(screen, sx-sr*0.7, sy+sr*0.2, sr*1.4, 1.5*p, stripeColor, false)
		}
		vector.StrokeCircle(screen, sx, sy, sr, 1.2*p, outlineColor, false)
	}

	// Head details
	headX := x + segments[0].sx*p
	headY := y + segments[0].sy*p
	headR := segments[0].r * p
	vector.DrawFilledCircle(screen, headX, headY, headR, bodyColor, false)
	vector.StrokeCircle(screen, headX, headY, headR, 1*p, outlineColor, false)

	// Big cute eye
	eyeX := headX - 0.5*p
	eyeY := headY - 2*p
	vector.DrawFilledCircle(screen, eyeX, eyeY, 3*p, color.RGBA{255, 255, 255, 255}, false)
	vector.StrokeCircle(screen, eyeX, eyeY, 3*p, 0.8*p, outlineColor, false)
	vector.DrawFilledCircle(screen, eyeX+0.3*p, eyeY+0.3*p, 1.5*p, color.RGBA{20, 20, 20, 255}, false)
	vector.DrawFilledCircle(screen, eyeX+0.8*p, eyeY-0.8*p, 0.6*p, color.RGBA{255, 255, 255, 255}, false)

	// Cute small smile
	vector.DrawFilledRect(screen, headX-0.5*p, headY+3*p, 2*p, 1*p, outlineColor, false)

	// Stubby legs - front pair
	vector.DrawFilledRect(screen, x+7*p, y+41*p, 2*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+11*p, y+41*p, 2*p, 4*p, legColor, false)

	// Stubby legs - back pair
	vector.DrawFilledRect(screen, x+40*p, y+41*p, 2*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+41*p, 2*p, 4*p, legColor, false)
}

// Draw all avatars marching across the top (including all avatars, reversed order)
func (is *IntroScreen) drawMarchingAvatars(screen *ebiten.Image) {
	y := float32(120) // Near top, above the trees
	spacing := float32(120)          // More spacing for bigger avatars
	scale := float32(2.0)            // Bigger avatars

	numAvatars := int(AvatarNumTypes)

	// Draw each avatar in reverse order
	for idx := 0; idx < numAvatars; idx++ {
		// Reverse the order: last avatar first
		i := AvatarType(numAvatars - 1 - idx)

		x := float32(is.marchX) + float32(idx)*spacing

		// Only draw if on screen
		if x > -100 && x < float32(screenWidth)+100 {
			// Add a little bounce to the march
			bounceY := y + float32(math.Sin(float64(is.frameCount+idx*10)*0.15))*5

			DrawAvatar(screen, i, x, bounceY, scale)

			// Draw name below avatar
			name := GetAvatarName(i)
			nameX := int(x + 50*scale/2 - float32(len(name)*3))
			nameY := int(bounceY + 50*scale + 10)
			ebitenutil.DebugPrintAt(screen, name, nameX, nameY)
		}
	}
}

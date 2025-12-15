package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AvatarType int

const (
	AvatarHuman AvatarType = iota
	AvatarDog
	AvatarCat
	AvatarRabbit
	AvatarGiraffe
	AvatarOwl
	AvatarMillipede
	AvatarPuppy
	AvatarTiger
	AvatarChimpanzee
	AvatarPlatypus
	AvatarLynx
	AvatarGator
	AvatarOcelot
	AvatarHen
	AvatarMouse
	AvatarCaribou
	AvatarPitbull
	AvatarLabrador
	AvatarGoldenDoodle
	AvatarShepherd
	AvatarWhiteCat
	AvatarGreyCat
	AvatarRibbonPuddles
	AvatarNumTypes
)

var avatarNames = []string{
	"Human", "Teddy", "Kaycat", "Zach Rabbit", "Kiraffe", "Owlive", "Milliepede", "Sweet Puppy Paw", "Tygler", "Chimpancici", "Papapus", "Kaitlynx", "Reagator", "Ocelivia", "Hen-ry", "Tomouse", "Karabou", "Valkyrie", "Eleanor", "Stella", "Huckleberry", "Winston", "Baxter", "Ribbon & Puddles",
}

func GetAvatarName(avatarType AvatarType) string {
	if avatarType >= 0 && avatarType < AvatarNumTypes {
		return avatarNames[avatarType]
	}
	return "Unknown"
}

func DrawAvatar(screen *ebiten.Image, avatarType AvatarType, x, y, scale float32) {
	switch avatarType {
	case AvatarHuman:
		DrawPlayer1Avatar(screen, x, y, scale)
	case AvatarDog:
		DrawPlayer2Avatar(screen, x, y, scale)
	case AvatarCat:
		DrawCatAvatar(screen, x, y, scale)
	case AvatarRabbit:
		DrawRabbitAvatar(screen, x, y, scale)
	case AvatarGiraffe:
		DrawGiraffeAvatar(screen, x, y, scale)
	case AvatarOwl:
		DrawOwlAvatar(screen, x, y, scale)
	case AvatarMillipede:
		DrawMillipedeAvatar(screen, x, y, scale)
	case AvatarPuppy:
		DrawPuppyAvatar(screen, x, y, scale)
	case AvatarTiger:
		DrawTigerAvatar(screen, x, y, scale)
	case AvatarChimpanzee:
		DrawChimpanzeeAvatar(screen, x, y, scale)
	case AvatarPlatypus:
		DrawPlatypusAvatar(screen, x, y, scale)
	case AvatarLynx:
		DrawLynxAvatar(screen, x, y, scale)
	case AvatarGator:
		DrawGatorAvatar(screen, x, y, scale)
	case AvatarOcelot:
		DrawOcelotAvatar(screen, x, y, scale)
	case AvatarHen:
		DrawHenAvatar(screen, x, y, scale)
	case AvatarMouse:
		DrawMouseAvatar(screen, x, y, scale)
	case AvatarCaribou:
		DrawCaribouAvatar(screen, x, y, scale)
	case AvatarPitbull:
		DrawPitbullAvatar(screen, x, y, scale)
	case AvatarLabrador:
		DrawLabradorAvatar(screen, x, y, scale)
	case AvatarGoldenDoodle:
		DrawGoldenDoodleAvatar(screen, x, y, scale)
	case AvatarShepherd:
		DrawShepherdAvatar(screen, x, y, scale)
	case AvatarWhiteCat:
		DrawWhiteCatAvatar(screen, x, y, scale)
	case AvatarGreyCat:
		DrawGreyCatAvatar(screen, x, y, scale)
	case AvatarRibbonPuddles:
		DrawRibbonPuddlesAvatar(screen, x, y, scale)
	}
}

// Draw Cat avatar
func DrawCatAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body/chest
	chestColor := color.RGBA{160, 160, 160, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, chestColor, false)

	// Neck
	neckColor := color.RGBA{150, 150, 150, 255}
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, neckColor, false)

	// Head shape (more triangular for cat)
	headColor := color.RGBA{140, 140, 140, 255}
	vector.DrawFilledRect(screen, x+14*p, y+20*p, 22*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+16*p, 18*p, 6*p, headColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+14*p, 14*p, 4*p, headColor, false)

	// Ears (triangular)
	earColor := color.RGBA{120, 120, 120, 255}
	// Left ear
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+10*p+float32(i)*p/2, y+8*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	// Right ear
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+30*p-float32(i)*p/2, y+8*p+float32(i)*p, w*p, 2*p, earColor, false)
	}

	// Inner ears
	innerEarColor := color.RGBA{255, 180, 200, 255}
	vector.DrawFilledRect(screen, x+13*p, y+12*p, 4*p, 4*p, innerEarColor, false)
	vector.DrawFilledRect(screen, x+33*p, y+12*p, 4*p, 4*p, innerEarColor, false)

	// Eyes (cat-like)
	eyeColor := color.RGBA{100, 200, 100, 255}
	vector.DrawFilledRect(screen, x+18*p, y+22*p, 5*p, 4*p, eyeColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+22*p, 5*p, 4*p, eyeColor, false)

	// Pupils (vertical slits)
	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+20*p, y+23*p, 1*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+23*p, 1*p, 3*p, pupilColor, false)

	// Nose (pink triangle)
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+28*p, 4*p, 2*p, noseColor, false)
	vector.DrawFilledRect(screen, x+24*p, y+27*p, 2*p, 1*p, noseColor, false)

	// Mouth
	mouthColor := color.RGBA{80, 60, 50, 255}
	vector.DrawFilledRect(screen, x+25*p, y+30*p, 1*p, 3*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+23*p, y+32*p, 5*p, 1*p, mouthColor, false)

	// Whiskers
	whiskerColor := color.RGBA{60, 60, 60, 255}
	vector.DrawFilledRect(screen, x+8*p, y+26*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+8*p, y+28*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+26*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+28*p, 8*p, 1*p, whiskerColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{200, 100, 200, 255}, false)
}

// Draw Rabbit avatar
func DrawRabbitAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{230, 220, 210, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+18*p, y+38*p, 14*p, 6*p, bodyColor, false)

	// Head (round)
	headColor := color.RGBA{240, 230, 220, 255}
	vector.DrawFilledRect(screen, x+14*p, y+22*p, 22*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+20*p, 18*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+36*p, 18*p, 4*p, headColor, false)

	// Long ears
	earColor := color.RGBA{230, 220, 210, 255}
	// Left ear
	vector.DrawFilledRect(screen, x+14*p, y+4*p, 8*p, 18*p, earColor, false)
	vector.DrawFilledRect(screen, x+15*p, y+2*p, 6*p, 4*p, earColor, false)
	// Right ear
	vector.DrawFilledRect(screen, x+28*p, y+4*p, 8*p, 18*p, earColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+2*p, 6*p, 4*p, earColor, false)

	// Inner ears
	innerEarColor := color.RGBA{255, 200, 220, 255}
	vector.DrawFilledRect(screen, x+16*p, y+8*p, 4*p, 10*p, innerEarColor, false)
	vector.DrawFilledRect(screen, x+30*p, y+8*p, 4*p, 10*p, innerEarColor, false)

	// Eyes (big and round)
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledRect(screen, x+17*p, y+24*p, 7*p, 6*p, eyeWhite, false)
	vector.DrawFilledRect(screen, x+26*p, y+24*p, 7*p, 6*p, eyeWhite, false)

	eyeColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledRect(screen, x+18*p, y+25*p, 5*p, 5*p, eyeColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+25*p, 5*p, 5*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+19*p, y+26*p, 3*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+28*p, y+26*p, 3*p, 3*p, pupilColor, false)

	// Light reflection
	vector.DrawFilledRect(screen, x+20*p, y+26*p, 2*p, 2*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, x+29*p, y+26*p, 2*p, 2*p, color.RGBA{255, 255, 255, 255}, false)

	// Nose (pink)
	noseColor := color.RGBA{255, 180, 200, 255}
	vector.DrawFilledRect(screen, x+23*p, y+31*p, 4*p, 3*p, noseColor, false)

	// Mouth (like Y shape)
	mouthColor := color.RGBA{100, 80, 70, 255}
	vector.DrawFilledRect(screen, x+25*p, y+33*p, 1*p, 2*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+23*p, y+34*p, 2*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+26*p, y+34*p, 2*p, 1*p, mouthColor, false)

	// Cheek fluff
	fluffColor := color.RGBA{250, 240, 230, 255}
	vector.DrawFilledRect(screen, x+10*p, y+26*p, 4*p, 8*p, fluffColor, false)
	vector.DrawFilledRect(screen, x+36*p, y+26*p, 4*p, 8*p, fluffColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 200, 150, 255}, false)
}

// Draw Giraffe avatar
func DrawGiraffeAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Very long neck
	neckColor := color.RGBA{255, 200, 120, 255}
	vector.DrawFilledRect(screen, x+20*p, y+25*p, 10*p, 20*p, neckColor, false)

	// Body (at bottom)
	bodyColor := color.RGBA{255, 200, 120, 255}
	vector.DrawFilledRect(screen, x+15*p, y+42*p, 20*p, 8*p, bodyColor, false)

	// Head (small compared to neck)
	headColor := color.RGBA{255, 180, 100, 255}
	vector.DrawFilledRect(screen, x+16*p, y+10*p, 18*p, 15*p, headColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+8*p, 14*p, 4*p, headColor, false)

	// Spots pattern
	spotColor := color.RGBA{139, 90, 43, 255}
	// Neck spots
	vector.DrawFilledCircle(screen, x+23*p, y+30*p, 3*p, spotColor, false)
	vector.DrawFilledCircle(screen, x+26*p, y+35*p, 3*p, spotColor, false)
	vector.DrawFilledCircle(screen, x+22*p, y+40*p, 3*p, spotColor, false)
	// Head spots
	vector.DrawFilledCircle(screen, x+20*p, y+14*p, 2*p, spotColor, false)
	vector.DrawFilledCircle(screen, x+28*p, y+16*p, 2*p, spotColor, false)

	// Horns (ossicones)
	hornColor := color.RGBA{139, 90, 43, 255}
	vector.DrawFilledRect(screen, x+20*p, y+6*p, 3*p, 5*p, hornColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+6*p, 3*p, 5*p, hornColor, false)
	// Horn tips
	vector.DrawFilledCircle(screen, x+21.5*p, y+6*p, 2*p, hornColor, false)
	vector.DrawFilledCircle(screen, x+28.5*p, y+6*p, 2*p, hornColor, false)

	// Ears
	vector.DrawFilledRect(screen, x+14*p, y+12*p, 4*p, 6*p, headColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+12*p, 4*p, 6*p, headColor, false)

	// Eyes
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+21*p, y+16*p, 3*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+29*p, y+16*p, 3*p, eyeWhite, false)

	eyeColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+21*p, y+16*p, 2*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+29*p, y+16*p, 2*p, eyeColor, false)

	// Nose
	noseColor := color.RGBA{100, 80, 60, 255}
	vector.DrawFilledRect(screen, x+23*p, y+20*p, 4*p, 3*p, noseColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 200, 120, 255}, false)
}

// Draw Millipede avatar - cute green inchworm style
func DrawMillipedeAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Colors - cute sage green with yellow stripes
	bodyColor := color.RGBA{158, 194, 145, 255}   // Sage green
	stripeColor := color.RGBA{220, 200, 130, 255} // Yellow/cream stripe
	outlineColor := color.RGBA{70, 90, 60, 255}   // Dark green outline
	legColor := color.RGBA{130, 170, 120, 255}    // Green for legs

	// Segment positions forming a smooth arch curve (like the reference image)
	// Head on left low, curves up and over, tail on right low
	type segment struct {
		x, y, r float32
	}

	segments := []segment{
		{x: 5, y: 38, r: 4},    // Head
		{x: 9, y: 38, r: 3.5},  // Straight section
		{x: 13, y: 38, r: 3.5}, // Straight section
		{x: 16, y: 35, r: 3.5}, // Start curving up
		{x: 19, y: 28, r: 3.5}, // Going up
		{x: 21, y: 21, r: 3.5}, // Going up more
		{x: 23, y: 15, r: 3.5}, // Top left
		{x: 26, y: 12, r: 3.5}, // Top center
		{x: 29, y: 15, r: 3.5}, // Top right
		{x: 31, y: 21, r: 3.5}, // Coming down
		{x: 33, y: 28, r: 3.5}, // Coming down more
		{x: 36, y: 35, r: 3.5}, // End curve
		{x: 39, y: 38, r: 3.5}, // Straight section
		{x: 43, y: 38, r: 3.5}, // Straight section
		{x: 47, y: 38, r: 3},   // Tail
	}

	// Draw segments back to front so head overlaps properly
	for i := len(segments) - 1; i >= 0; i-- {
		seg := segments[i]
		sx := x + seg.x*p
		sy := y + seg.y*p
		sr := seg.r * p

		// Main body segment
		vector.DrawFilledCircle(screen, sx, sy, sr, bodyColor, false)

		// Yellow stripe bands (two lines per segment, skip head and tail)
		if i > 0 && i < len(segments)-1 {
			vector.DrawFilledRect(screen, sx-sr*0.7, sy-sr*0.4, sr*1.4, 1.5*p, stripeColor, false)
			vector.DrawFilledRect(screen, sx-sr*0.7, sy+sr*0.2, sr*1.4, 1.5*p, stripeColor, false)
		}

		// Outline
		vector.StrokeCircle(screen, sx, sy, sr, 1.2*p, outlineColor, false)
	}

	// Head details
	headX := x + segments[0].x*p
	headY := y + segments[0].y*p
	headR := segments[0].r * p

	// Redraw head on top
	vector.DrawFilledCircle(screen, headX, headY, headR, bodyColor, false)
	vector.StrokeCircle(screen, headX, headY, headR, 1*p, outlineColor, false)

	// Big cute eye
	eyeX := headX - 0.5*p
	eyeY := headY - 2*p
	vector.DrawFilledCircle(screen, eyeX, eyeY, 3*p, color.RGBA{255, 255, 255, 255}, false)
	vector.StrokeCircle(screen, eyeX, eyeY, 3*p, 0.8*p, outlineColor, false)

	// Pupil
	vector.DrawFilledCircle(screen, eyeX+0.3*p, eyeY+0.3*p, 1.5*p, color.RGBA{20, 20, 20, 255}, false)

	// Eye shine
	vector.DrawFilledCircle(screen, eyeX+0.8*p, eyeY-0.8*p, 0.6*p, color.RGBA{255, 255, 255, 255}, false)

	// Cute small smile
	vector.DrawFilledRect(screen, headX-0.5*p, headY+3*p, 2*p, 1*p, outlineColor, false)

	// Stubby legs - front pair (under head/straight section)
	vector.DrawFilledRect(screen, x+7*p, y+41*p, 2*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+11*p, y+41*p, 2*p, 4*p, legColor, false)

	// Stubby legs - back pair (under tail/straight section)
	vector.DrawFilledRect(screen, x+40*p, y+41*p, 2*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+41*p, 2*p, 4*p, legColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, bodyColor, false)
}

// Draw Owl avatar
func DrawOwlAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{139, 90, 60, 255}
	vector.DrawFilledRect(screen, x+12*p, y+38*p, 26*p, 12*p, bodyColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+36*p, 22*p, 4*p, bodyColor, false)

	// Belly pattern
	bellyColor := color.RGBA{200, 180, 160, 255}
	vector.DrawFilledRect(screen, x+18*p, y+40*p, 14*p, 10*p, bellyColor, false)

	// Head (round and wide)
	headColor := color.RGBA{120, 80, 50, 255}
	vector.DrawFilledRect(screen, x+8*p, y+18*p, 34*p, 20*p, headColor, false)
	vector.DrawFilledRect(screen, x+10*p, y+16*p, 30*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+14*p, 26*p, 4*p, headColor, false)

	// Ear tufts
	// Left tuft
	for i := 0; i < 6; i++ {
		h := float32(6 - i)
		vector.DrawFilledRect(screen, x+10*p+float32(i)*p, y+10*p+float32(i)*p, 3*p, h*p, headColor, false)
	}
	// Right tuft
	for i := 0; i < 6; i++ {
		h := float32(6 - i)
		vector.DrawFilledRect(screen, x+37*p-float32(i)*p, y+10*p+float32(i)*p, 3*p, h*p, headColor, false)
	}

	// Face disc
	faceColor := color.RGBA{180, 160, 140, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 16*p, faceColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, faceColor, false)

	// Large eyes (owl characteristic)
	// Eye backgrounds
	eyeBgColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 7*p, eyeBgColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 7*p, eyeBgColor, false)

	// Yellow eyes
	eyeColor := color.RGBA{255, 200, 80, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 5*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 5*p, eyeColor, false)

	// Pupils
	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 2*p, pupilColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 2*p, pupilColor, false)

	// Beak
	beakColor := color.RGBA{255, 180, 100, 255}
	// Triangle beak
	for i := 0; i < 5; i++ {
		w := float32(5 - i)
		vector.DrawFilledRect(screen, x+25*p-w*p/2, y+30*p+float32(i)*p, w*p, 1*p, beakColor, false)
	}

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{120, 80, 50, 255}, false)
}

// Draw Puppy avatar
func DrawPuppyAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{220, 180, 140, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, bodyColor, false)

	// Head (round and cute)
	headColor := color.RGBA{230, 190, 150, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 18*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+36*p, 22*p, 3*p, headColor, false)

	// Floppy ears
	earColor := color.RGBA{200, 160, 120, 255}
	// Left ear
	vector.DrawFilledRect(screen, x+8*p, y+20*p, 8*p, 14*p, earColor, false)
	vector.DrawFilledRect(screen, x+6*p, y+24*p, 4*p, 10*p, earColor, false)
	// Right ear
	vector.DrawFilledRect(screen, x+34*p, y+20*p, 8*p, 14*p, earColor, false)
	vector.DrawFilledRect(screen, x+40*p, y+24*p, 4*p, 10*p, earColor, false)

	// Big puppy eyes
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+26*p, 5*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+31*p, y+26*p, 5*p, eyeWhite, false)

	eyeColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+26*p, 4*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+26*p, 4*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+27*p, 2*p, pupilColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+27*p, 2*p, pupilColor, false)

	// Sparkle in eyes
	vector.DrawFilledCircle(screen, x+20*p, y+25*p, 1.5*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(screen, x+32*p, y+25*p, 1.5*p, color.RGBA{255, 255, 255, 255}, false)

	// Nose (pink)
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+31*p, 4*p, 3*p, noseColor, false)

	// Mouth (smiling)
	mouthColor := color.RGBA{100, 80, 70, 255}
	vector.DrawFilledRect(screen, x+25*p, y+33*p, 1*p, 2*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+22*p, y+34*p, 3*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+26*p, y+34*p, 3*p, 1*p, mouthColor, false)

	// Tail (visible)
	tailColor := color.RGBA{210, 170, 130, 255}
	vector.DrawFilledRect(screen, x+38*p, y+40*p, 8*p, 4*p, tailColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+38*p, 4*p, 6*p, tailColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 180, 200, 255}, false)
}

// Draw Tiger avatar
func DrawTigerAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{255, 140, 50, 255}
	vector.DrawFilledRect(screen, x+8*p, y+42*p, 34*p, 8*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+16*p, y+36*p, 18*p, 8*p, bodyColor, false)

	// Head (wide for tiger)
	headColor := color.RGBA{255, 150, 60, 255}
	vector.DrawFilledRect(screen, x+10*p, y+20*p, 30*p, 18*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+18*p, 26*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+36*p, 22*p, 3*p, headColor, false)

	// Ears (rounded triangular)
	earColor := color.RGBA{255, 130, 40, 255}
	vector.DrawFilledCircle(screen, x+14*p, y+16*p, 5*p, earColor, false)
	vector.DrawFilledCircle(screen, x+36*p, y+16*p, 5*p, earColor, false)

	// Stripes
	stripeColor := color.RGBA{40, 40, 40, 255}
	// Head stripes
	vector.DrawFilledRect(screen, x+10*p, y+22*p, 2*p, 6*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+24*p, 2*p, 5*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+24*p, 2*p, 5*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+38*p, y+22*p, 2*p, 6*p, stripeColor, false)
	// Body stripes
	vector.DrawFilledRect(screen, x+12*p, y+43*p, 3*p, 6*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+20*p, y+44*p, 3*p, 5*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+28*p, y+43*p, 3*p, 6*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+36*p, y+44*p, 3*p, 5*p, stripeColor, false)

	// White muzzle area
	muzzleColor := color.RGBA{255, 240, 230, 255}
	vector.DrawFilledRect(screen, x+18*p, y+28*p, 14*p, 10*p, muzzleColor, false)

	// Eyes
	eyeColor := color.RGBA{100, 200, 100, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+24*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+24*p, 3*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+17*p, y+23*p, 2*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+31*p, y+23*p, 2*p, 3*p, pupilColor, false)

	// Nose (pink)
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+31*p, 4*p, 3*p, noseColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 140, 50, 255}, false)
}

// Draw Chimpanzee avatar
func DrawChimpanzeeAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{80, 60, 50, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, bodyColor, false)

	// Arms (longer)
	vector.DrawFilledRect(screen, x+6*p, y+38*p, 10*p, 12*p, bodyColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+38*p, 10*p, 12*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, bodyColor, false)

	// Head
	headColor := color.RGBA{100, 80, 70, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 18*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, headColor, false)

	// Face (lighter)
	faceColor := color.RGBA{180, 160, 140, 255}
	vector.DrawFilledRect(screen, x+14*p, y+22*p, 22*p, 14*p, faceColor, false)

	// Large ears
	earColor := color.RGBA{180, 160, 140, 255}
	vector.DrawFilledCircle(screen, x+10*p, y+26*p, 6*p, earColor, false)
	vector.DrawFilledCircle(screen, x+40*p, y+26*p, 6*p, earColor, false)

	// Eyes
	eyeColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+26*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+26*p, 3*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+26*p, 2*p, pupilColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+26*p, 2*p, pupilColor, false)

	// Nose
	noseColor := color.RGBA{120, 100, 80, 255}
	vector.DrawFilledRect(screen, x+23*p, y+30*p, 4*p, 2*p, noseColor, false)

	// Mouth
	mouthColor := color.RGBA{100, 80, 60, 255}
	vector.DrawFilledRect(screen, x+20*p, y+33*p, 10*p, 2*p, mouthColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{100, 80, 70, 255}, false)
}

// Draw Platypus avatar
func DrawPlatypusAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body (brown)
	bodyColor := color.RGBA{139, 90, 60, 255}
	vector.DrawFilledRect(screen, x+8*p, y+38*p, 34*p, 12*p, bodyColor, false)

	// Tail (beaver-like)
	tailColor := color.RGBA{100, 70, 50, 255}
	vector.DrawFilledRect(screen, x+38*p, y+40*p, 10*p, 8*p, tailColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+42*p, 4*p, 4*p, tailColor, false)

	// Head
	headColor := color.RGBA{160, 110, 70, 255}
	vector.DrawFilledRect(screen, x+10*p, y+22*p, 20*p, 16*p, headColor, false)

	// Bill (distinctive platypus feature)
	billColor := color.RGBA{200, 140, 80, 255}
	vector.DrawFilledRect(screen, x+4*p, y+26*p, 16*p, 8*p, billColor, false)
	vector.DrawFilledRect(screen, x+2*p, y+28*p, 18*p, 4*p, billColor, false)

	// Nostrils on bill
	nostrilColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+6*p, y+29*p, 1*p, nostrilColor, false)
	vector.DrawFilledCircle(screen, x+10*p, y+29*p, 1*p, nostrilColor, false)

	// Eyes
	eyeColor := color.RGBA{40, 30, 20, 255}
	vector.DrawFilledCircle(screen, x+16*p, y+26*p, 2*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+24*p, y+26*p, 2*p, eyeColor, false)

	// Feet (webbed)
	feetColor := color.RGBA{200, 140, 80, 255}
	vector.DrawFilledRect(screen, x+10*p, y+46*p, 8*p, 4*p, feetColor, false)
	vector.DrawFilledRect(screen, x+22*p, y+46*p, 8*p, 4*p, feetColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{139, 90, 60, 255}, false)
}

// Draw Lynx avatar
func DrawLynxAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{200, 180, 160, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, bodyColor, false)

	// Head
	headColor := color.RGBA{210, 190, 170, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, headColor, false)

	// Distinctive ear tufts
	earColor := color.RGBA{190, 170, 150, 255}
	// Left ear with tuft
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+10*p+float32(i)*p/2, y+10*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	// Right ear with tuft
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+32*p-float32(i)*p/2, y+10*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	// Black tufts on ears
	tuftColor := color.RGBA{40, 40, 40, 255}
	vector.DrawFilledRect(screen, x+12*p, y+8*p, 3*p, 6*p, tuftColor, false)
	vector.DrawFilledRect(screen, x+35*p, y+8*p, 3*p, 6*p, tuftColor, false)

	// Face pattern
	facePattern := color.RGBA{230, 210, 190, 255}
	vector.DrawFilledRect(screen, x+16*p, y+24*p, 18*p, 12*p, facePattern, false)

	// Spots
	spotColor := color.RGBA{100, 80, 60, 255}
	vector.DrawFilledCircle(screen, x+14*p, y+30*p, 2*p, spotColor, false)
	vector.DrawFilledCircle(screen, x+36*p, y+30*p, 2*p, spotColor, false)
	vector.DrawFilledCircle(screen, x+20*p, y+34*p, 1.5*p, spotColor, false)
	vector.DrawFilledCircle(screen, x+30*p, y+34*p, 1.5*p, spotColor, false)

	// Eyes (cat-like)
	eyeColor := color.RGBA{100, 200, 100, 255}
	vector.DrawFilledRect(screen, x+17*p, y+24*p, 6*p, 4*p, eyeColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+24*p, 6*p, 4*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+19*p, y+25*p, 2*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+25*p, 2*p, 3*p, pupilColor, false)

	// Nose
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+29*p, 4*p, 2*p, noseColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{200, 180, 160, 255}, false)
}

// Draw Gator avatar
func DrawGatorAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{50, 100, 50, 255}
	vector.DrawFilledRect(screen, x+8*p, y+40*p, 34*p, 10*p, bodyColor, false)

	// Tail
	tailColor := color.RGBA{40, 80, 40, 255}
	vector.DrawFilledRect(screen, x+38*p, y+42*p, 10*p, 6*p, tailColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+44*p, 4*p, 2*p, tailColor, false)

	// Head (long snout)
	headColor := color.RGBA{60, 120, 60, 255}
	vector.DrawFilledRect(screen, x+14*p, y+24*p, 22*p, 16*p, headColor, false)

	// Snout
	snoutColor := color.RGBA{50, 100, 50, 255}
	vector.DrawFilledRect(screen, x+4*p, y+28*p, 16*p, 8*p, snoutColor, false)
	vector.DrawFilledRect(screen, x+2*p, y+30*p, 18*p, 4*p, snoutColor, false)

	// Teeth
	teethColor := color.RGBA{255, 255, 255, 255}
	for i := 0; i < 4; i++ {
		vector.DrawFilledRect(screen, x+5*p+float32(i)*4*p, y+33*p, 2*p, 3*p, teethColor, false)
	}

	// Nostrils
	nostrilColor := color.RGBA{30, 60, 30, 255}
	vector.DrawFilledCircle(screen, x+6*p, y+30*p, 1.5*p, nostrilColor, false)
	vector.DrawFilledCircle(screen, x+10*p, y+30*p, 1.5*p, nostrilColor, false)

	// Eyes (on top of head)
	eyeColor := color.RGBA{255, 200, 100, 255}
	vector.DrawFilledCircle(screen, x+20*p, y+26*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+30*p, y+26*p, 3*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+19*p, y+25*p, 2*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+25*p, 2*p, 3*p, pupilColor, false)

	// Scales pattern
	scaleColor := color.RGBA{40, 80, 40, 255}
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			vector.DrawFilledCircle(screen, x+14*p+float32(i)*6*p, y+34*p+float32(j)*4*p, 2*p, scaleColor, false)
		}
	}

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{50, 100, 50, 255}, false)
}

// Draw Ocelot avatar
func DrawOcelotAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{200, 160, 100, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, bodyColor, false)

	// Head
	headColor := color.RGBA{210, 170, 110, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, headColor, false)

	// Ears
	earColor := color.RGBA{190, 150, 90, 255}
	for i := 0; i < 6; i++ {
		w := float32(6 - i)
		vector.DrawFilledRect(screen, x+12*p+float32(i)*p/2, y+12*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	for i := 0; i < 6; i++ {
		w := float32(6 - i)
		vector.DrawFilledRect(screen, x+32*p-float32(i)*p/2, y+12*p+float32(i)*p, w*p, 2*p, earColor, false)
	}

	// Rosette pattern (distinctive ocelot spots)
	spotColor := color.RGBA{80, 60, 40, 255}
	// Head spots
	vector.DrawFilledCircle(screen, x+16*p, y+22*p, 2*p, spotColor, false)
	vector.StrokeCircle(screen, x+16*p, y+22*p, 3*p, 1, spotColor, false)
	vector.DrawFilledCircle(screen, x+34*p, y+22*p, 2*p, spotColor, false)
	vector.StrokeCircle(screen, x+34*p, y+22*p, 3*p, 1, spotColor, false)
	// Body rosettes
	vector.StrokeCircle(screen, x+15*p, y+44*p, 4*p, 1.5, spotColor, false)
	vector.DrawFilledCircle(screen, x+15*p, y+44*p, 1*p, spotColor, false)
	vector.StrokeCircle(screen, x+25*p, y+45*p, 4*p, 1.5, spotColor, false)
	vector.DrawFilledCircle(screen, x+25*p, y+45*p, 1*p, spotColor, false)
	vector.StrokeCircle(screen, x+35*p, y+44*p, 4*p, 1.5, spotColor, false)
	vector.DrawFilledCircle(screen, x+35*p, y+44*p, 1*p, spotColor, false)

	// White muzzle area
	muzzleColor := color.RGBA{240, 220, 200, 255}
	vector.DrawFilledRect(screen, x+18*p, y+28*p, 14*p, 8*p, muzzleColor, false)

	// Eyes
	eyeColor := color.RGBA{180, 150, 80, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+25*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+25*p, 3*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+17*p, y+24*p, 2*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+31*p, y+24*p, 2*p, 3*p, pupilColor, false)

	// Nose
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+30*p, 4*p, 2*p, noseColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{200, 160, 100, 255}, false)
}

// Draw Hen avatar
func DrawHenAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{220, 180, 140, 255}
	vector.DrawFilledCircle(screen, x+25*p, y+38*p, 12*p, bodyColor, false)

	// Neck
	neckColor := color.RGBA{230, 190, 150, 255}
	vector.DrawFilledRect(screen, x+20*p, y+30*p, 10*p, 10*p, neckColor, false)

	// Head
	headColor := color.RGBA{240, 200, 160, 255}
	vector.DrawFilledCircle(screen, x+25*p, y+20*p, 8*p, headColor, false)

	// Comb (red crest)
	combColor := color.RGBA{220, 50, 50, 255}
	vector.DrawFilledCircle(screen, x+23*p, y+12*p, 3*p, combColor, false)
	vector.DrawFilledCircle(screen, x+25*p, y+10*p, 4*p, combColor, false)
	vector.DrawFilledCircle(screen, x+27*p, y+12*p, 3*p, combColor, false)

	// Wattle
	wattleColor := color.RGBA{200, 40, 40, 255}
	vector.DrawFilledCircle(screen, x+22*p, y+26*p, 3*p, wattleColor, false)

	// Beak
	beakColor := color.RGBA{255, 180, 100, 255}
	// Triangle beak
	for i := 0; i < 4; i++ {
		w := float32(4 - i)
		vector.DrawFilledRect(screen, x+17*p, y+20*p+float32(i)*p, w*p, 1*p, beakColor, false)
	}

	// Eye
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+28*p, y+20*p, 3*p, eyeWhite, false)
	eyeColor := color.RGBA{40, 40, 40, 255}
	vector.DrawFilledCircle(screen, x+28*p, y+20*p, 2*p, eyeColor, false)

	// Wing
	wingColor := color.RGBA{200, 160, 120, 255}
	vector.DrawFilledRect(screen, x+30*p, y+32*p, 10*p, 12*p, wingColor, false)
	// Wing feathers
	featherColor := color.RGBA{180, 140, 100, 255}
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(screen, x+32*p, y+34*p+float32(i)*3*p, 6*p, 2*p, featherColor, false)
	}

	// Tail feathers
	tailColor := color.RGBA{180, 140, 100, 255}
	vector.DrawFilledRect(screen, x+35*p, y+36*p, 8*p, 3*p, tailColor, false)
	vector.DrawFilledRect(screen, x+36*p, y+34*p, 7*p, 3*p, tailColor, false)
	vector.DrawFilledRect(screen, x+37*p, y+38*p, 6*p, 3*p, tailColor, false)

	// Feet
	feetColor := color.RGBA{255, 180, 100, 255}
	vector.DrawFilledRect(screen, x+20*p, y+46*p, 3*p, 4*p, feetColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+46*p, 3*p, 4*p, feetColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{220, 180, 140, 255}, false)
}

// Draw Mouse avatar
func DrawMouseAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body (small and round)
	bodyColor := color.RGBA{160, 140, 120, 255}
	vector.DrawFilledRect(screen, x+15*p, y+40*p, 20*p, 10*p, bodyColor, false)

	// Neck
	vector.DrawFilledRect(screen, x+20*p, y+36*p, 10*p, 6*p, bodyColor, false)

	// Head (small and cute)
	headColor := color.RGBA{170, 150, 130, 255}
	vector.DrawFilledRect(screen, x+16*p, y+22*p, 18*p, 14*p, headColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+20*p, 14*p, 4*p, headColor, false)

	// Large round ears (distinctive mouse feature)
	earColor := color.RGBA{150, 130, 110, 255}
	vector.DrawFilledCircle(screen, x+16*p, y+18*p, 8*p, earColor, false)
	vector.DrawFilledCircle(screen, x+34*p, y+18*p, 8*p, earColor, false)

	// Inner ears
	innerEarColor := color.RGBA{255, 200, 200, 255}
	vector.DrawFilledCircle(screen, x+16*p, y+18*p, 5*p, innerEarColor, false)
	vector.DrawFilledCircle(screen, x+34*p, y+18*p, 5*p, innerEarColor, false)

	// Eyes (big and round)
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+20*p, y+26*p, 4*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+30*p, y+26*p, 4*p, eyeWhite, false)

	eyeColor := color.RGBA{40, 40, 40, 255}
	vector.DrawFilledCircle(screen, x+20*p, y+26*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+30*p, y+26*p, 3*p, eyeColor, false)

	// Light reflection
	vector.DrawFilledCircle(screen, x+21*p, y+25*p, 1.5*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(screen, x+31*p, y+25*p, 1.5*p, color.RGBA{255, 255, 255, 255}, false)

	// Nose (small pink)
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledCircle(screen, x+25*p, y+32*p, 2*p, noseColor, false)

	// Whiskers (long and thin)
	whiskerColor := color.RGBA{80, 70, 60, 255}
	vector.DrawFilledRect(screen, x+8*p, y+30*p, 10*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+8*p, y+32*p, 10*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+30*p, 10*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+32*p, 10*p, 1*p, whiskerColor, false)

	// Long thin tail
	tailColor := color.RGBA{140, 120, 100, 255}
	vector.DrawFilledRect(screen, x+33*p, y+42*p, 14*p, 2*p, tailColor, false)
	vector.DrawFilledRect(screen, x+44*p, y+40*p, 2*p, 6*p, tailColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{160, 140, 120, 255}, false)
}

// Draw Caribou avatar
func DrawCaribouAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body
	bodyColor := color.RGBA{139, 115, 85, 255}
	vector.DrawFilledRect(screen, x+8*p, y+38*p, 34*p, 12*p, bodyColor, false)

	// Neck
	neckColor := color.RGBA{150, 125, 95, 255}
	vector.DrawFilledRect(screen, x+18*p, y+28*p, 14*p, 12*p, neckColor, false)

	// Head
	headColor := color.RGBA{160, 135, 105, 255}
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 12*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+16*p, 18*p, 4*p, headColor, false)

	// Snout
	snoutColor := color.RGBA{180, 155, 125, 255}
	vector.DrawFilledRect(screen, x+18*p, y+24*p, 14*p, 8*p, snoutColor, false)

	// Large branching antlers (distinctive caribou feature)
	antlerColor := color.RGBA{220, 200, 180, 255}

	// Left antler main branch
	vector.DrawFilledRect(screen, x+16*p, y+8*p, 3*p, 10*p, antlerColor, false)
	// Left antler branches
	vector.DrawFilledRect(screen, x+14*p, y+10*p, 6*p, 2*p, antlerColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+8*p, 4*p, 2*p, antlerColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+6*p, 3*p, 4*p, antlerColor, false)

	// Right antler main branch
	vector.DrawFilledRect(screen, x+31*p, y+8*p, 3*p, 10*p, antlerColor, false)
	// Right antler branches
	vector.DrawFilledRect(screen, x+30*p, y+10*p, 6*p, 2*p, antlerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+8*p, 4*p, 2*p, antlerColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+6*p, 3*p, 4*p, antlerColor, false)

	// Ears
	earColor := color.RGBA{140, 120, 90, 255}
	vector.DrawFilledRect(screen, x+14*p, y+14*p, 5*p, 6*p, earColor, false)
	vector.DrawFilledRect(screen, x+31*p, y+14*p, 5*p, 6*p, earColor, false)

	// Eyes
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+20*p, y+22*p, 3*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+30*p, y+22*p, 3*p, eyeWhite, false)

	eyeColor := color.RGBA{60, 50, 40, 255}
	vector.DrawFilledCircle(screen, x+20*p, y+22*p, 2*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+30*p, y+22*p, 2*p, eyeColor, false)

	// Nose
	noseColor := color.RGBA{100, 80, 60, 255}
	vector.DrawFilledRect(screen, x+23*p, y+28*p, 4*p, 3*p, noseColor, false)

	// White chest fur
	chestColor := color.RGBA{240, 230, 220, 255}
	vector.DrawFilledRect(screen, x+20*p, y+32*p, 10*p, 8*p, chestColor, false)

	// Legs/hooves
	legColor := color.RGBA{120, 100, 70, 255}
	vector.DrawFilledRect(screen, x+12*p, y+46*p, 4*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+20*p, y+46*p, 4*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+26*p, y+46*p, 4*p, 4*p, legColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+46*p, 4*p, 4*p, legColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{139, 115, 85, 255}, false)
}

// Draw Pitbull avatar (Valkyrie)
func DrawPitbullAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Muscular body (lighter for visibility)
	bodyColor := color.RGBA{70, 70, 75, 255}
	vector.DrawFilledRect(screen, x+8*p, y+42*p, 34*p, 8*p, bodyColor, false)

	// Chest (muscular)
	vector.DrawFilledRect(screen, x+12*p, y+38*p, 26*p, 6*p, bodyColor, false)

	// Thick neck
	vector.DrawFilledRect(screen, x+16*p, y+34*p, 18*p, 6*p, bodyColor, false)

	// Head (broad and strong)
	headColor := color.RGBA{80, 80, 85, 255}
	vector.DrawFilledRect(screen, x+12*p, y+20*p, 26*p, 14*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+18*p, 22*p, 4*p, headColor, false)

	// Strong jaw/muzzle
	muzzleColor := color.RGBA{90, 90, 95, 255}
	vector.DrawFilledRect(screen, x+16*p, y+28*p, 18*p, 8*p, muzzleColor, false)

	// Short ears (cropped look)
	earColor := color.RGBA{60, 60, 65, 255}
	vector.DrawFilledRect(screen, x+14*p, y+16*p, 6*p, 6*p, earColor, false)
	vector.DrawFilledRect(screen, x+30*p, y+16*p, 6*p, 6*p, earColor, false)

	// Eyes (strong gaze)
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+24*p, 3*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+31*p, y+24*p, 3*p, eyeWhite, false)

	eyeColor := color.RGBA{100, 80, 60, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+24*p, 2*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+24*p, 2*p, eyeColor, false)

	// Nose (black)
	noseColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+23*p, y+32*p, 4*p, 3*p, noseColor, false)

	// White chest patch
	chestPatch := color.RGBA{200, 200, 200, 255}
	vector.DrawFilledRect(screen, x+20*p, y+40*p, 10*p, 6*p, chestPatch, false)

	// Frame (bright to stand out)
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{100, 150, 255, 255}, false)
}

// Draw Labrador avatar (Eleanor)
func DrawLabradorAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body (bigger and more substantial)
	bodyColor := color.RGBA{65, 65, 70, 255}
	vector.DrawFilledRect(screen, x+8*p, y+40*p, 34*p, 10*p, bodyColor, false)

	// Chest (broad)
	vector.DrawFilledRect(screen, x+12*p, y+36*p, 26*p, 6*p, bodyColor, false)

	// Thick neck
	vector.DrawFilledRect(screen, x+16*p, y+32*p, 18*p, 6*p, bodyColor, false)

	// Head (larger, broad and friendly - typical lab)
	headColor := color.RGBA{75, 75, 80, 255}
	vector.DrawFilledRect(screen, x+10*p, y+18*p, 30*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+16*p, 26*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+32*p, 26*p, 3*p, headColor, false)

	// Floppy ears (characteristic labrador ears - bigger)
	earColor := color.RGBA{55, 55, 60, 255}
	// Left ear
	vector.DrawFilledRect(screen, x+6*p, y+20*p, 10*p, 16*p, earColor, false)
	vector.DrawFilledRect(screen, x+4*p, y+24*p, 6*p, 12*p, earColor, false)
	// Right ear
	vector.DrawFilledRect(screen, x+34*p, y+20*p, 10*p, 16*p, earColor, false)
	vector.DrawFilledRect(screen, x+40*p, y+24*p, 6*p, 12*p, earColor, false)

	// Friendly eyes (wider set for broad head)
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+17*p, y+24*p, 4*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+33*p, y+24*p, 4*p, eyeWhite, false)

	eyeColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+17*p, y+24*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+33*p, y+24*p, 3*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledCircle(screen, x+17*p, y+25*p, 1.5*p, pupilColor, false)
	vector.DrawFilledCircle(screen, x+33*p, y+25*p, 1.5*p, pupilColor, false)

	// Eye sparkle
	vector.DrawFilledCircle(screen, x+18*p, y+23*p, 1*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(screen, x+34*p, y+23*p, 1*p, color.RGBA{255, 255, 255, 255}, false)

	// Broad muzzle
	muzzleColor := color.RGBA{85, 85, 90, 255}
	vector.DrawFilledRect(screen, x+16*p, y+28*p, 18*p, 8*p, muzzleColor, false)

	// Nose (larger for lab)
	noseColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+22*p, y+30*p, 6*p, 4*p, noseColor, false)

	// Friendly mouth
	mouthColor := color.RGBA{60, 60, 65, 255}
	vector.DrawFilledRect(screen, x+25*p, y+33*p, 1*p, 2*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+21*p, y+34*p, 4*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+26*p, y+34*p, 4*p, 1*p, mouthColor, false)

	// Frame (bright to stand out)
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{180, 100, 255, 255}, false)
}

// Draw GoldenDoodle avatar (Stella)
func DrawGoldenDoodleAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Fluffy body
	bodyColor := color.RGBA{255, 215, 170, 255}
	vector.DrawFilledRect(screen, x+8*p, y+40*p, 34*p, 10*p, bodyColor, false)
	// Fluffy texture
	vector.DrawFilledCircle(screen, x+12*p, y+42*p, 4*p, bodyColor, false)
	vector.DrawFilledCircle(screen, x+22*p, y+44*p, 4*p, bodyColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+42*p, 4*p, bodyColor, false)

	// Fluffy neck
	vector.DrawFilledRect(screen, x+16*p, y+34*p, 18*p, 8*p, bodyColor, false)
	vector.DrawFilledCircle(screen, x+18*p, y+36*p, 3*p, bodyColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+36*p, 3*p, bodyColor, false)

	// Round fluffy head (teddy bear look)
	headColor := color.RGBA{255, 220, 180, 255}
	vector.DrawFilledRect(screen, x+10*p, y+18*p, 30*p, 18*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+16*p, 26*p, 4*p, headColor, false)
	vector.DrawFilledRect(screen, x+12*p, y+34*p, 26*p, 4*p, headColor, false)
	// Fluffy cheeks
	vector.DrawFilledCircle(screen, x+12*p, y+28*p, 5*p, headColor, false)
	vector.DrawFilledCircle(screen, x+38*p, y+28*p, 5*p, headColor, false)

	// Floppy fluffy ears
	earColor := color.RGBA{240, 200, 160, 255}
	// Left ear
	vector.DrawFilledCircle(screen, x+12*p, y+20*p, 6*p, earColor, false)
	vector.DrawFilledCircle(screen, x+10*p, y+26*p, 5*p, earColor, false)
	// Right ear
	vector.DrawFilledCircle(screen, x+38*p, y+20*p, 6*p, earColor, false)
	vector.DrawFilledCircle(screen, x+40*p, y+26*p, 5*p, earColor, false)

	// Big cute eyes
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+24*p, 4*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+31*p, y+24*p, 4*p, eyeWhite, false)

	eyeColor := color.RGBA{80, 60, 40, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+24*p, 3*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+24*p, 3*p, eyeColor, false)

	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledCircle(screen, x+19*p, y+25*p, 1.5*p, pupilColor, false)
	vector.DrawFilledCircle(screen, x+31*p, y+25*p, 1.5*p, pupilColor, false)

	// Big sparkles
	vector.DrawFilledCircle(screen, x+20*p, y+23*p, 1.5*p, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(screen, x+32*p, y+23*p, 1.5*p, color.RGBA{255, 255, 255, 255}, false)

	// Nose (black)
	noseColor := color.RGBA{40, 40, 40, 255}
	vector.DrawFilledRect(screen, x+23*p, y+30*p, 4*p, 3*p, noseColor, false)

	// Happy mouth
	mouthColor := color.RGBA{100, 80, 60, 255}
	vector.DrawFilledRect(screen, x+25*p, y+32*p, 1*p, 2*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+21*p, y+33*p, 4*p, 1*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+26*p, y+33*p, 4*p, 1*p, mouthColor, false)

	// Fluffy tail
	tailColor := color.RGBA{250, 210, 165, 255}
	vector.DrawFilledCircle(screen, x+40*p, y+42*p, 5*p, tailColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 215, 170, 255}, false)
}

// Draw Shepherd avatar (Huckleberry)
func DrawShepherdAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body (lighter for visibility)
	bodyColor := color.RGBA{70, 70, 75, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, bodyColor, false)

	// Neck (thick and strong)
	vector.DrawFilledRect(screen, x+16*p, y+36*p, 18*p, 8*p, bodyColor, false)

	// Head (noble and strong)
	headColor := color.RGBA{80, 80, 85, 255}
	vector.DrawFilledRect(screen, x+12*p, y+22*p, 26*p, 14*p, headColor, false)
	vector.DrawFilledRect(screen, x+14*p, y+20*p, 22*p, 4*p, headColor, false)

	// Pointed ears (distinctive shepherd characteristic)
	earColor := color.RGBA{60, 60, 65, 255}
	// Left ear
	for i := 0; i < 10; i++ {
		w := float32(10 - i)
		vector.DrawFilledRect(screen, x+12*p+float32(i)*p/2, y+10*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	// Right ear
	for i := 0; i < 10; i++ {
		w := float32(10 - i)
		vector.DrawFilledRect(screen, x+33*p-float32(i)*p/2, y+10*p+float32(i)*p, w*p, 2*p, earColor, false)
	}

	// Inner ears
	innerEarColor := color.RGBA{90, 90, 95, 255}
	vector.DrawFilledRect(screen, x+15*p, y+14*p, 3*p, 5*p, innerEarColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+14*p, 3*p, 5*p, innerEarColor, false)

	// Alert, intelligent eyes
	eyeWhite := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 3*p, eyeWhite, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 3*p, eyeWhite, false)

	eyeColor := color.RGBA{120, 80, 40, 255}
	vector.DrawFilledCircle(screen, x+18*p, y+26*p, 2*p, eyeColor, false)
	vector.DrawFilledCircle(screen, x+32*p, y+26*p, 2*p, eyeColor, false)

	// Snout/muzzle
	muzzleColor := color.RGBA{90, 90, 95, 255}
	vector.DrawFilledRect(screen, x+16*p, y+30*p, 18*p, 8*p, muzzleColor, false)

	// Black nose
	noseColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+23*p, y+32*p, 4*p, 3*p, noseColor, false)

	// White belly/chest
	chestPatch := color.RGBA{220, 220, 220, 255}
	vector.DrawFilledRect(screen, x+18*p, y+38*p, 14*p, 8*p, chestPatch, false)

	// Frame (bright to stand out)
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 150, 100, 255}, false)
}

// Draw White Cat avatar (Winston)
func DrawWhiteCatAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body/chest
	chestColor := color.RGBA{250, 250, 250, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, chestColor, false)

	// Neck
	neckColor := color.RGBA{245, 245, 245, 255}
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, neckColor, false)

	// Head shape
	headColor := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledRect(screen, x+14*p, y+20*p, 22*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+16*p, 18*p, 6*p, headColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+14*p, 14*p, 4*p, headColor, false)

	// Ears (triangular)
	earColor := color.RGBA{240, 240, 240, 255}
	// Left ear
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+10*p+float32(i)*p/2, y+8*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	// Right ear
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+30*p-float32(i)*p/2, y+8*p+float32(i)*p, w*p, 2*p, earColor, false)
	}

	// Inner ears (pink)
	innerEarColor := color.RGBA{255, 200, 220, 255}
	vector.DrawFilledRect(screen, x+13*p, y+12*p, 4*p, 4*p, innerEarColor, false)
	vector.DrawFilledRect(screen, x+33*p, y+12*p, 4*p, 4*p, innerEarColor, false)

	// Eyes (blue for white cat)
	eyeColor := color.RGBA{100, 180, 220, 255}
	vector.DrawFilledRect(screen, x+18*p, y+22*p, 5*p, 4*p, eyeColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+22*p, 5*p, 4*p, eyeColor, false)

	// Pupils (vertical slits)
	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+20*p, y+23*p, 1*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+23*p, 1*p, 3*p, pupilColor, false)

	// Nose (pink)
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+28*p, 4*p, 2*p, noseColor, false)
	vector.DrawFilledRect(screen, x+24*p, y+27*p, 2*p, 1*p, noseColor, false)

	// Mouth
	mouthColor := color.RGBA{200, 180, 180, 255}
	vector.DrawFilledRect(screen, x+25*p, y+30*p, 1*p, 3*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+23*p, y+32*p, 5*p, 1*p, mouthColor, false)

	// Whiskers
	whiskerColor := color.RGBA{200, 200, 200, 255}
	vector.DrawFilledRect(screen, x+8*p, y+26*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+8*p, y+28*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+26*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+28*p, 8*p, 1*p, whiskerColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{200, 200, 240, 255}, false)
}

// Draw Grey Cat avatar (Baxter)
func DrawGreyCatAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Body/chest
	chestColor := color.RGBA{130, 130, 130, 255}
	vector.DrawFilledRect(screen, x+10*p, y+42*p, 30*p, 8*p, chestColor, false)

	// Neck
	neckColor := color.RGBA{120, 120, 120, 255}
	vector.DrawFilledRect(screen, x+18*p, y+36*p, 14*p, 8*p, neckColor, false)

	// Head shape
	headColor := color.RGBA{140, 140, 140, 255}
	vector.DrawFilledRect(screen, x+14*p, y+20*p, 22*p, 16*p, headColor, false)
	vector.DrawFilledRect(screen, x+16*p, y+16*p, 18*p, 6*p, headColor, false)
	vector.DrawFilledRect(screen, x+18*p, y+14*p, 14*p, 4*p, headColor, false)

	// Ears (triangular)
	earColor := color.RGBA{110, 110, 110, 255}
	// Left ear
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+10*p+float32(i)*p/2, y+8*p+float32(i)*p, w*p, 2*p, earColor, false)
	}
	// Right ear
	for i := 0; i < 8; i++ {
		w := float32(8 - i)
		vector.DrawFilledRect(screen, x+30*p-float32(i)*p/2, y+8*p+float32(i)*p, w*p, 2*p, earColor, false)
	}

	// Inner ears
	innerEarColor := color.RGBA{255, 180, 200, 255}
	vector.DrawFilledRect(screen, x+13*p, y+12*p, 4*p, 4*p, innerEarColor, false)
	vector.DrawFilledRect(screen, x+33*p, y+12*p, 4*p, 4*p, innerEarColor, false)

	// Eyes (green/yellow)
	eyeColor := color.RGBA{150, 200, 100, 255}
	vector.DrawFilledRect(screen, x+18*p, y+22*p, 5*p, 4*p, eyeColor, false)
	vector.DrawFilledRect(screen, x+27*p, y+22*p, 5*p, 4*p, eyeColor, false)

	// Pupils (vertical slits)
	pupilColor := color.RGBA{20, 20, 20, 255}
	vector.DrawFilledRect(screen, x+20*p, y+23*p, 1*p, 3*p, pupilColor, false)
	vector.DrawFilledRect(screen, x+29*p, y+23*p, 1*p, 3*p, pupilColor, false)

	// Nose (pink)
	noseColor := color.RGBA{255, 150, 180, 255}
	vector.DrawFilledRect(screen, x+23*p, y+28*p, 4*p, 2*p, noseColor, false)
	vector.DrawFilledRect(screen, x+24*p, y+27*p, 2*p, 1*p, noseColor, false)

	// Mouth
	mouthColor := color.RGBA{80, 60, 50, 255}
	vector.DrawFilledRect(screen, x+25*p, y+30*p, 1*p, 3*p, mouthColor, false)
	vector.DrawFilledRect(screen, x+23*p, y+32*p, 5*p, 1*p, mouthColor, false)

	// Whiskers (lighter grey)
	whiskerColor := color.RGBA{90, 90, 90, 255}
	vector.DrawFilledRect(screen, x+8*p, y+26*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+8*p, y+28*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+26*p, 8*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, x+34*p, y+28*p, 8*p, 1*p, whiskerColor, false)

	// Tabby stripes
	stripeColor := color.RGBA{80, 80, 80, 255}
	vector.DrawFilledRect(screen, x+16*p, y+22*p, 2*p, 5*p, stripeColor, false)
	vector.DrawFilledRect(screen, x+32*p, y+22*p, 2*p, 5*p, stripeColor, false)

	// Frame
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{140, 140, 180, 255}, false)
}

// Draw Ribbon & Puddles avatar - two orange kitten sisters
func DrawRibbonPuddlesAvatar(screen *ebiten.Image, x, y, scale float32) {
	p := float32(1) * scale

	// Orange color palette
	orangeFur := color.RGBA{255, 140, 50, 255}
	darkOrange := color.RGBA{220, 100, 30, 255}
	lightOrange := color.RGBA{255, 180, 100, 255}
	pinkNose := color.RGBA{255, 150, 180, 255}
	pinkInnerEar := color.RGBA{255, 180, 200, 255}
	eyeColor := color.RGBA{100, 180, 100, 255}
	pupilColor := color.RGBA{20, 20, 20, 255}
	whiskerColor := color.RGBA{200, 120, 60, 255}

	// === Left kitten (Ribbon) - slightly larger ===
	lx := x + 2*p

	// Body
	vector.DrawFilledRect(screen, lx+2*p, y+38*p, 18*p, 10*p, orangeFur, false)

	// Neck
	vector.DrawFilledRect(screen, lx+6*p, y+32*p, 10*p, 8*p, darkOrange, false)

	// Head
	vector.DrawFilledRect(screen, lx+4*p, y+18*p, 14*p, 14*p, orangeFur, false)
	vector.DrawFilledRect(screen, lx+5*p, y+15*p, 12*p, 5*p, orangeFur, false)

	// Ears
	for i := 0; i < 5; i++ {
		w := float32(5 - i)
		vector.DrawFilledRect(screen, lx+3*p+float32(i)*p/2, y+10*p+float32(i)*p, w*p, 2*p, darkOrange, false)
	}
	for i := 0; i < 5; i++ {
		w := float32(5 - i)
		vector.DrawFilledRect(screen, lx+14*p-float32(i)*p/2, y+10*p+float32(i)*p, w*p, 2*p, darkOrange, false)
	}

	// Inner ears
	vector.DrawFilledRect(screen, lx+5*p, y+13*p, 2*p, 2*p, pinkInnerEar, false)
	vector.DrawFilledRect(screen, lx+15*p, y+13*p, 2*p, 2*p, pinkInnerEar, false)

	// Eyes
	vector.DrawFilledRect(screen, lx+6*p, y+20*p, 4*p, 3*p, eyeColor, false)
	vector.DrawFilledRect(screen, lx+12*p, y+20*p, 4*p, 3*p, eyeColor, false)

	// Pupils
	vector.DrawFilledRect(screen, lx+7*p, y+21*p, 1*p, 2*p, pupilColor, false)
	vector.DrawFilledRect(screen, lx+13*p, y+21*p, 1*p, 2*p, pupilColor, false)

	// Nose
	vector.DrawFilledRect(screen, lx+9*p, y+25*p, 3*p, 2*p, pinkNose, false)

	// Mouth
	vector.DrawFilledRect(screen, lx+10*p, y+27*p, 1*p, 2*p, darkOrange, false)

	// Whiskers
	vector.DrawFilledRect(screen, lx+1*p, y+24*p, 5*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, lx+16*p, y+24*p, 5*p, 1*p, whiskerColor, false)

	// Tabby stripes on forehead
	vector.DrawFilledRect(screen, lx+10*p, y+16*p, 2*p, 3*p, darkOrange, false)

	// === Right kitten (Puddles) - slightly smaller ===
	rx := x + 26*p

	// Body
	vector.DrawFilledRect(screen, rx+2*p, y+40*p, 16*p, 8*p, lightOrange, false)

	// Neck
	vector.DrawFilledRect(screen, rx+5*p, y+35*p, 8*p, 6*p, orangeFur, false)

	// Head
	vector.DrawFilledRect(screen, rx+3*p, y+22*p, 12*p, 13*p, lightOrange, false)
	vector.DrawFilledRect(screen, rx+4*p, y+19*p, 10*p, 5*p, lightOrange, false)

	// Ears
	for i := 0; i < 4; i++ {
		w := float32(4 - i)
		vector.DrawFilledRect(screen, rx+2*p+float32(i)*p/2, y+15*p+float32(i)*p, w*p, 2*p, orangeFur, false)
	}
	for i := 0; i < 4; i++ {
		w := float32(4 - i)
		vector.DrawFilledRect(screen, rx+12*p-float32(i)*p/2, y+15*p+float32(i)*p, w*p, 2*p, orangeFur, false)
	}

	// Inner ears
	vector.DrawFilledRect(screen, rx+4*p, y+17*p, 2*p, 2*p, pinkInnerEar, false)
	vector.DrawFilledRect(screen, rx+12*p, y+17*p, 2*p, 2*p, pinkInnerEar, false)

	// Eyes
	vector.DrawFilledRect(screen, rx+5*p, y+24*p, 3*p, 3*p, eyeColor, false)
	vector.DrawFilledRect(screen, rx+10*p, y+24*p, 3*p, 3*p, eyeColor, false)

	// Pupils
	vector.DrawFilledRect(screen, rx+6*p, y+25*p, 1*p, 2*p, pupilColor, false)
	vector.DrawFilledRect(screen, rx+11*p, y+25*p, 1*p, 2*p, pupilColor, false)

	// Nose
	vector.DrawFilledRect(screen, rx+7*p, y+28*p, 3*p, 2*p, pinkNose, false)

	// Mouth
	vector.DrawFilledRect(screen, rx+8*p, y+30*p, 1*p, 2*p, orangeFur, false)

	// Whiskers
	vector.DrawFilledRect(screen, rx+1*p, y+27*p, 4*p, 1*p, whiskerColor, false)
	vector.DrawFilledRect(screen, rx+13*p, y+27*p, 4*p, 1*p, whiskerColor, false)

	// Tabby stripes on forehead
	vector.DrawFilledRect(screen, rx+8*p, y+20*p, 2*p, 3*p, orangeFur, false)

	// Frame (warm orange-ish border)
	vector.StrokeRect(screen, x, y, 50*p, 50*p, 2, color.RGBA{255, 160, 100, 255}, false)
}

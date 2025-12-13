package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type HomeScreen struct {
	gameButtons    []*Button
	goOnlineButton *Button
}

func NewHomeScreen() *HomeScreen {
	hs := &HomeScreen{
		gameButtons: make([]*Button, 4),
	}

	// Create game selection buttons - 2x2 grid (4 games)
	buttonWidth := 280.0
	buttonHeight := 90.0
	startX := float64(screenWidth/2) - buttonWidth - 20
	startY := 250.0
	spacingX := buttonWidth + 40.0
	spacingY := buttonHeight + 20.0

	hs.gameButtons[0] = &Button{
		x:       startX,
		y:       startY,
		width:   buttonWidth,
		height:  buttonHeight,
		text:    "YAHTZEE",
		enabled: true,
	}

	hs.gameButtons[1] = &Button{
		x:       startX + spacingX,
		y:       startY,
		width:   buttonWidth,
		height:  buttonHeight,
		text:    "SANTORINI",
		enabled: true,
	}

	hs.gameButtons[2] = &Button{
		x:       startX,
		y:       startY + spacingY,
		width:   buttonWidth,
		height:  buttonHeight,
		text:    "CONNECT FOUR",
		enabled: true,
	}

	hs.gameButtons[3] = &Button{
		x:       startX + spacingX,
		y:       startY + spacingY,
		width:   buttonWidth,
		height:  buttonHeight,
		text:    "MEMORY MATCH",
		enabled: true,
	}

	// Add "Go Online" button at the bottom
	hs.goOnlineButton = &Button{
		x:       float64(screenWidth/2) - 150,
		y:       float64(screenHeight) - 80,
		width:   300,
		height:  60,
		text:    "GO ONLINE",
		enabled: true,
	}

	return hs
}

func (hs *HomeScreen) Update(gr *GameRoom) error {
	// Update hover states
	x, y := ebiten.CursorPosition()
	for _, btn := range hs.gameButtons {
		btn.hovered = btn.Contains(x, y)
	}

	// Update Go Online button hover state (only show if offline)
	if !gr.isOnlineMode && hs.goOnlineButton != nil {
		hs.goOnlineButton.hovered = hs.goOnlineButton.Contains(x, y)
	}

	// Handle button clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if hs.gameButtons[0].hovered {
			gr.SwitchToGame(NewYahtzeeGame())
		} else if hs.gameButtons[1].hovered {
			gr.SwitchToGame(NewSantoriniGame())
		} else if hs.gameButtons[2].hovered {
			gr.SwitchToGame(NewConnectFourGame())
		} else if hs.gameButtons[3].hovered {
			gr.SwitchToGame(NewMemoryGame())
		} else if !gr.isOnlineMode && hs.goOnlineButton != nil && hs.goOnlineButton.hovered {
			// Try to reconnect
			gr.TryGoOnline()
		}
	}

	return nil
}

func (hs *HomeScreen) Draw(screen *ebiten.Image, gr *GameRoom) {
	// Draw background
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)

	// Draw title
	titleWidth := float32(500)
	titleX := float32(screenWidth/2) - titleWidth/2
	titleY := float32(80)

	vector.DrawFilledRect(screen, titleX, titleY, titleWidth, 80, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, titleY, titleWidth, 80, 3, color.RGBA{100, 150, 220, 255}, false)

	// Title text
	title1 := "OLIVE & MILLIE'S"
	title2 := "GAME ROOM"
	title1X := int(titleX + (titleWidth-float32(len(title1)*6))/2)
	title2X := int(titleX + (titleWidth-float32(len(title2)*6))/2)

	ebitenutil.DebugPrintAt(screen, title1, title1X, int(titleY+20))
	ebitenutil.DebugPrintAt(screen, title1, title1X+1, int(titleY+20))
	ebitenutil.DebugPrintAt(screen, title2, title2X, int(titleY+45))
	ebitenutil.DebugPrintAt(screen, title2, title2X+1, int(titleY+45))

	// Draw avatars next to title
	DrawPlayer1Avatar(screen, titleX-80, titleY+15, 1.0)
	DrawPlayer2Avatar(screen, titleX+titleWidth+30, titleY+15, 1.0)

	// Draw game selection buttons
	for _, btn := range hs.gameButtons {
		hs.drawGameButton(screen, btn)
	}

	// Draw "Go Online" button if in offline mode
	if !gr.isOnlineMode && hs.goOnlineButton != nil {
		DrawButton(screen, hs.goOnlineButton)
	}

	// Instructions
	instructionText := "Select a game to play!"
	if !gr.isOnlineMode {
		instructionText = "Offline Mode - Click 'Go Online' to connect"
	}
	instructionX := screenWidth/2 - len(instructionText)*3
	vector.DrawFilledRect(screen, float32(instructionX-10), 550, float32(len(instructionText)*6+20), 30, color.RGBA{30, 50, 80, 200}, false)
	ebitenutil.DebugPrintAt(screen, instructionText, instructionX, 560)
}

func (hs *HomeScreen) drawGameButton(screen *ebiten.Image, btn *Button) {
	btnColor := color.RGBA{100, 150, 220, 255}
	borderColor := color.RGBA{70, 120, 190, 255}
	if btn.hovered {
		btnColor = color.RGBA{130, 180, 240, 255}
		borderColor = color.RGBA{100, 150, 220, 255}
	}

	// Shadow
	vector.DrawFilledRect(screen, float32(btn.x+4), float32(btn.y+4), float32(btn.width), float32(btn.height), color.RGBA{0, 0, 0, 60}, false)

	// Button background
	vector.DrawFilledRect(screen, float32(btn.x), float32(btn.y), float32(btn.width), float32(btn.height), btnColor, false)
	vector.StrokeRect(screen, float32(btn.x), float32(btn.y), float32(btn.width), float32(btn.height), 3, borderColor, false)

	// Button text
	textX := int(btn.x + btn.width/2 - float64(len(btn.text))*3)
	textY := int(btn.y + btn.height/2 - 5)
	ebitenutil.DebugPrintAt(screen, btn.text, textX, textY)
	ebitenutil.DebugPrintAt(screen, btn.text, textX+1, textY) // Bold
	ebitenutil.DebugPrintAt(screen, btn.text, textX, textY+1) // Bold
}

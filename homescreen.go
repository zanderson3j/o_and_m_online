package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type HomeScreen struct {
	gameButtons []*Button
	retryButton *Button
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

	// Add "Retry" button (only shown when connection fails)
	hs.retryButton = &Button{
		x:       float64(screenWidth/2) - 100,
		y:       float64(screenHeight) - 100,
		width:   200,
		height:  60,
		text:    "RETRY CONNECTION",
		enabled: true,
	}

	return hs
}

func (hs *HomeScreen) Update(gr *GameRoom) error {
	// Only show home screen if not connected
	if gr.connectionState != StateConnected {
		x, y := ebiten.CursorPosition()

		// Update retry button hover state (only when connection failed)
		if gr.connectionState == StateFailed && hs.retryButton != nil {
			hs.retryButton.hovered = hs.retryButton.Contains(x, y)
		}

		// Handle retry button click
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if gr.connectionState == StateFailed && hs.retryButton != nil && hs.retryButton.hovered {
				log.Println("Retry button clicked - attempting to reconnect")
				gr.TryGoOnline()
			}
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

	// Draw connection status message
	var statusText string
	var statusColor color.RGBA

	switch gr.connectionState {
	case StateConnecting:
		statusText = "Connecting to server..."
		statusColor = color.RGBA{100, 150, 220, 255}
	case StateFailed:
		statusText = "Connection Failed"
		statusColor = color.RGBA{200, 80, 80, 255}
	case StateConnected:
		statusText = "Connected!"
		statusColor = color.RGBA{80, 200, 80, 255}
	}

	// Draw status box
	statusWidth := float32(len(statusText)*6 + 40)
	statusX := float32(screenWidth/2) - statusWidth/2
	statusY := float32(400)
	vector.DrawFilledRect(screen, statusX, statusY, statusWidth, 50, color.RGBA{30, 50, 80, 200}, false)
	vector.StrokeRect(screen, statusX, statusY, statusWidth, 50, 3, statusColor, false)

	// Draw status text
	textX := int(statusX + statusWidth/2 - float32(len(statusText)*6)/2)
	ebitenutil.DebugPrintAt(screen, statusText, textX, int(statusY+17))
	ebitenutil.DebugPrintAt(screen, statusText, textX+1, int(statusY+17)) // Bold

	// Draw error message if failed
	if gr.connectionState == StateFailed && gr.connectionError != "" {
		errorX := screenWidth/2 - len(gr.connectionError)*3
		ebitenutil.DebugPrintAt(screen, gr.connectionError, errorX, int(statusY+60))
	}

	// Draw retry button if connection failed
	if gr.connectionState == StateFailed && hs.retryButton != nil {
		DrawButton(screen, hs.retryButton)
	}
	
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

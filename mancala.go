package main

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	mancala_pits    = 6
	mancala_pitSize = 80
	mancala_storeW  = 100
	mancala_storeH  = 200
)

type MancalaMove struct {
	PitIndex int `json:"pit_index"`
}

type MancalaGame struct {
	pits          [2][mancala_pits]int // Each player has 6 pits
	stores        [2]int               // Each player's store (mancala)
	currentPlayer int
	winner        int // 0 = no winner, 1 = player 1, 2 = player 2, 3 = tie
	gameOver      bool
	boardOffsetX  float32
	boardOffsetY  float32
	networkClient *NetworkClient
	myPlayerNum   int // 0 or 1
}

func NewMancalaGame() *MancalaGame {
	return NewMancalaGameWithNetwork(nil, 0)
}

func NewMancalaGameWithNetwork(nc *NetworkClient, playerNum int) *MancalaGame {
	g := &MancalaGame{
		currentPlayer: 0,
		winner:        0,
		gameOver:      false,
		boardOffsetX:  150,
		boardOffsetY:  300,
		networkClient: nc,
		myPlayerNum:   playerNum,
	}

	// Initialize each pit with 4 stones
	for i := 0; i < mancala_pits; i++ {
		g.pits[0][i] = 4
		g.pits[1][i] = 4
	}

	// Register network handler for opponent moves
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move MancalaMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.makeMove(move.PitIndex)
			}
		})
	}

	return g
}

func (g *MancalaGame) Reset() {
	*g = *NewMancalaGame()
}

func (g *MancalaGame) Update(gr *GameRoom) error {
	if IsLogoClicked() {
		gr.ReturnHome()
		return nil
	}

	if g.gameOver {
		return nil
	}

	// Only allow input if it's my turn (or if no network client)
	isMyTurn := g.networkClient == nil || g.currentPlayer == g.myPlayerNum

	if isMyTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		pitIndex := g.handleClick(float32(x), float32(y))

		// Send move to opponent if a valid move was made
		if pitIndex >= 0 && g.networkClient != nil {
			move := MancalaMove{PitIndex: pitIndex}
			g.networkClient.SendGameMove(move)
		}
	}

	return nil
}

func (g *MancalaGame) handleClick(x, y float32) int {
	// Check if clicked on current player's pits
	for i := 0; i < mancala_pits; i++ {
		pitX, pitY := g.getPitPosition(g.currentPlayer, i)

		if x >= pitX && x <= pitX+mancala_pitSize && y >= pitY && y <= pitY+mancala_pitSize {
			if g.pits[g.currentPlayer][i] > 0 {
				g.makeMove(i)
				return i // Return pit index for network sync
			}
			return -1
		}
	}
	return -1 // No valid pit clicked
}

func (g *MancalaGame) getPitPosition(player, pit int) (float32, float32) {
	if player == 0 {
		// Player 0's pits are on bottom row, left to right
		x := g.boardOffsetX + mancala_storeW + float32(pit*mancala_pitSize) + 20
		y := g.boardOffsetY + mancala_storeH - mancala_pitSize - 20
		return x, y
	} else {
		// Player 1's pits are on top row, right to left
		x := g.boardOffsetX + mancala_storeW + float32((mancala_pits-1-pit)*mancala_pitSize) + 20
		y := g.boardOffsetY + 20
		return x, y
	}
}

func (g *MancalaGame) makeMove(pitIndex int) {
	stones := g.pits[g.currentPlayer][pitIndex]
	g.pits[g.currentPlayer][pitIndex] = 0

	currentPlayer := g.currentPlayer
	currentPit := pitIndex

	// Distribute stones
	for stones > 0 {
		// Move to next position
		currentPit++

		// Handle wrapping within current player's side
		if currentPit >= mancala_pits {
			// Add to current player's store
			if currentPlayer == g.currentPlayer {
				g.stores[g.currentPlayer]++
				stones--
				if stones == 0 {
					// Landed in own store, get another turn
					g.checkGameOver()
					return
				}
			}
			// Switch to opponent's side
			currentPlayer = 1 - currentPlayer
			currentPit = 0
		}

		if stones > 0 {
			g.pits[currentPlayer][currentPit]++
			stones--
		}
	}

	// Check for capture: if last stone landed in empty pit on own side
	if currentPlayer == g.currentPlayer && g.pits[currentPlayer][currentPit] == 1 {
		oppositePit := mancala_pits - 1 - currentPit
		if g.pits[1-currentPlayer][oppositePit] > 0 {
			// Capture own stone + opposite stones
			g.stores[g.currentPlayer] += 1 + g.pits[1-currentPlayer][oppositePit]
			g.pits[currentPlayer][currentPit] = 0
			g.pits[1-currentPlayer][oppositePit] = 0
		}
	}

	// Switch player
	g.currentPlayer = 1 - g.currentPlayer

	g.checkGameOver()
}

func (g *MancalaGame) checkGameOver() {
	// Check if either side is empty
	player0Empty := true
	player1Empty := true

	for i := 0; i < mancala_pits; i++ {
		if g.pits[0][i] > 0 {
			player0Empty = false
		}
		if g.pits[1][i] > 0 {
			player1Empty = false
		}
	}

	if player0Empty || player1Empty {
		// Game over - collect remaining stones
		for i := 0; i < mancala_pits; i++ {
			g.stores[0] += g.pits[0][i]
			g.stores[1] += g.pits[1][i]
			g.pits[0][i] = 0
			g.pits[1][i] = 0
		}

		g.gameOver = true

		if g.stores[0] > g.stores[1] {
			g.winner = 1
		} else if g.stores[1] > g.stores[0] {
			g.winner = 2
		} else {
			g.winner = 3 // Tie
		}
	}
}

func (g *MancalaGame) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)
	g.drawGameInfo(screen)
	g.drawBoard(screen)
	g.drawPlayerInfo(screen)

	if g.gameOver {
		g.drawWinner(screen)
	}
}

func (g *MancalaGame) drawGameInfo(screen *ebiten.Image) {
	titleWidth := float32(200)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "MANCALA"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	infoWidth := float32(300)
	infoX := float32(screenWidth/2) - infoWidth/2
	vector.DrawFilledRect(screen, infoX, 70, infoWidth, 50, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, infoX, 70, infoWidth, 50, 2, color.RGBA{100, 150, 220, 255}, false)

	var phaseText string
	if g.gameOver {
		phaseText = "Game Over"
	} else {
		phaseText = fmt.Sprintf("Player %d's Turn", g.currentPlayer+1)
	}

	phaseTextX := int(infoX + (infoWidth-float32(len(phaseText)*6))/2)
	ebitenutil.DebugPrintAt(screen, phaseText, phaseTextX, 90)
}

func (g *MancalaGame) drawBoard(screen *ebiten.Image) {
	// Draw stores (mancalas)
	// Player 1's store (left side)
	storeX1 := g.boardOffsetX
	storeY := g.boardOffsetY
	vector.DrawFilledRect(screen, storeX1, storeY, mancala_storeW, mancala_storeH, color.RGBA{100, 150, 220, 255}, false)
	vector.StrokeRect(screen, storeX1, storeY, mancala_storeW, mancala_storeH, 2, color.RGBA{70, 120, 190, 255}, false)

	// Player 0's store (right side)
	storeX2 := g.boardOffsetX + mancala_storeW + float32(mancala_pits*mancala_pitSize) + 40
	vector.DrawFilledRect(screen, storeX2, storeY, mancala_storeW, mancala_storeH, color.RGBA{100, 150, 220, 255}, false)
	vector.StrokeRect(screen, storeX2, storeY, mancala_storeW, mancala_storeH, 2, color.RGBA{70, 120, 190, 255}, false)

	// Draw store counts with kodama spirits
	g.drawStones(screen, storeX1+mancala_storeW/2, storeY+mancala_storeH/2, g.stores[1])
	g.drawStones(screen, storeX2+mancala_storeW/2, storeY+mancala_storeH/2, g.stores[0])

	// Draw pits
	for player := 0; player < 2; player++ {
		for i := 0; i < mancala_pits; i++ {
			x, y := g.getPitPosition(player, i)

			pitColor := color.RGBA{180, 200, 230, 255}
			if player == g.currentPlayer && !g.gameOver && g.pits[player][i] > 0 {
				pitColor = color.RGBA{200, 220, 255, 255} // Highlight available moves
			}

			vector.DrawFilledRect(screen, x, y, mancala_pitSize, mancala_pitSize, pitColor, false)
			vector.StrokeRect(screen, x, y, mancala_pitSize, mancala_pitSize, 2, color.RGBA{100, 150, 220, 255}, false)

			// Draw stones
			g.drawStones(screen, x+mancala_pitSize/2, y+mancala_pitSize/2, g.pits[player][i])
		}
	}
}

func (g *MancalaGame) drawStones(screen *ebiten.Image, centerX, centerY float32, count int) {
	if count == 0 {
		return
	}

	// Draw kodama-like stones
	stoneColor := color.RGBA{200, 255, 220, 200}

	// Draw count as text
	countText := fmt.Sprintf("%d", count)
	textX := int(centerX) - len(countText)*3
	textY := int(centerY) - 5

	// Small kodama representation
	if count <= 3 {
		for i := 0; i < count; i++ {
			offset := float32((i - count/2) * 15)
			vector.DrawFilledCircle(screen, centerX+offset, centerY, 8, stoneColor, false)
		}
	} else {
		// Just show number for larger counts
		vector.DrawFilledCircle(screen, centerX, centerY, 12, stoneColor, false)
		ebitenutil.DebugPrintAt(screen, countText, textX, textY)
	}
}

func (g *MancalaGame) drawPlayerInfo(screen *ebiten.Image) {
	cardWidth := float32(320)
	edgeSpacing := float32(60)
	gapBetween := screenWidth - 2*edgeSpacing - 2*cardWidth

	y := float32(600)
	for i := 0; i < 2; i++ {
		var x float32
		if i == 0 {
			x = edgeSpacing
		} else {
			x = edgeSpacing + cardWidth + gapBetween
		}

		panelColor := color.RGBA{30, 50, 80, 255}
		var borderColor color.RGBA
		if i == 0 {
			borderColor = color.RGBA{100, 150, 220, 255}
		} else {
			borderColor = color.RGBA{200, 160, 120, 255}
		}

		vector.DrawFilledRect(screen, x, y, 320, 100, panelColor, false)
		vector.StrokeRect(screen, x, y, 320, 100, 2, borderColor, false)

		if i == 0 {
			DrawPlayer1Avatar(screen, x+10, y+10, 1.5)
		} else {
			DrawPlayer2Avatar(screen, x+10, y+10, 1.5)
		}

		playerName := fmt.Sprintf("Player %d", i+1)
		ebitenutil.DebugPrintAt(screen, playerName, int(x+90), int(y+30))
		ebitenutil.DebugPrintAt(screen, playerName, int(x+91), int(y+30))

		// Show store count
		storeText := fmt.Sprintf("Stones: %d", g.stores[i])
		ebitenutil.DebugPrintAt(screen, storeText, int(x+90), int(y+50))

		if g.currentPlayer == i && !g.gameOver {
			ebitenutil.DebugPrintAt(screen, "Current Turn", int(x+90), int(y+70))
		}
	}
}

func (g *MancalaGame) drawWinner(screen *ebiten.Image) {
	bannerWidth := float32(450)
	bannerHeight := float32(60)
	bannerX := (screenWidth - bannerWidth) / 2
	bannerY := (screenHeight - bannerHeight) / 2

	vector.DrawFilledRect(screen, bannerX, bannerY, bannerWidth, bannerHeight, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, bannerX, bannerY, bannerWidth, bannerHeight, 3, color.RGBA{100, 150, 220, 255}, false)

	starColor := color.RGBA{150, 200, 255, 200}
	vector.DrawFilledRect(screen, bannerX-20, bannerY+10, 12, 12, starColor, false)
	vector.DrawFilledRect(screen, bannerX+bannerWidth+8, bannerY+10, 12, 12, starColor, false)
	vector.DrawFilledRect(screen, bannerX-20, bannerY+40, 10, 10, starColor, false)
	vector.DrawFilledRect(screen, bannerX+bannerWidth+8, bannerY+40, 10, 10, starColor, false)

	var winnerText string
	if g.winner == 3 {
		winnerText = "TIE GAME!"
	} else {
		winnerText = fmt.Sprintf("WINNER: Player %d", g.winner)
	}

	winnerTextX := int(bannerX + (bannerWidth-float32(len(winnerText)*6))/2)
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX, int(bannerY+15))
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX+1, int(bannerY+15))

	scoreText := fmt.Sprintf("%d - %d", g.stores[0], g.stores[1])
	scoreTextX := int(bannerX + (bannerWidth-float32(len(scoreText)*6))/2)
	ebitenutil.DebugPrintAt(screen, scoreText, scoreTextX, int(bannerY+35))
	ebitenutil.DebugPrintAt(screen, scoreText, scoreTextX+1, int(bannerY+35))
}

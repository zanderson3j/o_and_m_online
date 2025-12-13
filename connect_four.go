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
	cf_rows    = 6
	cf_cols    = 7
	cf_cellSize = 65
)

type ConnectFourMove struct {
	Column int `json:"column"`
}

type ConnectFourPlayer struct {
	id     int
	name   string
	avatar AvatarType
}

type ConnectFourGame struct {
	board         [cf_rows][cf_cols]int // 0 = empty, 1 = player 1, 2 = player 2
	currentPlayer int
	winner        int // 0 = no winner, 1 = player 1, 2 = player 2
	boardOffsetX  float32
	boardOffsetY  float32
	hoveredCol    int
	networkClient *NetworkClient
	myPlayerNum   int // 1 or 2 (determined by join order)
	players       []*ConnectFourPlayer
}

func NewConnectFourGame() *ConnectFourGame {
	return NewConnectFourGameWithNetwork(nil, 1)
}

func NewConnectFourGameWithNetwork(nc *NetworkClient, playerNum int) *ConnectFourGame {
	return NewConnectFourGameWithPlayers(nc, playerNum, nil)
}

func NewConnectFourGameWithPlayers(nc *NetworkClient, playerNum int, playerData []map[string]interface{}) *ConnectFourGame {
	boardWidth := float32(cf_cols * cf_cellSize)
	boardHeight := float32(cf_rows * cf_cellSize)

	topSpace := float32(120)
	bottomSpace := float32(600)
	availableHeight := bottomSpace - topSpace
	boardCenterY := topSpace + (availableHeight-boardHeight)/2

	g := &ConnectFourGame{
		currentPlayer: 1,
		winner:        0,
		boardOffsetX:  (screenWidth - boardWidth) / 2,
		boardOffsetY:  boardCenterY,
		hoveredCol:    -1,
		networkClient: nc,
		myPlayerNum:   playerNum + 1, // Connect Four uses 1/2
		players:       make([]*ConnectFourPlayer, 2),
	}

	// Initialize players with server data
	for i := 0; i < 2; i++ {
		name := fmt.Sprintf("Player %d", i+1)
		avatar := i % int(AvatarNumTypes)

		if playerData != nil && i < len(playerData) {
			if n, ok := playerData[i]["name"].(string); ok {
				name = n
			}
			if a, ok := playerData[i]["avatar"].(float64); ok {
				avatar = int(a)
			}
		}

		g.players[i] = &ConnectFourPlayer{
			id:     i + 1,
			name:   name,
			avatar: AvatarType(avatar),
		}
	}

	// Register network handler for opponent moves
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move ConnectFourMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.dropPiece(move.Column)
			}
		})
	}

	return g
}

func (g *ConnectFourGame) Reset() {
	*g = *NewConnectFourGame()
}

func (g *ConnectFourGame) Update(gr *GameRoom) error {
	if IsLogoClicked() {
		gr.ReturnHome()
		return nil
	}

	// Update hovered column
	x, y := ebiten.CursorPosition()
	boardX := int((float32(x) - g.boardOffsetX) / cf_cellSize)
	boardY := int((float32(y) - g.boardOffsetY) / cf_cellSize)

	if boardX >= 0 && boardX < cf_cols && boardY >= 0 && boardY < cf_rows {
		g.hoveredCol = boardX
	} else {
		g.hoveredCol = -1
	}

	if g.winner != 0 {
		return nil
	}

	// Only allow input if it's my turn (or if no network client)
	isMyTurn := g.networkClient == nil || g.currentPlayer == g.myPlayerNum

	if isMyTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.hoveredCol >= 0 {
			g.dropPiece(g.hoveredCol)

			// Send move to opponent
			if g.networkClient != nil {
				move := ConnectFourMove{Column: g.hoveredCol}
				g.networkClient.SendGameMove(move)
			}
		}
	}

	return nil
}

func (g *ConnectFourGame) dropPiece(col int) {
	// Find lowest empty row in this column
	for row := cf_rows - 1; row >= 0; row-- {
		if g.board[row][col] == 0 {
			g.board[row][col] = g.currentPlayer
			if g.checkWin(row, col) {
				g.winner = g.currentPlayer
			} else {
				// Switch players
				g.currentPlayer = 3 - g.currentPlayer // Toggle between 1 and 2
			}
			return
		}
	}
}

func (g *ConnectFourGame) checkWin(row, col int) bool {
	player := g.board[row][col]

	// Check horizontal
	count := 1
	// Check left
	for c := col - 1; c >= 0 && g.board[row][c] == player; c-- {
		count++
	}
	// Check right
	for c := col + 1; c < cf_cols && g.board[row][c] == player; c++ {
		count++
	}
	if count >= 4 {
		return true
	}

	// Check vertical
	count = 1
	// Check down
	for r := row + 1; r < cf_rows && g.board[r][col] == player; r++ {
		count++
	}
	if count >= 4 {
		return true
	}

	// Check diagonal (top-left to bottom-right)
	count = 1
	// Check up-left
	for r, c := row-1, col-1; r >= 0 && c >= 0 && g.board[r][c] == player; r, c = r-1, c-1 {
		count++
	}
	// Check down-right
	for r, c := row+1, col+1; r < cf_rows && c < cf_cols && g.board[r][c] == player; r, c = r+1, c+1 {
		count++
	}
	if count >= 4 {
		return true
	}

	// Check diagonal (top-right to bottom-left)
	count = 1
	// Check up-right
	for r, c := row-1, col+1; r >= 0 && c < cf_cols && g.board[r][c] == player; r, c = r-1, c+1 {
		count++
	}
	// Check down-left
	for r, c := row+1, col-1; r < cf_rows && c >= 0 && g.board[r][c] == player; r, c = r+1, c-1 {
		count++
	}
	if count >= 4 {
		return true
	}

	return false
}

func (g *ConnectFourGame) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)
	g.drawGameInfo(screen)
	g.drawBoard(screen)
	g.drawPieces(screen)
	g.drawPlayerInfo(screen)

	if g.winner != 0 {
		g.drawWinner(screen)
	}
}

func (g *ConnectFourGame) drawGameInfo(screen *ebiten.Image) {
	titleWidth := float32(250)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "CONNECT FOUR"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	infoWidth := float32(300)
	infoX := float32(screenWidth/2) - infoWidth/2
	vector.DrawFilledRect(screen, infoX, 70, infoWidth, 50, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, infoX, 70, infoWidth, 50, 2, color.RGBA{100, 150, 220, 255}, false)

	var phaseText string
	if g.winner == 0 {
		phaseText = fmt.Sprintf("Player %d's Turn", g.currentPlayer)
	} else {
		phaseText = "Game Over"
	}

	phaseTextX := int(infoX + (infoWidth-float32(len(phaseText)*6))/2)
	ebitenutil.DebugPrintAt(screen, phaseText, phaseTextX, 90)
}

func (g *ConnectFourGame) drawBoard(screen *ebiten.Image) {
	boardColor := color.RGBA{100, 150, 220, 255}

	// Draw board background
	boardWidth := float32(cf_cols * cf_cellSize)
	boardHeight := float32(cf_rows * cf_cellSize)
	vector.DrawFilledRect(screen, g.boardOffsetX, g.boardOffsetY, boardWidth, boardHeight, boardColor, false)

	// Draw holes
	holeRadius := float32(25)
	for row := 0; row < cf_rows; row++ {
		for col := 0; col < cf_cols; col++ {
			x := g.boardOffsetX + float32(col*cf_cellSize) + cf_cellSize/2
			y := g.boardOffsetY + float32(row*cf_cellSize) + cf_cellSize/2
			vector.DrawFilledCircle(screen, x, y, holeRadius, color.RGBA{40, 70, 110, 255}, false)
		}
	}

	// Highlight hovered column
	if g.hoveredCol >= 0 && g.winner == 0 {
		x := g.boardOffsetX + float32(g.hoveredCol*cf_cellSize)
		vector.StrokeRect(screen, x, g.boardOffsetY, cf_cellSize, boardHeight, 3, color.RGBA{255, 255, 100, 200}, false)
	}
}

func (g *ConnectFourGame) drawPieces(screen *ebiten.Image) {
	for row := 0; row < cf_rows; row++ {
		for col := 0; col < cf_cols; col++ {
			if g.board[row][col] != 0 {
				x := g.boardOffsetX + float32(col*cf_cellSize) + cf_cellSize/2
				y := g.boardOffsetY + float32(row*cf_cellSize) + cf_cellSize/2

				var pieceColor color.RGBA
				if g.board[row][col] == 1 {
					pieceColor = color.RGBA{255, 100, 100, 255} // Red for player 1
				} else {
					pieceColor = color.RGBA{255, 220, 100, 255} // Yellow for player 2
				}

				pieceRadius := float32(23)
				vector.DrawFilledCircle(screen, x, y, pieceRadius, pieceColor, false)
				vector.StrokeCircle(screen, x, y, pieceRadius, 2, color.RGBA{200, 200, 200, 255}, false)
			}
		}
	}
}

func (g *ConnectFourGame) drawPlayerInfo(screen *ebiten.Image) {
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

		player := g.players[i]
		DrawAvatar(screen, player.avatar, x+10, y+10, 1.5)

		playerName := player.name
		ebitenutil.DebugPrintAt(screen, playerName, int(x+90), int(y+30))
		ebitenutil.DebugPrintAt(screen, playerName, int(x+91), int(y+30))

		if g.currentPlayer == i+1 {
			ebitenutil.DebugPrintAt(screen, "Current Turn", int(x+90), int(y+50))
		}
	}
}

func (g *ConnectFourGame) drawWinner(screen *ebiten.Image) {
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

	winnerText := fmt.Sprintf("WINNER: %s", g.players[g.winner-1].name)
	winnerTextX := int(bannerX + (bannerWidth-float32(len(winnerText)*6))/2)
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX, int(bannerY+25))
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX+1, int(bannerY+25))
}

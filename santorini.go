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
	boardSize = 5
	cellSize  = 80
)

type Worker struct {
	x, y     int
	playerID int
}

type Cell struct {
	level int
}

type SantoriniPlayer struct {
	id      int
	name    string
	avatar  AvatarType
	workers [2]*Worker
}

type SantoriniMove struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Phase  string `json:"phase"`
	Worker int    `json:"worker"` // 0 or 1
}

type SantoriniGame struct {
	board          [boardSize][boardSize]*Cell
	players        [2]*SantoriniPlayer
	currentPlayer  int
	selectedWorker *Worker
	gamePhase      string
	winner         *SantoriniPlayer
	placementCount int
	boardOffsetX   float32
	boardOffsetY   float32
	networkClient  *NetworkClient
	myPlayerNum    int
}

func NewSantoriniGame() *SantoriniGame {
	return NewSantoriniGameWithNetwork(nil, 0)
}

func NewSantoriniGameWithNetwork(nc *NetworkClient, playerNum int) *SantoriniGame {
	return NewSantoriniGameWithPlayers(nc, playerNum, nil)
}

func NewSantoriniGameWithPlayers(nc *NetworkClient, playerNum int, playerData []map[string]interface{}) *SantoriniGame {
	boardWidth := float32(boardSize * cellSize)
	boardHeight := float32(boardSize * cellSize)

	topSpace := float32(120)
	bottomSpace := float32(600)
	availableHeight := bottomSpace - topSpace
	boardCenterY := topSpace + (availableHeight-boardHeight)/2

	g := &SantoriniGame{
		currentPlayer:  0,
		gamePhase:      "place",
		placementCount: 0,
		boardOffsetX:   (screenWidth - boardWidth) / 2,
		boardOffsetY:   boardCenterY,
		networkClient:  nc,
		myPlayerNum:    playerNum,
	}

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			g.board[i][j] = &Cell{level: 0}
		}
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

		g.players[i] = &SantoriniPlayer{id: i, name: name, avatar: AvatarType(avatar), workers: [2]*Worker{}}
	}

	// Register network handler
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move SantoriniMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.applyMove(move)
			}
		})
	}

	return g
}

func (g *SantoriniGame) Reset() {
	*g = *NewSantoriniGame()
}

func (g *SantoriniGame) Update(gr *GameRoom) error {
	if IsLogoClicked() {
		gr.ReturnHome()
		return nil
	}

	// Only allow input if it's my turn (or no network client)
	isMyTurn := g.networkClient == nil || g.currentPlayer == g.myPlayerNum

	if isMyTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		boardX := int((float32(x) - g.boardOffsetX) / cellSize)
		boardY := int((float32(y) - g.boardOffsetY) / cellSize)

		if boardX >= 0 && boardX < boardSize && boardY >= 0 && boardY < boardSize {
			move := g.handleClick(boardX, boardY)

			// Send move to opponent if valid
			if move != nil && g.networkClient != nil {
				g.networkClient.SendGameMove(move)
			}
		}
	}

	return nil
}

func (g *SantoriniGame) handleClick(x, y int) *SantoriniMove {
	var workerIdx int
	if g.selectedWorker != nil {
		// Find selected worker index
		for i, w := range g.players[g.currentPlayer].workers {
			if w == g.selectedWorker {
				workerIdx = i
				break
			}
		}
	}

	move := &SantoriniMove{
		X:      x,
		Y:      y,
		Phase:  g.gamePhase,
		Worker: workerIdx,
	}

	switch g.gamePhase {
	case "place":
		if g.handlePlacement(x, y) {
			return move
		}
	case "select":
		if g.handleSelection(x, y) {
			return move
		}
	case "move":
		if g.handleMove(x, y) {
			return move
		}
	case "build":
		if g.handleBuild(x, y) {
			return move
		}
	}
	return nil
}

func (g *SantoriniGame) handlePlacement(x, y int) bool {
	if g.isOccupied(x, y) {
		return false
	}
	player := g.players[g.currentPlayer]
	workerIndex := g.placementCount / 2
	if workerIndex < 2 {
		player.workers[workerIndex] = &Worker{x: x, y: y, playerID: g.currentPlayer}
		g.placementCount++
		if g.placementCount == 4 {
			g.gamePhase = "select"
			g.currentPlayer = 0
		} else {
			g.currentPlayer = (g.currentPlayer + 1) % 2
		}
		return true
	}
	return false
}

func (g *SantoriniGame) handleSelection(x, y int) bool {
	player := g.players[g.currentPlayer]
	for _, worker := range player.workers {
		if worker != nil && worker.x == x && worker.y == y {
			g.selectedWorker = worker
			g.gamePhase = "move"
			return true
		}
	}
	return false
}

func (g *SantoriniGame) handleMove(x, y int) bool {
	if g.selectedWorker == nil {
		return false
	}
	if !g.isValidMove(g.selectedWorker, x, y) {
		return false
	}
	oldLevel := g.board[g.selectedWorker.y][g.selectedWorker.x].level
	g.selectedWorker.x = x
	g.selectedWorker.y = y
	newLevel := g.board[y][x].level
	if oldLevel < 3 && newLevel == 3 {
		g.winner = g.players[g.currentPlayer]
		g.gamePhase = "gameover"
		return true
	}
	g.gamePhase = "build"
	return true
}

func (g *SantoriniGame) handleBuild(x, y int) bool {
	if g.selectedWorker == nil {
		return false
	}
	if !g.isValidBuild(g.selectedWorker, x, y) {
		return false
	}
	g.board[y][x].level++
	g.selectedWorker = nil
	g.currentPlayer = (g.currentPlayer + 1) % 2
	g.gamePhase = "select"
	return true
}

func (g *SantoriniGame) applyMove(move SantoriniMove) {
	switch move.Phase {
	case "place":
		g.handlePlacement(move.X, move.Y)
	case "select":
		// Select the worker
		player := g.players[g.currentPlayer]
		if move.Worker < len(player.workers) && player.workers[move.Worker] != nil {
			g.selectedWorker = player.workers[move.Worker]
			g.gamePhase = "move"
		}
	case "move":
		g.handleMove(move.X, move.Y)
	case "build":
		g.handleBuild(move.X, move.Y)
	}
}

func (g *SantoriniGame) isOccupied(x, y int) bool {
	for _, player := range g.players {
		for _, worker := range player.workers {
			if worker != nil && worker.x == x && worker.y == y {
				return true
			}
		}
	}
	return false
}

func (g *SantoriniGame) isValidMove(worker *Worker, x, y int) bool {
	dx := abs(worker.x - x)
	dy := abs(worker.y - y)
	if dx > 1 || dy > 1 || (dx == 0 && dy == 0) {
		return false
	}
	if g.isOccupied(x, y) {
		return false
	}
	if g.board[y][x].level == 4 {
		return false
	}
	currentLevel := g.board[worker.y][worker.x].level
	targetLevel := g.board[y][x].level
	if targetLevel > currentLevel+1 {
		return false
	}
	return true
}

func (g *SantoriniGame) isValidBuild(worker *Worker, x, y int) bool {
	dx := abs(worker.x - x)
	dy := abs(worker.y - y)
	if dx > 1 || dy > 1 {
		return false
	}
	if g.isOccupied(x, y) {
		return false
	}
	if g.board[y][x].level >= 4 {
		return false
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *SantoriniGame) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)
	g.drawGameInfo(screen)
	g.drawBoard(screen)
	g.drawWorkers(screen)
	g.drawPlayerInfo(screen)
	if g.gamePhase == "gameover" {
		g.drawWinner(screen)
	}
}

func (g *SantoriniGame) drawGameInfo(screen *ebiten.Image) {
	titleWidth := float32(200)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "SANTORINI"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	infoWidth := float32(300)
	infoX := float32(screenWidth/2) - infoWidth/2
	vector.DrawFilledRect(screen, infoX, 70, infoWidth, 50, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, infoX, 70, infoWidth, 50, 2, color.RGBA{100, 150, 220, 255}, false)

	var phaseText string
	switch g.gamePhase {
	case "place":
		phaseText = fmt.Sprintf("%s: Place Worker", g.players[g.currentPlayer].name)
	case "select":
		phaseText = fmt.Sprintf("%s: Select Worker", g.players[g.currentPlayer].name)
	case "move":
		phaseText = fmt.Sprintf("%s: Move Worker", g.players[g.currentPlayer].name)
	case "build":
		phaseText = fmt.Sprintf("%s: Build", g.players[g.currentPlayer].name)
	}

	phaseTextX := int(infoX + (infoWidth-float32(len(phaseText)*6))/2)
	ebitenutil.DebugPrintAt(screen, phaseText, phaseTextX, 90)
}

func (g *SantoriniGame) drawBoard(screen *ebiten.Image) {
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			x := g.boardOffsetX + float32(j)*cellSize
			y := g.boardOffsetY + float32(i)*cellSize

			cellColor := color.RGBA{220, 200, 170, 255}
			vector.DrawFilledRect(screen, x, y, cellSize, cellSize, cellColor, false)
			vector.StrokeRect(screen, x, y, cellSize, cellSize, 2, color.RGBA{100, 80, 60, 255}, false)

			cell := g.board[i][j]
			g.drawBuilding(screen, x, y, cell.level)

			if g.gamePhase == "move" && g.selectedWorker != nil {
				if g.isValidMove(g.selectedWorker, j, i) {
					vector.StrokeRect(screen, x+2, y+2, cellSize-4, cellSize-4, 3, color.RGBA{100, 255, 100, 200}, false)
				}
			} else if g.gamePhase == "build" && g.selectedWorker != nil {
				if g.isValidBuild(g.selectedWorker, j, i) {
					vector.StrokeRect(screen, x+2, y+2, cellSize-4, cellSize-4, 3, color.RGBA{100, 150, 255, 200}, false)
				}
			}
		}
	}
}

func (g *SantoriniGame) drawBuilding(screen *ebiten.Image, x, y float32, level int) {
	if level == 0 {
		return
	}
	baseColors := []color.RGBA{
		{200, 200, 200, 255},
		{180, 180, 180, 255},
		{160, 160, 160, 255},
		{100, 150, 255, 255},
	}
	centerX := x + cellSize/2
	centerY := y + cellSize/2
	for l := 0; l < level && l < 4; l++ {
		size := float32(60 - l*8)
		offset := float32(l * 8)
		buildX := centerX - size/2
		buildY := centerY - size/2 - offset
		vector.DrawFilledRect(screen, buildX, buildY, size, size, baseColors[l], false)
		vector.StrokeRect(screen, buildX, buildY, size, size, 1, color.RGBA{100, 100, 100, 255}, false)
	}
}

func (g *SantoriniGame) drawWorkers(screen *ebiten.Image) {
	for _, player := range g.players {
		for _, worker := range player.workers {
			if worker != nil {
				x := g.boardOffsetX + float32(worker.x)*cellSize + cellSize/2
				y := g.boardOffsetY + float32(worker.y)*cellSize + cellSize/2

				DrawAvatar(screen, player.avatar, x-25, y-25, 1.0)

				if g.selectedWorker == worker {
					vector.StrokeRect(screen, x-27, y-27, 54, 54, 3, color.RGBA{255, 255, 100, 255}, false)
				}
			}
		}
	}
}

func (g *SantoriniGame) drawPlayerInfo(screen *ebiten.Image) {
	cardWidth := float32(320)
	edgeSpacing := float32(60) // spacing from edges

	// Calculate center gap: total width - 2 edge spacings - 2 card widths
	gapBetween := screenWidth - 2*edgeSpacing - 2*cardWidth

	for i, player := range g.players {
		var x float32
		if i == 0 {
			x = edgeSpacing
		} else {
			x = edgeSpacing + cardWidth + gapBetween
		}
		y := float32(600)
		panelColor := color.RGBA{30, 50, 80, 255}
		var borderColor color.RGBA
		if i == 0 {
			borderColor = color.RGBA{100, 150, 220, 255}
		} else {
			borderColor = color.RGBA{200, 160, 120, 255}
		}
		vector.DrawFilledRect(screen, x, y, 320, 100, panelColor, false)
		vector.StrokeRect(screen, x, y, 320, 100, 2, borderColor, false)
		DrawAvatar(screen, player.avatar, x+10, y+10, 1.5)
		ebitenutil.DebugPrintAt(screen, player.name, int(x+90), int(y+30))
		ebitenutil.DebugPrintAt(screen, player.name, int(x+91), int(y+30))
		if g.currentPlayer == i {
			ebitenutil.DebugPrintAt(screen, "Current Turn", int(x+90), int(y+50))
		}
	}
}

func (g *SantoriniGame) drawWinner(screen *ebiten.Image) {
	// Center the winner banner
	bannerWidth := float32(450)
	bannerHeight := float32(60)
	bannerX := (screenWidth - bannerWidth) / 2
	bannerY := float32(350)

	vector.DrawFilledRect(screen, bannerX, bannerY, bannerWidth, bannerHeight, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, bannerX, bannerY, bannerWidth, bannerHeight, 3, color.RGBA{100, 150, 220, 255}, false)

	starColor := color.RGBA{150, 200, 255, 200}
	vector.DrawFilledRect(screen, bannerX-20, bannerY+10, 12, 12, starColor, false)
	vector.DrawFilledRect(screen, bannerX+bannerWidth+8, bannerY+10, 12, 12, starColor, false)
	vector.DrawFilledRect(screen, bannerX-20, bannerY+40, 10, 10, starColor, false)
	vector.DrawFilledRect(screen, bannerX+bannerWidth+8, bannerY+40, 10, 10, starColor, false)

	winnerText := fmt.Sprintf("WINNER: %s", g.winner.name)
	winnerTextX := int(bannerX + (bannerWidth-float32(len(winnerText)*6))/2)
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX, int(bannerY+20))
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX+1, int(bannerY+20))
}

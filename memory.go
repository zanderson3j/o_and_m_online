package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	mem_cardWidth  = 80
	mem_cardHeight = 100
	mem_gridCols   = 6
	mem_gridRows   = 4
	mem_totalPairs = 12
)

type CardType int

const (
	CardTotoro CardType = iota
	CardNoFace
	CardCalcifer
	CardKiki
	CardChihiro
	CardHowl
	CardPonyo
	CardSophie
	CardHaku
	CardSan
	CardAshitaka
	CardSatsuki
)

type Card struct {
	cardType  CardType
	x         float32
	y         float32
	flipped   bool
	matched   bool
	index     int
}

type MemoryMove struct {
	CardIndex int `json:"card_index"`
}

type MemoryPlayer struct {
	name   string
	avatar AvatarType
	score  int
}

type MemoryGame struct {
	cards          []*Card
	flippedIndices []int
	currentPlayer  int
	players        []*MemoryPlayer
	winner         int
	gameOver       bool
	flipDelay      int
	networkClient  *NetworkClient
	myPlayerNum    int
	numPlayers     int
}

func NewMemoryGame() *MemoryGame {
	return NewMemoryGameWithNetwork(nil, 0)
}

func NewMemoryGameWithPlayers(nc *NetworkClient, playerNum int, playerData []map[string]interface{}) *MemoryGame {
	numPlayers := len(playerData)
	if numPlayers == 0 {
		numPlayers = 2
	}

	g := &MemoryGame{
		currentPlayer:  0,
		players:        make([]*MemoryPlayer, numPlayers),
		winner:         -1,
		gameOver:       false,
		flipDelay:      0,
		flippedIndices: make([]int, 0),
		networkClient:  nc,
		myPlayerNum:    playerNum,
		numPlayers:     numPlayers,
	}

	// Initialize players from server data
	for i := 0; i < numPlayers; i++ {
		name := fmt.Sprintf("Player %d", i+1)
		avatar := i % int(AvatarNumTypes)
		
		if i < len(playerData) {
			if n, ok := playerData[i]["name"].(string); ok {
				name = n
			}
			if a, ok := playerData[i]["avatar"].(float64); ok {
				avatar = int(a)
			}
		}
		
		g.players[i] = &MemoryPlayer{
			name:   name,
			avatar: AvatarType(avatar),
			score:  0,
		}
	}

	// Setup game board
	g.setupBoard(nc)

	// Register network handler for opponent moves
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move MemoryMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.flipCard(move.CardIndex)
			}
		})
	}

	return g
}

func NewMemoryGameWithNetwork(nc *NetworkClient, playerNum int) *MemoryGame {
	// Default to 2 players for backward compatibility
	g := &MemoryGame{
		currentPlayer:  0,
		players:        make([]*MemoryPlayer, 2),
		winner:         -1,
		gameOver:       false,
		flipDelay:      0,
		flippedIndices: make([]int, 0),
		networkClient:  nc,
		myPlayerNum:    playerNum,
		numPlayers:     2,
	}

	// Initialize default players
	for i := 0; i < 2; i++ {
		g.players[i] = &MemoryPlayer{
			name:   fmt.Sprintf("Player %d", i+1),
			avatar: AvatarType(i),
			score:  0,
		}
	}

	// Setup game board
	g.setupBoard(nc)

	// Register network handler for opponent moves
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move MemoryMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.flipCard(move.CardIndex)
			}
		})
	}

	return g
}

func (g *MemoryGame) setupBoard(nc *NetworkClient) {
	// Create pairs of cards
	cardTypes := make([]CardType, mem_totalPairs*2)
	for i := 0; i < mem_totalPairs; i++ {
		cardTypes[i*2] = CardType(i)
		cardTypes[i*2+1] = CardType(i)
	}

	// Shuffle cards with fixed seed for network games (both players have same layout)
	if nc != nil {
		rand.Seed(12345) // Fixed seed so both players have same board
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	rand.Shuffle(len(cardTypes), func(i, j int) {
		cardTypes[i], cardTypes[j] = cardTypes[j], cardTypes[i]
	})

	// Create card grid
	g.cards = make([]*Card, mem_gridRows*mem_gridCols)
	startX := float32(screenWidth/2) - float32(mem_gridCols*mem_cardWidth+5*(mem_gridCols-1))/2
	startY := float32(140)

	idx := 0
	for row := 0; row < mem_gridRows; row++ {
		for col := 0; col < mem_gridCols; col++ {
			g.cards[idx] = &Card{
				cardType: cardTypes[idx],
				x:        startX + float32(col*(mem_cardWidth+5)),
				y:        startY + float32(row*(mem_cardHeight+5)),
				flipped:  false,
				matched:  false,
				index:    idx,
			}
			idx++
		}
	}
}

func (g *MemoryGame) Reset() {
	*g = *NewMemoryGame()
}

func (g *MemoryGame) Update(gr *GameRoom) error {
	if IsLogoClicked() {
		gr.ReturnHome()
		return nil
	}

	if g.gameOver {
		return nil
	}

	// Handle flip delay (auto-unflip after mismatch)
	if g.flipDelay > 0 {
		g.flipDelay--
		if g.flipDelay == 0 {
			// Unflip non-matched cards
			for _, idx := range g.flippedIndices {
				if !g.cards[idx].matched {
					g.cards[idx].flipped = false
				}
			}
			g.flippedIndices = make([]int, 0)
		}
		return nil
	}

	// Only allow input if it's my turn (or if no network client)
	isMyTurn := g.networkClient == nil || g.currentPlayer == g.myPlayerNum

	// Handle card clicks
	if isMyTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		for _, card := range g.cards {
			if card.matched || card.flipped {
				continue
			}

			if float32(mx) >= card.x && float32(mx) <= card.x+mem_cardWidth &&
				float32(my) >= card.y && float32(my) <= card.y+mem_cardHeight {

				g.flipCard(card.index)

				// Send move to opponent
				if g.networkClient != nil {
					move := MemoryMove{CardIndex: card.index}
					g.networkClient.SendGameMove(move)
				}

				break
			}
		}
	}

	return nil
}

func (g *MemoryGame) flipCard(cardIndex int) {
	card := g.cards[cardIndex]
	if card.matched || card.flipped {
		return
	}

	card.flipped = true
	g.flippedIndices = append(g.flippedIndices, card.index)

	// Check if two cards are flipped
	if len(g.flippedIndices) == 2 {
		card1 := g.cards[g.flippedIndices[0]]
		card2 := g.cards[g.flippedIndices[1]]

		if card1.cardType == card2.cardType {
			// Match!
			card1.matched = true
			card2.matched = true
			g.players[g.currentPlayer].score++
			g.flippedIndices = make([]int, 0)

			// Check if game is over
			allMatched := true
			for _, c := range g.cards {
				if !c.matched {
					allMatched = false
					break
				}
			}
			if allMatched {
				g.gameOver = true
				// Find winner
				maxScore := -1
				g.winner = -1
				for i, player := range g.players {
					if player.score > maxScore {
						maxScore = player.score
						g.winner = i
					} else if player.score == maxScore {
						g.winner = -1 // Tie
					}
				}
			}
		} else {
			// No match, switch to next player and set delay
			g.currentPlayer = (g.currentPlayer + 1) % g.numPlayers
			g.flipDelay = 60 // 1 second
		}
	}
}

func (g *MemoryGame) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)
	g.drawGameInfo(screen)
	g.drawCards(screen)
	g.drawPlayerInfo(screen)

	if g.gameOver {
		g.drawWinner(screen)
	}
}

func (g *MemoryGame) drawGameInfo(screen *ebiten.Image) {
	titleWidth := float32(250)
	titleX := float32(screenWidth/2) - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleText := "MEMORY MATCH"
	titleTextX := int(titleX + (titleWidth-float32(len(titleText)*6))/2)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, titleText, titleTextX+1, 32)

	infoWidth := float32(300)
	infoX := float32(screenWidth/2) - infoWidth/2
	vector.DrawFilledRect(screen, infoX, 70, infoWidth, 50, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, infoX, 70, infoWidth, 50, 2, color.RGBA{100, 150, 220, 255}, false)

	var turnText string
	if g.gameOver {
		turnText = "Game Over!"
	} else {
		turnText = fmt.Sprintf("%s's Turn", g.players[g.currentPlayer].name)
	}

	turnTextX := int(infoX + (infoWidth-float32(len(turnText)*6))/2)
	ebitenutil.DebugPrintAt(screen, turnText, turnTextX, 90)
}

func (g *MemoryGame) drawCards(screen *ebiten.Image) {
	for _, card := range g.cards {
		if card.flipped || card.matched {
			g.drawCardFace(screen, card)
		} else {
			g.drawCardBack(screen, card)
		}
	}
}

func (g *MemoryGame) drawCardBack(screen *ebiten.Image, card *Card) {
	// Card background
	backColor := color.RGBA{30, 50, 80, 255}
	vector.DrawFilledRect(screen, card.x, card.y, mem_cardWidth, mem_cardHeight, backColor, false)
	vector.StrokeRect(screen, card.x, card.y, mem_cardWidth, mem_cardHeight, 2, color.RGBA{100, 150, 220, 255}, false)

	// Simple pattern on back
	patternColor := color.RGBA{50, 80, 120, 255}
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			x := card.x + 15 + float32(i*25)
			y := card.y + 15 + float32(j*25)
			vector.DrawFilledCircle(screen, x, y, 5, patternColor, false)
		}
	}
}

func (g *MemoryGame) drawCardFace(screen *ebiten.Image, card *Card) {
	// Card background
	faceColor := color.RGBA{240, 240, 220, 255}
	vector.DrawFilledRect(screen, card.x, card.y, mem_cardWidth, mem_cardHeight, faceColor, false)
	vector.StrokeRect(screen, card.x, card.y, mem_cardWidth, mem_cardHeight, 2, color.RGBA{100, 150, 220, 255}, false)

	// Draw character based on card type
	cx := card.x + mem_cardWidth/2
	cy := card.y + mem_cardHeight/2

	switch card.cardType {
	case CardTotoro:
		g.drawTotoro(screen, cx, cy)
	case CardNoFace:
		g.drawNoFace(screen, cx, cy)
	case CardCalcifer:
		g.drawCalcifer(screen, cx, cy)
	case CardKiki:
		g.drawKiki(screen, cx, cy)
	case CardChihiro:
		g.drawChihiro(screen, cx, cy)
	case CardHowl:
		g.drawHowl(screen, cx, cy)
	case CardPonyo:
		g.drawPonyo(screen, cx, cy)
	case CardSophie:
		g.drawSophie(screen, cx, cy)
	case CardHaku:
		g.drawHaku(screen, cx, cy)
	case CardSan:
		g.drawSan(screen, cx, cy)
	case CardAshitaka:
		g.drawAshitaka(screen, cx, cy)
	case CardSatsuki:
		g.drawSatsuki(screen, cx, cy)
	}
}

// Character drawing functions
func (g *MemoryGame) drawTotoro(screen *ebiten.Image, cx, cy float32) {
	// Gray body
	bodyColor := color.RGBA{120, 120, 120, 255}
	vector.DrawFilledCircle(screen, cx, cy, 25, bodyColor, false)
	// Belly
	bellyColor := color.RGBA{200, 200, 200, 255}
	vector.DrawFilledCircle(screen, cx, cy+5, 15, bellyColor, false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-8, cy-8, 4, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx+8, cy-8, 4, color.RGBA{0, 0, 0, 255}, false)
	// Ears
	vector.DrawFilledCircle(screen, cx-18, cy-20, 8, bodyColor, false)
	vector.DrawFilledCircle(screen, cx+18, cy-20, 8, bodyColor, false)
}

func (g *MemoryGame) drawNoFace(screen *ebiten.Image, cx, cy float32) {
	// Black body
	bodyColor := color.RGBA{40, 40, 50, 255}
	vector.DrawFilledCircle(screen, cx, cy, 28, bodyColor, false)
	// White mask
	maskColor := color.RGBA{240, 240, 240, 255}
	vector.DrawFilledCircle(screen, cx, cy-5, 18, maskColor, false)
	// Eyes (dots)
	vector.DrawFilledCircle(screen, cx-6, cy-10, 2, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-10, 2, color.RGBA{0, 0, 0, 255}, false)
	// Mouth line
	vector.StrokeLine(screen, cx-8, cy, cx+8, cy, 2, color.RGBA{0, 0, 0, 255}, false)
}

func (g *MemoryGame) drawCalcifer(screen *ebiten.Image, cx, cy float32) {
	// Orange/red flame body
	flameColors := []color.RGBA{
		{255, 150, 50, 255},
		{255, 100, 30, 255},
		{255, 200, 100, 255},
	}
	// Flame shape
	vector.DrawFilledCircle(screen, cx, cy, 20, flameColors[0], false)
	vector.DrawFilledCircle(screen, cx-8, cy-15, 12, flameColors[1], false)
	vector.DrawFilledCircle(screen, cx+8, cy-15, 12, flameColors[2], false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-6, cy-5, 3, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-5, 3, color.RGBA{0, 0, 0, 255}, false)
	// Smile
	vector.StrokeLine(screen, cx-8, cy+5, cx+8, cy+5, 2, color.RGBA{0, 0, 0, 255}, false)
}

func (g *MemoryGame) drawKiki(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 20, skinColor, false)
	// Black hair with red bow
	hairColor := color.RGBA{30, 20, 20, 255}
	vector.DrawFilledCircle(screen, cx, cy-15, 18, hairColor, false)
	// Red bow
	bowColor := color.RGBA{200, 50, 50, 255}
	vector.DrawFilledRect(screen, cx-15, cy-25, 10, 8, bowColor, false)
	vector.DrawFilledRect(screen, cx+5, cy-25, 10, 8, bowColor, false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-7, cy-3, 3, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx+7, cy-3, 3, color.RGBA{0, 0, 0, 255}, false)
}

func (g *MemoryGame) drawChihiro(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Brown hair
	hairColor := color.RGBA{100, 60, 40, 255}
	vector.DrawFilledCircle(screen, cx, cy-12, 20, hairColor, false)
	// Eyes (large anime eyes)
	eyeColor := color.RGBA{80, 50, 30, 255}
	vector.DrawFilledCircle(screen, cx-6, cy-2, 4, eyeColor, false)
	vector.DrawFilledCircle(screen, cx+6, cy-2, 4, eyeColor, false)
	// White highlights
	vector.DrawFilledCircle(screen, cx-5, cy-3, 2, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(screen, cx+7, cy-3, 2, color.RGBA{255, 255, 255, 255}, false)
}

func (g *MemoryGame) drawHowl(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Blonde hair
	hairColor := color.RGBA{220, 180, 100, 255}
	vector.DrawFilledCircle(screen, cx-10, cy-15, 15, hairColor, false)
	vector.DrawFilledCircle(screen, cx+10, cy-15, 15, hairColor, false)
	// Blue eyes
	eyeColor := color.RGBA{100, 150, 220, 255}
	vector.DrawFilledCircle(screen, cx-6, cy-2, 3, eyeColor, false)
	vector.DrawFilledCircle(screen, cx+6, cy-2, 3, eyeColor, false)
}

func (g *MemoryGame) drawPonyo(screen *ebiten.Image, cx, cy float32) {
	// Skin (pinkish)
	skinColor := color.RGBA{255, 200, 180, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Red/orange hair
	hairColor := color.RGBA{255, 100, 80, 255}
	vector.DrawFilledCircle(screen, cx, cy-12, 20, hairColor, false)
	// Large eyes
	vector.DrawFilledCircle(screen, cx-6, cy-3, 5, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-3, 5, color.RGBA{0, 0, 0, 255}, false)
	// Eye highlights
	vector.DrawFilledCircle(screen, cx-5, cy-5, 2, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledCircle(screen, cx+7, cy-5, 2, color.RGBA{255, 255, 255, 255}, false)
}

func (g *MemoryGame) drawSophie(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Gray hair (old Sophie)
	hairColor := color.RGBA{180, 180, 180, 255}
	vector.DrawFilledCircle(screen, cx, cy-12, 20, hairColor, false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-6, cy-2, 3, color.RGBA{100, 120, 140, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-2, 3, color.RGBA{100, 120, 140, 255}, false)
}

func (g *MemoryGame) drawHaku(screen *ebiten.Image, cx, cy float32) {
	// Dragon form - white/blue
	dragonColor := color.RGBA{200, 220, 255, 255}
	vector.DrawFilledCircle(screen, cx, cy, 25, dragonColor, false)
	// Green eyes
	eyeColor := color.RGBA{100, 200, 150, 255}
	vector.DrawFilledCircle(screen, cx-8, cy-5, 4, eyeColor, false)
	vector.DrawFilledCircle(screen, cx+8, cy-5, 4, eyeColor, false)
	// Scales/details
	scaleColor := color.RGBA{180, 200, 240, 255}
	vector.DrawFilledCircle(screen, cx-10, cy+10, 6, scaleColor, false)
	vector.DrawFilledCircle(screen, cx+10, cy+10, 6, scaleColor, false)
}

func (g *MemoryGame) drawSan(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Dark brown hair
	hairColor := color.RGBA{60, 40, 30, 255}
	vector.DrawFilledCircle(screen, cx, cy-12, 20, hairColor, false)
	// Red face paint marks
	paintColor := color.RGBA{180, 40, 40, 255}
	vector.StrokeLine(screen, cx-12, cy-5, cx-8, cy-5, 2, paintColor, false)
	vector.StrokeLine(screen, cx+8, cy-5, cx+12, cy-5, 2, paintColor, false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-6, cy-2, 3, color.RGBA{80, 50, 30, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-2, 3, color.RGBA{80, 50, 30, 255}, false)
}

func (g *MemoryGame) drawAshitaka(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Black hair
	hairColor := color.RGBA{30, 20, 20, 255}
	vector.DrawFilledCircle(screen, cx, cy-12, 20, hairColor, false)
	// Headband
	bandColor := color.RGBA{200, 180, 150, 255}
	vector.DrawFilledRect(screen, cx-20, cy-18, 40, 4, bandColor, false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-6, cy-2, 3, color.RGBA{80, 50, 30, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-2, 3, color.RGBA{80, 50, 30, 255}, false)
}

func (g *MemoryGame) drawSatsuki(screen *ebiten.Image, cx, cy float32) {
	// Skin
	skinColor := color.RGBA{255, 220, 190, 255}
	vector.DrawFilledCircle(screen, cx, cy, 18, skinColor, false)
	// Dark hair
	hairColor := color.RGBA{40, 30, 30, 255}
	vector.DrawFilledCircle(screen, cx, cy-12, 20, hairColor, false)
	// Yellow dress hint
	dressColor := color.RGBA{255, 220, 100, 255}
	vector.DrawFilledRect(screen, cx-15, cy+15, 30, 8, dressColor, false)
	// Eyes
	vector.DrawFilledCircle(screen, cx-6, cy-2, 4, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledCircle(screen, cx+6, cy-2, 4, color.RGBA{0, 0, 0, 255}, false)
}

func (g *MemoryGame) drawPlayerInfo(screen *ebiten.Image) {
	// Dynamic layout based on number of players
	// Memory board: starts at y=140, 4 rows of 100px cards + 3*5px spacing = 415px total
	// Board ends at y=555, so player info starts at y=570 for safety
	availableHeight := 190.0  // From y=570 to y=760
	availableWidth := float64(screenWidth) - 40  // Leave margins
	spacing := 5.0
	
	var cols, rows int
	var panelWidth, panelHeight float64
	
	// Determine optimal grid layout (with less vertical space)
	switch {
	case g.numPlayers <= 2:
		cols, rows = g.numPlayers, 1
		panelHeight = 80.0
	case g.numPlayers <= 4:
		cols, rows = g.numPlayers, 1
		panelHeight = 70.0
	case g.numPlayers <= 6:
		cols, rows = 3, 2
		panelHeight = 60.0
	case g.numPlayers <= 10:
		cols, rows = 5, 2
		panelHeight = 55.0
	case g.numPlayers <= 15:
		cols, rows = 5, 3
		panelHeight = 50.0
	default: // 16-20 players
		cols, rows = 5, 4
		panelHeight = 45.0
	}
	
	// Adjust rows based on actual players
	actualRows := (g.numPlayers + cols - 1) / cols
	if actualRows < rows {
		rows = actualRows
	}
	
	// Calculate panel width based on available space and columns
	panelWidth = (availableWidth - float64(cols-1)*spacing) / float64(cols)
	if panelWidth > 180 {
		panelWidth = 180  // Cap max width
	}
	
	// Calculate actual height based on available space
	maxPanelHeight := (availableHeight - float64(rows-1)*spacing) / float64(rows)
	if panelHeight > maxPanelHeight {
		panelHeight = maxPanelHeight
	}
	
	// Center the grid
	totalWidth := float64(cols)*panelWidth + float64(cols-1)*spacing
	startX := (float64(screenWidth) - totalWidth) / 2
	startY := 570.0  // Below the memory board with margin
	
	// Draw player panels
	for i, player := range g.players {
		col := i % cols
		row := i / cols
		x := startX + float64(col)*(panelWidth+spacing)
		y := startY + float64(row)*(panelHeight+spacing)
		g.drawPlayerPanel(screen, player, i, x, y, panelWidth, panelHeight)
	}
}

func (g *MemoryGame) drawPlayerPanel(screen *ebiten.Image, player *MemoryPlayer, index int, x, y, width, height float64) {
	panelColor := color.RGBA{30, 50, 80, 255}
	borderColor := color.RGBA{100, 150, 220, 255}
	
	// Highlight current player
	if index == g.currentPlayer && !g.gameOver {
		borderColor = color.RGBA{255, 200, 100, 255}
		panelColor = color.RGBA{50, 70, 100, 255}
	}
	
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), panelColor, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(width), float32(height), 2, borderColor, false)
	
	// Dynamic sizing based on panel height
	avatarSize := height * 0.6
	if avatarSize > 50 {
		avatarSize = 50  // Cap max size
	}
	avatarScale := float32(avatarSize / 50.0)  // Base avatar is 50x50
	
	// Position avatar with consistent padding
	avatarX := float32(x + 5)
	avatarY := float32(y + (height-avatarSize)/2)
	DrawAvatar(screen, player.avatar, avatarX, avatarY, avatarScale)
	
	// Text positioning - after avatar with padding
	textX := int(x + avatarSize + 10)
	nameY := int(y + height*0.3)
	scoreY := int(y + height*0.6)
	
	// Use smaller font for very compact layouts
	if height < 70 {
		// For very small panels, put text on single line
		combinedText := fmt.Sprintf("%s: %d pairs", player.name, player.score)
		ebitenutil.DebugPrintAt(screen, combinedText, textX, int(y+height/2-4))
		if index == g.currentPlayer && !g.gameOver {
			ebitenutil.DebugPrintAt(screen, combinedText, textX+1, int(y+height/2-4))
		}
	} else {
		// Normal two-line display
		ebitenutil.DebugPrintAt(screen, player.name, textX, nameY)
		if index == g.currentPlayer && !g.gameOver {
			ebitenutil.DebugPrintAt(screen, player.name, textX+1, nameY)
		}
		
		pairsText := fmt.Sprintf("Pairs: %d", player.score)
		ebitenutil.DebugPrintAt(screen, pairsText, textX, scoreY)
		
		// Show "Your turn!" for current player if space allows
		if index == g.currentPlayer && !g.gameOver && height >= 90 {
			ebitenutil.DebugPrintAt(screen, "Your turn!", textX, int(y+height*0.85))
		}
	}
}

func (g *MemoryGame) drawWinner(screen *ebiten.Image) {
	bannerWidth := float32(450)
	bannerHeight := float32(80)
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
	if g.winner == -1 {
		winnerText = "IT'S A TIE!"
	} else {
		winnerText = fmt.Sprintf("WINNER: %s", g.players[g.winner].name)
	}

	winnerTextX := int(bannerX + (bannerWidth-float32(len(winnerText)*6))/2)
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX, int(bannerY+20))
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX+1, int(bannerY+20))
	
	if g.winner != -1 {
		scoreText := fmt.Sprintf("Score: %d pairs!", g.players[g.winner].score)
		scoreTextX := int(bannerX + (bannerWidth-float32(len(scoreText)*6))/2)
		ebitenutil.DebugPrintAt(screen, scoreText, scoreTextX, int(bannerY+45))
	}
}

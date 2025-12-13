package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ScoreCategory int

const (
	Ones ScoreCategory = iota
	Twos
	Threes
	Fours
	Fives
	Sixes
	ThreeOfKind
	FourOfKind
	FullHouse
	SmallStraight
	LargeStraight
	Yahtzee
	Chance
	NumCategories
)

var categoryNames = []string{
	"Ones", "Twos", "Threes", "Fours", "Fives", "Sixes",
	"3 of a Kind", "4 of a Kind", "Full House",
	"Small Straight", "Large Straight", "Yahtzee", "Chance",
}

type Die struct {
	value  int
	held   bool
	x, y   float64
	width  float64
	height float64
}

type YahtzeePlayer struct {
	name       string
	avatar     AvatarType
	scores     [NumCategories]*int
	totalScore int
}

func (p *YahtzeePlayer) calculateTotal() int {
	total := 0
	upperTotal := 0
	for i := Ones; i <= Sixes; i++ {
		if p.scores[i] != nil {
			upperTotal += *p.scores[i]
		}
	}
	total += upperTotal
	if upperTotal >= 63 {
		total += 35
	}
	for i := ThreeOfKind; i < NumCategories; i++ {
		if p.scores[i] != nil {
			total += *p.scores[i]
		}
	}
	return total
}

type YahtzeeMove struct {
	Action   string `json:"action"` // "roll", "hold", "score"
	DiceIdx  int    `json:"dice_idx,omitempty"`
	Category int    `json:"category,omitempty"`
	DiceVals [5]int `json:"dice_vals,omitempty"`
}

type YahtzeeGame struct {
	dice          [5]*Die
	players       []*YahtzeePlayer
	currentPlayer int
	rollsLeft     int
	rollButton    *Button
	scoreButtons  [NumCategories]*Button
	newGameButton *Button
	rng           *rand.Rand
	networkClient *NetworkClient
	myPlayerNum   int
	numPlayers    int
	playerAvatars []AvatarType
}

func NewYahtzeeGame() *YahtzeeGame {
	return NewYahtzeeGameWithNetwork(nil, 0)
}

func NewYahtzeeGameWithPlayers(nc *NetworkClient, playerNum int, playerData []map[string]interface{}) *YahtzeeGame {
	numPlayers := len(playerData)
	if numPlayers == 0 {
		numPlayers = 2
	}

	g := &YahtzeeGame{
		players:       make([]*YahtzeePlayer, numPlayers),
		currentPlayer: 0,
		rollsLeft:     3,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		networkClient: nc,
		myPlayerNum:   playerNum,
		numPlayers:    numPlayers,
		playerAvatars: make([]AvatarType, numPlayers),
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

		g.players[i] = &YahtzeePlayer{
			name:   name,
			avatar: AvatarType(avatar),
		}
		g.playerAvatars[i] = AvatarType(avatar)
	}

	// Register network handler
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move YahtzeeMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.applyMove(move)
			}
		})
	}

	// Setup UI elements
	g.setupUI()

	return g
}

func NewYahtzeeGameWithNetwork(nc *NetworkClient, playerNum int) *YahtzeeGame {
	// Default to 2 players for offline mode
	numPlayers := 2
	if nc != nil {
		// In network mode, we'll get the actual player count from the server
		// For now, default to 2 but this will be updated when game starts
		numPlayers = 2
	}

	g := &YahtzeeGame{
		players:       make([]*YahtzeePlayer, numPlayers),
		currentPlayer: 0,
		rollsLeft:     3,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		networkClient: nc,
		myPlayerNum:   playerNum,
		numPlayers:    numPlayers,
		playerAvatars: make([]AvatarType, numPlayers),
	}

	// Initialize players with default names and avatars
	for i := 0; i < numPlayers; i++ {
		g.players[i] = &YahtzeePlayer{
			name:   fmt.Sprintf("Player %d", i+1),
			avatar: AvatarType(i % int(AvatarNumTypes)),
		}
		g.playerAvatars[i] = AvatarType(i % int(AvatarNumTypes))
	}

	// Register network handler
	if nc != nil {
		nc.RegisterHandler(MsgGameMove, func(msg Message) {
			var move YahtzeeMove
			if err := json.Unmarshal(msg.Data, &move); err == nil {
				g.applyMove(move)
			}
		})
	}

	// Setup UI elements
	g.setupUI()

	return g
}

func (g *YahtzeeGame) setupUI() {
	diceY := 150.0
	diceSpacing := 100.0
	diceStartX := 150.0
	for i := 0; i < 5; i++ {
		g.dice[i] = &Die{
			value:  0,
			held:   false,
			x:      diceStartX + float64(i)*diceSpacing,
			y:      diceY,
			width:  80,
			height: 80,
		}
	}

	thirdDieX := diceStartX + 2*diceSpacing
	g.rollButton = &Button{
		x:       thirdDieX + 40 - 75,
		y:       280,
		width:   150,
		height:  40,
		text:    "Roll Dice",
		enabled: true,
	}

	scoreX := 720.0
	scoreY := 70.0
	scoreSpacing := 38.0
	for i := 0; i < int(NumCategories); i++ {
		g.scoreButtons[i] = &Button{
			x:       scoreX,
			y:       scoreY + float64(i)*scoreSpacing,
			width:   250,
			height:  33,
			text:    categoryNames[i],
			enabled: false,
		}
	}

	// Position will be calculated dynamically in drawWinner to be centered
	g.newGameButton = &Button{
		x:       (screenWidth - 200) / 2,
		y:       0, // Will be set dynamically when drawing
		width:   200,
		height:  50,
		text:    "New Game",
		enabled: false,
	}
}

func (g *YahtzeeGame) Reset() {
	*g = *NewYahtzeeGame()
}

func (g *YahtzeeGame) Update(gr *GameRoom) error {
	// Check if logo clicked
	if IsLogoClicked() {
		gr.ReturnHome()
		return nil
	}

	// Only allow input if it's my turn (or no network client)
	isMyTurn := g.networkClient == nil || g.currentPlayer == g.myPlayerNum

	if isMyTurn && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if g.rollsLeft < 3 && g.rollsLeft > 0 {
			for i, die := range g.dice {
				if float64(x) >= die.x && float64(x) <= die.x+die.width &&
					float64(y) >= die.y && float64(y) <= die.y+die.height {
					die.held = !die.held

					// Send hold action
					if g.networkClient != nil {
						move := YahtzeeMove{
							Action:  "hold",
							DiceIdx: i,
						}
						g.networkClient.SendGameMove(move)
					}
				}
			}
		}

		if g.rollButton.enabled && g.rollButton.Contains(x, y) {
			g.rollDice()

			// Send roll action with dice values
			if g.networkClient != nil {
				var diceVals [5]int
				for i, die := range g.dice {
					diceVals[i] = die.value
				}
				move := YahtzeeMove{
					Action:   "roll",
					DiceVals: diceVals,
				}
				g.networkClient.SendGameMove(move)
			}
		}

		for i, btn := range g.scoreButtons {
			if btn.enabled && btn.Contains(x, y) {
				g.scoreCategory(ScoreCategory(i))

				// Send score action
				if g.networkClient != nil {
					move := YahtzeeMove{
						Action:   "score",
						Category: i,
					}
					g.networkClient.SendGameMove(move)
				}
			}
		}

		if g.newGameButton.enabled && g.newGameButton.Contains(x, y) {
			g.Reset()
		}
	}

	x, y := ebiten.CursorPosition()
	g.rollButton.hovered = g.rollButton.Contains(x, y)
	for _, btn := range g.scoreButtons {
		btn.hovered = btn.Contains(x, y)
	}
	g.newGameButton.hovered = g.newGameButton.Contains(x, y)

	return nil
}

func (g *YahtzeeGame) rollDice() {
	if g.rollsLeft <= 0 {
		return
	}
	for _, die := range g.dice {
		if !die.held {
			die.value = g.rng.Intn(6) + 1
		}
	}
	g.rollsLeft--
	g.enableScoreButtons()
	if g.rollsLeft == 0 {
		g.rollButton.enabled = false
	}
}

func (g *YahtzeeGame) enableScoreButtons() {
	player := g.players[g.currentPlayer]
	for i := 0; i < int(NumCategories); i++ {
		if player.scores[i] == nil {
			g.scoreButtons[i].enabled = true
		}
	}
}

func (g *YahtzeeGame) disableScoreButtons() {
	for i := 0; i < int(NumCategories); i++ {
		g.scoreButtons[i].enabled = false
	}
}

func (g *YahtzeeGame) scoreCategory(category ScoreCategory) {
	player := g.players[g.currentPlayer]
	if player.scores[category] != nil {
		return
	}
	score := g.calculateScore(category)
	player.scores[category] = &score
	player.totalScore = player.calculateTotal()
	g.disableScoreButtons()
	g.nextTurn()
}

func (g *YahtzeeGame) nextTurn() {
	allScored := true
	for _, player := range g.players {
		for i := 0; i < int(NumCategories); i++ {
			if player.scores[i] == nil {
				allScored = false
				break
			}
		}
		if !allScored {
			break
		}
	}
	if allScored {
		g.newGameButton.enabled = true
		return
	}
	g.currentPlayer = (g.currentPlayer + 1) % len(g.players)
	g.rollsLeft = 3
	g.rollButton.enabled = true
	for _, die := range g.dice {
		die.held = false
		die.value = 0
	}
}

func (g *YahtzeeGame) applyMove(move YahtzeeMove) {
	switch move.Action {
	case "hold":
		if move.DiceIdx >= 0 && move.DiceIdx < 5 {
			g.dice[move.DiceIdx].held = !g.dice[move.DiceIdx].held
		}
	case "roll":
		// Apply the dice values from opponent's roll
		for i, val := range move.DiceVals {
			if i < 5 {
				g.dice[i].value = val
			}
		}
		g.rollsLeft--
		g.enableScoreButtons()
		if g.rollsLeft == 0 {
			g.rollButton.enabled = false
		}
	case "score":
		if move.Category >= 0 && move.Category < int(NumCategories) {
			g.scoreCategory(ScoreCategory(move.Category))
		}
	}
}

func (g *YahtzeeGame) calculateScore(category ScoreCategory) int {
	values := make([]int, 5)
	for i, die := range g.dice {
		values[i] = die.value
	}
	sort.Ints(values)
	counts := make(map[int]int)
	sum := 0
	for _, v := range values {
		counts[v]++
		sum += v
	}
	switch category {
	case Ones, Twos, Threes, Fours, Fives, Sixes:
		target := int(category) + 1
		return counts[target] * target
	case ThreeOfKind:
		for _, count := range counts {
			if count >= 3 {
				return sum
			}
		}
		return 0
	case FourOfKind:
		for _, count := range counts {
			if count >= 4 {
				return sum
			}
		}
		return 0
	case FullHouse:
		hasThree, hasTwo := false, false
		for _, count := range counts {
			if count == 3 {
				hasThree = true
			}
			if count == 2 {
				hasTwo = true
			}
		}
		if hasThree && hasTwo {
			return 25
		}
		return 0
	case SmallStraight:
		straights := [][]int{{1, 2, 3, 4}, {2, 3, 4, 5}, {3, 4, 5, 6}}
		for _, straight := range straights {
			found := true
			for _, v := range straight {
				if counts[v] == 0 {
					found = false
					break
				}
			}
			if found {
				return 30
			}
		}
		return 0
	case LargeStraight:
		if (values[0] == 1 && values[1] == 2 && values[2] == 3 && values[3] == 4 && values[4] == 5) ||
			(values[0] == 2 && values[1] == 3 && values[2] == 4 && values[3] == 5 && values[4] == 6) {
			return 40
		}
		return 0
	case Yahtzee:
		for _, count := range counts {
			if count == 5 {
				return 50
			}
		}
		return 0
	case Chance:
		return sum
	}
	return 0
}

func (g *YahtzeeGame) Draw(screen *ebiten.Image, gr *GameRoom) {
	DrawForestBackground(screen)
	DrawKodamaSpirits(screen)
	DrawOMLogo(screen)

	diceStartX := 150.0
	diceSpacing := 100.0
	thirdDieCenterX := float32(diceStartX + 2*diceSpacing + 40)

	titleWidth := float32(220)
	titleX := thirdDieCenterX - titleWidth/2
	vector.DrawFilledRect(screen, titleX, 15, titleWidth, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, titleX, 15, titleWidth, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	titleTextX := int(titleX + (titleWidth-42)/2)
	ebitenutil.DebugPrintAt(screen, "YAHTZEE", titleTextX, 32)
	ebitenutil.DebugPrintAt(screen, "YAHTZEE", titleTextX+1, 32)

	player := g.players[g.currentPlayer]
	playerInfoWidth := float32(270)
	playerInfoX := thirdDieCenterX - playerInfoWidth/2
	vector.DrawFilledRect(screen, playerInfoX, 70, playerInfoWidth, 50, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, playerInfoX, 70, playerInfoWidth, 50, 2, color.RGBA{100, 150, 220, 255}, false)

	turnText := fmt.Sprintf("%s's Turn", player.name)
	turnTextX := int(playerInfoX + (playerInfoWidth-float32(len(turnText)*6))/2)
	ebitenutil.DebugPrintAt(screen, turnText, turnTextX, 82)

	rollsText := fmt.Sprintf("Rolls Left: %d", g.rollsLeft)
	rollsTextX := int(playerInfoX + (playerInfoWidth-float32(len(rollsText)*6))/2)
	ebitenutil.DebugPrintAt(screen, rollsText, rollsTextX, 100)

	for _, die := range g.dice {
		var dieColor color.Color = color.RGBA{255, 255, 255, 255}
		var borderColor color.Color = color.RGBA{100, 150, 220, 255}
		if die.held {
			dieColor = color.RGBA{200, 230, 255, 255}
			borderColor = color.RGBA{70, 120, 200, 255}
		}
		vector.DrawFilledRect(screen, float32(die.x+2), float32(die.y+2), float32(die.width), float32(die.height), color.RGBA{0, 0, 0, 20}, false)
		vector.DrawFilledRect(screen, float32(die.x), float32(die.y), float32(die.width), float32(die.height), dieColor, false)
		vector.StrokeRect(screen, float32(die.x), float32(die.y), float32(die.width), float32(die.height), 2, borderColor, false)
		g.drawDieDots(screen, die)
	}

	DrawButton(screen, g.rollButton)

	vector.DrawFilledRect(screen, 690, 15, 310, 45, color.RGBA{30, 50, 80, 255}, false)
	vector.StrokeRect(screen, 690, 15, 310, 45, 2, color.RGBA{100, 150, 220, 255}, false)
	scorecardTextX := 690 + (310-54)/2
	ebitenutil.DebugPrintAt(screen, "SCORECARD", scorecardTextX, 32)
	ebitenutil.DebugPrintAt(screen, "SCORECARD", scorecardTextX+1, 32)

	for i, btn := range g.scoreButtons {
		g.drawScoreButton(screen, btn, ScoreCategory(i))
	}

	g.drawScoreSummary(screen)

	if g.newGameButton.enabled {
		g.drawWinner(screen)
	}
}

func (g *YahtzeeGame) drawDieDots(screen *ebiten.Image, die *Die) {
	if die.value == 0 {
		return
	}
	dotRadius := float32(6)
	cx := float32(die.x + die.width/2)
	cy := float32(die.y + die.height/2)
	offset := float32(20)
	dotColor := color.RGBA{40, 40, 40, 255}
	switch die.value {
	case 1:
		vector.DrawFilledCircle(screen, cx, cy, dotRadius, dotColor, false)
	case 2:
		vector.DrawFilledCircle(screen, cx-offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy+offset, dotRadius, dotColor, false)
	case 3:
		vector.DrawFilledCircle(screen, cx-offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx, cy, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy+offset, dotRadius, dotColor, false)
	case 4:
		vector.DrawFilledCircle(screen, cx-offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx-offset, cy+offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy+offset, dotRadius, dotColor, false)
	case 5:
		vector.DrawFilledCircle(screen, cx-offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx, cy, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx-offset, cy+offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy+offset, dotRadius, dotColor, false)
	case 6:
		vector.DrawFilledCircle(screen, cx-offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy-offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx-offset, cy, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx-offset, cy+offset, dotRadius, dotColor, false)
		vector.DrawFilledCircle(screen, cx+offset, cy+offset, dotRadius, dotColor, false)
	}
}

func (g *YahtzeeGame) drawScoreButton(screen *ebiten.Image, btn *Button, category ScoreCategory) {
	player := g.players[g.currentPlayer]
	btnColor := color.RGBA{40, 60, 90, 255}
	borderColor := color.RGBA{100, 150, 220, 255}
	if player.scores[category] != nil {
		btnColor = color.RGBA{60, 100, 150, 255}
		borderColor = color.RGBA{120, 170, 230, 255}
	} else if !btn.enabled {
		btnColor = color.RGBA{30, 45, 70, 255}
		borderColor = color.RGBA{70, 100, 140, 255}
	} else if btn.hovered {
		btnColor = color.RGBA{50, 80, 120, 255}
		borderColor = color.RGBA{130, 180, 240, 255}
	}
	vector.DrawFilledRect(screen, float32(btn.x), float32(btn.y), float32(btn.width), float32(btn.height), btnColor, false)
	vector.StrokeRect(screen, float32(btn.x), float32(btn.y), float32(btn.width), float32(btn.height), 1, borderColor, false)
	textX := int(btn.x + 8)
	textY := int(btn.y + btn.height/2 - 5)
	if player.scores[category] != nil {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: %d", btn.text, *player.scores[category]), textX, textY)
	} else if btn.enabled {
		potentialScore := g.calculateScore(category)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: (%d)", btn.text, potentialScore), textX, textY)
	} else {
		ebitenutil.DebugPrintAt(screen, btn.text, textX, textY)
	}
}

func (g *YahtzeeGame) drawScoreSummary(screen *ebiten.Image) {
	// Dynamic layout based on number of players
	// Available space: from y=340 (below roll button) to y=700 (leaving space at bottom)
	availableHeight := 360.0
	// Scorecard starts at x=720, so we have from x=0 to x=700 available (with margins)
	availableWidth := 680.0 // Leave 20px margin before scorecard
	spacing := 5.0

	var cols, rows int
	var panelWidth, panelHeight float64

	// Determine optimal grid layout
	switch {
	case g.numPlayers <= 2:
		cols, rows = g.numPlayers, 1
		panelHeight = 120.0
	case g.numPlayers <= 4:
		cols, rows = g.numPlayers, 1
		panelHeight = 100.0
	case g.numPlayers <= 6:
		cols, rows = 3, 2
		panelHeight = 80.0
	case g.numPlayers <= 10:
		cols, rows = 5, 2
		panelHeight = 70.0
	case g.numPlayers <= 15:
		cols, rows = 5, 3
		panelHeight = 65.0
	default: // 16-20 players
		cols, rows = 5, 4
		panelHeight = 60.0
	}

	// Adjust rows based on actual players
	actualRows := (g.numPlayers + cols - 1) / cols
	if actualRows < rows {
		rows = actualRows
	}

	// Calculate panel width based on available space and columns
	panelWidth = (availableWidth - float64(cols-1)*spacing) / float64(cols)
	if panelWidth > 180 {
		panelWidth = 180 // Cap max width
	}

	// Calculate actual height based on available space
	maxPanelHeight := (availableHeight - float64(rows-1)*spacing) / float64(rows)
	if panelHeight > maxPanelHeight {
		panelHeight = maxPanelHeight
	}

	// Center the grid within available space (left side of screen)
	totalWidth := float64(cols)*panelWidth + float64(cols-1)*spacing
	// Center within the left 700 pixels (before scorecard)
	startX := (700.0 - totalWidth) / 2
	if startX < 10 {
		startX = 10 // Minimum margin
	}
	startY := 340.0 // Below the roll button

	// Draw player panels
	for i, player := range g.players {
		col := i % cols
		row := i / cols
		x := startX + float64(col)*(panelWidth+spacing)
		y := startY + float64(row)*(panelHeight+spacing)
		g.drawPlayerPanel(screen, player, i, x, y, panelWidth, panelHeight)
	}
}

func (g *YahtzeeGame) drawPlayerPanel(screen *ebiten.Image, player *YahtzeePlayer, index int, x, y, width, height float64) {
	panelColor := color.RGBA{30, 50, 80, 255}
	borderColor := color.RGBA{100, 150, 220, 255}

	// Highlight current player
	if index == g.currentPlayer {
		borderColor = color.RGBA{255, 220, 100, 255}
		panelColor = color.RGBA{50, 70, 100, 255}
	}

	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), panelColor, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(width), float32(height), 2, borderColor, false)

	// Dynamic sizing based on panel height
	// Avatar should take up about 60% of height, leaving room for text
	avatarSize := height * 0.6
	if avatarSize > 50 {
		avatarSize = 50 // Cap max size
	}
	avatarScale := float32(avatarSize / 50.0) // Base avatar is 50x50

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
		combinedText := fmt.Sprintf("%s: %d", player.name, player.totalScore)
		ebitenutil.DebugPrintAt(screen, combinedText, textX, int(y+height/2-4))
		if index == g.currentPlayer {
			ebitenutil.DebugPrintAt(screen, combinedText, textX+1, int(y+height/2-4))
		}
	} else {
		// Normal two-line display
		ebitenutil.DebugPrintAt(screen, player.name, textX, nameY)
		if index == g.currentPlayer {
			ebitenutil.DebugPrintAt(screen, player.name, textX+1, nameY)
		}

		scoreText := fmt.Sprintf("Score: %d", player.totalScore)
		ebitenutil.DebugPrintAt(screen, scoreText, textX, scoreY)
		if index == g.currentPlayer {
			ebitenutil.DebugPrintAt(screen, scoreText, textX+1, scoreY)
		}
	}
}

func (g *YahtzeeGame) drawWinner(screen *ebiten.Image) {
	winner := g.players[0]
	for _, player := range g.players {
		if player.totalScore > winner.totalScore {
			winner = player
		}
	}

	// Center the winner banner vertically and horizontally
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

	winnerText := fmt.Sprintf("WINNER: %s", winner.name)
	winnerTextX := int(bannerX + (bannerWidth-float32(len(winnerText)*6))/2)
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX, int(bannerY+20))
	ebitenutil.DebugPrintAt(screen, winnerText, winnerTextX+1, int(bannerY+20))

	scoreText := fmt.Sprintf("Score: %d points!", winner.totalScore)
	scoreTextX := int(bannerX + (bannerWidth-float32(len(scoreText)*6))/2)
	ebitenutil.DebugPrintAt(screen, scoreText, scoreTextX, int(bannerY+40))
	ebitenutil.DebugPrintAt(screen, scoreText, scoreTextX+1, int(bannerY+40))
}

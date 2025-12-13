// Memory Match game implementation
class MemoryGame extends GameInterface {
    constructor(gameRoom) {
        super();
        this.gameRoom = gameRoom;
        this.backButton = new Button(20, 20, 100, 40, "BACK");
        
        // Game state
        this.cards = [];
        this.flipped = [];
        this.matched = [];
        this.currentPlayer = 0;
        this.players = [];
        this.scores = [];
        this.gameOver = false;
        this.gameState = null;
        this.hoveredCard = -1;
        
        console.log("Memory game created");
    }
    
    update(inputState) {
        // Update button hover
        this.backButton.hovered = this.backButton.contains(inputState.mouseX, inputState.mouseY);
        
        // Calculate grid position (6x4 grid for 24 cards)
        const cardSize = 80;
        const cardSpacing = 10;
        const gridWidth = 6 * cardSize + 5 * cardSpacing;
        const gridHeight = 4 * cardSize + 3 * cardSpacing;
        const gridX = (SCREEN_WIDTH - gridWidth) / 2;
        const gridY = (SCREEN_HEIGHT - gridHeight) / 2 + 50;
        
        // Update hovered card
        this.hoveredCard = -1;
        for (let i = 0; i < 24; i++) {
            const row = Math.floor(i / 6);
            const col = i % 6;
            const x = gridX + col * (cardSize + cardSpacing);
            const y = gridY + row * (cardSize + cardSpacing);
            
            if (inputState.mouseX >= x && inputState.mouseX < x + cardSize &&
                inputState.mouseY >= y && inputState.mouseY < y + cardSize) {
                this.hoveredCard = i;
                break;
            }
        }
        
        // Handle clicks
        if (inputState.mouseJustPressed) {
            if (this.backButton.contains(inputState.mouseX, inputState.mouseY)) {
                this.gameRoom.returnToHome();
                return;
            }
            
            // Handle card clicks
            if (this.hoveredCard >= 0 && !this.gameOver && this.isMyTurn()) {
                this.flipCard(this.hoveredCard);
            }
        }
    }
    
    draw(screen) {
        screen.clear();
        drawForestBackground(screen);
        
        // Draw game title
        screen.drawText("MEMORY MATCH", SCREEN_WIDTH / 2 - 60, 50, 
            { r: 255, g: 255, b: 255, a: 255 }, "24px monospace");
        
        // Draw back button
        drawButton(screen, this.backButton);
        
        // Draw game state
        if (this.gameState) {
            this.drawCards(screen);
            this.drawScores(screen);
            
            // Draw turn indicator or winner
            if (this.gameOver) {
                const winner = this.getWinner();
                if (winner) {
                    screen.drawText(`${winner.name} Wins!`, SCREEN_WIDTH / 2 - 50, 100,
                        { r: 255, g: 220, b: 100, a: 255 }, "20px monospace");
                } else {
                    screen.drawText("It's a tie!", SCREEN_WIDTH / 2 - 50, 100,
                        { r: 255, g: 255, b: 100, a: 255 }, "20px monospace");
                }
            } else if (this.players.length > 0 && this.currentPlayer < this.players.length) {
                const currentPlayerName = this.players[this.currentPlayer].name || "Player " + (this.currentPlayer + 1);
                screen.drawText(`${currentPlayerName}'s Turn`, SCREEN_WIDTH / 2 - 60, 100,
                    { r: 255, g: 255, b: 100, a: 255 }, "18px monospace");
            }
        } else {
            screen.drawText("Waiting for game to start...", SCREEN_WIDTH / 2 - 120, 200,
                { r: 200, g: 200, b: 200, a: 255 }, "16px monospace");
        }
    }
    
    drawCards(screen) {
        const cardSize = 80;
        const cardSpacing = 10;
        const gridWidth = 6 * cardSize + 5 * cardSpacing;
        const gridHeight = 4 * cardSize + 3 * cardSpacing;
        const gridX = (SCREEN_WIDTH - gridWidth) / 2;
        const gridY = (SCREEN_HEIGHT - gridHeight) / 2 + 50;
        
        for (let i = 0; i < 24; i++) {
            const row = Math.floor(i / 6);
            const col = i % 6;
            const x = gridX + col * (cardSize + cardSpacing);
            const y = gridY + row * (cardSize + cardSpacing);
            
            // Skip matched cards
            if (this.matched.includes(i)) {
                continue;
            }
            
            // Card background
            let cardColor = { r: 100, g: 100, b: 200, a: 255 };
            if (i === this.hoveredCard && this.isMyTurn()) {
                cardColor = { r: 130, g: 130, b: 230, a: 255 };
            }
            
            screen.fillRect(x, y, cardSize, cardSize, cardColor);
            screen.strokeRect(x, y, cardSize, cardSize, 2, { r: 255, g: 255, b: 255, a: 255 });
            
            // Draw card content if flipped
            if (this.flipped.includes(i) && this.cards[i] !== undefined) {
                // Draw avatar or symbol based on card value
                const avatarType = this.cards[i] % AVATAR_TYPES.NUM_TYPES;
                AvatarSystem.drawAvatar(screen, avatarType, x + cardSize / 2, y + cardSize / 2, 0.8);
            } else {
                // Draw card back pattern
                screen.drawText("?", x + cardSize / 2 - 6, y + cardSize / 2 + 6,
                    { r: 255, g: 255, b: 255, a: 255 }, "24px monospace");
            }
        }
    }
    
    drawScores(screen) {
        // Draw player scores on the right side
        const scoreX = SCREEN_WIDTH - 200;
        let scoreY = 150;
        
        screen.drawText("Scores:", scoreX, scoreY, { r: 255, g: 255, b: 255, a: 255 }, "18px monospace");
        scoreY += 30;
        
        for (let i = 0; i < this.players.length && i < this.scores.length; i++) {
            const player = this.players[i];
            const score = this.scores[i];
            const isCurrentPlayer = i === this.currentPlayer;
            
            // Highlight current player
            const color = isCurrentPlayer ? 
                { r: 255, g: 255, b: 100, a: 255 } : 
                { r: 200, g: 200, b: 200, a: 255 };
            
            screen.drawText(`${player.name}: ${score}`, scoreX, scoreY, color, "14px monospace");
            scoreY += 25;
        }
    }
    
    isMyTurn() {
        const myPlayerIndex = this.players.findIndex(p => p.id === this.gameRoom.playerID);
        return myPlayerIndex >= 0 && this.currentPlayer === myPlayerIndex;
    }
    
    flipCard(index) {
        if (this.flipped.includes(index) || this.matched.includes(index)) {
            return;
        }
        
        // Send flip action to server
        this.gameRoom.sendGameMove({
            action: "flip",
            card_index: index
        });
    }
    
    getWinner() {
        if (!this.gameOver || this.scores.length === 0) return null;
        
        let maxScore = Math.max(...this.scores);
        let winners = [];
        
        for (let i = 0; i < this.scores.length; i++) {
            if (this.scores[i] === maxScore) {
                winners.push(this.players[i]);
            }
        }
        
        return winners.length === 1 ? winners[0] : null;
    }
    
    handleGameState(data) {
        this.gameState = data;
        
        if (data.cards) {
            this.cards = data.cards;
        }
        if (data.flipped) {
            this.flipped = data.flipped;
        }
        if (data.matched) {
            this.matched = data.matched;
        }
        if (data.players) {
            this.players = data.players;
        }
        if (data.scores) {
            this.scores = data.scores;
        }
        if (data.current_player !== undefined) {
            this.currentPlayer = data.current_player;
        }
        if (data.game_over !== undefined) {
            this.gameOver = data.game_over;
        }
        
        console.log("Received Memory game state:", data);
    }
    
    reset() {
        this.cards = [];
        this.flipped = [];
        this.matched = [];
        this.currentPlayer = 0;
        this.players = [];
        this.scores = [];
        this.gameOver = false;
        this.gameState = null;
        this.hoveredCard = -1;
    }
}
// Connect Four game implementation
class ConnectFourGame extends GameInterface {
    constructor(gameRoom) {
        super();
        this.gameRoom = gameRoom;
        this.backButton = new Button(20, 20, 100, 40, "BACK");
        
        // Game state
        this.board = Array(6).fill(null).map(() => Array(7).fill(0));
        this.currentPlayer = 1;
        this.players = [];
        this.gameOver = false;
        this.winner = 0;
        this.gameState = null;
        this.hoveredColumn = -1;
        
        console.log("Connect Four game created");
    }
    
    update(inputState) {
        // Update button hover
        this.backButton.hovered = this.backButton.contains(inputState.mouseX, inputState.mouseY);
        
        // Calculate board position
        const boardWidth = 420; // 7 columns * 60 pixels
        const boardHeight = 360; // 6 rows * 60 pixels
        const boardX = (SCREEN_WIDTH - boardWidth) / 2;
        const boardY = (SCREEN_HEIGHT - boardHeight) / 2 + 50;
        const cellSize = 60;
        
        // Update hovered column
        if (inputState.mouseX >= boardX && inputState.mouseX < boardX + boardWidth &&
            inputState.mouseY >= boardY && inputState.mouseY < boardY + boardHeight) {
            this.hoveredColumn = Math.floor((inputState.mouseX - boardX) / cellSize);
        } else {
            this.hoveredColumn = -1;
        }
        
        // Handle clicks
        if (inputState.mouseJustPressed) {
            if (this.backButton.contains(inputState.mouseX, inputState.mouseY)) {
                this.gameRoom.returnToHome();
                return;
            }
            
            // Handle column clicks
            if (this.hoveredColumn >= 0 && !this.gameOver && this.isMyTurn()) {
                this.makeMove(this.hoveredColumn);
            }
        }
    }
    
    draw(screen) {
        screen.clear();
        drawForestBackground(screen);
        
        // Draw game title
        screen.drawText("CONNECT FOUR", SCREEN_WIDTH / 2 - 60, 50, 
            { r: 255, g: 255, b: 255, a: 255 }, "24px monospace");
        
        // Draw back button
        drawButton(screen, this.backButton);
        
        // Draw game state
        if (this.gameState) {
            this.drawBoard(screen);
            
            // Draw turn indicator or winner
            if (this.gameOver) {
                if (this.winner > 0 && this.winner <= this.players.length) {
                    const winnerName = this.players[this.winner - 1].name || "Player " + this.winner;
                    screen.drawText(`${winnerName} Wins!`, SCREEN_WIDTH / 2 - 50, 100,
                        { r: 255, g: 220, b: 100, a: 255 }, "20px monospace");
                } else {
                    screen.drawText("Draw!", SCREEN_WIDTH / 2 - 30, 100,
                        { r: 255, g: 255, b: 100, a: 255 }, "20px monospace");
                }
            } else if (this.players.length > 0 && this.currentPlayer > 0 && this.currentPlayer <= this.players.length) {
                const currentPlayerName = this.players[this.currentPlayer - 1].name || "Player " + this.currentPlayer;
                screen.drawText(`${currentPlayerName}'s Turn`, SCREEN_WIDTH / 2 - 60, 100,
                    { r: 255, g: 255, b: 100, a: 255 }, "18px monospace");
            }
        } else {
            screen.drawText("Waiting for game to start...", SCREEN_WIDTH / 2 - 120, 200,
                { r: 200, g: 200, b: 200, a: 255 }, "16px monospace");
        }
    }
    
    drawBoard(screen) {
        const cellSize = 60;
        const boardWidth = 7 * cellSize;
        const boardHeight = 6 * cellSize;
        const boardX = (SCREEN_WIDTH - boardWidth) / 2;
        const boardY = (SCREEN_HEIGHT - boardHeight) / 2 + 50;
        
        // Draw board background
        screen.fillRect(boardX - 10, boardY - 10, boardWidth + 20, boardHeight + 20,
            { r: 0, g: 0, b: 139, a: 255 });
        
        // Draw hover indicator
        if (this.hoveredColumn >= 0 && !this.gameOver && this.isMyTurn()) {
            const hoverX = boardX + this.hoveredColumn * cellSize;
            screen.fillRect(hoverX, boardY - 30, cellSize, 20,
                { r: 255, g: 255, b: 100, a: 128 });
        }
        
        // Draw cells
        for (let row = 0; row < 6; row++) {
            for (let col = 0; col < 7; col++) {
                const x = boardX + col * cellSize + cellSize / 2;
                const y = boardY + row * cellSize + cellSize / 2;
                
                // Draw hole
                let color = { r: 50, g: 50, b: 50, a: 255 }; // Empty
                
                if (this.board[row][col] === 1) {
                    color = { r: 255, g: 0, b: 0, a: 255 }; // Player 1 - Red
                } else if (this.board[row][col] === 2) {
                    color = { r: 255, g: 255, b: 0, a: 255 }; // Player 2 - Yellow
                }
                
                screen.fillCircle(x, y, cellSize / 2 - 5, color);
            }
        }
    }
    
    isMyTurn() {
        // Check if it's the current player's turn
        const myPlayerIndex = this.players.findIndex(p => p.id === this.gameRoom.playerID);
        return myPlayerIndex >= 0 && this.currentPlayer === myPlayerIndex + 1;
    }
    
    makeMove(column) {
        // Send move to server
        this.gameRoom.sendGameMove({
            column: column
        });
    }
    
    handleGameState(data) {
        this.gameState = data;
        
        if (data.board) {
            this.board = data.board;
        }
        if (data.players) {
            this.players = data.players;
        }
        if (data.current_player !== undefined) {
            this.currentPlayer = data.current_player;
        }
        if (data.game_over !== undefined) {
            this.gameOver = data.game_over;
        }
        if (data.winner !== undefined) {
            this.winner = data.winner;
        }
        
        console.log("Received Connect Four game state:", data);
    }
    
    reset() {
        this.board = Array(6).fill(null).map(() => Array(7).fill(0));
        this.currentPlayer = 1;
        this.players = [];
        this.gameOver = false;
        this.winner = 0;
        this.gameState = null;
        this.hoveredColumn = -1;
    }
}
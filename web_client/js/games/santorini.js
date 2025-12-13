// Santorini game implementation
class SantoriniGame extends GameInterface {
    constructor(gameRoom) {
        super();
        this.gameRoom = gameRoom;
        this.backButton = new Button(20, 20, 100, 40, "BACK");
        
        // Game state
        this.board = [];
        this.currentPlayer = 0;
        this.players = [];
        this.selectedWorker = null;
        this.gamePhase = "PLACE_WORKER";
        this.gameState = null;
        
        console.log("Santorini game created");
    }
    
    update(inputState) {
        // Update button hover
        this.backButton.hovered = this.backButton.contains(inputState.mouseX, inputState.mouseY);
        
        // Handle clicks
        if (inputState.mouseJustPressed) {
            if (this.backButton.contains(inputState.mouseX, inputState.mouseY)) {
                this.gameRoom.returnToHome();
            }
            
            // TODO: Handle board clicks for worker placement/movement
        }
    }
    
    draw(screen) {
        screen.clear();
        drawForestBackground(screen);
        
        // Draw game title
        screen.drawText("SANTORINI", SCREEN_WIDTH / 2 - 50, 50, 
            { r: 255, g: 255, b: 255, a: 255 }, "24px monospace");
        
        // Draw back button
        drawButton(screen, this.backButton);
        
        // Draw game state
        if (this.gameState) {
            // TODO: Draw board, workers, current player turn
            this.drawBoard(screen);
            
            // Draw turn indicator
            if (this.players.length > 0 && this.currentPlayer < this.players.length) {
                const currentPlayerName = this.players[this.currentPlayer].name || "Player " + (this.currentPlayer + 1);
                screen.drawText(`${currentPlayerName}'s Turn`, SCREEN_WIDTH / 2 - 60, 100,
                    { r: 255, g: 255, b: 100, a: 255 }, "18px monospace");
            }
        } else {
            screen.drawText("Waiting for game to start...", SCREEN_WIDTH / 2 - 120, 200,
                { r: 200, g: 200, b: 200, a: 255 }, "16px monospace");
        }
    }
    
    drawBoard(screen) {
        // Draw 5x5 board centered on screen
        const boardSize = 400;
        const cellSize = boardSize / 5;
        const boardX = (SCREEN_WIDTH - boardSize) / 2;
        const boardY = (SCREEN_HEIGHT - boardSize) / 2 + 50;
        
        // Draw grid
        for (let i = 0; i < 5; i++) {
            for (let j = 0; j < 5; j++) {
                const x = boardX + j * cellSize;
                const y = boardY + i * cellSize;
                
                // Draw cell
                screen.strokeRect(x, y, cellSize, cellSize, 2, { r: 100, g: 100, b: 100, a: 255 });
                
                // TODO: Draw building levels and workers based on game state
            }
        }
    }
    
    handleGameState(data) {
        this.gameState = data;
        // TODO: Update game state from server data
        if (data.board) {
            this.board = data.board;
        }
        if (data.players) {
            this.players = data.players;
        }
        if (data.current_player !== undefined) {
            this.currentPlayer = data.current_player;
        }
        if (data.phase) {
            this.gamePhase = data.phase;
        }
        console.log("Received Santorini game state:", data);
    }
    
    reset() {
        this.board = [];
        this.currentPlayer = 0;
        this.players = [];
        this.selectedWorker = null;
        this.gamePhase = "PLACE_WORKER";
        this.gameState = null;
    }
}
// Yahtzee game implementation
class YahtzeeGame extends GameInterface {
    constructor(gameRoom) {
        super();
        this.gameRoom = gameRoom;
        this.backButton = new Button(20, 20, 100, 40, "BACK");
        
        // Game state
        this.dice = [0, 0, 0, 0, 0];
        this.kept = [false, false, false, false, false];
        this.rollsLeft = 3;
        this.currentPlayer = 0;
        this.players = [];
        this.gameState = null;
        
        console.log("Yahtzee game created");
    }
    
    update(inputState) {
        // Update button hover
        this.backButton.hovered = this.backButton.contains(inputState.mouseX, inputState.mouseY);
        
        // Handle clicks
        if (inputState.mouseJustPressed) {
            if (this.backButton.contains(inputState.mouseX, inputState.mouseY)) {
                this.gameRoom.returnToHome();
            }
            
            // TODO: Handle dice clicks, roll button, score selection
        }
    }
    
    draw(screen) {
        screen.clear();
        drawForestBackground(screen);
        
        // Draw game title
        screen.drawText("YAHTZEE", SCREEN_WIDTH / 2 - 40, 50, 
            { r: 255, g: 255, b: 255, a: 255 }, "24px monospace");
        
        // Draw back button
        drawButton(screen, this.backButton);
        
        // Draw game state
        if (this.gameState) {
            // TODO: Draw dice, scores, current player, etc.
            screen.drawText("Game in progress...", SCREEN_WIDTH / 2 - 80, 200,
                { r: 200, g: 200, b: 200, a: 255 }, "16px monospace");
        } else {
            screen.drawText("Waiting for game to start...", SCREEN_WIDTH / 2 - 120, 200,
                { r: 200, g: 200, b: 200, a: 255 }, "16px monospace");
        }
    }
    
    handleGameState(data) {
        this.gameState = data;
        // TODO: Update game state from server data
        console.log("Received game state:", data);
    }
    
    reset() {
        this.dice = [0, 0, 0, 0, 0];
        this.kept = [false, false, false, false, false];
        this.rollsLeft = 3;
        this.currentPlayer = 0;
        this.gameState = null;
    }
}
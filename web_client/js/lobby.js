// LobbyScreen implementation
class LobbyScreen {
    constructor(gameRoom) {
        this.gameRoom = gameRoom;
        this.buttons = [];
        this.chatMessages = [];
        this.chatInput = "";
        
        this.initializeButtons();
    }
    
    initializeButtons() {
        // Back button
        this.backButton = new Button(20, 20, 100, 40, "BACK");
        
        // Start game button (only visible for host)
        this.startButton = new Button(SCREEN_WIDTH / 2 - 100, SCREEN_HEIGHT - 100, 200, 50, "START GAME");
        
        this.buttons = [this.backButton, this.startButton];
    }
    
    update(inputState) {
        // Update button hover states
        for (const btn of this.buttons) {
            btn.hovered = btn.contains(inputState.mouseX, inputState.mouseY);
        }
        
        // Only enable start button for host with enough players
        this.startButton.enabled = this.gameRoom.isHost && 
                                  this.gameRoom.players.length >= this.getMinPlayers();
        
        // Handle clicks
        if (inputState.mouseJustPressed) {
            if (this.backButton.contains(inputState.mouseX, inputState.mouseY)) {
                this.gameRoom.returnToHome();
            } else if (this.startButton.enabled && 
                      this.startButton.contains(inputState.mouseX, inputState.mouseY)) {
                this.gameRoom.startGame();
            }
        }
    }
    
    draw(screen) {
        // Clear screen
        screen.clear();
        
        // Draw background
        drawForestBackground(screen);
        
        // Draw room info
        if (this.gameRoom.roomInfo) {
            screen.drawText(this.gameRoom.roomInfo.name, SCREEN_WIDTH / 2 - 100, 50,
                { r: 255, g: 255, b: 255, a: 255 }, "24px monospace");
            
            const gameTypeText = this.gameRoom.roomInfo.game_type.replace('_', ' ').toUpperCase();
            screen.drawText(gameTypeText, SCREEN_WIDTH / 2 - 80, 80,
                { r: 200, g: 200, b: 200, a: 255 }, "18px monospace");
        }
        
        // Draw player list
        screen.drawText("Players:", 50, 150, { r: 255, g: 255, b: 255, a: 255 }, "18px monospace");
        
        for (let i = 0; i < this.gameRoom.players.length; i++) {
            const player = this.gameRoom.players[i];
            const y = 180 + i * 40;
            
            // Draw avatar
            const avatarType = player.avatar || 0;
            this.drawAvatar(screen, avatarType, 80, y + 15, 0.8);
            
            // Draw player name
            const name = AVATAR_NAMES[avatarType];
            const isHost = i === 0;
            const hostText = isHost ? " (HOST)" : "";
            const isYou = player.id === this.gameRoom.playerID ? " (YOU)" : "";
            
            screen.drawText(`${name}${hostText}${isYou}`, 120, y + 20,
                { r: 255, g: 255, b: 255, a: 255 }, "16px monospace");
        }
        
        // Draw player count
        const minPlayers = this.getMinPlayers();
        const maxPlayers = this.getMaxPlayers();
        const currentPlayers = this.gameRoom.players.length;
        const countText = `${currentPlayers}/${maxPlayers} players (min: ${minPlayers})`;
        screen.drawText(countText, SCREEN_WIDTH - 200, 150,
            { r: 200, g: 200, b: 200, a: 255 }, "14px monospace");
        
        // Draw waiting message if not enough players
        if (currentPlayers < minPlayers) {
            const waitText = `Waiting for ${minPlayers - currentPlayers} more player${minPlayers - currentPlayers > 1 ? 's' : ''}...`;
            screen.drawText(waitText, SCREEN_WIDTH / 2 - 100, SCREEN_HEIGHT - 150,
                { r: 255, g: 200, b: 100, a: 255 }, "16px monospace");
        }
        
        // Draw buttons
        drawButton(screen, this.backButton);
        if (this.gameRoom.isHost) {
            drawButton(screen, this.startButton);
        }
    }
    
    drawAvatar(screen, avatarType, x, y, scale = 1.0) {
        // Simplified avatar drawing - matches HomeScreen for consistency
        const size = 20 * scale;
        const color = this.getAvatarColor(avatarType);
        
        screen.fillCircle(x, y, size, color);
        
        // Draw simple face
        screen.fillCircle(x - 5 * scale, y - 3 * scale, 2 * scale, { r: 0, g: 0, b: 0, a: 255 });
        screen.fillCircle(x + 5 * scale, y - 3 * scale, 2 * scale, { r: 0, g: 0, b: 0, a: 255 });
    }
    
    getAvatarColor(avatarType) {
        // Same color mapping as HomeScreen
        const colors = [
            { r: 255, g: 200, b: 150, a: 255 }, // Human - peach
            { r: 139, g: 90, b: 43, a: 255 },   // Dog - brown
            { r: 255, g: 165, b: 0, a: 255 },   // Cat - orange
            { r: 200, g: 200, b: 200, a: 255 }, // Rabbit - gray
            { r: 255, g: 220, b: 100, a: 255 }, // Giraffe - yellow
            { r: 139, g: 69, b: 19, a: 255 },   // Owl - dark brown
            { r: 100, g: 50, b: 0, a: 255 },    // Millipede - dark brown
            { r: 210, g: 180, b: 140, a: 255 }, // Puppy - tan
            { r: 255, g: 140, b: 0, a: 255 },   // Tiger - orange
            { r: 101, g: 67, b: 33, a: 255 },   // Chimpanzee - dark brown
            { r: 139, g: 90, b: 43, a: 255 },   // Platypus - brown
            { r: 205, g: 133, b: 63, a: 255 },  // Lynx - tan
            { r: 0, g: 100, b: 0, a: 255 },     // Gator - dark green
            { r: 255, g: 200, b: 100, a: 255 }, // Ocelot - light orange
            { r: 255, g: 255, b: 200, a: 255 }  // Hen - light yellow
        ];
        
        return colors[avatarType] || colors[0];
    }
    
    getMinPlayers() {
        if (!this.gameRoom.roomInfo) return 2;
        
        switch (this.gameRoom.roomInfo.game_type) {
            case 'santorini':
            case 'connect_four':
                return 2;
            case 'yahtzee':
            case 'memory':
                return 1;
            default:
                return 2;
        }
    }
    
    getMaxPlayers() {
        if (!this.gameRoom.roomInfo) return 20;
        
        switch (this.gameRoom.roomInfo.game_type) {
            case 'santorini':
            case 'connect_four':
                return 2;
            case 'yahtzee':
            case 'memory':
                return 20;
            default:
                return 20;
        }
    }
}
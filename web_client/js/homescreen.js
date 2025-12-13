// HomeScreen implementation
class HomeScreen {
    constructor(gameRoom) {
        this.gameRoom = gameRoom;
        this.gameButtons = [];
        this.avatarButtons = [];
        this.currentAvatar = 0;
        this.sparkles = [];
        this.lastSparkleTime = Date.now();
        this.sparkleEmitRate = 100; // milliseconds
        
        this.initializeButtons();
        this.initializeSparkles();
    }
    
    initializeButtons() {
        // Game selection buttons (2x2 grid)
        const gameTypes = [
            { name: "YAHTZEE", type: "yahtzee" },
            { name: "SANTORINI", type: "santorini" },
            { name: "CONNECT FOUR", type: "connect_four" },
            { name: "MEMORY MATCH", type: "memory" }
        ];
        
        const buttonWidth = 200;
        const buttonHeight = 50;
        const startX = SCREEN_WIDTH / 2 - 210;
        const startY = 300;
        const spacing = 20;
        
        for (let i = 0; i < gameTypes.length; i++) {
            const row = Math.floor(i / 2);
            const col = i % 2;
            const x = startX + col * (buttonWidth + spacing);
            const y = startY + row * (buttonHeight + spacing);
            
            const btn = new Button(x, y, buttonWidth, buttonHeight, gameTypes[i].name);
            btn.gameType = gameTypes[i].type;
            this.gameButtons.push(btn);
        }
        
        // Avatar selection buttons
        const avatarY = 500;
        const avatarSize = 40;
        const totalAvatars = AVATAR_TYPES.NUM_TYPES;
        const avatarSpacing = 5;
        const totalWidth = totalAvatars * avatarSize + (totalAvatars - 1) * avatarSpacing;
        const avatarStartX = (SCREEN_WIDTH - totalWidth) / 2;
        
        for (let i = 0; i < totalAvatars; i++) {
            const x = avatarStartX + i * (avatarSize + avatarSpacing);
            const btn = new Button(x, avatarY, avatarSize, avatarSize, "");
            btn.avatarType = i;
            this.avatarButtons.push(btn);
        }
    }
    
    initializeSparkles() {
        // Create initial sparkles around the logo
        for (let i = 0; i < 5; i++) {
            this.sparkles.push(this.createSparkle());
        }
    }
    
    createSparkle() {
        return {
            x: 75 + Math.random() * 60 - 30,
            y: 45 + Math.random() * 40 - 20,
            size: Math.random() * 3 + 2,
            lifetime: 0,
            maxLifetime: Math.random() * 1000 + 1000,
            speed: Math.random() * 0.5 + 0.1
        };
    }
    
    update(inputState) {
        // Update sparkles
        const now = Date.now();
        if (now - this.lastSparkleTime > this.sparkleEmitRate && this.sparkles.length < 10) {
            this.sparkles.push(this.createSparkle());
            this.lastSparkleTime = now;
        }
        
        // Update sparkle positions and lifetimes
        this.sparkles = this.sparkles.filter(sparkle => {
            sparkle.lifetime += 16; // ~60fps
            sparkle.y -= sparkle.speed;
            return sparkle.lifetime < sparkle.maxLifetime;
        });
        
        // Update button hover states
        for (const btn of this.gameButtons) {
            btn.hovered = btn.contains(inputState.mouseX, inputState.mouseY);
        }
        for (const btn of this.avatarButtons) {
            btn.hovered = btn.contains(inputState.mouseX, inputState.mouseY);
        }
        
        // Handle clicks
        if (inputState.mouseJustPressed) {
            console.log(`Click detected at (${inputState.mouseX}, ${inputState.mouseY})`);
            
            // Check logo click
            if (isLogoClicked(inputState.mouseX, inputState.mouseY)) {
                console.log("Logo clicked!");
            }
            
            // Check game button clicks
            for (const btn of this.gameButtons) {
                if (btn.contains(inputState.mouseX, inputState.mouseY)) {
                    console.log(`Game button clicked: ${btn.gameType}`);
                    this.gameRoom.selectGame(btn.gameType);
                    break;
                }
            }
            
            // Check avatar button clicks
            for (let i = 0; i < this.avatarButtons.length; i++) {
                const btn = this.avatarButtons[i];
                if (btn.contains(inputState.mouseX, inputState.mouseY)) {
                    this.currentAvatar = i;
                    this.gameRoom.setAvatar(i);
                    console.log(`Avatar selected: ${AVATAR_NAMES[i]}`);
                    break;
                }
            }
        }
    }
    
    draw(screen) {
        // Clear screen
        screen.clear();
        
        // Draw forest background
        drawForestBackground(screen);
        drawKodamaSpirits(screen);
        
        // Draw logo
        drawOMLogo(screen);
        
        // Draw sparkles
        for (const sparkle of this.sparkles) {
            const alpha = 1 - (sparkle.lifetime / sparkle.maxLifetime);
            const color = { 
                r: COLORS.SPARKLE.r, 
                g: COLORS.SPARKLE.g, 
                b: COLORS.SPARKLE.b, 
                a: Math.floor(alpha * 255) 
            };
            screen.fillCircle(sparkle.x, sparkle.y, sparkle.size, color);
        }
        
        // Draw title
        screen.drawText("Olive & Millie's Game Room", SCREEN_WIDTH / 2 - 120, 150, 
            { r: 255, g: 255, b: 255, a: 255 }, "20px monospace");
        
        // Draw connection status
        const statusText = this.gameRoom.isOnline() ? "ONLINE" : "OFFLINE";
        const statusColor = this.gameRoom.isOnline() ? 
            { r: 100, g: 255, b: 100, a: 255 } : 
            { r: 255, g: 100, b: 100, a: 255 };
        screen.drawText(statusText, SCREEN_WIDTH / 2 - 30, 180, statusColor, "16px monospace");
        
        // Draw instructions
        screen.drawText("Select Your Avatar:", SCREEN_WIDTH / 2 - 80, 480, 
            { r: 200, g: 200, b: 200, a: 255 }, "14px monospace");
        
        // Draw game buttons
        for (const btn of this.gameButtons) {
            drawButton(screen, btn);
        }
        
        // Draw avatar buttons
        for (let i = 0; i < this.avatarButtons.length; i++) {
            const btn = this.avatarButtons[i];
            
            // Highlight selected avatar
            if (i === this.currentAvatar) {
                screen.fillRect(btn.x - 2, btn.y - 2, btn.width + 4, btn.height + 4, 
                    { r: 255, g: 220, b: 100, a: 255 });
            }
            
            // Draw avatar
            this.drawAvatar(screen, i, btn.x + btn.width / 2, btn.y + btn.height / 2, 0.8);
            
            // Draw hover effect
            if (btn.hovered) {
                screen.strokeRect(btn.x - 1, btn.y - 1, btn.width + 2, btn.height + 2, 2,
                    { r: 255, g: 255, b: 255, a: 200 });
            }
        }
        
        // Draw selected avatar name
        screen.drawText(AVATAR_NAMES[this.currentAvatar], SCREEN_WIDTH / 2 - 50, 560,
            { r: 255, g: 255, b: 255, a: 255 }, "16px monospace");
    }
    
    drawAvatar(screen, avatarType, x, y, scale = 1.0) {
        // Simplified avatar drawing for web
        // In a full implementation, this would match the Go version's avatar drawing
        const size = 20 * scale;
        const color = this.getAvatarColor(avatarType);
        
        // Draw simple circle for now
        screen.fillCircle(x, y, size, color);
        
        // Draw simple face
        screen.fillCircle(x - 5 * scale, y - 3 * scale, 2 * scale, { r: 0, g: 0, b: 0, a: 255 });
        screen.fillCircle(x + 5 * scale, y - 3 * scale, 2 * scale, { r: 0, g: 0, b: 0, a: 255 });
    }
    
    getAvatarColor(avatarType) {
        // Simple color mapping for avatars
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
}
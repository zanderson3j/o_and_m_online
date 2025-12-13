// Avatar drawing system
class AvatarSystem {
    static drawAvatar(screen, avatarType, x, y, scale = 1.0, playerName = null) {
        // This is a placeholder for the full avatar drawing system
        // In a complete implementation, this would match the Go version's detailed avatar drawings
        
        const baseSize = 30 * scale;
        
        switch (avatarType) {
            case AVATAR_TYPES.HUMAN:
                this.drawHuman(screen, x, y, scale);
                break;
            case AVATAR_TYPES.DOG:
                this.drawDog(screen, x, y, scale);
                break;
            case AVATAR_TYPES.CAT:
                this.drawCat(screen, x, y, scale);
                break;
            case AVATAR_TYPES.RABBIT:
                this.drawRabbit(screen, x, y, scale);
                break;
            default:
                // For now, draw a simple colored circle for other avatars
                const color = this.getAvatarColor(avatarType);
                screen.fillCircle(x, y, baseSize / 2, color);
                
                // Simple eyes
                screen.fillCircle(x - 5 * scale, y - 3 * scale, 2 * scale, { r: 0, g: 0, b: 0, a: 255 });
                screen.fillCircle(x + 5 * scale, y - 3 * scale, 2 * scale, { r: 0, g: 0, b: 0, a: 255 });
        }
        
        // Draw player name if provided
        if (playerName) {
            screen.drawText(playerName, x - playerName.length * 3, y + baseSize / 2 + 10,
                { r: 255, g: 255, b: 255, a: 255 }, "12px monospace");
        }
    }
    
    static drawHuman(screen, x, y, scale) {
        const size = 30 * scale;
        
        // Head
        screen.fillCircle(x, y - size * 0.3, size * 0.4, { r: 255, g: 200, b: 150, a: 255 });
        
        // Body
        screen.fillRect(x - size * 0.3, y - size * 0.1, size * 0.6, size * 0.6,
            { r: 100, g: 150, b: 200, a: 255 });
        
        // Eyes
        screen.fillCircle(x - size * 0.15, y - size * 0.35, size * 0.06, { r: 0, g: 0, b: 0, a: 255 });
        screen.fillCircle(x + size * 0.15, y - size * 0.35, size * 0.06, { r: 0, g: 0, b: 0, a: 255 });
        
        // Smile
        screen.ctx.strokeStyle = "rgba(0, 0, 0, 1)";
        screen.ctx.lineWidth = 2 * scale;
        screen.ctx.beginPath();
        screen.ctx.arc(x, y - size * 0.25, size * 0.15, 0, Math.PI);
        screen.ctx.stroke();
    }
    
    static drawDog(screen, x, y, scale) {
        const size = 30 * scale;
        
        // Body
        screen.fillCircle(x, y, size * 0.5, { r: 139, g: 90, b: 43, a: 255 });
        
        // Head
        screen.fillCircle(x, y - size * 0.4, size * 0.35, { r: 139, g: 90, b: 43, a: 255 });
        
        // Ears
        screen.fillCircle(x - size * 0.3, y - size * 0.5, size * 0.2, { r: 100, g: 60, b: 20, a: 255 });
        screen.fillCircle(x + size * 0.3, y - size * 0.5, size * 0.2, { r: 100, g: 60, b: 20, a: 255 });
        
        // Eyes
        screen.fillCircle(x - size * 0.12, y - size * 0.4, size * 0.05, { r: 0, g: 0, b: 0, a: 255 });
        screen.fillCircle(x + size * 0.12, y - size * 0.4, size * 0.05, { r: 0, g: 0, b: 0, a: 255 });
        
        // Nose
        screen.fillCircle(x, y - size * 0.25, size * 0.08, { r: 0, g: 0, b: 0, a: 255 });
    }
    
    static drawCat(screen, x, y, scale) {
        const size = 30 * scale;
        
        // Body
        screen.fillCircle(x, y, size * 0.45, { r: 255, g: 165, b: 0, a: 255 });
        
        // Head
        screen.fillCircle(x, y - size * 0.35, size * 0.3, { r: 255, g: 165, b: 0, a: 255 });
        
        // Ears (triangular)
        screen.ctx.fillStyle = "rgba(255, 140, 0, 1)";
        screen.ctx.beginPath();
        screen.ctx.moveTo(x - size * 0.25, y - size * 0.4);
        screen.ctx.lineTo(x - size * 0.35, y - size * 0.6);
        screen.ctx.lineTo(x - size * 0.15, y - size * 0.5);
        screen.ctx.fill();
        
        screen.ctx.beginPath();
        screen.ctx.moveTo(x + size * 0.25, y - size * 0.4);
        screen.ctx.lineTo(x + size * 0.35, y - size * 0.6);
        screen.ctx.lineTo(x + size * 0.15, y - size * 0.5);
        screen.ctx.fill();
        
        // Eyes
        screen.fillCircle(x - size * 0.1, y - size * 0.35, size * 0.04, { r: 0, g: 0, b: 0, a: 255 });
        screen.fillCircle(x + size * 0.1, y - size * 0.35, size * 0.04, { r: 0, g: 0, b: 0, a: 255 });
        
        // Whiskers
        screen.ctx.strokeStyle = "rgba(0, 0, 0, 0.8)";
        screen.ctx.lineWidth = 1 * scale;
        
        // Left whiskers
        screen.ctx.beginPath();
        screen.ctx.moveTo(x - size * 0.15, y - size * 0.25);
        screen.ctx.lineTo(x - size * 0.4, y - size * 0.25);
        screen.ctx.stroke();
        
        // Right whiskers
        screen.ctx.beginPath();
        screen.ctx.moveTo(x + size * 0.15, y - size * 0.25);
        screen.ctx.lineTo(x + size * 0.4, y - size * 0.25);
        screen.ctx.stroke();
    }
    
    static drawRabbit(screen, x, y, scale) {
        const size = 30 * scale;
        
        // Body
        screen.fillCircle(x, y, size * 0.4, { r: 200, g: 200, b: 200, a: 255 });
        
        // Head
        screen.fillCircle(x, y - size * 0.3, size * 0.25, { r: 200, g: 200, b: 200, a: 255 });
        
        // Ears (long)
        screen.fillCircle(x - size * 0.15, y - size * 0.6, size * 0.1, { r: 200, g: 200, b: 200, a: 255 });
        screen.fillCircle(x + size * 0.15, y - size * 0.6, size * 0.1, { r: 200, g: 200, b: 200, a: 255 });
        screen.fillRect(x - size * 0.2, y - size * 0.5, size * 0.1, size * 0.3, { r: 200, g: 200, b: 200, a: 255 });
        screen.fillRect(x + size * 0.1, y - size * 0.5, size * 0.1, size * 0.3, { r: 200, g: 200, b: 200, a: 255 });
        
        // Eyes
        screen.fillCircle(x - size * 0.08, y - size * 0.3, size * 0.04, { r: 0, g: 0, b: 0, a: 255 });
        screen.fillCircle(x + size * 0.08, y - size * 0.3, size * 0.04, { r: 0, g: 0, b: 0, a: 255 });
        
        // Nose
        screen.fillCircle(x, y - size * 0.2, size * 0.05, { r: 255, g: 150, b: 150, a: 255 });
    }
    
    static getAvatarColor(avatarType) {
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
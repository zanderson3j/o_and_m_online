// Shared drawing functions matching Go implementation

function drawButton(screen, btn) {
    const btnColor = !btn.enabled ? COLORS.BUTTON_DISABLED : 
                     btn.hovered ? COLORS.BUTTON_HOVER : COLORS.BUTTON;
    const borderColor = !btn.enabled ? { r: 100, g: 120, b: 140, a: 255 } :
                        btn.hovered ? { r: 100, g: 150, b: 220, a: 255 } : COLORS.BORDER;
    
    // Soft shadow
    screen.fillRect(btn.x + 2, btn.y + 2, btn.width, btn.height, { r: 0, g: 0, b: 0, a: 40 });
    
    // Button background
    screen.fillRect(btn.x, btn.y, btn.width, btn.height, btnColor);
    screen.strokeRect(btn.x, btn.y, btn.width, btn.height, 2, borderColor);
    
    // Button text (centered)
    const textMetrics = screen.measureText(btn.text);
    const textX = btn.x + (btn.width - textMetrics.width) / 2;
    const textY = btn.y + btn.height / 2 + 4; // Adjust for baseline
    screen.drawText(btn.text, textX, textY);
}

function drawForestBackground(screen) {
    // Simple background for web performance
    screen.fillRect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT, { r: 30, g: 55, b: 95, a: 255 });
    
    // Ground
    screen.fillRect(0, 400, SCREEN_WIDTH, SCREEN_HEIGHT - 400, COLORS.MID_FOREST);
    
    // Simple trees
    for (let i = 0; i < 10; i++) {
        const x = i * 110 + 20;
        const y = 350 + (i % 3) * 20;
        
        // Tree trunk
        screen.fillRect(x + 25, y + 80, 20, 60, { r: 50, g: 30, b: 20, a: 255 });
        
        // Tree top (triangle)
        screen.ctx.fillStyle = `rgba(${COLORS.DARK_FOREST.r}, ${COLORS.DARK_FOREST.g}, ${COLORS.DARK_FOREST.b}, 1)`;
        screen.ctx.beginPath();
        screen.ctx.moveTo(x + 35, y);
        screen.ctx.lineTo(x, y + 80);
        screen.ctx.lineTo(x + 70, y + 80);
        screen.ctx.closePath();
        screen.ctx.fill();
    }
}

function drawKodamaSpirits(screen) {
    // Skip for web performance, or draw simple versions
    const positions = [
        [100, 400], [250, 350], [450, 380], [650, 420],
        [850, 360], [150, 500], [550, 480], [900, 520]
    ];
    
    for (const [x, y] of positions) {
        // Simple glow
        screen.fillCircle(x, y, 12, COLORS.KODAMA_GLOW);
        // Kodama body
        screen.fillCircle(x, y, 8, COLORS.KODAMA);
        // Eyes
        screen.fillRect(x - 3, y - 2, 2, 3, { r: 40, g: 40, b: 40, a: 255 });
        screen.fillRect(x + 1, y - 2, 2, 3, { r: 40, g: 40, b: 40, a: 255 });
    }
}

function drawOMLogo(screen) {
    const x = 15;
    const y = 15;
    const scale = 3;
    
    // Simplified O&M logo
    // O
    screen.fillRect(x + 4 * scale, y + 3 * scale, 2 * scale, 14 * scale, COLORS.O_COLOR);
    screen.fillRect(x + 12 * scale, y + 3 * scale, 2 * scale, 14 * scale, COLORS.O_COLOR);
    screen.fillRect(x + 6 * scale, y + 3 * scale, 6 * scale, 2 * scale, COLORS.O_COLOR);
    screen.fillRect(x + 6 * scale, y + 15 * scale, 6 * scale, 2 * scale, COLORS.O_COLOR);
    
    // &
    screen.fillRect(x + 19 * scale, y + 8 * scale, 8 * scale, 2 * scale, { r: 200, g: 200, b: 200, a: 255 });
    screen.fillRect(x + 22 * scale, y + 4 * scale, 2 * scale, 12 * scale, { r: 200, g: 200, b: 200, a: 255 });
    
    // M
    screen.fillRect(x + 31 * scale, y + 3 * scale, 2 * scale, 14 * scale, COLORS.O_COLOR);
    screen.fillRect(x + 41 * scale, y + 3 * scale, 2 * scale, 14 * scale, COLORS.O_COLOR);
    screen.fillRect(x + 33 * scale, y + 3 * scale, 8 * scale, 2 * scale, COLORS.O_COLOR);
    screen.fillRect(x + 36 * scale, y + 7 * scale, 2 * scale, 4 * scale, COLORS.O_COLOR);
}

function isLogoClicked(x, y) {
    // Logo is at position (15, 15) with scale 3, roughly 120x60 pixels
    return x >= 15 && x <= 135 && y >= 15 && y <= 75;
}

// Helper function to get current timestamp
function getTimestamp() {
    return new Date().toISOString();
}

// Debug text function
function debugPrint(screen, text, x, y, color = { r: 255, g: 255, b: 255, a: 255 }) {
    screen.drawText(text, x, y, color, "12px monospace");
}
// Type definitions matching Go structures

class Button {
    constructor(x, y, width, height, text, enabled = true) {
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.text = text;
        this.enabled = enabled;
        this.hovered = false;
    }
    
    contains(x, y) {
        return x >= this.x && x <= this.x + this.width && 
               y >= this.y && y <= this.y + this.height;
    }
}

class Message {
    constructor(type, data = null) {
        this.type = type;
        this.data = data;
        this.timestamp = new Date().toISOString();
        this.player_id = null; // Set by server
        this.room_id = null; // Set by context
        this.game_type = null; // Set by context
    }
}

class RoomInfo {
    constructor(data) {
        this.id = data.id;
        this.name = data.name;
        this.game_type = data.game_type;
        this.players = data.players;
        this.max_players = data.max_players;
        this.started = data.started;
    }
}

// Base class for all games
class GameInterface {
    update(gameRoom) {
        throw new Error("update() must be implemented");
    }
    
    draw(screen, gameRoom) {
        throw new Error("draw() must be implemented");
    }
    
    reset() {
        throw new Error("reset() must be implemented");
    }
    
    handleClick(x, y) {
        // Override in subclasses
    }
    
    handleMouseMove(x, y) {
        // Override in subclasses
    }
}

// Input state
class InputState {
    constructor() {
        this.mouseX = 0;
        this.mouseY = 0;
        this.mousePressed = false;
        this.mouseJustPressed = false;
        this.mouseJustReleased = false;
        this.lastMousePressed = false;
    }
    
    update(pressed) {
        this.mouseJustPressed = pressed && !this.lastMousePressed;
        this.mouseJustReleased = !pressed && this.lastMousePressed;
        this.mousePressed = pressed;
        this.lastMousePressed = pressed;
    }
}

// Drawing context wrapper
class Screen {
    constructor(canvas, ctx) {
        this.canvas = canvas;
        this.ctx = ctx;
        this.width = canvas.width;
        this.height = canvas.height;
    }
    
    clear() {
        this.ctx.clearRect(0, 0, this.width, this.height);
    }
    
    fillRect(x, y, width, height, color) {
        this.ctx.fillStyle = `rgba(${color.r}, ${color.g}, ${color.b}, ${color.a / 255})`;
        this.ctx.fillRect(x, y, width, height);
    }
    
    strokeRect(x, y, width, height, lineWidth, color) {
        this.ctx.strokeStyle = `rgba(${color.r}, ${color.g}, ${color.b}, ${color.a / 255})`;
        this.ctx.lineWidth = lineWidth;
        this.ctx.strokeRect(x, y, width, height);
    }
    
    fillCircle(x, y, radius, color) {
        this.ctx.fillStyle = `rgba(${color.r}, ${color.g}, ${color.b}, ${color.a / 255})`;
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, Math.PI * 2);
        this.ctx.fill();
    }
    
    strokeCircle(x, y, radius, lineWidth, color) {
        this.ctx.strokeStyle = `rgba(${color.r}, ${color.g}, ${color.b}, ${color.a / 255})`;
        this.ctx.lineWidth = lineWidth;
        this.ctx.beginPath();
        this.ctx.arc(x, y, radius, 0, Math.PI * 2);
        this.ctx.stroke();
    }
    
    drawText(text, x, y, color = { r: 255, g: 255, b: 255, a: 255 }, font = "12px monospace") {
        this.ctx.font = font;
        this.ctx.fillStyle = `rgba(${color.r}, ${color.g}, ${color.b}, ${color.a / 255})`;
        this.ctx.fillText(text, x, y);
    }
    
    measureText(text, font = "12px monospace") {
        this.ctx.font = font;
        return this.ctx.measureText(text);
    }
}
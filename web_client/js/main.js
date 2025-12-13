// Main game controller
class GameRoom {
    constructor() {
        this.canvas = document.getElementById('gameCanvas');
        this.ctx = this.canvas.getContext('2d');
        this.screen = new Screen(this.canvas, this.ctx);
        this.inputState = new InputState();
        
        this.networkClient = null;
        this.currentScreen = null;
        this.homeScreen = null;
        this.lobbyScreen = null;
        this.gameScreen = null;
        
        this.roomInfo = null;
        this.players = [];
        this.playerID = null;
        this.isHost = false;
        
        this.animationFrameId = null;
        this.lastFrameTime = 0;
        this.targetFPS = 60;
        this.frameInterval = 1000 / this.targetFPS;
        
        this.setupInputHandlers();
    }
    
    async start() {
        console.log('Starting game...');
        
        // Initialize home screen first
        this.homeScreen = new HomeScreen(this);
        this.currentScreen = this.homeScreen;
        
        // Try to connect to server
        try {
            this.networkClient = new NetworkClient(SERVER_URL);
            await this.networkClient.connect();
            
            // Register message handlers
            this.setupNetworkHandlers();
            
            console.log('Connected to server!');
        } catch (error) {
            console.error('Failed to connect to server:', error);
            // Continue in offline mode
        }
        
        // Start game loop
        this.gameLoop(0);
    }
    
    setupInputHandlers() {
        // Mouse move
        this.canvas.addEventListener('mousemove', (e) => {
            const rect = this.canvas.getBoundingClientRect();
            this.inputState.mouseX = e.clientX - rect.left;
            this.inputState.mouseY = e.clientY - rect.top;
        });
        
        // Mouse down
        this.canvas.addEventListener('mousedown', (e) => {
            if (e.button === 0) { // Left click only
                this.inputState.mousePressed = true;
                this.inputState.mouseJustPressed = true;
            }
        });
        
        // Mouse up
        this.canvas.addEventListener('mouseup', (e) => {
            if (e.button === 0) {
                this.inputState.mousePressed = false;
                this.inputState.mouseJustReleased = true;
            }
        });
        
        // Touch events for mobile
        this.canvas.addEventListener('touchstart', (e) => {
            e.preventDefault();
            const touch = e.touches[0];
            const rect = this.canvas.getBoundingClientRect();
            this.inputState.mouseX = touch.clientX - rect.left;
            this.inputState.mouseY = touch.clientY - rect.top;
            this.inputState.mousePressed = true;
            this.inputState.mouseJustPressed = true;
        });
        
        this.canvas.addEventListener('touchmove', (e) => {
            e.preventDefault();
            const touch = e.touches[0];
            const rect = this.canvas.getBoundingClientRect();
            this.inputState.mouseX = touch.clientX - rect.left;
            this.inputState.mouseY = touch.clientY - rect.top;
        });
        
        this.canvas.addEventListener('touchend', (e) => {
            e.preventDefault();
            this.inputState.mousePressed = false;
            this.inputState.mouseJustReleased = true;
        });
        
        // Prevent context menu on right click
        this.canvas.addEventListener('contextmenu', (e) => {
            e.preventDefault();
        });
    }
    
    setupNetworkHandlers() {
        // Connected message
        this.networkClient.registerHandler('connected', (msg) => {
            this.playerID = msg.PlayerID || msg.player_id;
            console.log(`Connected as player: ${this.playerID}`);
        });
        
        // Room created
        this.networkClient.registerHandler('room_created', (msg) => {
            const roomID = msg.RoomID || msg.room_id;
            console.log(`Room created: ${roomID}`);
            this.joinRoom(roomID);
        });
        
        // Player joined (including self)
        this.networkClient.registerHandler('player_joined', (msg) => {
            const playerID = msg.PlayerID || msg.player_id;
            const roomID = msg.RoomID || msg.room_id;
            console.log(`Player ${playerID} joined room ${roomID}`);
            
            // Check if it's us who joined
            if (playerID === this.playerID) {
                // We joined a room!
                this.networkClient.currentRoom = roomID;
                
                // Create basic room info from available data
                this.roomInfo = {
                    id: roomID,
                    name: `Room ${roomID}`,
                    game_type: msg.GameType || msg.game_type || this.pendingGameType,
                    players: this.players,
                    max_players: 20,
                    started: false
                };
                
                console.log(`We joined room: ${msg.room_id}`);
                this.switchToLobby();
            }
        });
        
        // Player list update
        this.networkClient.registerHandler(MSG_PLAYER_LIST, (msg) => {
            if (msg.data && msg.data.players) {
                this.players = msg.data.players;
                this.isHost = this.players.length > 0 && this.players[0].id === this.playerID;
                console.log(`Updated player list: ${this.players.length} players`);
            }
        });
        
        // Game started
        this.networkClient.registerHandler(MSG_START_GAME, (msg) => {
            console.log('Game started!');
            if (this.roomInfo) {
                this.switchToGame();
            }
        });
        
        // Game state update
        this.networkClient.registerHandler(MSG_GAME_STATE, (msg) => {
            if (this.gameScreen && this.gameScreen.update) {
                this.gameScreen.handleGameState(msg.data);
            }
        });
        
        // Error messages
        this.networkClient.registerHandler(MSG_ERROR, (msg) => {
            if (msg.data && msg.data.error) {
                console.error(`Server error: ${msg.data.error}`);
                alert(`Error: ${msg.data.error}`);
            }
        });
    }
    
    gameLoop(timestamp) {
        // Calculate delta time
        const deltaTime = timestamp - this.lastFrameTime;
        
        // Only update if enough time has passed
        if (deltaTime >= this.frameInterval) {
            // Update current screen
            if (this.currentScreen) {
                this.currentScreen.update(this.inputState);
                this.currentScreen.draw(this.screen);
            }
            
            // Reset just pressed/released states for next frame
            this.inputState.mouseJustPressed = false;
            this.inputState.mouseJustReleased = false;
            
            this.lastFrameTime = timestamp - (deltaTime % this.frameInterval);
        }
        
        // Continue loop
        this.animationFrameId = requestAnimationFrame((t) => this.gameLoop(t));
    }
    
    isOnline() {
        return this.networkClient && this.networkClient.isConnected();
    }
    
    async createRoom(gameType) {
        if (!this.isOnline()) {
            alert('You must be online to create a room!');
            return;
        }
        
        // Store the game type for when we receive the response
        this.pendingGameType = gameType;
        
        const roomName = `${AVATAR_NAMES[this.homeScreen.currentAvatar]}'s ${gameType.replace('_', ' ').toUpperCase()} Room`;
        
        try {
            await this.networkClient.createRoom(gameType, roomName);
        } catch (error) {
            console.error('Failed to create room:', error);
            alert('Failed to create room. Please try again.');
        }
    }
    
    async joinRoom(roomID) {
        if (!this.isOnline()) {
            alert('You must be online to join a room!');
            return;
        }
        
        try {
            await this.networkClient.joinRoom(roomID);
        } catch (error) {
            console.error('Failed to join room:', error);
            alert('Failed to join room. Please try again.');
        }
    }
    
    async setAvatar(avatarType) {
        if (this.isOnline()) {
            try {
                await this.networkClient.setAvatar(avatarType);
            } catch (error) {
                console.error('Failed to set avatar:', error);
            }
        }
    }
    
    switchToLobby() {
        if (!this.lobbyScreen) {
            this.lobbyScreen = new LobbyScreen(this);
        }
        this.currentScreen = this.lobbyScreen;
    }
    
    switchToGame() {
        // Create appropriate game screen based on game type
        switch (this.roomInfo.game_type) {
            case 'yahtzee':
                this.gameScreen = new YahtzeeGame(this);
                break;
            case 'santorini':
                this.gameScreen = new SantoriniGame(this);
                break;
            case 'connect_four':
                this.gameScreen = new ConnectFourGame(this);
                break;
            case 'memory':
                this.gameScreen = new MemoryGame(this);
                break;
            default:
                console.error(`Unknown game type: ${this.roomInfo.game_type}`);
                return;
        }
        
        this.currentScreen = this.gameScreen;
    }
    
    returnToHome() {
        // Leave room if in one
        if (this.roomInfo && this.isOnline()) {
            this.networkClient.leaveRoom();
        }
        
        this.roomInfo = null;
        this.players = [];
        this.isHost = false;
        this.gameScreen = null;
        this.currentScreen = this.homeScreen;
    }
    
    async startGame() {
        if (!this.isHost) {
            alert('Only the host can start the game!');
            return;
        }
        
        if (!this.isOnline()) {
            alert('You must be online to start the game!');
            return;
        }
        
        try {
            await this.networkClient.startGame();
        } catch (error) {
            console.error('Failed to start game:', error);
            alert('Failed to start game. Please try again.');
        }
    }
    
    async sendGameMove(moveData) {
        if (!this.isOnline()) {
            console.error('Cannot send move - not online');
            return;
        }
        
        try {
            await this.networkClient.sendGameMove(moveData);
        } catch (error) {
            console.error('Failed to send game move:', error);
        }
    }
}
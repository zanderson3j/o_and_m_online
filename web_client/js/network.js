// NetworkClient for WebSocket communication
class NetworkClient {
    constructor(serverURL) {
        this.serverURL = serverURL;
        this.ws = null;
        this.playerID = null;
        this.currentRoom = null;
        this.rooms = [];
        this.msgHandlers = new Map();
        this.connected = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 2000; // 2 seconds
    }
    
    async connect() {
        return new Promise((resolve, reject) => {
            console.log(`Attempting to connect to: ${this.serverURL}`);
            
            try {
                this.ws = new WebSocket(this.serverURL);
                
                this.ws.onopen = () => {
                    console.log('WebSocket connected!');
                    this.connected = true;
                    this.reconnectAttempts = 0;
                    resolve();
                };
                
                this.ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                    if (!this.connected) {
                        reject(new Error('Failed to connect'));
                    }
                };
                
                this.ws.onmessage = (event) => {
                    try {
                        const msg = JSON.parse(event.data);
                        console.log('Received message:', msg);
                        this.handleMessage(msg);
                    } catch (error) {
                        console.error('Error parsing message:', error);
                    }
                };
                
                this.ws.onclose = (event) => {
                    console.log(`WebSocket closed: code=${event.code}, reason=${event.reason}`);
                    this.connected = false;
                    
                    // Attempt to reconnect if not a normal closure
                    if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                        this.reconnectAttempts++;
                        console.log(`Reconnecting in ${this.reconnectDelay}ms... (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
                        setTimeout(() => this.connect(), this.reconnectDelay);
                    }
                };
                
                // Timeout connection attempt after 5 seconds
                setTimeout(() => {
                    if (!this.connected) {
                        this.ws.close();
                        reject(new Error('Connection timeout'));
                    }
                }, 5000);
                
            } catch (error) {
                reject(error);
            }
        });
    }
    
    handleMessage(msg) {
        // Handle special messages
        switch (msg.type) {
            case 'connected':
                this.playerID = msg.PlayerID || msg.player_id;
                console.log(`Connected as player ${this.playerID}`);
                break;
                
            case MSG_ROOM_LIST:
                if (msg.data && msg.data.rooms) {
                    this.rooms = msg.data.rooms;
                    console.log(`Updated room list: ${this.rooms.length} rooms`);
                }
                break;
                
            case MSG_ERROR:
                if (msg.data && msg.data.error) {
                    console.error(`Server error: ${msg.data.error}`);
                }
                break;
        }
        
        // Call registered handler if exists
        const handler = this.msgHandlers.get(msg.type);
        if (handler) {
            handler(msg);
        }
    }
    
    registerHandler(msgType, handler) {
        this.msgHandlers.set(msgType, handler);
    }
    
    sendMessage(msg) {
        if (!this.connected) {
            console.error(`Cannot send message ${msg.type} - not connected!`);
            return Promise.reject(new Error('Not connected'));
        }
        
        // Add timestamp if not present
        if (!msg.timestamp) {
            msg.timestamp = new Date().toISOString();
        }
        
        console.log(`Sending message type=${msg.type}`);
        
        try {
            this.ws.send(JSON.stringify(msg));
            return Promise.resolve();
        } catch (error) {
            console.error('Error sending message:', error);
            return Promise.reject(error);
        }
    }
    
    createRoom(gameType, roomName) {
        return this.sendMessage({
            type: MSG_CREATE_ROOM,
            data: {
                game_type: gameType,
                room_name: roomName
            }
        });
    }
    
    joinRoom(roomID) {
        return this.sendMessage({
            type: MSG_JOIN_ROOM,
            data: {
                room_id: roomID
            }
        });
    }
    
    leaveRoom() {
        return this.sendMessage({
            type: MSG_LEAVE_ROOM
        });
    }
    
    startGame() {
        return this.sendMessage({
            type: MSG_START_GAME
        });
    }
    
    sendGameMove(moveData) {
        return this.sendMessage({
            type: MSG_GAME_MOVE,
            data: moveData
        });
    }
    
    setAvatar(avatarType) {
        return this.sendMessage({
            type: MSG_SET_AVATAR,
            data: {
                avatar: avatarType
            }
        });
    }
    
    getRooms() {
        return this.rooms;
    }
    
    getPlayerID() {
        return this.playerID;
    }
    
    getCurrentRoom() {
        return this.currentRoom;
    }
    
    isConnected() {
        return this.connected;
    }
    
    close() {
        if (this.ws) {
            this.ws.close();
            this.connected = false;
        }
    }
}
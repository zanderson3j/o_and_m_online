// Game constants matching Go implementation
const SCREEN_WIDTH = 1024;
const SCREEN_HEIGHT = 768;
// const SERVER_URL = "wss://o-and-m-online.onrender.com/ws";
const SERVER_URL = "ws://localhost:8080/ws"; // For local testing

// Message types matching Go
const MSG_JOIN_LOBBY = "join_lobby";
const MSG_LEAVE_LOBBY = "leave_lobby";
const MSG_CREATE_ROOM = "create_room";
const MSG_JOIN_ROOM = "join_room";
const MSG_LEAVE_ROOM = "leave_room";
const MSG_START_GAME = "start_game";
const MSG_GAME_MOVE = "game_move";
const MSG_GAME_STATE = "game_state";
const MSG_PLAYER_LIST = "player_list";
const MSG_ROOM_LIST = "room_list";
const MSG_ERROR = "error";
const MSG_CHAT = "chat";
const MSG_SET_AVATAR = "set_avatar";
const MSG_PLAYER_UPDATE = "player_update";

// Avatar types
const AVATAR_TYPES = {
    HUMAN: 0,
    DOG: 1,
    CAT: 2,
    RABBIT: 3,
    GIRAFFE: 4,
    OWL: 5,
    MILLIPEDE: 6,
    PUPPY: 7,
    TIGER: 8,
    CHIMPANZEE: 9,
    PLATYPUS: 10,
    LYNX: 11,
    GATOR: 12,
    OCELOT: 13,
    HEN: 14,
    NUM_TYPES: 15
};

// Avatar names matching Go
const AVATAR_NAMES = [
    "Human", "Teddy", "Kaycat", "Zach Rabbit", "Kiraffe", "Owlive", "Milliepede", 
    "Sweet Puppy Paw", "Tygler", "Chimpancici", "Papapus", "Kaitlynx", 
    "Reagator", "Ocelivia", "Hen-ry"
];

// Colors
const COLORS = {
    // UI Colors
    BUTTON: { r: 100, g: 150, b: 220, a: 255 },
    BUTTON_HOVER: { r: 130, g: 180, b: 240, a: 255 },
    BUTTON_DISABLED: { r: 140, g: 160, b: 180, a: 255 },
    BORDER: { r: 70, g: 120, b: 190, a: 255 },
    PANEL_BG: { r: 30, g: 50, b: 80, a: 255 },
    
    // Forest background
    SKY_TOP: { r: 20, g: 40, b: 80, a: 255 },
    SKY_BOTTOM: { r: 40, g: 70, b: 110, a: 255 },
    DARK_FOREST: { r: 15, g: 35, b: 60, a: 255 },
    MID_FOREST: { r: 20, g: 80, b: 80, a: 255 },
    FRONT_FOREST: { r: 30, g: 120, b: 100, a: 255 },
    
    // Kodama
    KODAMA: { r: 200, g: 255, b: 220, a: 200 },
    KODAMA_GLOW: { r: 150, g: 255, b: 200, a: 100 },
    
    // Logo colors
    O_COLOR: { r: 255, g: 220, b: 150, a: 255 },
    SPARKLE: { r: 255, g: 255, b: 200, a: 255 }
};
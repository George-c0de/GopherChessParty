export const config = {
    apiBaseUrl: process.env.REACT_APP_API_URL || 'http://127.0.0.1:8000/v1',
    wsBaseUrl: process.env.REACT_APP_WS_URL || 'ws://127.0.0.1:8000/v1',
    endpoints: {
        login: '/login',
        register: '/register',
        user: '/users/me',
        chess: '/chess/',
        game: {
            move: '/api/game/move',
            search: '/ws/search',
            play: '/ws/game'
        },
        auth: {
            refresh: '/refresh',
            validate: '/validate'
        }
    }
} as const; 
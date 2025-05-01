import React, { useState, useEffect, useCallback, useRef } from 'react';
import styled, { keyframes, css } from 'styled-components';
import { useNavigate, useParams } from 'react-router-dom';
import { config } from '../config';

const pieces = {
    white: {
        king: '♔', queen: '♕', rook: '♖', bishop: '♗', knight: '♘', pawn: '♙'
    },
    black: {
        king: '♚', queen: '♛', rook: '♜', bishop: '♝', knight: '♞', pawn: '♟'
    }
};

interface Piece {
    color: 'white' | 'black';
    type: 'king' | 'queen' | 'rook' | 'bishop' | 'knight' | 'pawn';
}

interface Position {
    row: number;
    col: number;
}

const BoardContainer = styled.div`
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f5f5f5;
    position: relative;
`;

const Overlay = styled.div<{ isGameStarted: boolean }>`
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1;
    opacity: ${props => props.isGameStarted ? 0 : 1};
    pointer-events: ${props => props.isGameStarted ? 'none' : 'auto'};
    transition: opacity 0.3s ease;
`;

const StartButton = styled.button`
    padding: 15px 30px;
    font-size: 20px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s;

    &:hover {
        background-color: #45a049;
    }

    &:disabled {
        background-color: #cccccc;
        cursor: not-allowed;
    }
`;

const StatusMessage = styled.div`
    color: white;
    font-size: 18px;
    margin-top: 10px;
    text-align: center;
`;

const Board = styled.table`
    border-collapse: collapse;
    background-color: #fff;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
`;

const Cell = styled.td<{ isLight: boolean; isSelected: boolean }>`
    width: 60px;
    height: 60px;
    text-align: center;
    font-size: 40px;
    cursor: pointer;
    background-color: ${props => props.isLight ? '#f0d9b5' : '#b58863'};
    border: 1px solid #999;
    transition: background-color 0.2s;

    &:hover {
        background-color: ${props => props.isLight ? '#e8d0a9' : '#a67b5b'};
    }

    ${props => props.isSelected && `
        background-color: #7b61ff;
        &:hover {
            background-color: #6b51ef;
        }
    `}
`;

const SideNumber = styled.td`
    width: 30px;
    text-align: center;
    font-weight: bold;
    color: #666;
`;

const CapturesContainer = styled.div`
    position: fixed;
    top: 20px;
    right: 20px;
    background: white;
    padding: 15px;
    border-radius: 5px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
`;

const LogoutButton = styled.button`
    position: fixed;
    top: 20px;
    left: 20px;
    padding: 10px 20px;
    background-color: #dc3545;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.2s;
    z-index: 2;

    &:hover {
        background-color: #c82333;
    }
`;

const PlayerInfo = styled.div`
    position: fixed;
    background: white;
    padding: 15px;
    border-radius: 5px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    gap: 10px;
`;

const WhitePlayerInfo = styled(PlayerInfo)`
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
`;

const BlackPlayerInfo = styled(PlayerInfo)`
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
`;

const PlayerAvatar = styled.div`
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: #e0e0e0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    color: #666;
`;

interface User {
    id: string;
    name: string;
    email: string;
}

interface GameInfo {
    ID: string;
    CreatedAt: string;
    Result: string;
    Status: string;
    WhiteUser: User;
    BlackUser: User;
}

const fadeIn = keyframes`
    from {
        opacity: 0;
        transform: scale(0.8);
    }
    to {
        opacity: 1;
        transform: scale(1);
    }
`;

const GameStatusOverlay = styled.div`
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10;
    animation: ${fadeIn} 0.3s ease-out;
`;

const GameStatusContent = styled.div`
    background: white;
    padding: 30px;
    border-radius: 10px;
    text-align: center;
    box-shadow: 0 0 20px rgba(0, 0, 0, 0.2);
`;

const GameStatusTitle = styled.h2`
    margin: 0 0 20px 0;
    font-size: 24px;
    color: #333;
`;

const GameStatusResult = styled.div`
    font-size: 36px;
    font-weight: bold;
    margin: 20px 0;
    color: #2c3e50;
`;

const GameStatusButton = styled.button`
    padding: 10px 20px;
    font-size: 16px;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
        background-color: #2980b9;
    }
`;

const getResultText = (result: string) => {
    switch (result) {
        case '1-0':
            return 'Победа белых!';
        case '0-1':
            return 'Победа черных!';
        case '1-1':
            return 'Ничья!';
        default:
            return 'Игра в процессе';
    }
};

const getResultEmoji = (result: string) => {
    switch (result) {
        case '1-0':
            return '🏆';
        case '0-1':
            return '🏆';
        case '1-1':
            return '🤝';
        default:
            return '⚔️';
    }
};

const HistoryControls = styled.div`
    position: fixed;
    bottom: 20px;
    left: 20px;
    display: flex;
    gap: 10px;
    background: white;
    padding: 10px;
    border-radius: 5px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
`;

const HistoryButton = styled.button`
    padding: 8px 16px;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
        background-color: #2980b9;
    }

    &:disabled {
        background-color: #bdc3c7;
        cursor: not-allowed;
    }
`;

interface Move {
    from: Position;
    to: Position;
    piece: Piece | null;
    capturedPiece: Piece | null;
    board: (Piece | null)[][];
    whiteCaptures: string[];
    blackCaptures: string[];
}

const movePiece = keyframes`
    from {
        transform: translate(0, 0);
    }
    to {
        transform: translate(var(--moveX), var(--moveY));
    }
`;

const AnimatedCell = styled(Cell)<{ isMoving: boolean; moveX: number; moveY: number }>`
    position: relative;
    ${props => props.isMoving && css`
        &::before {
            content: attr(data-piece);
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            animation: ${movePiece} 0.3s ease-out forwards;
            --moveX: ${props.moveX}px;
            --moveY: ${props.moveY}px;
            z-index: 2;
        }
    `}
`;

interface ChessMove {
    from: string;
    to: string;
    promotion?: string;
}

interface GameState {
    fen: string;
    moves: string[];
    status: 'playing' | 'finished' | 'aborted';
    result?: string;
}

const NewDesignButton = styled.button`
    position: fixed;
    top: 20px;
    left: 140px;
    padding: 10px 20px;
    background-color: #28a745;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.2s;
    z-index: 2;

    &:hover {
        background-color: #218838;
    }
`;

export const ChessBoard: React.FC = () => {
    const { gameId } = useParams();
    const [board, setBoard] = useState<(Piece | null)[][]>([]);
    const [selectedCell, setSelectedCell] = useState<Position | null>(null);
    const [currentTurn, setCurrentTurn] = useState<'white' | 'black'>('white');
    const [whiteCaptures, setWhiteCaptures] = useState<string[]>([]);
    const [blackCaptures, setBlackCaptures] = useState<string[]>([]);
    const [isGameStarted, setIsGameStarted] = useState(!!gameId);
    const [isSearching, setIsSearching] = useState(false);
    const [searchStatus, setSearchStatus] = useState('');
    const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
    const [ws, setWs] = useState<WebSocket | null>(null);
    const [showGameStatus, setShowGameStatus] = useState(false);
    const [moveHistory, setMoveHistory] = useState<Move[]>(() => {
        if (gameId) {
            const savedHistory = localStorage.getItem(`game_history_${gameId}`);
            return savedHistory ? JSON.parse(savedHistory) : [];
        }
        return [];
    });
    const [currentMoveIndex, setCurrentMoveIndex] = useState(() => {
        if (gameId) {
            const savedIndex = localStorage.getItem(`game_current_move_${gameId}`);
            return savedIndex ? parseInt(savedIndex) : -1;
        }
        return -1;
    });
    const [animatingPiece, setAnimatingPiece] = useState<{
        from: Position;
        to: Position;
        piece: string;
    } | null>(null);
    const [gameState, setGameState] = useState<GameState>({
        fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1',
        moves: [],
        status: 'playing'
    });
    const [isConnected, setIsConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();
    const [currentUserId, setCurrentUserId] = useState<string | null>(null);
    const [playerColor, setPlayerColor] = useState<'white' | 'black' | null>(null);
    const [reconnectAttempts, setReconnectAttempts] = useState(0);
    const MAX_RECONNECT_ATTEMPTS = 5;
    const RECONNECT_DELAY = 3000;
    const CONNECTION_TIMEOUT = 5000;

    // Refs для управления WebSocket соединением
    const wsRef = useRef<WebSocket | null>(null);
    const isConnectingRef = useRef(false);
    const pingIntervalRef = useRef<NodeJS.Timeout | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const connectionTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    // Добавляем refs для хранения мутабельных значений
    const boardRef = useRef(board);
    const selectedCellRef = useRef(selectedCell);
    const gameInfoRef = useRef(gameInfo);

    // Обновляем refs при изменении значений
    useEffect(() => {
        boardRef.current = board;
    }, [board]);

    useEffect(() => {
        selectedCellRef.current = selectedCell;
    }, [selectedCell]);

    useEffect(() => {
        gameInfoRef.current = gameInfo;
    }, [gameInfo]);

    // Добавляем функцию для определения цвета игрока
    const determinePlayerColor = useCallback((gameInfo: GameInfo) => {
        const userId = localStorage.getItem('userId');
        if (!userId) return null;

        if (gameInfo.WhiteUser.id === userId) return 'white';
        if (gameInfo.BlackUser.id === userId) return 'black';
        return null;
    }, []);

    // Разделяем useEffect для инициализации игры и WebSocket
    useEffect(() => {
        if (!gameId) return;

        const token = localStorage.getItem('authToken');
        const userId = localStorage.getItem('userId');
        if (!token || !userId) {
            navigate('/login');
            return;
        }

        // Проверяем, что это текущая игра
        const currentGameId = localStorage.getItem('currentGameId');
        if (currentGameId !== gameId) {
            navigate('/game');
            return;
        }

        // Получаем информацию об игре
        const fetchGameInfo = async () => {
            try {
                const response = await fetch(`${config.apiBaseUrl}/chess/${gameId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });

                if (!response.ok) {
                    throw new Error('Failed to fetch game info');
                }

                const data = await response.json();
                setGameInfo(data);
                
                // Определяем цвет игрока
                const color = determinePlayerColor(data);
                setPlayerColor(color);
                
                // Устанавливаем начальный ход
                setCurrentTurn('white');
            } catch (error) {
                console.error('Error fetching game info:', error);
            }
        };

        // Проверяем, не в режиме ли разработки
        const isDev = process.env.NODE_ENV === 'development';
        if (!isDev) {
            fetchGameInfo();
        } else {
            // В режиме разработки используем setTimeout для предотвращения дублирования запросов
            const timeoutId = setTimeout(() => {
                fetchGameInfo();
            }, 0);
            return () => clearTimeout(timeoutId);
        }
    }, [gameId, navigate, determinePlayerColor]);

    // Функция для очистки всех WebSocket ресурсов
    const cleanupWebSocket = useCallback(() => {
        if (connectionTimeoutRef.current) {
            clearTimeout(connectionTimeoutRef.current);
            connectionTimeoutRef.current = null;
        }
        if (wsRef.current) {
            if (wsRef.current.readyState === WebSocket.OPEN) {
                wsRef.current.close(1000, "Normal closure");
            }
            wsRef.current = null;
        }
        if (pingIntervalRef.current) {
            clearInterval(pingIntervalRef.current);
            pingIntervalRef.current = null;
        }
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }
        setIsConnected(false);
        isConnectingRef.current = false;
    }, []);

    // Функция для создания WebSocket соединения
    const createWebSocket = useCallback((token: string) => {
        if (!gameId || isConnectingRef.current) return null;
        
        cleanupWebSocket();
        isConnectingRef.current = true;
        
        // Исправляем URL на правильный формат v1/ws/game/gameID
        const wsUrl = `${config.wsBaseUrl}/ws/game/${gameId}?token=${encodeURIComponent('Bearer ' + token)}`;
        console.log('Connecting to WebSocket:', wsUrl);
        
        const socket = new WebSocket(wsUrl);

        // Увеличиваем таймаут соединения
        connectionTimeoutRef.current = setTimeout(() => {
            if (socket.readyState !== WebSocket.OPEN) {
                console.log('WebSocket connection timeout');
                socket.close(1000, "Connection timeout");
                cleanupWebSocket();
                setError('Не удалось установить соединение с сервером');
            }
        }, CONNECTION_TIMEOUT);

        socket.onopen = () => {
            console.log('WebSocket соединение установлено');
            if (connectionTimeoutRef.current) {
                clearTimeout(connectionTimeoutRef.current);
                connectionTimeoutRef.current = null;
            }
            setIsConnected(true);
            setIsGameStarted(true); // Устанавливаем флаг начала игры
            setError(null);
            setReconnectAttempts(0);
            isConnectingRef.current = false;
            
            // Устанавливаем интервал для пинга
            if (pingIntervalRef.current) {
                clearInterval(pingIntervalRef.current);
            }
            pingIntervalRef.current = setInterval(() => {
                if (socket.readyState === WebSocket.OPEN) {
                    try {
                        socket.send("0000");
                    } catch (error) {
                        console.error('Ошибка отправки пинга:', error);
                    }
                }
            }, 30000);
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            setError('Ошибка соединения с сервером');
            cleanupWebSocket();
        };

        socket.onclose = (event) => {
            console.log('WebSocket closed:', event.code, event.reason);
            cleanupWebSocket();
            
            if (event.code !== 1000) { // Если соединение закрыто не нормально
                setError('Соединение с сервером потеряно');
                
                // Пробуем переподключиться
                if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
                    reconnectTimeoutRef.current = setTimeout(() => {
                        setReconnectAttempts(prev => prev + 1);
                        createWebSocket(token);
                    }, RECONNECT_DELAY);
                } else {
                    setError('Не удалось восстановить соединение');
                }
            }
        };

        return socket;
    }, [gameId, reconnectAttempts, cleanupWebSocket]);

    // Единый useEffect для управления WebSocket соединением
    useEffect(() => {
        if (!gameId || !gameInfo) return;

        const token = localStorage.getItem('authToken');
        if (!token) {
            setError('Нет токена авторизации');
            return;
        }

        // Проверяем, не в режиме ли разработки
        const isDev = process.env.NODE_ENV === 'development';
        if (!isDev) {
            const socket = createWebSocket(token);
            if (socket) {
                wsRef.current = socket;
                setWs(socket);
            }
        } else {
            // В режиме разработки используем setTimeout для предотвращения дублирования соединений
            const timeoutId = setTimeout(() => {
                const socket = createWebSocket(token);
                if (socket) {
                    wsRef.current = socket;
                    setWs(socket);
                }
            }, 0);
            return () => clearTimeout(timeoutId);
        }

        return () => {
            cleanupWebSocket();
        };
    }, [gameId, gameInfo, createWebSocket, cleanupWebSocket]);

    // Очистка при размонтировании компонента
    useEffect(() => {
        return () => {
            cleanupWebSocket();
        };
    }, [cleanupWebSocket]);

    // Отдельный useEffect для сохранения истории
    useEffect(() => {
        if (gameId) {
            localStorage.setItem(`game_history_${gameId}`, JSON.stringify(moveHistory));
            localStorage.setItem(`game_current_move_${gameId}`, currentMoveIndex.toString());
        }
    }, [moveHistory, currentMoveIndex, gameId]);

    // Отдельный useEffect для очистки истории при завершении игры
    useEffect(() => {
        if (gameInfo && (gameInfo.Status === 'finished' || gameInfo.Status === 'aborted')) {
            if (gameId) {
                localStorage.setItem(`game_history_${gameId}_final`, JSON.stringify(moveHistory));
            }
        }
    }, [gameInfo, gameId, moveHistory]);

    useEffect(() => {
        initBoard();
        if (gameId) {
            // Получаем информацию об игре
            const fetchGameInfo = async () => {
                try {
                    const token = localStorage.getItem('authToken');
                    if (!token) {
                        navigate('/login');
                        return;
                    }

                    const response = await fetch(`${config.apiBaseUrl}/chess/${gameId}`, {
                        headers: {
                            'Authorization': `Bearer ${token}`
                        }
                    });

                    if (!response.ok) {
                        throw new Error('Failed to fetch game info');
                    }

                    const data = await response.json();
                    setGameInfo(data);
                } catch (error) {
                    console.error('Error fetching game info:', error);
                }
            };

            fetchGameInfo();

            // Подключаемся к существующей игре
            const token = localStorage.getItem('authToken');
            if (!token) {
                navigate('/login');
                return;
            }

            const wsUrl = `${config.wsBaseUrl}${config.endpoints.game.play}/${gameId}?token=${encodeURIComponent('Bearer ' + token)}`;
            const ws = new WebSocket(wsUrl);

            ws.onopen = () => {
                console.log('WebSocket соединение установлено для игры:', gameId);
            };

            ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    console.log('Received WebSocket message:', data);
                    
                    // Игнорируем ответы на пинги
                    if (data.move === "0000") return;
                    
                    if (data.move) {
                        const move = data.move;
                        if (typeof move === 'string' && move.length === 4) {
                            const fromNotation = move.slice(0, 2);
                            const toNotation = move.slice(2, 4);
                            const [fromRow, fromCol] = convertNotationToPosition(fromNotation);
                            const [toRow, toCol] = convertNotationToPosition(toNotation);
                            
                            // Создаем глубокую копию доски
                            const newBoard = boardRef.current.map(row => [...row]);
                            
                            // Перемещаем фигуру
                            newBoard[toRow][toCol] = newBoard[fromRow][fromCol];
                            newBoard[fromRow][fromCol] = null;
                            
                            // Обновляем состояние доски
                            setBoard(newBoard);
                            
                            // Обновляем чей ход после получения хода
                            setCurrentTurn(prevTurn => {
                                const nextTurn = prevTurn === 'white' ? 'black' : 'white';
                                console.log('Turn changed to:', nextTurn);
                                return nextTurn;
                            });
                            setError(null);
                        }
                    }
                    
                    if (data.ok === false) {
                        if (data.status === 'GameFinished') {
                            setShowGameStatus(true);
                            const currentGameInfo = gameInfoRef.current;
                            if (currentGameInfo) {
                                setGameInfo({
                                    ...currentGameInfo,
                                    Status: 'finished',
                                    Result: data.result || '1-1'
                                });
                            }
                        } else {
                            setError(data.message || 'Недопустимый ход');
                            // Возвращаем фигуру на исходную позицию
                            const currentSelectedCell = selectedCellRef.current;
                            if (currentSelectedCell) {
                                const newBoard = boardRef.current.map(row => [...row]);
                                setBoard(newBoard);
                            }
                            fetchGameStatus();
                        }
                    } else if (data.ok === true) {
                        // Обновляем чей ход после успешного хода
                        setCurrentTurn(prevTurn => {
                            const nextTurn = prevTurn === 'white' ? 'black' : 'white';
                            console.log('Turn changed to:', nextTurn);
                            return nextTurn;
                        });
                        setError(null);
                    }
                } catch (error) {
                    console.error('Ошибка при обработке сообщения:', error);
                    setError('Ошибка при обработке сообщения от сервера');
                }
            };

            ws.onerror = (error) => {
                console.error('WebSocket ошибка:', error);
            };

            return () => {
                ws.close();
            };
        }
    }, [gameId, navigate]);

    useEffect(() => {
        if (gameInfo) {
            if (gameInfo.Status === 'finished' || gameInfo.Status === 'aborted') {
                setShowGameStatus(true);
            }
        }
    }, [gameInfo]);

    const initBoard = () => {
        const newBoard: (Piece | null)[][] = Array(8).fill(null).map(() => Array(8).fill(null));

        // Черные фигуры
        newBoard[0][0] = { color: 'black', type: 'rook' };
        newBoard[0][1] = { color: 'black', type: 'knight' };
        newBoard[0][2] = { color: 'black', type: 'bishop' };
        newBoard[0][3] = { color: 'black', type: 'queen' };
        newBoard[0][4] = { color: 'black', type: 'king' };
        newBoard[0][5] = { color: 'black', type: 'bishop' };
        newBoard[0][6] = { color: 'black', type: 'knight' };
        newBoard[0][7] = { color: 'black', type: 'rook' };
        for (let j = 0; j < 8; j++) {
            newBoard[1][j] = { color: 'black', type: 'pawn' };
        }

        // Белые фигуры
        newBoard[7][0] = { color: 'white', type: 'rook' };
        newBoard[7][1] = { color: 'white', type: 'knight' };
        newBoard[7][2] = { color: 'white', type: 'bishop' };
        newBoard[7][3] = { color: 'white', type: 'queen' };
        newBoard[7][4] = { color: 'white', type: 'king' };
        newBoard[7][5] = { color: 'white', type: 'bishop' };
        newBoard[7][6] = { color: 'white', type: 'knight' };
        newBoard[7][7] = { color: 'white', type: 'rook' };
        for (let j = 0; j < 8; j++) {
            newBoard[6][j] = { color: 'white', type: 'pawn' };
        }

        setBoard(newBoard);
    };

    const handleStartGame = () => {
        setIsSearching(true);
        setSearchStatus('Поиск игры...');
        
        const token = localStorage.getItem('authToken');
        const userId = localStorage.getItem('userId');
        console.log('authToken:', localStorage.getItem('authToken'));
        console.log('userId:', localStorage.getItem('userId'));
        if (!token || !userId) {
            setSearchStatus('Ошибка: нет токена авторизации');
            setIsSearching(false);
            console.log('Нет токена авторизации или ID пользователя');
            navigate('/login');
            return;
        }

        

        // Передаем токен в URL как query-параметр
        const wsUrl = `${config.wsBaseUrl}${config.endpoints.game.search}?token=${encodeURIComponent('Bearer ' + token)}`;
        console.log('Открываю WebSocket:', wsUrl);
        const ws = new WebSocket(wsUrl);
        
        ws.onopen = () => {
            console.log('WebSocket соединение установлено');
        };

        ws.onmessage = (event) => {
            console.log('Получено сообщение от WebSocket:', event.data);
            try {
                const data = JSON.parse(event.data);
                const id = data.gameId || data.gemId || data.gameID; // поддержка всех вариантов
                if (id) {
                    console.log('Переход на:', `/game/${id}`);
                    setSearchStatus('Игра найдена! Перенаправление...');
                    // Сохраняем ID игры перед переходом
                    localStorage.setItem('currentGameId', id);
                    // Закрываем WebSocket только после сохранения ID
                    setTimeout(() => {
                        ws.close();
                        navigate(`/game/${id}`);
                    }, 100);
                } else {
                    console.log('В сообщении нет gameId/gemId/gameID:', data);
                }
            } catch (error) {
                console.error('Ошибка при обработке сообщения:', error);
                setSearchStatus('Ошибка при получении данных игры');
                setIsSearching(false);
            }
        };

        ws.onerror = (error) => {
            console.error('WebSocket ошибка:', error);
            setSearchStatus('Ошибка соединения');
            setIsSearching(false);
        };

        ws.onclose = (event) => {
            console.log('WebSocket соединение закрыто', event);
            // Не сбрасываем состояние поиска при нормальном закрытии
            if (event.code !== 1000) {
                setIsSearching(false);
            }
        };
    };

    // Функция для конвертации FEN в массив доски
    const fenToBoard = useCallback((fen: string) => {
        const [position] = fen.split(' ');
        const rows = position.split('/');
        const newBoard: (Piece | null)[][] = Array(8).fill(null).map(() => Array(8).fill(null));
        
        rows.forEach((row, i) => {
            let col = 0;
            for (const char of row) {
                if (char >= '1' && char <= '8') {
                    col += parseInt(char);
                } else {
                    const color = char === char.toUpperCase() ? 'white' : 'black';
                    const type = char.toLowerCase() as 'k' | 'q' | 'r' | 'b' | 'n' | 'p';
                    const pieceType = {
                        'k': 'king',
                        'q': 'queen',
                        'r': 'rook',
                        'b': 'bishop',
                        'n': 'knight',
                        'p': 'pawn'
                    }[type];
                    
                    newBoard[i][col] = { color, type: pieceType as Piece['type'] };
                    col++;
                }
            }
        });
        
        return newBoard;
    }, []);

    // Добавляем функцию для получения статуса игры
    const fetchGameStatus = async () => {
        try {
            const token = localStorage.getItem('authToken');
            if (!token || !gameId) return;

            const response = await fetch(`${config.apiBaseUrl}/chess/${gameId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch game status');
            }

            const data = await response.json();
            setGameInfo(data);

            // Если игра завершена, показываем оверлей
            if (data.Status === 'finished' || data.Status === 'aborted') {
                setShowGameStatus(true);
            }
        } catch (error) {
            console.error('Error fetching game status:', error);
            setError('Ошибка при получении статуса игры');
        }
    };

    const handleOpponentMove = (from: string, to: string) => {
        // Обновляем состояние доски после хода противника
        const newBoard = [...board];
        const [fromRow, fromCol] = convertNotationToPosition(from);
        const [toRow, toCol] = convertNotationToPosition(to);
        
        newBoard[toRow][toCol] = newBoard[fromRow][fromCol];
        newBoard[fromRow][fromCol] = null;
        
        setBoard(newBoard);
        setCurrentTurn('white');
    };

    const convertNotationToPosition = (notation: string): [number, number] => {
        const col = notation.charCodeAt(0) - 'a'.charCodeAt(0);
        const row = 8 - parseInt(notation[1]);
        // Если игрок играет черными, переворачиваем координаты
        if (playerColor === 'black') {
            return [7 - row, 7 - col];
        }
        return [row, col];
    };

    const convertPositionToNotation = (row: number, col: number): string => {
        // Если игрок играет черными, переворачиваем координаты
        let displayRow = row;
        let displayCol = col;
        if (playerColor === 'black') {
            displayRow = 7 - row;
            displayCol = 7 - col;
        }
        const colChar = String.fromCharCode('a'.charCodeAt(0) + displayCol);
        const rowNum = 8 - displayRow;
        return `${colChar}${rowNum}`;
    };

    const handleCellClick = async (row: number, col: number) => {
        if (!gameId) {
            setError('Игра не найдена');
            return;
        }

        if (!isConnected) {
            setError('Нет подключения к серверу');
            return;
        }

        // Проверяем, что это ход текущего игрока
        if (currentTurn !== playerColor) {
            setError(`Сейчас ход ${currentTurn === 'white' ? 'белых' : 'черных'}`);
            return;
        }

        if (!selectedCell) {
            // Выбор фигуры для хода
            if (board[row][col] && board[row][col]?.color === playerColor) {
                setSelectedCell({ row, col });
                setError(null);
            } else if (board[row][col]) {
                setError(`Вы играете за ${playerColor === 'white' ? 'белых' : 'черных'}`);
            }
        } else {
            // Ход выбранной фигурой
            const from = selectedCell;
            const fromNotation = convertPositionToNotation(from.row, from.col);
            const toNotation = convertPositionToNotation(row, col);
            const moveNotation = `${fromNotation}${toNotation}`;

            try {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    // Создаем временную копию доски для предварительного отображения хода
                    const tempBoard = board.map(row => [...row]);
                    tempBoard[row][col] = tempBoard[from.row][from.col];
                    tempBoard[from.row][from.col] = null;
                    setBoard(tempBoard);
                    
                    setSelectedCell(null);
                    console.log('Sending move:', moveNotation);
                    ws.send(moveNotation);
                } else {
                    setError('Нет подключения к серверу');
                }
            } catch (error) {
                console.error('Ошибка при выполнении хода:', error);
                setError('Ошибка при выполнении хода');
                // Возвращаем доску в исходное состояние при ошибке
                const originalBoard = board.map(row => [...row]);
                setBoard(originalBoard);
                setSelectedCell(null);
            }
        }
    };

    const handleUndo = () => {
        if (currentMoveIndex > 0) {
            const previousMove = moveHistory[currentMoveIndex - 1];
            const currentMove = moveHistory[currentMoveIndex];
            
            // Анимируем перемещение фигуры
            if (currentMove.piece) {
                const fromCell = document.querySelector(`[data-row="${currentMove.to.row}"][data-col="${currentMove.to.col}"]`);
                const toCell = document.querySelector(`[data-row="${currentMove.from.row}"][data-col="${currentMove.from.col}"]`);
                
                if (fromCell && toCell) {
                    const fromRect = fromCell.getBoundingClientRect();
                    const toRect = toCell.getBoundingClientRect();
                    
                    setAnimatingPiece({
                        from: currentMove.to,
                        to: currentMove.from,
                        piece: pieces[currentMove.piece.color][currentMove.piece.type]
                    });

                    // Ждем окончания анимации
                    setTimeout(() => {
                        setBoard(previousMove.board);
                        setWhiteCaptures(previousMove.whiteCaptures);
                        setBlackCaptures(previousMove.blackCaptures);
                        setCurrentMoveIndex(currentMoveIndex - 1);
                        setCurrentTurn(currentTurn === 'white' ? 'black' : 'white');
                        setAnimatingPiece(null);
                    }, 300);
                }
            } else {
                setBoard(previousMove.board);
                setWhiteCaptures(previousMove.whiteCaptures);
                setBlackCaptures(previousMove.blackCaptures);
                setCurrentMoveIndex(currentMoveIndex - 1);
                setCurrentTurn(currentTurn === 'white' ? 'black' : 'white');
            }
        }
    };

    const handleRedo = () => {
        if (currentMoveIndex < moveHistory.length - 1) {
            const nextMove = moveHistory[currentMoveIndex + 1];
            const currentMove = moveHistory[currentMoveIndex];
            
            // Анимируем перемещение фигуры
            if (nextMove.piece) {
                const fromCell = document.querySelector(`[data-row="${nextMove.from.row}"][data-col="${nextMove.from.col}"]`);
                const toCell = document.querySelector(`[data-row="${nextMove.to.row}"][data-col="${nextMove.to.col}"]`);
                
                if (fromCell && toCell) {
                    const fromRect = fromCell.getBoundingClientRect();
                    const toRect = toCell.getBoundingClientRect();
                    
                    setAnimatingPiece({
                        from: nextMove.from,
                        to: nextMove.to,
                        piece: pieces[nextMove.piece.color][nextMove.piece.type]
                    });

                    // Ждем окончания анимации
                    setTimeout(() => {
                        setBoard(nextMove.board);
                        setWhiteCaptures(nextMove.whiteCaptures);
                        setBlackCaptures(nextMove.blackCaptures);
                        setCurrentMoveIndex(currentMoveIndex + 1);
                        setCurrentTurn(currentTurn === 'white' ? 'black' : 'white');
                        setAnimatingPiece(null);
                    }, 300);
                }
            } else {
                setBoard(nextMove.board);
                setWhiteCaptures(nextMove.whiteCaptures);
                setBlackCaptures(nextMove.blackCaptures);
                setCurrentMoveIndex(currentMoveIndex + 1);
                setCurrentTurn(currentTurn === 'white' ? 'black' : 'white');
            }
        }
    };

    const handleLogout = () => {
        cleanupWebSocket();
        localStorage.removeItem('authToken');
        localStorage.removeItem('userId');
        navigate('/login');
    };

    const handleNewGame = () => {
        cleanupWebSocket();
        setBoard([]);
        setSelectedCell(null);
        setCurrentTurn('white');
        setWhiteCaptures([]);
        setBlackCaptures([]);
        setIsGameStarted(false);
        setIsSearching(false);
        setSearchStatus('');
        setGameInfo(null);
        setShowGameStatus(false);
        setMoveHistory([]);
        setCurrentMoveIndex(-1);
        
        if (gameId) {
            localStorage.removeItem(`game_history_${gameId}`);
            localStorage.removeItem(`game_current_move_${gameId}`);
        }
        
        navigate('/game');
    };

    const handleGoHome = () => {
        navigate('/');
    };

    // Добавляем компонент для отображения ошибок
    const ErrorMessage = styled.div`
        position: fixed;
        top: 80px;
        left: 50%;
        transform: translateX(-50%);
        background-color: #ff4444;
        color: white;
        padding: 10px 20px;
        border-radius: 5px;
        z-index: 1000;
        animation: fadeIn 0.3s ease-in;
    `;

    return (
        <BoardContainer>
            <LogoutButton onClick={handleLogout}>Выйти</LogoutButton>
            <NewDesignButton onClick={() => navigate('/new-design')}>New Design</NewDesignButton>
            <LogoutButton style={{ left: 160, backgroundColor: '#3498db' }} onClick={handleGoHome}>На главную</LogoutButton>
            {error && <ErrorMessage>{error}</ErrorMessage>}
            {gameInfo && (
                <>
                    <WhitePlayerInfo>
                        <PlayerAvatar>W</PlayerAvatar>
                        <div>
                            <div>{gameInfo.WhiteUser.name}</div>
                            <div style={{ fontSize: '12px', color: '#666' }}>Белые</div>
                        </div>
                    </WhitePlayerInfo>
                    <BlackPlayerInfo>
                        <PlayerAvatar>B</PlayerAvatar>
                        <div>
                            <div>{gameInfo.BlackUser.name}</div>
                            <div style={{ fontSize: '12px', color: '#666' }}>Черные</div>
                        </div>
                    </BlackPlayerInfo>
                </>
            )}
            <Board>
                <tbody>
                    {board.map((row, i) => {
                        // Если игрок играет черными, переворачиваем доску
                        const displayRow = playerColor === 'black' ? 7 - i : i;
                        return (
                            <tr key={i}>
                                <SideNumber>{playerColor === 'black' ? i + 1 : 8 - i}</SideNumber>
                                {row.map((cell, j) => {
                                    // Если игрок играет черными, переворачиваем колонки
                                    const displayCol = playerColor === 'black' ? 7 - j : j;
                                    return (
                                        <AnimatedCell
                                            key={`${i}-${j}`}
                                            isLight={(displayRow + displayCol) % 2 === 0}
                                            isSelected={selectedCell?.row === i && selectedCell?.col === j}
                                            onClick={() => handleCellClick(i, j)}
                                            isMoving={animatingPiece?.from.row === i && animatingPiece?.from.col === j}
                                            moveX={animatingPiece ? (animatingPiece.to.col - animatingPiece.from.col) * 60 : 0}
                                            moveY={animatingPiece ? (animatingPiece.to.row - animatingPiece.from.row) * 60 : 0}
                                            data-row={i}
                                            data-col={j}
                                            data-piece={cell ? pieces[cell.color][cell.type] : ''}
                                        >
                                            {cell && pieces[cell.color][cell.type]}
                                        </AnimatedCell>
                                    );
                                })}
                            </tr>
                        );
                    })}
                </tbody>
            </Board>
            <CapturesContainer>
                <h3>Захваченные фигуры:</h3>
                <p>Белые: {whiteCaptures.join(' ')}</p>
                <p>Черные: {blackCaptures.join(' ')}</p>
            </CapturesContainer>
            {!gameId && (
                <Overlay isGameStarted={isGameStarted}>
                    <div>
                        <StartButton 
                            onClick={handleStartGame}
                            disabled={isSearching}
                        >
                            {isSearching ? 'Поиск игры...' : 'Начать игру'}
                        </StartButton>
                        {searchStatus && <StatusMessage>{searchStatus}</StatusMessage>}
                    </div>
                </Overlay>
            )}
            {showGameStatus && gameInfo && (
                <GameStatusOverlay>
                    <GameStatusContent>
                        <GameStatusTitle>
                            {gameInfo.Status === 'finished' ? 'Игра завершена' : 'Игра прервана'}
                        </GameStatusTitle>
                        <GameStatusResult>
                            {getResultEmoji(gameInfo.Result)} {getResultText(gameInfo.Result)}
                        </GameStatusResult>
                        <GameStatusButton onClick={handleNewGame}>
                            Новая игра
                        </GameStatusButton>
                    </GameStatusContent>
                </GameStatusOverlay>
            )}
            <HistoryControls>
                <HistoryButton 
                    onClick={handleUndo}
                    disabled={currentMoveIndex <= 0}
                >
                    ⬅️ Назад
                </HistoryButton>
                <HistoryButton 
                    onClick={handleRedo}
                    disabled={currentMoveIndex >= moveHistory.length - 1}
                >
                    Вперед ➡️
                </HistoryButton>
            </HistoryControls>
        </BoardContainer>
    );
}; 
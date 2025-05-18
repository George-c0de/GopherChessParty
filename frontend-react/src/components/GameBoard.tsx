import React, { useState, useEffect, useRef, useCallback } from 'react';
import styled from 'styled-components';
import { useParams, useNavigate } from 'react-router-dom';
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

export const GameBoard: React.FC = () => {
    const { gameId } = useParams();
    const navigate = useNavigate();
    const [board, setBoard] = useState<(Piece | null)[][]>([]);
    const [selectedCell, setSelectedCell] = useState<Position | null>(null);
    const [currentTurn, setCurrentTurn] = useState<'white' | 'black'>('white');
    const [whiteCaptures, setWhiteCaptures] = useState<string[]>([]);
    const [blackCaptures, setBlackCaptures] = useState<string[]>([]);
    const [ws, setWs] = useState<WebSocket | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Добавляем функции для конвертации нотации
    const convertNotationToPosition = (notation: string): [number, number] => {
        const col = notation.charCodeAt(0) - 'a'.charCodeAt(0);
        const row = 8 - parseInt(notation[1]);
        return [row, col];
    };

    const convertPositionToNotation = (row: number, col: number): string => {
        const colChar = String.fromCharCode('a'.charCodeAt(0) + col);
        const rowNum = 8 - row;
        return `${colChar}${rowNum}`;
    };

    // Refs для управления WebSocket соединением
    const wsRef = useRef<WebSocket | null>(null);
    const isConnectingRef = useRef(false);

    // Функция для очистки WebSocket соединения
    const cleanupWebSocket = useCallback(() => {
        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }
        setIsConnected(false);
        isConnectingRef.current = false;
    }, []);

    // Инициализация WebSocket соединения
    useEffect(() => {
        if (!gameId) return;

        const token = localStorage.getItem('access_token');
        if (!token) {
            setError('Нет токена авторизации');
            navigate('/login');
            return;
        }

        if (isConnectingRef.current) return;
        isConnectingRef.current = true;

        const wsUrl = `${config.wsBaseUrl}${config.endpoints.game.play}/${gameId}?token=${encodeURIComponent('Bearer ' + token)}`;
        const socket = new WebSocket(wsUrl);

        socket.onopen = () => {
            console.log('WebSocket соединение установлено');
            setIsConnected(true);
            setError(null);
            isConnectingRef.current = false;
        };

        socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                console.log('Received WebSocket message:', data);
                
                if (data.move) {
                    const move = data.move;
                    if (typeof move === 'string' && move.length === 4) {
                        const fromNotation = move.slice(0, 2);
                        const toNotation = move.slice(2, 4);
                        const [fromRow, fromCol] = convertNotationToPosition(fromNotation);
                        const [toRow, toCol] = convertNotationToPosition(toNotation);
                        const newBoard = [...board];
                        newBoard[toRow][toCol] = newBoard[fromRow][fromCol];
                        newBoard[fromRow][fromCol] = null;
                        setBoard(newBoard);
                        setCurrentTurn(prevTurn => prevTurn === 'white' ? 'black' : 'white');
                        setError(null);
                    }
                }
            } catch (error) {
                console.error('Ошибка при обработке сообщения:', error);
                setError('Ошибка при обработке сообщения от сервера');
            }
        };

        socket.onclose = () => {
            console.log('WebSocket соединение закрыто');
            cleanupWebSocket();
        };

        socket.onerror = (error) => {
            console.error('WebSocket ошибка:', error);
            setError('Ошибка соединения с сервером');
            cleanupWebSocket();
        };

        wsRef.current = socket;
        setWs(socket);

        return () => {
            cleanupWebSocket();
        };
    }, [gameId, navigate]);

    // Очистка при размонтировании компонента
    useEffect(() => {
        return () => {
            cleanupWebSocket();
        };
    }, [cleanupWebSocket]);

    const handleLogout = () => {
        cleanupWebSocket();
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('userId');
        navigate('/login');
    };

    useEffect(() => {
        initBoard();
    }, []);

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

    const handleCellClick = async (row: number, col: number) => {
        if (!selectedCell) {
            if (board[row][col] && board[row][col]?.color === currentTurn) {
                setSelectedCell({ row, col });
            }
        } else {
            const from = selectedCell;
            try {
                const response = await fetch(`${config.apiBaseUrl}${config.endpoints.game.move}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                    },
                    body: JSON.stringify({
                        from: from,
                        to: { row, col },
                        gameId // <-- передаем gameId
                    })
                });

                if (!response.ok) {
                    throw new Error('Ошибка при выполнении хода');
                }

                const newBoard = [...board];
                if (board[row][col]) {
                    if (board[row][col]?.color !== currentTurn) {
                        const capturedPiece = board[row][col];
                        if (capturedPiece) {
                            if (currentTurn === 'white') {
                                setWhiteCaptures([...whiteCaptures, pieces[capturedPiece.color][capturedPiece.type]]);
                            } else {
                                setBlackCaptures([...blackCaptures, pieces[capturedPiece.color][capturedPiece.type]]);
                            }
                        }
                    } else {
                        setSelectedCell(null);
                        return;
                    }
                }
                newBoard[row][col] = newBoard[from.row][from.col];
                newBoard[from.row][from.col] = null;
                setBoard(newBoard);
                setCurrentTurn(currentTurn === 'white' ? 'black' : 'white');
                setSelectedCell(null);
            } catch (error) {
                setSelectedCell(null);
            }
        }
    };

    return (
        <BoardContainer>
            <button
                style={{
                    position: 'absolute',
                    top: 20,
                    left: 20,
                    padding: '10px 20px',
                    background: '#e74c3c',
                    color: '#fff',
                    border: 'none',
                    borderRadius: '5px',
                    cursor: 'pointer',
                    zIndex: 2
                }}
                onClick={handleLogout}
            >
                Выйти
            </button>
            <Board>
                <tbody>
                    {board.map((row, i) => (
                        <tr key={i}>
                            <SideNumber>{8 - i}</SideNumber>
                            {row.map((cell, j) => (
                                <Cell
                                    key={`${i}-${j}`}
                                    isLight={(i + j) % 2 === 0}
                                    isSelected={selectedCell?.row === i && selectedCell?.col === j}
                                    onClick={() => handleCellClick(i, j)}
                                >
                                    {cell && pieces[cell.color][cell.type]}
                                </Cell>
                            ))}
                        </tr>
                    ))}
                </tbody>
            </Board>
            <CapturesContainer>
                <h3>Захваченные фигуры:</h3>
                <p>Белые: {whiteCaptures.join(' ')}</p>
                <p>Черные: {blackCaptures.join(' ')}</p>
            </CapturesContainer>
        </BoardContainer>
    );
}; 
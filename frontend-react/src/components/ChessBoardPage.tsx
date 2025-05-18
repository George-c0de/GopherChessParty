import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { Chess } from 'chess.js';
import { config } from '../config';
import { Chessboard } from 'react-chessboard';
import Header from './Header';

// Theme types
declare module 'styled-components' {
  export interface DefaultTheme {
    colors: {
      background: string;
      text: string;
    };
  }
}

// styled-components
const BoardContainer = styled.div`
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #f0f0f0;
`;

const MainContent = styled.main`
  flex: 1;
  padding-top: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const BoardWrapper = styled.div`
  display: flex;
  gap: 20px;
  align-items: flex-start;
  position: relative;
  margin-top: 20px;
`;

const SideResultContainer = styled.div`
  min-width: 220px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  margin-left: 10px;
`;

const ChessBoardContainer = styled.div`
  position: relative;
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  margin: 0 auto;
  min-width: 600px;
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const MoveHistoryContainer = styled.div`
  background: white;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0,0,0,0.1);
  width: 200px;
  max-height: 600px;
  display: flex;
  flex-direction: column;
`;

const MoveHistoryTitle = styled.h3`
  margin: 0 0 10px 0;
  color: #333;
  text-align: center;
`;

const MoveHistoryList = styled.div`
  flex-grow: 1;
  overflow-y: auto;
  margin-bottom: 10px;
`;

const MoveItem = styled.div<{ isCurrent: boolean }>`
  padding: 5px;
  border-radius: 4px;
  margin: 2px 0;
  background: ${props => props.isCurrent ? '#e3f2fd' : 'transparent'};
  color: ${props => props.isCurrent ? '#1976d2' : '#333'};
  font-weight: ${props => props.isCurrent ? 'bold' : 'normal'};
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #f5f5f5;
  }
`;

const MoveControls = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 10px;
  margin-top: 10px;
`;

const ControlButton = styled.button`
  flex: 1;
  padding: 8px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.2);
  }

  &:disabled {
    background: #ccc;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }
`;

const PlayerInfo = styled.div`
  background: rgba(255, 255, 255, 0.95);
  padding: 12px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 200px;
  transition: all 0.3s ease;
  align-self: center;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(0,0,0,0.15);
  }
`;

const WhitePlayerInfo = styled(PlayerInfo)`
  border: 2px solid #ebecd0;
`;

const BlackPlayerInfo = styled(PlayerInfo)`
  border: 2px solid #779556;
`;

const PlayerAvatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  color: white;
  font-size: 18px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
`;

const PlayerName = styled.div`
  font-weight: bold;
  color: #2c3e50;
  font-size: 16px;
`;

const PlayerColor = styled.div`
  font-size: 12px;
  color: #7f8c8d;
  margin-top: 2px;
`;

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
  cursor: pointer;
  animation: fadeOut 3s forwards;
  animation-delay: 2s;

  @keyframes fadeOut {
    from { opacity: 1; }
    to { opacity: 0; }
  }
`;

const StartGameButton = styled.button`
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: 16px 32px;
  background: linear-gradient(135deg, #4a90e2, #357abd);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 18px;
  font-weight: bold;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  z-index: 100;

  &:hover {
    transform: translate(-50%, -50%) scale(1.05);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
  }

  &:active {
    transform: translate(-50%, -50%) scale(0.95);
  }
`;

const NewGameButton = styled(StartGameButton)`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: linear-gradient(135deg, #2ecc71 0%, #27ae60 100%);
  z-index: 20;
  padding: 15px 30px;
  font-size: 20px;
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);

  &:hover {
    transform: translate(-50%, -50%) translateY(-2px);
    box-shadow: 0 10px 20px rgba(0,0,0,0.25);
  }

  &:disabled {
    background: #cccccc;
    cursor: not-allowed;
    transform: translate(-50%, -50%);
    box-shadow: none;
  }
`;

const GameStatus = styled.div`
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  background-color: white;
  padding: 12px 24px;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 10;
  font-size: 16px;
  color: #333;
`;

const VictoryResult = styled.div`
  margin-top: 0;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(46, 204, 113, 0.3);
  font-size: 26px;
  font-weight: 800;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  animation: subtleAppear 0.5s ease-in-out;
  border: 2px solid #27ae60;
`;

const DefeatResult = styled.div`
  margin-top: 0;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #7f8c8d 0%, #95a5a6 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(149, 165, 166, 0.3);
  font-size: 26px;
  font-weight: 800;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  animation: subtleAppear 0.5s ease-in-out;
  border: 2px solid #7f8c8d;
`;

const DrawResult = styled.div`
  margin-top: 0;
  margin-bottom: 20px;
  background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(41, 128, 185, 0.3);
  font-size: 26px;
  font-weight: 800;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  animation: subtleAppear 0.5s ease-in-out;
  border: 2px solid #2c3e50;
`;

const WelcomeMessage = styled.div`
  text-align: center;
  margin-bottom: 20px;
  h1 {
    font-size: 32px;
    color: #2c3e50;
    margin-bottom: 10px;
  }
  p {
    font-size: 18px;
    color: #7f8c8d;
  }
`;

const MutedChessboard = styled.div`
  opacity: 0.3;
  pointer-events: none;
  position: relative;
`;

const CenteredStartButton = styled(StartGameButton)`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 20;
  opacity: 1;
  pointer-events: auto;
  padding: 15px 30px;
  font-size: 20px;
  background: linear-gradient(135deg, #2ecc71 0%, #27ae60 100%);
  box-shadow: 0 8px 16px rgba(0,0,0,0.2);

  &:hover {
    transform: translate(-50%, -50%) translateY(-2px);
    box-shadow: 0 10px 20px rgba(0,0,0,0.25);
  }
`;

const CheckAnimation = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 0, 0, 0.2);
  animation: pulse 1s infinite;
  pointer-events: none;
  z-index: 5;

  @keyframes pulse {
    0% { opacity: 0.2; }
    50% { opacity: 0.4; }
    100% { opacity: 0.2; }
  }
`;

const MateAnimation = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  pointer-events: none;
  z-index: 10;
  animation: fadeIn 0.5s ease-out;

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
`;

const MateText = styled.div`
  color: white;
  font-size: 48px;
  font-weight: bold;
  text-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  animation: scaleIn 0.5s ease-out;

  @keyframes scaleIn {
    from { transform: scale(0.5); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }
`;

interface ChessBoardWrapperProps {
  playerColor: Color | null;
}

const ChessBoardWrapper = styled.div<ChessBoardWrapperProps>`
  position: relative;
  width: 600px;
  height: 600px;
`;

const ButtonIcon = styled.span`
  font-size: 18px;
`;

const DefeatMessage = styled.div`
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
  z-index: 100;
  min-width: 300px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
`;

// Types
interface GameInfo {
  ID: string;
  CreatedAt: string;
  Result: string;
  Status: string;
  WhiteUser: { id: string; name: string; email: string };
  BlackUser: { id: string; name: string; email: string };
  HistoryMove?: string[];
}
interface WSMessage {
  historyMove?: string[];
  currentMove?: number;
  ok?: boolean;
  message?: string;
  move?: string;
  error?: string;
  result?: string;
  event?: string;
}

type Color = 'white' | 'black';

// Custom hook encapsulating game logic
const useChessGame = (gameId: string | undefined) => {
  const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
  const [position, setPosition] = useState('start');
  const [moveHistory, setMoveHistory] = useState<string[]>([]);
  const [currentTurn, setCurrentTurn] = useState<Color>('white');
  const [playerColor, setPlayerColor] = useState<Color | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSearching, setIsSearching] = useState(false);
  const [searchStatus, setSearchStatus] = useState('');
  const [gameResult, setGameResult] = useState<string | null>(null);
  const navigate = useNavigate();
  const [currentMoveIndex, setCurrentMoveIndex] = useState<number>(-1);
  const [temporaryPosition, setTemporaryPosition] = useState<string>('start');
  const [isCheck, setIsCheck] = useState(false);
  const [isMate, setIsMate] = useState(false);
  const [isConnected, setIsConnected] = useState(false);

  const chessRef = React.useRef(new Chess());
  const socketRef = React.useRef<WebSocket | null>(null);

  // Initialize game: fetch info
  useEffect(() => {
    async function init() {
      if (!gameId) {
        setIsLoading(false);
        return;
      }

      const accessToken = localStorage.getItem('access_token');
      const refreshToken = localStorage.getItem('refresh_token');
      const userId = localStorage.getItem('userId');
      
      if (!accessToken || !refreshToken || !userId) {
        navigate('/login');
        return;
      }

      try {
        setIsLoading(true);

        // Проверяем срок действия токена
        const tokenData = JSON.parse(atob(accessToken.split('.')[1]));
        const isTokenExpired = tokenData.exp * 1000 < Date.now();

        let currentAccessToken = accessToken;

        // Если токен истек, обновляем его
        if (isTokenExpired) {
          const refreshRes = await fetch(`${config.apiBaseUrl}${config.endpoints.auth.refresh}`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({ refresh_token: refreshToken })
          });

          if (!refreshRes.ok) {
            navigate('/login');
            return;
          }

          const refreshData = await refreshRes.json();
          currentAccessToken = refreshData.access_token;
          localStorage.setItem('access_token', currentAccessToken);
        }

        // Получаем информацию об игре
        const res = await fetch(`${config.apiBaseUrl}/chess/${gameId}`, {
          headers: { 
            Authorization: `Bearer ${currentAccessToken}`,
            'Content-Type': 'application/json'
          }
        });

        if (!res.ok) {
          throw new Error('Failed to fetch game info');
        }

        const data = await res.json();
        console.log('=== Game Info from REST API ===');
        console.log('Game Data:', data);
        console.log('White User:', data.WhiteUser);
        console.log('Black User:', data.BlackUser);
        console.log('Game Status:', data.Status);
        console.log('============================');
        
        // Проверяем структуру данных
        if (!data) {
          console.error('Data is null or undefined');
          setError('Ошибка получения данных игры');
          return;
        }

        if (!data.WhiteUser) {
          console.error('WhiteUser is missing');
          setError('Ошибка получения данных игры');
          return;
        }

        if (!data.BlackUser) {
          console.error('BlackUser is missing');
          setError('Ошибка получения данных игры');
          return;
        }

        if (!data.WhiteUser.id || !data.BlackUser.id) {
          console.error('User IDs are missing');
          setError('Ошибка получения данных игры');
          return;
        }
        
        setGameInfo(data);
        
        const userColor = data.WhiteUser.id === userId ? 'white' : 
                         data.BlackUser.id === userId ? 'black' : null;
        setPlayerColor(userColor);

        if (!userColor) {
          console.error('User color not found. User ID:', userId);
          console.error('White User ID:', data.WhiteUser.id);
          console.error('Black User ID:', data.BlackUser.id);
          setError('Вы не являетесь участником этой игры');
          return;
        }

        // Если игра в статусе waiting или playing, открываем WebSocket соединение
        if (data.Status === 'waiting' || data.Status === 'playing') {
          const ws = new WebSocket(`${config.wsBaseUrl}${config.endpoints.game.play}/${gameId}?token=${encodeURIComponent('Bearer ' + currentAccessToken)}`);
          socketRef.current = ws;

          ws.onopen = () => {
            console.log('=== WebSocket Connected ===');
            console.log('User ID:', userId);
            console.log('Player Color:', userColor);
            console.log('Game Status:', data.Status);
            console.log('========================');
            setIsConnected(true);
          };

          ws.onmessage = (event) => {
            console.log('=== Received Move ===');
            console.log('Player Color:', userColor);
            console.log('Move:', event.data);
            
            // Обрабатываем только ходы
            if (typeof event.data === 'string' && event.data.length === 4) {
              const move = event.data;
              const from = move.substring(0, 2);
              const to = move.substring(2, 4);
              chessRef.current.move({ from, to });
              setPosition(chessRef.current.fen());
              setMoveHistory(prev => [...prev, move]);
              setCurrentTurn(prev => prev === 'white' ? 'black' : 'white');
            }
          };

          ws.onclose = () => {
            console.log('=== WebSocket Disconnected ===');
            console.log('User ID:', userId);
            console.log('Player Color:', userColor);
            console.log('===========================');
            setIsConnected(false);
          };
        }

      } catch (e) {
        console.error('Error fetching game info:', e);
        setError(e instanceof Error ? e.message : 'Ошибка загрузки игры');
      } finally {
        setIsLoading(false);
      }
    }

    init();
  }, [gameId, navigate]);

  const startNewGame = async () => {
    const accessToken = localStorage.getItem('access_token');
    if (!accessToken) {
      navigate('/login');
      return;
    }

    setIsSearching(true);
    setSearchStatus('Поиск противника...');

    try {
      const wsUrl = `${config.wsBaseUrl}/ws/search?token=${encodeURIComponent('Bearer ' + accessToken)}`;
      console.log('Connecting to WebSocket:', wsUrl);
      
      const ws = new WebSocket(wsUrl);
      
      ws.onopen = () => {
        console.log('WebSocket connection established');
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log('Received WebSocket message:', data);
          if (data.gameID) {
            setSearchStatus('Противник найден! Перенаправление...');
            localStorage.setItem('currentGameId', data.gameID);
            setTimeout(() => {
              ws.close();
              navigate(`/game/${data.gameID}`);
            }, 100);
          }
        } catch (error) {
          console.error('Ошибка при обработке сообщения:', error);
          setSearchStatus('Ошибка при поиске игры');
          setIsSearching(false);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setSearchStatus('Ошибка соединения');
        setIsSearching(false);
      };

      ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason);
        if (!localStorage.getItem('currentGameId')) {
          setIsSearching(false);
        }
      };
    } catch (error) {
      console.error('Ошибка при запуске игры:', error);
      setSearchStatus('Ошибка при запуске игры');
      setIsSearching(false);
    }
  };

  const handlePieceDrop = (from: string, to: string) => {
    if (!playerColor || currentTurn !== playerColor) {
      setError(`Сейчас ход ${currentTurn === 'white' ? 'белых' : 'черных'}`);
      return false;
    }

    try {
      const move = chessRef.current.move({ from, to });
      if (move === null) {
        setError('Недопустимый ход');
        return false;
      }

      if (socketRef.current?.readyState === WebSocket.OPEN) {
        const moveStr = `${from}${to}`;
        console.log('Sending move:', moveStr);
        socketRef.current.send(moveStr);
        setPosition(chessRef.current.fen());
        setMoveHistory(prev => [...prev, moveStr]);
        setCurrentTurn(prev => prev === 'white' ? 'black' : 'white');
        return true;
      } else {
        setError('Нет соединения с сервером');
        return false;
      }
    } catch (error) {
      console.error('Error making move:', error);
      setError('Ошибка при выполнении хода');
      return false;
    }
  };

  // WebSocket and move handling
  useEffect(() => {
    if (!gameId) return;

    const accessToken = localStorage.getItem('access_token');
    const userId = localStorage.getItem('userId');
    
    if (!accessToken || !userId) {
      navigate('/login');
      return;
    }

    const ws = new WebSocket(`${config.wsBaseUrl}${config.endpoints.game.play}/${gameId}?token=${encodeURIComponent('Bearer ' + accessToken)}`);
    socketRef.current = ws;

    ws.onopen = () => {
      console.log('=== WebSocket Connected ===');
      console.log('User ID:', userId);
      console.log('Player Color:', playerColor);
      console.log('========================');
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      console.log('=== Received Message ===');
      console.log('Player Color:', playerColor);
      console.log('Raw message:', event.data);
      
      try {
        const data = JSON.parse(event.data);
        console.log('Parsed data:', data);
        
        // Обработка информации о игре
        if (data.ok === true) {
          console.log('=== Game Info ===');
          console.log('White User ID:', data.WhiteUserId);
          console.log('Black User ID:', data.BlackUserId);
          console.log('Current Move:', data.currentMove);
          console.log('History Moves:', data.historyMove);
          console.log('Game Status:', data.status);
          console.log('Game Result:', data.result);
          console.log('================');
          
          // Обновляем информацию об игре
          const userId = localStorage.getItem('userId');
          if (userId) {
            const userColor = data.WhiteUserId === userId ? 'white' : 
                            data.BlackUserId === userId ? 'black' : null;
            setPlayerColor(userColor);
            
            // Обновляем gameInfo
            setGameInfo({
              ID: gameId || '',
              CreatedAt: new Date().toISOString(),
              Result: data.result,
              Status: data.status,
              WhiteUser: { id: data.WhiteUserId, name: 'White Player', email: '' },
              BlackUser: { id: data.BlackUserId, name: 'Black Player', email: '' }
            });
          }
          
          // Обновляем состояние игры
          if (data.historyMove && data.historyMove.length > 0) {
            chessRef.current.reset();
            data.historyMove.forEach((move: string) => {
              try {
                chessRef.current.move({ from: move.slice(0, 2), to: move.slice(2, 4) });
              } catch (e) {
                console.error('Error processing move:', move, e);
              }
            });
            setPosition(chessRef.current.fen());
            setMoveHistory(data.historyMove);
            setCurrentTurn((data.currentMove % 2 === 0) ? 'white' : 'black');
          }
          
          // Обновляем статус игры
          if (data.status === 'finished' || data.status === 'aborted') {
            setGameResult(data.result);
          }
        }
        
        // Обработка хода
        if (typeof event.data === 'string' && event.data.length === 4) {
          const move = event.data;
          const from = move.substring(0, 2);
          const to = move.substring(2, 4);
          chessRef.current.move({ from, to });
          setPosition(chessRef.current.fen());
          setMoveHistory(prev => [...prev, move]);
          setCurrentTurn(prev => prev === 'white' ? 'black' : 'white');
        }
      } catch (error) {
        console.error('Error processing message:', error);
      }
    };

    ws.onclose = () => {
      console.log('=== WebSocket Disconnected ===');
      console.log('User ID:', userId);
      console.log('Player Color:', playerColor);
      console.log('===========================');
      setIsConnected(false);
    };

    return () => {
      ws.close();
    };
  }, [gameId, navigate, playerColor]);

  useEffect(() => {
    if (gameInfo) {
      if (gameInfo.Result) {
        setGameResult(gameInfo.Result);
      }
      
      if (gameInfo.HistoryMove && gameInfo.HistoryMove.length > 0) {
        chessRef.current.reset();
        gameInfo.HistoryMove.forEach(move => {
          try {
            chessRef.current.move(move);
          } catch (e) {
            console.error('Invalid move in history:', move, e);
          }
        });
        setPosition(chessRef.current.fen());
      }
    }
  }, [gameInfo]);

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('userId');
    navigate('/login');
  };

  const handleMoveBack = () => {
    if (currentMoveIndex > 0) {
      const newIndex = currentMoveIndex - 1;
      setCurrentMoveIndex(newIndex);
      const chess = new Chess();
      moveHistory.slice(0, newIndex + 1).forEach(move => {
        chess.move({ from: move.slice(0, 2), to: move.slice(2, 4) });
      });
      setTemporaryPosition(chess.fen());
    }
  };

  const handleMoveForward = () => {
    if (currentMoveIndex < moveHistory.length - 1) {
      const newIndex = currentMoveIndex + 1;
      setCurrentMoveIndex(newIndex);
      const chess = new Chess();
      moveHistory.slice(0, newIndex + 1).forEach(move => {
        chess.move({ from: move.slice(0, 2), to: move.slice(2, 4) });
      });
      setTemporaryPosition(chess.fen());
    }
  };

  const handleMoveClick = (index: number) => {
    setCurrentMoveIndex(index);
    const chess = new Chess();
    moveHistory.slice(0, index + 1).forEach(move => {
      chess.move({ from: move.slice(0, 2), to: move.slice(2, 4) });
    });
    setTemporaryPosition(chess.fen());
  };

  // Update currentMoveIndex when new moves are received
  useEffect(() => {
    if (moveHistory.length > 0) {
      setCurrentMoveIndex(moveHistory.length - 1);
      setTemporaryPosition(position);
    }
  }, [moveHistory, position]);

  // Автоматически очищаем ошибку через 3 секунды
  useEffect(() => {
    if (error) {
      const timer = setTimeout(() => {
        setError(null);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [error]);

  return { 
    gameInfo, 
    position: temporaryPosition,
    moveHistory, 
    playerColor, 
    error, 
    setError,
    isLoading, 
    isSearching,
    searchStatus,
    handlePieceDrop,
    startNewGame,
    gameResult,
    handleLogout,
    currentMoveIndex,
    handleMoveBack,
    handleMoveForward,
    handleMoveClick,
    isCheck,
    isMate,
    isConnected
  };
};

// Component
const ChessBoardPage: React.FC = () => {
  const { gameId } = useParams();
  const navigate = useNavigate();
  const [showResult, setShowResult] = useState(false);
  const handleResultClick = () => setShowResult(false);

  const { 
    gameInfo, 
    position: temporaryPosition,
    moveHistory, 
    playerColor, 
    error, 
    setError,
    isLoading, 
    isSearching,
    searchStatus,
    handlePieceDrop,
    startNewGame,
    gameResult,
    handleLogout,
    currentMoveIndex,
    handleMoveBack,
    handleMoveForward,
    handleMoveClick,
    isCheck,
    isMate,
    isConnected
  } = useChessGame(gameId);

  // Определяем, кто пользователь
  const userId = localStorage.getItem('userId');
  const isWhitePlayer = gameInfo?.WhiteUser?.id === userId;
  const isBlackPlayer = gameInfo?.BlackUser?.id === userId;

  const getGameResultText = (result: string) => {
    const userId = localStorage.getItem('userId');
    if (!gameInfo || !userId) return 'Игра завершена';

    const isWhitePlayer = gameInfo.WhiteUser?.id === userId;
    const isBlackPlayer = gameInfo.BlackUser?.id === userId;

    switch (result) {
      case '1-0':
        if (isWhitePlayer) {
          return 'ПОБЕДА';
        } else if (isBlackPlayer) {
          return 'ПОРАЖЕНИЕ';
        }
        return 'Белые победили';
      case '0-1':
        if (isBlackPlayer) {
          return 'ПОБЕДА';
        } else if (isWhitePlayer) {
          return 'ПОРАЖЕНИЕ';
        }
        return 'Чёрные победили';
      case '1/2-1/2':
        return 'НИЧЬЯ';
      default:
        return 'Игра завершена';
    }
  };

  if (isLoading) {
    return (
      <BoardContainer>
        <div style={{ fontSize: '20px', color: '#333' }}>Загрузка...</div>
      </BoardContainer>
    );
  }

  return (
    <BoardContainer>
      <Header />
      <MainContent>
        {!gameId && (
          <>
            <WelcomeMessage>
              <h1>Добро пожаловать в GopherChess!</h1>
              <p>Начните новую игру и сразитесь с противником</p>
            </WelcomeMessage>
            <BoardWrapper>
              <MutedChessboard>
                <Chessboard
                  position="start"
                  boardWidth={600}
                  customBoardStyle={{
                    borderRadius: '4px',
                    boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
                  }}
                  customDarkSquareStyle={{ backgroundColor: '#779556' }}
                  customLightSquareStyle={{ backgroundColor: '#ebecd0' }}
                />
              </MutedChessboard>
              <CenteredStartButton 
                onClick={startNewGame}
                disabled={isSearching}
              >
                <ButtonIcon>♟️</ButtonIcon>
                {isSearching ? 'Поиск...' : 'Начать игру'}
              </CenteredStartButton>
            </BoardWrapper>
            {searchStatus && <GameStatus>{searchStatus}</GameStatus>}
          </>
        )}

        {error && (
          <ErrorMessageComponent 
            message={error} 
            onClose={() => setError(null)} 
          />
        )}

        {gameInfo && gameId && (
          <BoardWrapper>
            <MoveHistoryContainer>
              <MoveHistoryTitle>История ходов</MoveHistoryTitle>
              <MoveHistoryList>
                {moveHistory.map((move, index) => (
                  <MoveItem 
                    key={index}
                    isCurrent={index === currentMoveIndex}
                    onClick={() => handleMoveClick(index)}
                  >
                    {Math.floor(index/2) + 1}. {move}
                  </MoveItem>
                ))}
              </MoveHistoryList>
              <MoveControls>
                <ControlButton 
                  onClick={handleMoveBack}
                  disabled={currentMoveIndex <= 0}
                >
                  ⏮️ Назад
                </ControlButton>
                <ControlButton 
                  onClick={handleMoveForward}
                  disabled={currentMoveIndex >= moveHistory.length - 1}
                >
                  Вперед ⏭️
                </ControlButton>
              </MoveControls>
            </MoveHistoryContainer>
            <ChessBoardContainer style={{position: 'relative'}}>
              {/* Верхний блок — соперник */}
              {playerColor === 'white' ? (
                <BlackPlayerInfo>
                  <PlayerAvatar>B</PlayerAvatar>
                  <div>
                    <PlayerName>{gameInfo?.BlackUser?.name}</PlayerName>
                    <PlayerColor>Черные</PlayerColor>
                    <div style={{ fontSize: '12px', color: '#666' }}>{gameInfo?.BlackUser?.email}</div>
                  </div>
                </BlackPlayerInfo>
              ) : (
                <WhitePlayerInfo>
                  <PlayerAvatar>W</PlayerAvatar>
                  <div>
                    <PlayerName>{gameInfo?.WhiteUser?.name}</PlayerName>
                    <PlayerColor>Белые</PlayerColor>
                    <div style={{ fontSize: '12px', color: '#666' }}>{gameInfo?.WhiteUser?.email}</div>
                  </div>
                </WhitePlayerInfo>
              )}
              <ChessBoardWrapper playerColor={playerColor}>
                <Chessboard
                  position={temporaryPosition}
                  onPieceDrop={handlePieceDrop}
                  boardOrientation={playerColor || 'white'}
                  boardWidth={600}
                  customBoardStyle={{
                    borderRadius: '4px',
                    boxShadow: '0 2px 10px rgba(0,0,0,0.1)'
                  }}
                  customDarkSquareStyle={{ backgroundColor: '#779556' }}
                  customLightSquareStyle={{ backgroundColor: '#ebecd0' }}
                  showBoardNotation={true}
                />
                {gameResult && gameResult !== "0-0" && (
                  <NewGameButton 
                    onClick={startNewGame}
                    disabled={isSearching}
                  >
                    <ButtonIcon>♟️</ButtonIcon>
                    {isSearching ? 'Поиск...' : 'Новая игра'}
                  </NewGameButton>
                )}
              </ChessBoardWrapper>
              {/* Нижний блок — пользователь */}
              {playerColor === 'white' ? (
                <WhitePlayerInfo>
                  <PlayerAvatar>W</PlayerAvatar>
                  <div>
                    <PlayerName>{gameInfo?.WhiteUser?.name} <span style={{color:'#27ae60', fontWeight:'bold', fontSize:'14px', marginLeft:'6px'}}>Вы</span></PlayerName>
                    <PlayerColor>Белые</PlayerColor>
                    <div style={{ fontSize: '12px', color: '#666' }}>{gameInfo?.WhiteUser?.email}</div>
                  </div>
                </WhitePlayerInfo>
              ) : (
                <BlackPlayerInfo>
                  <PlayerAvatar>B</PlayerAvatar>
                  <div>
                    <PlayerName>{gameInfo?.BlackUser?.name} <span style={{color:'#27ae60', fontWeight:'bold', fontSize:'14px', marginLeft:'6px'}}>Вы</span></PlayerName>
                    <PlayerColor>Черные</PlayerColor>
                    <div style={{ fontSize: '12px', color: '#666' }}>{gameInfo?.BlackUser?.email}</div>
                  </div>
                </BlackPlayerInfo>
              )}
            </ChessBoardContainer>
            {showResult && gameResult && gameResult !== "0-0" && (
              <div 
                style={{ position: 'fixed', top: 40, left: '50%', transform: 'translateX(-50%)', zIndex: 2000, cursor: 'pointer' }}
                onClick={handleResultClick}
                title="Нажмите, чтобы скрыть"
              >
                {gameResult === '1-0' && isWhitePlayer && (
                  <VictoryResult>ПОБЕДА</VictoryResult>
                )}
                {gameResult === '1-0' && isBlackPlayer && (
                  <DefeatResult>ПОРАЖЕНИЕ</DefeatResult>
                )}
                {gameResult === '0-1' && isBlackPlayer && (
                  <VictoryResult>ПОБЕДА</VictoryResult>
                )}
                {gameResult === '0-1' && isWhitePlayer && (
                  <DefeatResult>ПОРАЖЕНИЕ</DefeatResult>
                )}
                {gameResult === '1/2-1/2' && (
                  <DrawResult>НИЧЬЯ</DrawResult>
                )}
              </div>
            )}
          </BoardWrapper>
        )}
      </MainContent>
    </BoardContainer>
  );
};

interface ErrorMessageProps {
  message: string;
  onClose: () => void;
}

const ErrorMessageComponent: React.FC<ErrorMessageProps> = ({ message, onClose }) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      onClose();
    }, 3000);
    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <ErrorMessage onClick={onClose}>
      {message}
    </ErrorMessage>
  );
};

export default ChessBoardPage;
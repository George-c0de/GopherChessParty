import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { Chess } from 'chess.js';
import { config } from '../config';
import { Chessboard } from 'react-chessboard';  // измените импорт

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
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f0f0;
  position: relative;
  padding: 20px;
`;

const MainContent = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  margin-top: 80px;
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

const BoardWrapper = styled.div`
  display: flex;
  gap: 20px;
  align-items: flex-start;
  position: relative;
  margin-top: 100px;
  margin-bottom: 100px;
`;

const ChessBoardContainer = styled.div`
  position: relative;
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  margin: 0 auto;
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
  position: absolute;
  background: rgba(255, 255, 255, 0.95);
  padding: 12px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 200px;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(0,0,0,0.15);
  }
`;

const WhitePlayerInfo = styled(PlayerInfo)`
  top: -80px;
  left: 50%;
  transform: translateX(-50%);
  border: 2px solid #ebecd0;
  z-index: 10;
`;

const BlackPlayerInfo = styled(PlayerInfo)`
  bottom: -80px;
  left: 50%;
  transform: translateX(-50%);
  border: 2px solid #779556;
  z-index: 10;
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
`;

const StartGameButton = styled.button`
  position: absolute;
  top: -50px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  background: linear-gradient(135deg, #2ecc71 0%, #27ae60 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
  transition: all 0.3s ease;
  z-index: 10;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 8px;

  &:hover {
    transform: translateX(-50%) translateY(-2px);
    box-shadow: 0 6px 12px rgba(0,0,0,0.15);
    background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  }

  &:active {
    transform: translateX(-50%) translateY(0);
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  &:disabled {
    background: #cccccc;
    cursor: not-allowed;
    transform: translateX(-50%);
    box-shadow: none;
  }
`;

const NewGameButton = styled(StartGameButton)`
  position: fixed;
  top: auto;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background-color: #2196F3;
  z-index: 10;

  &:hover {
    background-color: #1976D2;
  }

  &:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
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
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(46, 204, 113, 0.3);
  z-index: 1000;
  font-size: 26px;
  font-weight: 800;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  animation: subtleAppear 0.5s ease-in-out;
  border: 2px solid #27ae60;
`;

const DefeatResult = styled.div`
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #7f8c8d 0%, #95a5a6 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(149, 165, 166, 0.3);
  z-index: 1000;
  font-size: 26px;
  font-weight: 800;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  animation: subtleAppear 0.5s ease-in-out;
  border: 2px solid #7f8c8d;
`;

const DrawResult = styled.div`
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
  color: white;
  padding: 16px 32px;
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(41, 128, 185, 0.3);
  z-index: 1000;
  font-size: 26px;
  font-weight: 800;
  text-align: center;
  text-transform: uppercase;
  letter-spacing: 1px;
  animation: subtleAppear 0.5s ease-in-out;
  border: 2px solid #2c3e50;
`;

const NavigationButton = styled.button`
  padding: 12px 24px;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(0,0,0,0.15);
  }

  &:active {
    transform: translateY(0);
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
`;

const NavigationButtons = styled.div`
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  gap: 10px;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.9);
  padding: 10px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
`;

const LogoutButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);

  &:hover {
    background: linear-gradient(135deg, #c0392b 0%, #e74c3c 100%);
  }
`;

const ProfileButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #2ecc71 0%, #27ae60 100%);

  &:hover {
    background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  }
`;

const HistoryButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #9b59b6 0%, #8e44ad 100%);

  &:hover {
    background: linear-gradient(135deg, #8e44ad 0%, #9b59b6 100%);
  }
`;

const HomeButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #9b59b6 0%, #8e44ad 100%);

  &:hover {
    background: linear-gradient(135deg, #8e44ad 0%, #9b59b6 100%);
  }
`;

const ButtonIcon = styled.span`
  font-size: 18px;
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

  const chessRef = React.useRef(new Chess());
  const socketRef = React.useRef<WebSocket | null>(null);

  // Initialize game: fetch info
  useEffect(() => {
    async function init() {
      if (!gameId) {
        setIsLoading(false);
        return;
      }

      const token = localStorage.getItem('authToken');
      const userId = localStorage.getItem('userId');
      
      if (!token || !userId) {
        navigate('/login');
        return;
      }

      try {
        setIsLoading(true);
        const res = await fetch(`${config.apiBaseUrl}/chess/${gameId}`, {
          headers: { 
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!res.ok) {
          throw new Error('Failed to fetch game info');
        }

        const data: GameInfo = await res.json();
        setGameInfo(data);
        
        const userColor = data.WhiteUser.id === userId ? 'white' : 
                         data.BlackUser.id === userId ? 'black' : null;
        setPlayerColor(userColor);

        if (!userColor) {
          setError('Вы не являетесь участником этой игры');
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
    const token = localStorage.getItem('authToken');
    if (!token) {
      navigate('/login');
      return;
    }

    setIsSearching(true);
    setSearchStatus('Поиск противника...');

    try {
      const wsUrl = `${config.wsBaseUrl}/ws/search?token=${encodeURIComponent('Bearer ' + token)}`;
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
    if (!gameInfo) return;
    const token = localStorage.getItem('authToken')!;
    const ws = new WebSocket(`${config.wsBaseUrl}/ws/game/${gameId}?token=${encodeURIComponent('Bearer ' + token)}`);
    socketRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket connection established for game:', gameId);
      setError(null);
    };

    ws.onmessage = ({ data: msg }) => {
      try {
        console.log('Received WebSocket message:', msg);
        const data: WSMessage = JSON.parse(msg);
        
        // Handle game result first
        if (data.result && data.result !== "0-0") {
          console.log('Game result received:', data.result);
          setGameResult(data.result);
          return;
        }
        
        // initial history
        if (data.historyMove) {
          console.log('Processing initial history:', data.historyMove);
          chessRef.current.reset();
          data.historyMove.forEach((m: string) => {
            try {
              chessRef.current.move({ from: m.slice(0, 2), to: m.slice(2, 4) });
            } catch (moveError) {
              console.error('Error processing move:', m, moveError);
            }
          });
          setPosition(chessRef.current.fen());
          setMoveHistory(data.historyMove);
          setCurrentTurn((data.currentMove! % 2 === 0) ? 'white' : 'black');
          return;
        }
        
        // single move from opponent
        if (data.move) {
          console.log('Processing opponent move:', data.move);
          try {
            chessRef.current.move({ from: data.move.slice(0, 2), to: data.move.slice(2, 4) });
            setPosition(chessRef.current.fen());
            setMoveHistory(prev => [...prev, data.move!]);
            setCurrentTurn(prev => prev === 'white' ? 'black' : 'white');
          } catch (moveError) {
            console.error('Error processing opponent move:', data.move, moveError);
            setError('Ошибка при обработке хода противника');
          }
          return;
        }
        
        // server response with result
        if (data.ok === true && data.result && data.result !== "0-0") {
          console.log('Server response with result:', data.result);
          setGameResult(data.result);
          return;
        }
        
        // server response
        if (typeof data.ok === 'boolean') {
          if (!data.ok && data.message) {
            console.error('Server error message:', data.message);
            setError(data.message);
          }
        }
        
        if (data.error) {
          console.error('Server error:', data.error);
          setError(data.error);
        }
      } catch (parseError) {
        console.error('Error parsing WebSocket message:', parseError);
        console.error('Raw message:', msg);
        setError('Ошибка при обработке сообщения от сервера');
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      setError('Ошибка соединения с сервером');
    };

    ws.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason);
      if (event.code !== 1000) {
        setError('Соединение с сервером прервано');
      }
      socketRef.current = null;
    };

    return () => { 
      console.log('Cleaning up WebSocket connection');
      if (socketRef.current) {
        socketRef.current.close(1000, 'Component unmounting');
        socketRef.current = null;
      }
    };
  }, [gameInfo, gameId]);

  useEffect(() => {
    if (gameInfo) {
      if (gameInfo.Result) {
        setGameResult(gameInfo.Result);
      }
      
      if (gameInfo.HistoryMove && gameInfo.HistoryMove.length > 0) {
        const chess = new Chess();
        gameInfo.HistoryMove.forEach(move => {
          try {
            chess.move(move);
          } catch (e) {
            console.error('Invalid move in history:', move, e);
          }
        });
        setPosition(chess.fen());
      }
    }
  }, [gameInfo]);

  const handleLogout = () => {
    localStorage.removeItem('authToken');
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

  return { 
    gameInfo, 
    position: temporaryPosition,
    moveHistory, 
    playerColor, 
    error, 
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
    handleMoveClick
  };
};

// Component
const ChessBoardPage: React.FC = () => {
  const { gameId } = useParams();
  const navigate = useNavigate();
  const { 
    gameInfo, 
    position: temporaryPosition,
    moveHistory, 
    playerColor, 
    error, 
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
    handleMoveClick
  } = useChessGame(gameId);

  const getGameResultText = (result: string) => {
    const userId = localStorage.getItem('userId');
    if (!gameInfo || !userId) return 'Игра завершена';

    const isWhitePlayer = gameInfo.WhiteUser.id === userId;
    const isBlackPlayer = gameInfo.BlackUser.id === userId;

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
      <NavigationButtons>
        <LogoutButton onClick={handleLogout}>
          <ButtonIcon>🚪</ButtonIcon>
          Выйти
        </LogoutButton>
        <ProfileButton onClick={() => navigate('/profile')}>
          <ButtonIcon>👤</ButtonIcon>
          Профиль
        </ProfileButton>
        <HistoryButton onClick={() => navigate('/history')}>
          <ButtonIcon>📜</ButtonIcon>
          История
        </HistoryButton>
        <HomeButton onClick={() => {
          navigate('/game');
          window.location.reload();
        }}>
          <ButtonIcon>🏠</ButtonIcon>
          Главная
        </HomeButton>
      </NavigationButtons>

      {!gameId && (
        <MainContent>
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
        </MainContent>
      )}

      {error && <ErrorMessage>{error}</ErrorMessage>}

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
          <ChessBoardContainer>
            {gameInfo && (
              <>
                <WhitePlayerInfo>
                  <PlayerAvatar>W</PlayerAvatar>
                  <div>
                    <PlayerName>{gameInfo.WhiteUser.name}</PlayerName>
                    <PlayerColor>Белые</PlayerColor>
                  </div>
                </WhitePlayerInfo>
                <BlackPlayerInfo>
                  <PlayerAvatar>B</PlayerAvatar>
                  <div>
                    <PlayerName>{gameInfo.BlackUser.name}</PlayerName>
                    <PlayerColor>Черные</PlayerColor>
                  </div>
                </BlackPlayerInfo>
              </>
            )}
            <Chessboard
              position={temporaryPosition}
              onPieceDrop={handlePieceDrop}
              boardOrientation={playerColor || 'white'}
              boardWidth={600}
              customBoardStyle={{
                borderRadius: '4px',
                boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
              }}
              customDarkSquareStyle={{ backgroundColor: '#779556' }}
              customLightSquareStyle={{ backgroundColor: '#ebecd0' }}
            />
          </ChessBoardContainer>
        </BoardWrapper>
      )}

      {gameResult && gameResult !== "0-0" && (
        <>
          {gameResult === '1-0' && gameInfo?.WhiteUser.id === localStorage.getItem('userId') && (
            <VictoryResult>ПОБЕДА</VictoryResult>
          )}
          {gameResult === '1-0' && gameInfo?.BlackUser.id === localStorage.getItem('userId') && (
            <DefeatResult>ПОРАЖЕНИЕ</DefeatResult>
          )}
          {gameResult === '0-1' && gameInfo?.BlackUser.id === localStorage.getItem('userId') && (
            <VictoryResult>ПОБЕДА</VictoryResult>
          )}
          {gameResult === '0-1' && gameInfo?.WhiteUser.id === localStorage.getItem('userId') && (
            <DefeatResult>ПОРАЖЕНИЕ</DefeatResult>
          )}
          {gameResult === '1/2-1/2' && (
            <DrawResult>НИЧЬЯ</DrawResult>
          )}
          <CenteredStartButton 
            onClick={startNewGame}
            disabled={isSearching}
          >
            <ButtonIcon>♟️</ButtonIcon>
            {isSearching ? 'Поиск...' : 'Новая игра'}
          </CenteredStartButton>
        </>
      )}
    </BoardContainer>
  );
};

export default ChessBoardPage;
import React, { useState, useEffect, useRef, useCallback } from 'react';
import styled, { keyframes } from 'styled-components';
import { useNavigate } from 'react-router-dom';
import { config } from '../config';
import Header from './Header';
import { MainContent } from './MainContent';

const Container = styled.div`
  max-width: 800px;
  margin: 100px auto 60px;
  padding: 30px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
`;

const Title = styled.h2`
  color: #2c3e50;
  margin-bottom: 30px;
  text-align: center;
`;

const GameList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
`;

const fadeIn = keyframes`
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`;

const GameCard = styled.div`
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  transition: all 0.3s ease;
  cursor: pointer;
  animation: ${fadeIn} 0.5s ease-out;
  margin-bottom: 15px;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }
`;

const GameHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
`;

const GameId = styled.span`
  color: #6c757d;
  font-size: 14px;
`;

const GameStatus = styled.span<{ status: string }>`
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  background-color: ${props => {
    switch (props.status) {
      case 'finished':
        return '#28a745';
      case 'playing':
        return '#17a2b8';
      case 'aborted':
        return '#dc3545';
      case 'waiting':
        return '#ffc107';
      default:
        return '#6c757d';
    }
  }};
  color: white;
`;

const GameResult = styled.span<{ result: string, isCurrentUser: boolean }>`
  font-weight: 500;
  color: ${props => {
    if (props.result === '1-0') {
      return props.isCurrentUser ? '#28a745' : '#dc3545';
    } else if (props.result === '0-1') {
      return props.isCurrentUser ? '#28a745' : '#dc3545';
    } else if (props.result === '1/2-1/2') {
      return '#6c757d';
    }
    return '#6c757d';
  }};
`;

const GameInfo = styled.div`
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 14px;
  color: #6c757d;
`;

const PlayerInfo = styled.div`
  display: flex;
  flex-direction: column;
  gap: 5px;
`;

const PlayerName = styled.span`
  font-weight: 500;
  color: #2c3e50;
`;

const PlayerColor = styled.span`
  font-size: 12px;
  color: #6c757d;
`;

const ErrorMessage = styled.div`
  color: #e74c3c;
  text-align: center;
  margin-bottom: 20px;
`;

const LoadingMessage = styled.div`
  text-align: center;
  color: #6c757d;
  font-size: 18px;
`;

const NavigationButtons = styled.div`
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  gap: 10px;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.95);
  padding: 10px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  backdrop-filter: blur(5px);
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

const HomeButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);

  &:hover {
    background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
  }
`;

const ButtonIcon = styled.span`
  font-size: 18px;
`;

const LoadingIndicator = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  color: #6c757d;
  font-size: 16px;
`;

const LoadingSpinner = styled.div`
  width: 20px;
  height: 20px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 10px;

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;

const GameModal = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
`;

const GameModalContent = styled.div`
  background: white;
  padding: 30px;
  border-radius: 12px;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
`;

const GameModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
`;

const GameModalTitle = styled.h2`
  color: #2c3e50;
  margin: 0;
`;

const GameModalClose = styled.button`
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6c757d;
  padding: 0;
  
  &:hover {
    color: #2c3e50;
  }
`;

const GameModalPlayers = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
  gap: 20px;
`;

const GameModalPlayer = styled.div`
  flex: 1;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
  text-align: center;
`;

const GameModalPlayerName = styled.div`
  font-weight: bold;
  font-size: 18px;
  margin-bottom: 5px;
`;

const GameModalPlayerColor = styled.div`
  color: #6c757d;
  font-size: 14px;
`;

const GameModalInfo = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 20px;
`;

const GameModalInfoRow = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 8px;
`;

const GameModalInfoLabel = styled.span`
  color: #6c757d;
`;

const GameModalInfoValue = styled.span`
  font-weight: 500;
`;

const GameModalButton = styled.button`
  padding: 12px 24px;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
  transition: all 0.3s ease;
  margin-top: 20px;
  width: 100%;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  }

  &:active {
    transform: translateY(0);
    box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  }
`;

interface ChessGame {
  id: string;
  created_at: string;
  updated_at: string;
  status: string;
  result: string;
  white_player?: {
    id: string;
    name: string;
  };
  black_player?: {
    id: string;
    name: string;
  };
}

export const ChessHistory: React.FC = () => {
  const [games, setGames] = useState<ChessGame[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [currentUserId, setCurrentUserId] = useState<string | null>(null);
  const navigate = useNavigate();
  const [selectedGame, setSelectedGame] = useState<ChessGame | null>(null);

  useEffect(() => {
    const userId = localStorage.getItem('userId');
    setCurrentUserId(userId);
  }, []);

  const getPlayerName = (player: { id: string; name: string } | undefined) => {
    if (!player) return 'Неизвестный игрок';
    return player.id === currentUserId ? 'Вы' : player.name;
  };

  const fetchGames = async () => {
    const token = localStorage.getItem('authToken');
    if (!token) {
      navigate('/login');
      return;
    }

    try {
      const response = await fetch(`${config.apiBaseUrl}${config.endpoints.chess}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Failed to fetch games');
      }

      const data = await response.json();
      setGames(data.items || []);
    } catch (err) {
      setError('Ошибка при загрузке истории игр');
      console.error('Error fetching games:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchGames();
  }, []);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'finished':
        return 'Завершена';
      case 'playing':
        return 'В процессе';
      case 'aborted':
        return 'Прервана';
      case 'waiting':
        return 'Ожидание';
      default:
        return status;
    }
  };

  const getResultText = (result: string) => {
    switch (result) {
      case '1-0':
        return 'Победа белых';
      case '0-1':
        return 'Победа черных';
      case '1/2-1/2':
        return 'Ничья';
      default:
        return result;
    }
  };

  const handleGameClick = (game: ChessGame) => {
    setSelectedGame(game);
  };

  const handleCloseModal = () => {
    setSelectedGame(null);
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('userId');
    navigate('/login');
  };

  const handleGoHome = () => {
    navigate('/game');
  };

  const isCurrentUserWinner = (game: ChessGame) => {
    if (!currentUserId) return false;
    if (game.result === '1-0' && game.white_player?.id === currentUserId) return true;
    if (game.result === '0-1' && game.black_player?.id === currentUserId) return true;
    return false;
  };

  if (isLoading) {
    return (
      <>
        <Header />
        <MainContent>
          <Container>
            <LoadingMessage>Загрузка истории игр...</LoadingMessage>
          </Container>
        </MainContent>
      </>
    );
  }

  if (error) {
    return (
      <>
        <Header />
        <MainContent>
          <Container>
            <ErrorMessage>{error}</ErrorMessage>
          </Container>
        </MainContent>
      </>
    );
  }

  return (
    <>
      <Header />
      <MainContent>
        <Container>
          <Title>История игр</Title>
          <GameList>
            {games.map((game) => (
              <GameCard 
                key={game.id} 
                onClick={() => handleGameClick(game)}
              >
                <GameHeader>
                  <GameId>Игра #{game.id.slice(0, 8)}</GameId>
                  <GameStatus status={game.status}>{getStatusText(game.status)}</GameStatus>
                </GameHeader>
                <GameResult result={game.result} isCurrentUser={isCurrentUserWinner(game)}>
                  {getResultText(game.result)}
                </GameResult>
                <GameInfo>
                  <PlayerInfo>
                    <PlayerName>{getPlayerName(game.white_player)}</PlayerName>
                    <PlayerColor>Белые</PlayerColor>
                  </PlayerInfo>
                  <PlayerInfo>
                    <PlayerName>{getPlayerName(game.black_player)}</PlayerName>
                    <PlayerColor>Черные</PlayerColor>
                  </PlayerInfo>
                </GameInfo>
                <GameInfo>
                  <span>Создана: {formatDate(game.created_at)}</span>
                  <span>Обновлена: {formatDate(game.updated_at)}</span>
                </GameInfo>
              </GameCard>
            ))}
          </GameList>
        </Container>

        {selectedGame && (
          <GameModal>
            <GameModalContent>
              <GameModalHeader>
                <GameModalTitle>Игра #{selectedGame.id.slice(0, 8)}</GameModalTitle>
                <GameModalClose onClick={handleCloseModal}>×</GameModalClose>
              </GameModalHeader>
              
              <GameModalPlayers>
                <GameModalPlayer>
                  <GameModalPlayerName>{getPlayerName(selectedGame.white_player)}</GameModalPlayerName>
                  <GameModalPlayerColor>Белые</GameModalPlayerColor>
                </GameModalPlayer>
                <GameModalPlayer>
                  <GameModalPlayerName>{getPlayerName(selectedGame.black_player)}</GameModalPlayerName>
                  <GameModalPlayerColor>Черные</GameModalPlayerColor>
                </GameModalPlayer>
              </GameModalPlayers>

              <GameModalInfo>
                <GameModalInfoRow>
                  <GameModalInfoLabel>Статус:</GameModalInfoLabel>
                  <GameModalInfoValue>{getStatusText(selectedGame.status)}</GameModalInfoValue>
                </GameModalInfoRow>
                <GameModalInfoRow>
                  <GameModalInfoLabel>Результат:</GameModalInfoLabel>
                  <GameModalInfoValue style={{ 
                    color: isCurrentUserWinner(selectedGame) ? '#28a745' : '#dc3545'
                  }}>
                    {getResultText(selectedGame.result)}
                  </GameModalInfoValue>
                </GameModalInfoRow>
                <GameModalInfoRow>
                  <GameModalInfoLabel>Создана:</GameModalInfoLabel>
                  <GameModalInfoValue>{formatDate(selectedGame.created_at)}</GameModalInfoValue>
                </GameModalInfoRow>
                <GameModalInfoRow>
                  <GameModalInfoLabel>Обновлена:</GameModalInfoLabel>
                  <GameModalInfoValue>{formatDate(selectedGame.updated_at)}</GameModalInfoValue>
                </GameModalInfoRow>
              </GameModalInfo>

              <GameModalButton onClick={() => navigate(`/game/${selectedGame.id}`)}>
                Просмотреть игру
              </GameModalButton>
            </GameModalContent>
          </GameModal>
        )}
      </MainContent>
    </>
  );
}; 
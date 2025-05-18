import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import { config } from '../config';
import Header from './Header';
import { MainContent } from './MainContent';

const Container = styled.div`
  max-width: 600px;
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

const HomeButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);

  &:hover {
    background: linear-gradient(135deg, #2980b9 0%, #3498db 100%);
  }
`;

const HistoryButton = styled(NavigationButton)`
  background: linear-gradient(135deg, #2ecc71 0%, #27ae60 100%);

  &:hover {
    background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  }
`;

const ButtonIcon = styled.span`
  font-size: 18px;
`;

const UserInfo = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const InfoRow = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
`;

const InfoLabel = styled.span`
  color: #6c757d;
`;

const InfoValue = styled.span`
  font-weight: 500;
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

interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export const ProfilePage: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUserData = async () => {
      const token = localStorage.getItem('access_token');
      if (!token) {
        navigate('/login');
        return;
      }

      try {
        const response = await fetch(`${config.apiBaseUrl}${config.endpoints.user}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!response.ok) {
          throw new Error('Failed to fetch user data');
        }

        const data = await response.json();
        setUser(data.item);
      } catch (err) {
        setError('Ошибка при загрузке данных пользователя');
        console.error('Error fetching user data:', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserData();
  }, [navigate]);

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('userId');
    navigate('/login');
  };

  const handleGoHome = () => {
    navigate('/game');
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <>
      <Header />
      <MainContent>
        <Container>
          <Title>Личный кабинет</Title>
          {isLoading ? (
            <LoadingMessage>Загрузка данных пользователя...</LoadingMessage>
          ) : error ? (
            <ErrorMessage>{error}</ErrorMessage>
          ) : user ? (
            <UserInfo>
              <InfoRow>
                <InfoLabel>ID:</InfoLabel>
                <InfoValue>{user.id}</InfoValue>
              </InfoRow>
              <InfoRow>
                <InfoLabel>Имя:</InfoLabel>
                <InfoValue>{user.name}</InfoValue>
              </InfoRow>
              <InfoRow>
                <InfoLabel>Email:</InfoLabel>
                <InfoValue>{user.email}</InfoValue>
              </InfoRow>
              <InfoRow>
                <InfoLabel>Дата регистрации:</InfoLabel>
                <InfoValue>{formatDate(user.created_at)}</InfoValue>
              </InfoRow>
              <InfoRow>
                <InfoLabel>Последнее обновление:</InfoLabel>
                <InfoValue>{formatDate(user.updated_at)}</InfoValue>
              </InfoRow>
            </UserInfo>
          ) : null}
        </Container>
      </MainContent>
    </>
  );
}; 
import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import styled from 'styled-components';

const HeaderContainer = styled.header<{ visible: boolean }>`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  transform: translateY(${props => props.visible ? '0' : '-100%'});
  transition: transform 0.3s ease;
`;

const ToggleButton = styled.button<{ isVisible: boolean }>`
  padding: 8px 12px;
  background: transparent;
  color: #666;
  border: 1px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;

  &:hover {
    background: #f5f5f5;
    border-color: #ccc;
  }

  &:active {
    background: #eee;
  }
`;

const NavigationButtons = styled.div`
  display: flex;
  gap: 8px;
`;

const NavigationButton = styled.button`
  padding: 8px 12px;
  color: #666;
  background: transparent;
  border: 1px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;

  &:hover {
    background: #f5f5f5;
    border-color: #ccc;
  }

  &:active {
    background: #eee;
  }
`;

const LogoutButton = styled(NavigationButton)`
  color: #e74c3c;
  border-color: #e74c3c;
  &:hover {
    background: #fdf3f2;
  }
`;

const ProfileButton = styled(NavigationButton)`
  color: #2ecc71;
  border-color: #2ecc71;
  &:hover {
    background: #f2fdf6;
  }
`;

const HistoryButton = styled(NavigationButton)`
  color: #9b59b6;
  border-color: #9b59b6;
  &:hover {
    background: #f9f5fc;
  }
`;

const HomeButton = styled(NavigationButton)`
  color: #3498db;
  border-color: #3498db;
  &:hover {
    background: #f2f8fc;
  }
`;

const ButtonIcon = styled.span`
  font-size: 16px;
`;

const Header: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [isHovering, setIsHovering] = useState(false);
  const [scrollTimeout, setScrollTimeout] = useState<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      if (scrollTimeout) {
        clearTimeout(scrollTimeout);
      }

      if (isHovering) {
        setIsVisible(true);
        return;
      }

      if (currentScrollY > lastScrollY && currentScrollY > 100) {
        setIsVisible(false);
      } else {
        setIsVisible(true);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => {
      window.removeEventListener('scroll', handleScroll);
      if (scrollTimeout) {
        clearTimeout(scrollTimeout);
      }
    };
  }, [lastScrollY, isHovering, scrollTimeout]);

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('userId');
    navigate('/login');
  };

  const toggleHeader = () => {
    setIsVisible(!isVisible);
  };

  if (location.pathname === '/login' || location.pathname === '/register') {
    return null;
  }

  return (
    <>
      <HeaderContainer visible={isVisible || isHovering}>
        <ToggleButton 
          isVisible={isVisible || isHovering}
          onClick={toggleHeader}
        >
          <ButtonIcon>{isVisible ? '▼' : '▲'}</ButtonIcon>
          {isVisible ? 'Скрыть' : 'Показать'}
        </ToggleButton>
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
      </HeaderContainer>
    </>
  );
};

export default Header; 
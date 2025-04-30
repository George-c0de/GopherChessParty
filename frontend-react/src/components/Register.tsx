import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authApi } from '../services/api';
import {
    Container,
    FormGroup,
    Label,
    Input,
    Button,
    ErrorMessage,
    LinkContainer
} from './styled';

export const Register: React.FC = () => {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        if (password !== confirmPassword) {
            setError('Пароли не совпадают');
            return;
        }

        if (password.length < 6) {
            setError('Пароль должен содержать минимум 6 символов');
            return;
        }

        try {
            const response = await authApi.register({ name, email, password });
            localStorage.setItem('authToken', response.token);
            navigate('/game');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Ошибка регистрации');
        }
    };

    return (
        <Container>
            <h2>Регистрация</h2>
            <form onSubmit={handleSubmit}>
                <FormGroup>
                    <Label htmlFor="name">Имя пользователя:</Label>
                    <Input
                        type="text"
                        id="name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </FormGroup>
                <FormGroup>
                    <Label htmlFor="email">Email:</Label>
                    <Input
                        type="email"
                        id="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </FormGroup>
                <FormGroup>
                    <Label htmlFor="password">Пароль:</Label>
                    <Input
                        type="password"
                        id="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </FormGroup>
                <FormGroup>
                    <Label htmlFor="confirmPassword">Подтвердите пароль:</Label>
                    <Input
                        type="password"
                        id="confirmPassword"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                    />
                </FormGroup>
                <Button type="submit">Зарегистрироваться</Button>
                <ErrorMessage hidden={!error}>{error}</ErrorMessage>
            </form>
            <LinkContainer>
                Уже есть аккаунт? <Link to="/login">Войти</Link>
            </LinkContainer>
        </Container>
    );
}; 
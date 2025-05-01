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
import { jwtDecode } from "jwt-decode";

export const Login: React.FC = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        if (isSubmitting) return;
        
        setIsSubmitting(true);
        setError('');

        try {
            const response = await authApi.login({ email, password });
            let token = response.token;
            if (token && token.startsWith('Bearer ')) {
                token = token.slice(7);
            }
            localStorage.setItem('authToken', token);

            try {
                const payload: any = jwtDecode(token);
                if (payload && payload.id) {
                    localStorage.setItem('userId', payload.id);
                }
            } catch (e) {
                console.error('Error decoding token:', e);
            }

            navigate('/game');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Ошибка авторизации');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Container>
            <h2>Вход в систему</h2>
            <form onSubmit={handleSubmit}>
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
                <Button type="submit">Войти</Button>
                <ErrorMessage hidden={!error}>{error}</ErrorMessage>
            </form>
            <LinkContainer>
                Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
            </LinkContainer>
        </Container>
    );
}; 
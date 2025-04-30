import axios from 'axios';
import { config } from '../config';
import { AuthResponse, LoginRequest, RegisterRequest } from '../types/api';

const api = axios.create({
    baseURL: config.apiBaseUrl,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const authApi = {
    login: async (data: LoginRequest): Promise<AuthResponse> => {
        const response = await api.post<AuthResponse>(config.endpoints.login, data);
        return response.data;
    },

    register: async (data: RegisterRequest): Promise<AuthResponse> => {
        const response = await api.post<AuthResponse>(config.endpoints.register, data);
        return response.data;
    },
};

// Добавляем интерцептор для добавления токена к запросам
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}); 
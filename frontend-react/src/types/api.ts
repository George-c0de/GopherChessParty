export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
}

export interface AuthResponse {
    access_token: string;
    refresh_token: string;
    userId?: string;
    id?: string;
    user_id?: string;
}

export interface ErrorResponse {
    error: string;
} 
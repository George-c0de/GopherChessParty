import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Login } from './components/Login';
import { Register } from './components/Register';
import { ChessBoard } from './components/ChessBoard';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const token = localStorage.getItem('authToken');
    if (!token) {
        return <Navigate to="/login" replace />;
    }
    return <>{children}</>;
};

const App: React.FC = () => {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route
                    path="/game"
                    element={
                        <ProtectedRoute>
                            <ChessBoard />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/game/:gameId"
                    element={
                        <ProtectedRoute>
                            <ChessBoard />
                        </ProtectedRoute>
                    }
                />
                <Route path="/" element={<Navigate to="/game" replace />} />
            </Routes>
        </Router>
    );
};

export default App;

import React, { useState, useEffect } from 'react';
import { ToastProvider, useToast } from './components/Toast';
import { User } from './views/LandingView';
import { Router } from './router';
import api from './api';

function AppContent() {
  const [user, setUser] = useState<User | null>(() => {
    const stored = localStorage.getItem('sn_user');
    return stored ? JSON.parse(stored) : null;
  });
  const [hash, setHash] = useState(window.location.hash || '#/');
  const showToast = useToast();

  useEffect(() => {
    const handleHashChange = () => setHash(window.location.hash || '#/');
    window.addEventListener('hashchange', handleHashChange);
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, []);

  useEffect(() => {
    const handleUserUpdated = (e: Event) => {
      const customEvent = e as CustomEvent<User>;
      setUser(customEvent.detail);
    };

    const handleLoggedOut = () => {
      setUser(null);
      navigate('#/login');
      showToast("Session expired. Please log in again.", "error");
    };

    window.addEventListener('sn-user-updated', handleUserUpdated);
    window.addEventListener('sn-user-logged-out', handleLoggedOut);

    return () => {
      window.removeEventListener('sn-user-updated', handleUserUpdated);
      window.removeEventListener('sn-user-logged-out', handleLoggedOut);
    };
  }, []);

  const navigate = (toHash: string) => {
    window.location.hash = toHash;
  };

  const handleLoginSuccess = (userData: User) => {
    setUser(userData);
    localStorage.setItem('sn_user', JSON.stringify(userData));
    navigate('#/dashboard');
  };

  const handleLogout = () => {
    setUser(null);
    localStorage.removeItem('sn_user');
    navigate('#/');
    showToast("Logged out successfully", "success");
  };

  // Axios-based wrapper to compatibility-support old views
  const apiFetch = async (endpoint: string, options: RequestInit = {}) => {
    const method = (options.method || 'GET').toUpperCase();
    const headers = (options.headers as Record<string, string>) || {};
    let data: any = undefined;

    if (options.body) {
      try {
        data = JSON.parse(options.body as string);
      } catch (e) {
        data = options.body;
      }
    }

    try {
      const response = await api({
        url: endpoint,
        method,
        headers,
        data,
      });

      return {
        ok: true,
        status: response.status,
        json: async () => response.data,
        text: async () => JSON.stringify(response.data),
      } as Response;
    } catch (error: any) {
      if (error.response) {
        return {
          ok: false,
          status: error.response.status,
          json: async () => error.response.data,
          text: async () => JSON.stringify(error.response.data),
        } as Response;
      }
      throw error;
    }
  };

  return (
    <Router
      user={user}
      hash={hash}
      navigate={navigate}
      handleLoginSuccess={handleLoginSuccess}
      handleLogout={handleLogout}
      apiFetch={apiFetch}
    />
  );
}

export default function App() {
  return (
    <ToastProvider>
      <AppContent />
    </ToastProvider>
  );
}

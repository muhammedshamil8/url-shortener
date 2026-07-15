import React, { useEffect } from 'react';
import LandingView, { User } from './views/LandingView';
import LoginView from './views/LoginView';
import RegisterView from './views/RegisterView';
import DashboardView from './views/DashboardView';
import MyURLsView from './views/MyURLsView';
import AdminView from './views/AdminView';
import DashboardLayout from './components/DashboardLayout';

interface RouterProps {
  user: User | null;
  hash: string;
  navigate: (toHash: string) => void;
  handleLoginSuccess: (userData: User) => void;
  handleLogout: () => void;
  apiFetch: (endpoint: string, options?: RequestInit) => Promise<Response>;
}

export function Router({
  user,
  hash,
  navigate,
  handleLoginSuccess,
  handleLogout,
  apiFetch,
}: RouterProps) {
  const route = hash.split('?')[0];

  switch (route) {
    case '#/':
      return <LandingView user={user} apiFetch={apiFetch} navigate={navigate} />;
    case '#/login':
      return <LoginView onLoginSuccess={handleLoginSuccess} navigate={navigate} />;
    case '#/register':
      return <RegisterView navigate={navigate} />;
    case '#/dashboard':
      return user ? (
        <DashboardLayout user={user} onLogout={handleLogout} activeTab="dashboard" navigate={navigate}>
          <DashboardView user={user} apiFetch={apiFetch} />
        </DashboardLayout>
      ) : (
        <RedirectView to="#/login" />
      );
    case '#/my-urls':
      return user ? (
        <DashboardLayout user={user} onLogout={handleLogout} activeTab="my-urls" navigate={navigate}>
          <MyURLsView apiFetch={apiFetch} />
        </DashboardLayout>
      ) : (
        <RedirectView to="#/login" />
      );
    case '#/admin':
      return user && user.role === 'admin' ? (
        <DashboardLayout user={user} onLogout={handleLogout} activeTab="admin" navigate={navigate}>
          <AdminView apiFetch={apiFetch} />
        </DashboardLayout>
      ) : (
        <RedirectView to="#/" />
      );
    default:
      return <LandingView user={user} apiFetch={apiFetch} navigate={navigate} />;
  }
}

function RedirectView({ to }: { to: string }) {
  useEffect(() => {
    window.location.hash = to;
  }, [to]);
  return null;
}

import React, { ReactNode, useEffect } from 'react';
import { Link, LayoutDashboard, Link2, ShieldAlert, LogOut } from 'lucide-react';
import LandingView, { User } from './views/LandingView';
import LoginView from './views/LoginView';
import RegisterView from './views/RegisterView';
import DashboardView from './views/DashboardView';
import MyURLsView from './views/MyURLsView';
import AdminView from './views/AdminView';

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

interface DashboardLayoutProps {
  user: User;
  activeTab: string;
  onLogout: () => void;
  navigate: (toHash: string) => void;
  children: ReactNode;
}

function DashboardLayout({ user, activeTab, onLogout, navigate, children }: DashboardLayoutProps) {
  return (
    <div className="flex-1 flex flex-col md:flex-row ">
      {/* Sidebar */}
      <aside className="w-full md:w-64 glass border-r border-white/5 flex flex-col p-6">
        {/* Logo */}
        <div className="flex items-center gap-2 cursor-pointer mb-10" onClick={() => navigate('#/')}>
          <div className="p-1.5 bg-brand-600 rounded-lg text-white font-bold">
            <Link className="w-4 h-4" />
          </div>
          <span className="font-outfit text-lg font-bold tracking-tight">Snippy</span>
        </div>

        {/* Nav List */}
        <nav className="flex-1 space-y-1">
          <button
            onClick={() => navigate('#/dashboard')}
            className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition text-sm font-medium ${
              activeTab === 'dashboard' ? 'bg-brand-600 text-white' : 'text-gray-400 hover:bg-white/5 hover:text-white'
            }`}
          >
            <LayoutDashboard className="w-4 h-4" />
            Dashboard
          </button>

          <button
            onClick={() => navigate('#/my-urls')}
            className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition text-sm font-medium ${
              activeTab === 'my-urls' ? 'bg-brand-600 text-white' : 'text-gray-400 hover:bg-white/5 hover:text-white'
            }`}
          >
            <Link2 className="w-4 h-4" />
            My URLs
          </button>

          {user.role === 'admin' && (
            <button
              onClick={() => navigate('#/admin')}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition text-sm font-medium ${
                activeTab === 'admin' ? 'bg-brand-600 text-white' : 'text-gray-400 hover:bg-white/5 hover:text-white'
              }`}
            >
              <ShieldAlert className="w-4 h-4" />
              Admin Panel
            </button>
          )}
        </nav>

        {/* Profile widget */}
        <div className="pt-6 border-t border-white/5 flex flex-col gap-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-brand-900 border border-brand-500/30 flex items-center justify-center font-bold text-brand-300">
              {user.username.substring(0, 2).toUpperCase()}
            </div>
            <div className="overflow-hidden">
              <span className="font-semibold text-sm text-gray-200 block truncate">{user.username}</span>
              <span className="text-xs text-gray-500 block truncate">{user.email}</span>
            </div>
          </div>

          <button
            onClick={onLogout}
            className="w-full flex items-center gap-3 px-4 py-2.5 rounded-xl hover:bg-rose-950/20 text-rose-400 hover:text-rose-300 transition text-sm font-medium border border-rose-500/10"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </button>
        </div>
      </aside>

      {/* Main content body */}
      <main className="flex-1 p-6 md:p-10 overflow-y-auto max-w-7xl">
        {children}
      </main>
    </div>
  );
}

function RedirectView({ to }: { to: string }) {
  useEffect(() => {
    window.location.hash = to;
  }, [to]);
  return null;
}

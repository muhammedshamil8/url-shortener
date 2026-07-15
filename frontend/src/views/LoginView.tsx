import React, { useState, FormEvent } from 'react';
import { Link } from 'lucide-react';
import { useToast } from '../components/Toast';
import { User } from './LandingView';
import { API_BASE_URL } from '../config';

interface LoginViewProps {
  onLoginSuccess: (userData: User) => void;
  navigate: (toHash: string) => void;
}

export default function LoginView({ onLoginSuccess, navigate }: LoginViewProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const showToast = useToast();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      const data = await res.json();
      if (res.ok) {
        onLoginSuccess({
          accessToken: data.data.access_token,
          refreshToken: data.data.refresh_token,
          username: data.data.user.username,
          email: data.data.user.email,
          role: data.data.user.role,
        });
        showToast(`Welcome back, ${data.data.user.username}!`, "success");
      } else {
        showToast(data.error || "Failed to authenticate", "error");
      }
    } catch (e) {
      showToast("Something went wrong", "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex-1 flex items-center justify-center px-6 py-12">
      <div className="w-full max-w-md p-8 glass rounded-2xl shadow-2xl relative z-10">
        <div className="flex flex-col items-center mb-8">
          <div 
            className="p-3 bg-brand-600 rounded-2xl text-white font-bold cursor-pointer mb-3"
            onClick={() => navigate('#/')}
          >
            <Link className="w-6 h-6" />
          </div>
          <h2 className="font-outfit text-2xl font-bold">Sign in to Snippy</h2>
          <p className="text-gray-400 text-sm mt-1">Access your urls and analytics</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-xs text-gray-400 font-semibold uppercase tracking-wider block mb-1.5">Email Address</label>
            <input 
              type="email" 
              required
              placeholder="name@company.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-3 glass-input rounded-xl text-white text-sm"
            />
          </div>

          <div>
            <label className="text-xs text-gray-400 font-semibold uppercase tracking-wider block mb-1.5">Password</label>
            <input 
              type="password" 
              required
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 glass-input rounded-xl text-white text-sm"
            />
          </div>

          <button 
            type="submit" 
            disabled={loading}
            className="w-full bg-brand-600 hover:bg-brand-700 disabled:opacity-50 text-white font-medium py-3 rounded-xl transition duration-200 text-sm mt-6"
          >
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-400">
          Don't have an account? <span onClick={() => navigate('#/register')} className="text-brand-500 hover:underline cursor-pointer font-medium">Register</span>
        </div>
      </div>
    </div>
  );
}

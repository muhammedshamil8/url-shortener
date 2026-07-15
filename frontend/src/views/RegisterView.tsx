import React, { useState, FormEvent } from 'react';
import { Link } from 'lucide-react';
import { useToast } from '../components/Toast';
import api from '../api';

interface RegisterViewProps {
  navigate: (toHash: string) => void;
}

export default function RegisterView({ navigate }: RegisterViewProps) {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const showToast = useToast();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      await api.post('/api/v1/auth/register', { username, email, password });
      showToast("Account created successfully! Please log in.", "success");
      navigate('#/login');
    } catch (e: any) {
      const errMsg = e.response?.data?.error || "Failed to create account";
      showToast(errMsg, "error");
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
          <h2 className="font-outfit text-2xl font-bold">Create an account</h2>
          <p className="text-gray-400 text-sm mt-1">Get detailed analytics and management tools</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-xs text-gray-400 font-semibold uppercase tracking-wider block mb-1.5">Username</label>
            <input 
              type="text" 
              required
              placeholder="johndoe"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-3 glass-input rounded-xl text-white text-sm"
            />
          </div>

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
            {loading ? 'Registering...' : 'Register'}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-400">
          Already have an account? <span onClick={() => navigate('#/login')} className="text-brand-500 hover:underline cursor-pointer font-medium">Sign In</span>
        </div>
      </div>
    </div>
  );
}

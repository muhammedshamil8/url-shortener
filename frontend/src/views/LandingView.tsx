import React, { useState, FormEvent } from 'react';
import { Link2, LayoutDashboard, ArrowRight, Copy, ExternalLink, Link } from 'lucide-react';
import { useToast } from '../components/Toast';
import { API_BASE_URL } from '../config';

export interface User {
  username: string;
  email: string;
  role: string;
  accessToken: string;
  refreshToken: string;
}

interface LandingViewProps {
  user: User | null;
  apiFetch: (endpoint: string, options?: RequestInit) => Promise<Response>;
  navigate: (toHash: string) => void;
}

interface ShortenedURL {
  id: number;
  short_code: string;
  original_url: string;
  short_url: string;
}

export default function LandingView({ user, apiFetch, navigate }: LandingViewProps) {
  const apiBase = (API_BASE_URL || '').replace(/\/$/, '').replace(/\/api\/v1$/, '');
  const swaggerUrl = `${apiBase}/swagger/index.html`;

  const [urlInput, setUrlInput] = useState('');
  const [shortenedResult, setShortenedResult] = useState<ShortenedURL | null>(null);
  const [loading, setLoading] = useState(false);
  const showToast = useToast();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!urlInput.trim()) return;

    setLoading(true);
    setShortenedResult(null);
    try {
      const res = await apiFetch('/api/v1/shorten', {
        method: 'POST',
        body: JSON.stringify({ url: urlInput })
      });

      const data = await res.json();
      if (res.ok) {
        setShortenedResult(data.data);
        showToast("URL shortened successfully!", "success");
      } else {
        showToast(data.error || "Failed to shorten URL", "error");
      }
    } catch (err) {
      showToast("Something went wrong", "error");
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    showToast("Copied to clipboard!", "success");
  };

  return (
    <div className="flex-1 flex flex-col">
      {/* Header */}
      <header className="w-full max-w-7xl mx-auto px-6 py-5 flex items-center justify-between relative z-10">
        <div className="flex items-center gap-2 cursor-pointer" onClick={() => navigate('#/')}>
          <div className="p-2 bg-brand-600 rounded-lg text-white font-bold">
            <Link className="w-5 h-5" />
          </div>
          <span className="font-outfit text-xl font-bold tracking-tight">Snippy</span>
        </div>

        <div className="flex items-center gap-4">
          {user ? (
            <button
              onClick={() => navigate('#/dashboard')}
              className="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white font-medium rounded-xl transition duration-200 text-sm flex items-center gap-2"
            >
              <LayoutDashboard className="w-4 h-4" />
              Dashboard
            </button>
          ) : (
            <div className="flex items-center gap-4">
              <button onClick={() => navigate('#/login')} className="px-4 py-2 text-gray-400 hover:text-white font-medium transition text-sm">
                Sign In
              </button>
              <button
                onClick={() => navigate('#/register')}
                className="px-4 py-2 bg-white/10 hover:bg-white/15 text-white font-medium rounded-xl border border-white/10 transition text-sm"
              >
                Register
              </button>
            </div>
          )}
        </div>
      </header>

      {/* Hero Section */}
      <main className="flex-1 max-w-4xl mx-auto px-6 flex flex-col items-center justify-center text-center py-20 relative z-10">
        <h1 className="font-outfit text-4xl sm:text-6xl font-extrabold tracking-tight mb-6 leading-tight">
          Shorten links. <br />
          <span className="bg-gradient-to-r from-brand-500 to-violet-400 bg-clip-text text-transparent">Measure impact.</span>
        </h1>
        <p className="text-gray-400 text-lg max-w-xl mb-12">
          Create clean, memorable, and analytics-driven short URLs. Free, instant, and completely secure.
        </p>

        {/* Input Box */}
        <div className="w-full max-w-2xl glass p-3 rounded-2xl shadow-2xl mb-8">
          <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-2">
            <div className="flex-1 relative">
              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500">
                <Link2 className="w-5 h-5" />
              </span>
              <input
                type="url"
                placeholder="Paste your long URL here..."
                required
                value={urlInput}
                onChange={(e) => setUrlInput(e.target.value)}
                className="w-full pl-12 pr-4 py-3.5 glass-input rounded-xl text-white placeholder-gray-500 text-sm"
              />
            </div>
            <button
              type="submit"
              disabled={loading}
              className="bg-brand-600 hover:bg-brand-700 disabled:opacity-50 text-white font-medium px-8 py-3.5 rounded-xl transition duration-200 text-sm flex items-center justify-center gap-2"
            >
              {loading ? 'Shortening...' : 'Shorten URL'}
              <ArrowRight className="w-4 h-4" />
            </button>
          </form>
        </div>

        {/* Shortened Result Box */}
        {shortenedResult && (
          <div className="w-full max-w-2xl p-5 rounded-2xl glass border-emerald-500/20 bg-emerald-950/10 text-left animate-fade-in-up">
            <span className="text-xs font-semibold text-emerald-400 uppercase tracking-wider block mb-2">Shortened Result</span>
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <div className="flex-1 overflow-hidden pr-4">
                <a
                  href={shortenedResult.short_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-outfit text-xl font-bold text-brand-500 hover:underline break-all"
                >
                  {shortenedResult.short_url}
                </a>
                <span className="text-xs text-gray-500 block truncate mt-1">{shortenedResult.original_url}</span>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => copyToClipboard(shortenedResult.short_url)}
                  className="p-2.5 bg-white/5 hover:bg-white/10 rounded-xl transition border border-white/5"
                  title="Copy Link"
                >
                  <Copy className="w-4 h-4 text-gray-300" />
                </button>
                <a
                  href={shortenedResult.short_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="p-2.5 bg-brand-600 hover:bg-brand-700 rounded-xl transition text-white"
                  title="Open Link"
                >
                  <ExternalLink className="w-4 h-4" />
                </a>
              </div>
            </div>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="w-full max-w-7xl mx-auto px-6 py-6 border-t border-white/5 flex flex-col sm:flex-row items-center justify-between text-xs text-gray-500 relative z-10">
        <span>&copy; 2026 Snippy. Built with Go & React.</span>
        <div className="flex gap-4 mt-2 sm:mt-0">
          <a href="#/login" className="hover:text-gray-300">Admin Dashboard</a>
          <a href={swaggerUrl} target="_blank" rel="noopener noreferrer" className="hover:text-gray-300">API Documentation</a>
        </div>
      </footer>
    </div>
  );
}

import React, { useState, useEffect } from 'react';
import { Link2, BarChart3, PlusCircle, ArrowRight, Copy, ExternalLink } from 'lucide-react';
import { useToast } from '../components/Toast';

export default function DashboardView({ user, apiFetch }) {
  const [urlInput, setUrlInput] = useState('');
  const [shortenedResult, setShortenedResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState({ totalUrls: 0, totalClicks: 0 });
  const showToast = useToast();

  const fetchStats = async () => {
    try {
      const res = await apiFetch('/api/v1/my/urls');
      if (res.ok) {
        const data = await res.json();
        const list = data.data || [];
        const clicks = list.reduce((sum, item) => sum + (item.click_count || 0), 0);
        setStats({ totalUrls: list.length, totalClicks: clicks });
      }
    } catch (e) {}
  };

  useEffect(() => {
    fetchStats();
  }, []);

  const handleSubmit = async (e) => {
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
        showToast("URL shortened and linked to account!", "success");
        fetchStats();
      } else {
        showToast(data.error || "Failed to shorten URL", "error");
      }
    } catch (err) {
      showToast("Something went wrong", "error");
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    showToast("Copied to clipboard!", "success");
  };

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-outfit text-3xl font-extrabold tracking-tight">Welcome, {user.username}!</h1>
        <p className="text-gray-400 text-sm mt-1">Here is a quick overview of your shortened links.</p>
      </div>

      {/* Quick Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="glass p-6 rounded-2xl flex items-center gap-5">
          <div className="p-3 bg-brand-500/10 text-brand-400 rounded-xl">
            <Link2 className="w-6 h-6" />
          </div>
          <div>
            <span className="text-xs font-semibold text-gray-500 uppercase tracking-wider block">Total URLs</span>
            <span className="font-outfit text-2xl font-bold">{stats.totalUrls}</span>
          </div>
        </div>

        <div className="glass p-6 rounded-2xl flex items-center gap-5">
          <div className="p-3 bg-indigo-500/10 text-indigo-400 rounded-xl">
            <BarChart3 className="w-6 h-6" />
          </div>
          <div>
            <span className="text-xs font-semibold text-gray-500 uppercase tracking-wider block">Total Clicks</span>
            <span className="font-outfit text-2xl font-bold">{stats.totalClicks}</span>
          </div>
        </div>
      </div>

      {/* Shortener box */}
      <div className="glass p-6 rounded-2xl shadow-xl">
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <PlusCircle className="w-5 h-5 text-brand-500" />
          Shorten a new URL
        </h2>
        <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-2">
          <div className="flex-1 relative">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-500">
              <Link2 className="w-5 h-5" />
            </span>
            <input 
              type="url" 
              placeholder="https://example.com/very-long-url-path"
              required
              value={urlInput}
              onChange={(e) => setUrlInput(e.target.value)}
              className="w-full pl-12 pr-4 py-3 glass-input rounded-xl text-white placeholder-gray-500 text-sm"
            />
          </div>
          <button 
            type="submit" 
            disabled={loading}
            className="bg-brand-600 hover:bg-brand-700 disabled:opacity-50 text-white font-medium px-8 py-3 rounded-xl transition duration-200 text-sm flex items-center justify-center gap-2"
          >
            {loading ? 'Shortening...' : 'Shorten'}
            <ArrowRight className="w-4 h-4" />
          </button>
        </form>

        {shortenedResult && (
          <div className="mt-6 p-4 rounded-xl glass border-emerald-500/20 bg-emerald-950/5 animate-fade-in-up">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
              <div className="flex-1 overflow-hidden pr-4">
                <a 
                  href={shortenedResult.short_url} 
                  target="_blank" 
                  rel="noopener noreferrer" 
                  className="font-outfit text-lg font-bold text-brand-500 hover:underline break-all"
                >
                  {shortenedResult.short_url}
                </a>
                <span className="text-xs text-gray-500 block truncate mt-0.5">{shortenedResult.original_url}</span>
              </div>
              <div className="flex gap-2">
                <button 
                  onClick={() => copyToClipboard(shortenedResult.short_url)}
                  className="p-2 bg-white/5 hover:bg-white/10 rounded-lg transition border border-white/5"
                >
                  <Copy className="w-4 h-4 text-gray-300" />
                </button>
                <a 
                  href={shortenedResult.short_url} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="p-2 bg-brand-600 hover:bg-brand-700 rounded-lg transition text-white"
                >
                  <ExternalLink className="w-4 h-4" />
                </a>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

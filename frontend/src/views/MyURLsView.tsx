import React, { useState, useEffect } from 'react';
import { Link2Off, Copy, Trash2 } from 'lucide-react';
import { useToast } from '../components/Toast';

interface MyURLsViewProps {
  apiFetch: (endpoint: string, options?: RequestInit) => Promise<Response>;
}

interface ShortenedURL {
  id: number;
  short_code: string;
  original_url: string;
  click_count?: number;
  created_at: string;
}

export default function MyURLsView({ apiFetch }: MyURLsViewProps) {
  const [urls, setUrls] = useState<ShortenedURL[]>([]);
  const [loading, setLoading] = useState(true);
  const showToast = useToast();

  const fetchUrls = async () => {
    setLoading(true);
    try {
      const res = await apiFetch('/api/v1/my/urls');
      if (res.ok) {
        const data = await res.json();
        setUrls(data.data.urls || []);
      }
    } catch (e) {
      showToast("Failed to fetch URLs", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this URL?")) return;

    try {
      const res = await apiFetch(`/api/v1/my/urls/${id}`, { method: 'DELETE' });
      if (res.ok) {
        showToast("URL deleted successfully", "success");
        fetchUrls();
      } else {
        showToast("Failed to delete URL", "error");
      }
    } catch (e) {
      showToast("Something went wrong", "error");
    }
  };

  useEffect(() => {
    fetchUrls();
  }, []);

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    showToast("Copied to clipboard!", "success");
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-outfit text-3xl font-extrabold tracking-tight">My URLs</h1>
        <p className="text-gray-400 text-sm mt-1">Manage and track your custom shortened links.</p>
      </div>

      <div className="glass rounded-2xl overflow-hidden shadow-xl border border-white/5">
        {loading ? (
          <div className="p-12 text-center text-gray-500">Loading your URLs...</div>
        ) : urls.length === 0 ? (
          <div className="p-12 text-center text-gray-500">
            <Link2Off className="w-10 h-10 text-gray-600 mx-auto mb-3" />
            No URLs created yet.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-white/5 text-gray-400 font-semibold uppercase tracking-wider text-xs border-b border-white/5">
                <tr>
                  <th className="px-6 py-4">Short Code</th>
                  <th className="px-6 py-4">Original URL</th>
                  <th className="px-6 py-4 text-center">Clicks</th>
                  <th className="px-6 py-4">Created Date</th>
                  <th className="px-6 py-4 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {urls.map((u) => {
                  const absoluteShort = `${window.location.protocol}//${window.location.host}/${u.short_code}`;
                  return (
                    <tr key={u.id} className="hover:bg-white/5 transition duration-150">
                      <td className="px-6 py-4 font-semibold text-brand-400 font-outfit">
                        <a href={absoluteShort} target="_blank" rel="noopener noreferrer" className="hover:underline">
                          {u.short_code}
                        </a>
                      </td>
                      <td className="px-6 py-4 max-w-xs truncate text-gray-300" title={u.original_url}>
                        {u.original_url}
                      </td>
                      <td className="px-6 py-4 text-center font-semibold text-gray-200">
                        {u.click_count || 0}
                      </td>
                      <td className="px-6 py-4 text-gray-400 text-xs">
                        {new Date(u.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex justify-end gap-1.5">
                          <button 
                            onClick={() => copyToClipboard(absoluteShort)}
                            className="p-2 hover:bg-white/10 rounded-lg text-gray-400 hover:text-white transition"
                            title="Copy"
                          >
                            <Copy className="w-4 h-4" />
                          </button>
                          <button 
                            onClick={() => handleDelete(u.id)}
                            className="p-2 hover:bg-rose-950/30 rounded-lg text-gray-500 hover:text-rose-400 transition"
                            title="Delete"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

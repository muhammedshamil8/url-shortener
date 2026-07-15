import React, { useState, useEffect } from 'react';
import { Trash2, UserX } from 'lucide-react';
import { useToast } from '../components/Toast';

export default function AdminView({ apiFetch }) {
  const [activeSubTab, setActiveSubTab] = useState('urls');
  const [urls, setUrls] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const showToast = useToast();

  const fetchAdminData = async () => {
    setLoading(true);
    try {
      if (activeSubTab === 'urls') {
        const res = await apiFetch('/api/v1/admin/urls?limit=100');
        if (res.ok) {
          const data = await res.json();
          setUrls(data.data || []);
        }
      } else {
        const res = await apiFetch('/api/v1/admin/users');
        if (res.ok) {
          const data = await res.json();
          setUsers(data.users || []);
        }
      }
    } catch (e) {
      showToast("Failed to retrieve administration datasets", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteUrl = async (id) => {
    if (!confirm("Are you sure you want to delete this URL?")) return;
    try {
      const res = await apiFetch(`/api/v1/admin/urls/${id}`, { method: 'DELETE' });
      if (res.ok) {
        showToast("URL deleted successfully", "success");
        fetchAdminData();
      } else {
        showToast("Failed to delete URL", "error");
      }
    } catch (e) {
      showToast("Something went wrong", "error");
    }
  };

  const handleDeleteUser = async (id) => {
    if (!confirm("Deleting a user will cascade delete all of their shortened URLs. Proceed?")) return;
    try {
      const res = await apiFetch(`/api/v1/admin/users/${id}`, { method: 'DELETE' });
      if (res.ok) {
        showToast("User deleted successfully", "success");
        fetchAdminData();
      } else {
        showToast("Failed to delete user", "error");
      }
    } catch (e) {
      showToast("Something went wrong", "error");
    }
  };

  useEffect(() => {
    fetchAdminData();
  }, [activeSubTab]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-outfit text-3xl font-extrabold tracking-tight">Admin Control Panel</h1>
        <p className="text-gray-400 text-sm mt-1">Global audit logs, accounts and shortening metrics.</p>
      </div>

      {/* Subtab selection toggles */}
      <div className="flex border-b border-white/5 gap-6">
        <button 
          onClick={() => setActiveSubTab('urls')}
          className={`pb-3 font-semibold text-sm transition relative ${
            activeSubTab === 'urls' ? 'text-brand-500 border-b-2 border-brand-500' : 'text-gray-500 hover:text-white'
          }`}
        >
          All URLs
        </button>
        <button 
          onClick={() => setActiveSubTab('users')}
          className={`pb-3 font-semibold text-sm transition relative ${
            activeSubTab === 'users' ? 'text-brand-500 border-b-2 border-brand-500' : 'text-gray-500 hover:text-white'
          }`}
        >
          Registered Users
        </button>
      </div>

      <div className="glass rounded-2xl overflow-hidden shadow-xl border border-white/5">
        {loading ? (
          <div className="p-12 text-center text-gray-500">Retrieving audit dataset...</div>
        ) : activeSubTab === 'urls' ? (
          urls.length === 0 ? (
            <div className="p-12 text-center text-gray-500">No database links populated.</div>
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
                      <tr key={u.id} className="hover:bg-white/5 transition">
                        <td className="px-6 py-4 font-semibold text-brand-400 font-outfit">
                          <a href={absoluteShort} target="_blank" rel="noopener noreferrer" className="hover:underline">
                            {u.short_code}
                          </a>
                        </td>
                        <td className="px-6 py-4 max-w-sm truncate text-gray-300" title={u.original_url}>
                          {u.original_url}
                        </td>
                        <td className="px-6 py-4 text-center font-semibold text-gray-200">{u.click_count || 0}</td>
                        <td className="px-6 py-4 text-gray-400 text-xs">
                          {new Date(u.created_at).toLocaleDateString()}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <button 
                            onClick={() => handleDeleteUrl(u.id)}
                            className="p-2 hover:bg-rose-950/30 rounded-lg text-gray-500 hover:text-rose-400 transition"
                            title="Delete Link"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )
        ) : (
          users.length === 0 ? (
            <div className="p-12 text-center text-gray-500">No database user accounts found.</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="bg-white/5 text-gray-400 font-semibold uppercase tracking-wider text-xs border-b border-white/5">
                  <tr>
                    <th className="px-6 py-4">User ID</th>
                    <th className="px-6 py-4">Username</th>
                    <th className="px-6 py-4">Email</th>
                    <th className="px-6 py-4">Registered</th>
                    <th className="px-6 py-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-white/5">
                  {users.map((u) => (
                    <tr key={u.id} className="hover:bg-white/5 transition">
                      <td className="px-6 py-4 font-bold text-gray-400">{u.id}</td>
                      <td className="px-6 py-4 font-semibold text-gray-200">{u.username}</td>
                      <td className="px-6 py-4 text-gray-300">{u.email}</td>
                      <td className="px-6 py-4 text-gray-400 text-xs">
                        {new Date(u.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <button 
                          onClick={() => handleDeleteUser(u.id)}
                          className="p-2 hover:bg-rose-950/30 rounded-lg text-gray-500 hover:text-rose-400 transition"
                          title="Delete Account"
                        >
                          <UserX className="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
      </div>
    </div>
  );
}

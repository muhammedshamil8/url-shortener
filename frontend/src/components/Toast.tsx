import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';

export type ToastType = 'info' | 'success' | 'error';
export type ShowToast = (message: string, type?: ToastType) => void;

const ToastContext = createContext<ShowToast | null>(null);

interface ToastProviderProps {
  children: ReactNode;
}

export function ToastProvider({ children }: ToastProviderProps) {
  const [toast, setToast] = useState<{ message: string; type: ToastType } | null>(null);

  const showToast = useCallback((message: string, type: ToastType = 'info') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 4000);
  }, []);

  const colors = {
    success: 'border-emerald-500/30 bg-emerald-950/80 text-emerald-300',
    error: 'border-rose-500/30 bg-rose-950/80 text-rose-300',
    info: 'border-brand-500/30 bg-brand-950/80 text-brand-300',
  };

  return (
    <ToastContext.Provider value={showToast}>
      {children}
      {toast && (
        <div className="fixed bottom-6 right-6 z-50 animate-bounce">
          <div className={`flex items-center gap-3 px-4 py-3 rounded-xl border glass shadow-2xl ${colors[toast.type] || colors.info}`}>
            <span className="text-sm font-medium">{toast.message}</span>
          </div>
        </div>
      )}
    </ToastContext.Provider>
  );
}

export function useToast(): ShowToast {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
}

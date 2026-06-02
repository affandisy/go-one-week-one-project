import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/useAuthStore';
import { ShieldCheck, Lock, User } from 'lucide-react';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  const login = useAuthStore((state) => state.login);
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      const response = await login(username, password);
      if (response.success) {
        navigate('/'); // Arahkan ke Dasbor setelah sukses
      } else {
        setError(response.message);
      }
    } catch (err) {
      setError('Terjadi kesalahan pada peladen');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-white rounded-3xl shadow-xl overflow-hidden border border-slate-100">
        
        {/* Bagian Header Form */}
        <div className="bg-indigo-900 p-8 text-center">
          <div className="w-16 h-16 bg-indigo-800 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-inner">
            <ShieldCheck size={32} className="text-blue-400" />
          </div>
          <h1 className="text-3xl font-black text-white tracking-tight">PetCare<span className="text-blue-400">Pro</span></h1>
          <p className="text-indigo-300 text-sm mt-2 font-medium">Sistem Manajemen & Gizi Klinik</p>
        </div>

        {/* Bagian Input Form */}
        <div className="p-8">
          <h2 className="text-xl font-bold text-slate-800 mb-6 text-center">Masuk ke Akun Anda</h2>
          
          {error && (
            <div className="mb-6 p-4 bg-red-50 text-red-600 text-sm font-bold rounded-xl border border-red-100 text-center">
              {error}
            </div>
          )}

          <form onSubmit={handleLogin} className="space-y-5">
            <div>
              <label className="block text-sm font-bold text-slate-600 mb-2">Username</label>
              <div className="relative">
                <User className="absolute left-4 top-3.5 text-slate-400" size={18} />
                <input 
                  type="text" 
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="kasir / manajer"
                  className="w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 text-sm font-medium"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-600 mb-2">Kata Sandi</label>
              <div className="relative">
                <Lock className="absolute left-4 top-3.5 text-slate-400" size={18} />
                <input 
                  type="password" 
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="rahasia"
                  className="w-full pl-11 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 text-sm font-medium"
                  required
                />
              </div>
            </div>

            <button 
              type="submit" 
              disabled={isLoading}
              className="w-full py-4 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-300 text-white font-black text-lg rounded-xl shadow-lg shadow-indigo-200 transition-all active:scale-95 flex justify-center items-center gap-2 mt-4"
            >
              {isLoading ? 'Memproses...' : 'Masuk Sistem'}
            </button>
          </form>
          
          <div className="mt-6 text-center text-xs text-slate-400 font-medium">
            <p>Gunakan username <b>kasir</b> atau <b>manajer</b></p>
            <p>Kata sandi: <b>rahasia</b></p>
          </div>
        </div>
      </div>
    </div>
  );
}
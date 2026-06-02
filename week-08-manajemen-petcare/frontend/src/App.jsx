import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from 'react-router-dom';
import { Activity, ClipboardList, Syringe, Users, LogOut } from 'lucide-react';

import Dashboard from './pages/Dashboard';
import CashierView from './pages/CashierView';
import NutritionLog from './pages/NutritionLog';
import Login from './pages/Login';
import ProtectedRoute from './components/ProtectedRoute';
import { useAuthStore } from './store/useAuthStore';

const MainLayout = () => {
  const { userRole, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="flex h-screen bg-slate-50">
      <aside className="w-64 bg-indigo-900 text-white flex flex-col shadow-xl">
        <div className="p-6 border-b border-indigo-800">
          <h1 className="text-2xl font-black tracking-tight text-indigo-100">PetCare<span className="text-blue-400">Pro</span></h1>
          <p className="text-xs text-indigo-300 mt-1 font-medium">Sistem Operasional</p>
        </div>
        
        <nav className="flex-1 p-4 space-y-2">
          <Link to="/" className="flex items-center gap-3 p-3 rounded-xl hover:bg-indigo-800 transition-colors font-medium">
            <Activity size={20} className="text-blue-400" /> Dasbor
          </Link>
          <Link to="/cashier" className="flex items-center gap-3 p-3 rounded-xl hover:bg-indigo-800 transition-colors font-medium">
            <ClipboardList size={20} className="text-green-400" /> Kasir & Tagihan
          </Link>
          <Link to="/nutrition" className="flex items-center gap-3 p-3 rounded-xl hover:bg-indigo-800 transition-colors font-medium">
            <Syringe size={20} className="text-orange-400" /> Rekam Gizi
          </Link>
        </nav>
        
        <div className="p-4 border-t border-indigo-800">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-indigo-700 rounded-full flex justify-center items-center font-bold text-indigo-100">
                {userRole?.charAt(0) || 'U'}
              </div>
              <div>
                <p className="text-sm font-bold">{userRole}</p>
                <p className="text-xs text-indigo-300">Aktif</p>
              </div>
            </div>
            <button onClick={handleLogout} className="p-2 text-red-300 hover:text-red-100 hover:bg-red-900/30 rounded-lg transition-colors">
              <LogOut size={18} />
            </button>
          </div>
        </div>
      </aside>

      <main className="flex-1 overflow-y-auto">
        {/* Tempat halaman anak di-render */}
        <Outlet /> 
      </main>
    </div>
  );
};

import { Outlet } from 'react-router-dom';

function App() {
  return (
    <Router>
      <Routes>
        {/* Rute Publik */}
        <Route path="/login" element={<Login />} />

        {/* Rute Terproteksi */}
        <Route element={<ProtectedRoute />}>
          <Route element={<MainLayout />}>
            <Route path="/" element={<Dashboard />} />
            
            {/* Contoh pembatasan Role: Hanya Cashier & Manager yang bisa buka Kasir */}
            <Route element={<ProtectedRoute allowedRoles={['Cashier', 'Manager']} />}>
              <Route path="/cashier" element={<CashierView />} />
            </Route>

            <Route path="/nutrition" element={<NutritionLog />} />
          </Route>
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
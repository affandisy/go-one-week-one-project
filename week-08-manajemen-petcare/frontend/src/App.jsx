import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Activity, ClipboardList, Syringe, Users } from 'lucide-react';
import CashierView from './pages/CashierView';

// Placeholder komponen (akan kita buat di langkah selanjutnya)
const Dashboard = () => <div className="p-8"><h1 className="text-2xl font-bold">Dasbor Operasional</h1></div>;
const CashierView = () => <div className="p-8"><h1 className="text-2xl font-bold">Modul Kasir & Penagihan</h1></div>;
const NutritionLog = () => <div className="p-8"><h1 className="text-2xl font-bold">Rekam Gizi & Transformasi</h1></div>;

function App() {
  return (
    <Router>
      <div className="flex h-screen bg-slate-50">
        
        {/* Sidebar Navigasi */}
        <aside className="w-64 bg-indigo-900 text-white flex flex-col shadow-xl">
          <div className="p-6 border-b border-indigo-800">
            <h1 className="text-2xl font-black tracking-tight text-indigo-100">PetCare<span className="text-blue-400">Pro</span></h1>
            <p className="text-xs text-indigo-300 mt-1 font-medium">Manajemen & Gizi Klinik</p>
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
            <Link to="/master" className="flex items-center gap-3 p-3 rounded-xl hover:bg-indigo-800 transition-colors font-medium">
              <Users size={20} className="text-purple-400" /> Data Master
            </Link>
          </nav>
          
          <div className="p-4 border-t border-indigo-800">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-indigo-700 rounded-full flex justify-center items-center font-bold">M</div>
              <div>
                <p className="text-sm font-bold">Manajer Operasional</p>
                <button className="text-xs text-red-300 hover:text-red-100 font-medium">Keluar Sesi</button>
              </div>
            </div>
          </div>
        </aside>

        {/* Area Konten Utama */}
        <main className="flex-1 overflow-y-auto">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/cashier" element={<CashierView />} />
            <Route path="/nutrition" element={<NutritionLog />} />
          </Routes>
        </main>

      </div>
    </Router>
  );
}

export default App;
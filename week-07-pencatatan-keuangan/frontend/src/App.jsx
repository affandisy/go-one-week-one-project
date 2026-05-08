import { useState, useEffect } from 'react';
import InputModal from './components/InputModal';
import ReportView from './components/ReportView';
import HistoryView from './components/HistoryView';

const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard'); // dashboard | history | report
  const [wallet, setWallet] = useState(null);
  const [recentTrx, setRecentTrx] = useState([]);
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [txType, setTxType] = useState('expense');

  const loadDashboardData = async () => {
    try {
      const walletRes = await fetch('http://localhost:3000/api/v1/wallet');
      const walletData = await walletRes.json();
      setWallet(walletData.data);

      const trxRes = await fetch('http://localhost:3000/api/v1/transactions/recent');
      const trxData = await trxRes.json();
      setRecentTrx(trxData.data || []);
    } catch (error) { console.error("Gagal memuat data", error); }
  };

  useEffect(() => {
    if (activeTab === 'dashboard') loadDashboardData();
  }, [activeTab]);

  const openInput = (type) => { setTxType(type); setIsModalOpen(true); };

  return (
    <div className="max-w-md mx-auto h-screen bg-gray-50 shadow-2xl flex flex-col font-sans relative overflow-hidden">
      
      {/* KONTEN DINAMIS BERDASARKAN TAB */}
      {activeTab === 'dashboard' && (
        <div className="flex-1 overflow-y-auto pb-24">
          <div className="bg-blue-800 text-white p-6 rounded-b-3xl shadow-md z-10">
            <h1 className="text-xl font-bold tracking-wide mb-6">Keuangan Sederhana</h1>
            <p className="text-blue-200 text-sm font-medium">{wallet?.name || 'Memuat...'}</p>
            <h2 className="text-4xl font-black mt-1 mb-2">{formatRp(wallet?.balance)}</h2>
            <span className="bg-blue-700/50 text-xs px-3 py-1 rounded-full font-bold uppercase tracking-wider">Saldo Saat Ini</span>
          </div>

          <div className="p-6">
            <div className="flex gap-4 mb-8">
              <button onClick={() => openInput('expense')} className="flex-1 bg-red-100 text-red-700 py-4 rounded-2xl font-black text-lg hover:bg-red-200 active:scale-95 transition-all shadow-sm">
                - KELUAR
              </button>
              <button onClick={() => openInput('income')} className="flex-1 bg-green-100 text-green-700 py-4 rounded-2xl font-black text-lg hover:bg-green-200 active:scale-95 transition-all shadow-sm">
                + MASUK
              </button>
            </div>

            <div className="flex justify-between items-center mb-4">
              <h3 className="font-bold text-gray-800 text-lg">Riwayat Terakhir</h3>
            </div>
            
            <div className="space-y-3">
              {recentTrx.slice(0, 3).map((trx) => ( // Hanya tampilkan 3 di dasbor
                <div key={trx.id} className="bg-white p-4 rounded-2xl shadow-sm border border-gray-100 flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <span className="text-2xl">{trx.category?.icon}</span>
                    <div>
                      <p className="font-bold text-gray-800">{trx.category?.name}</p>
                      <p className="text-xs text-gray-400 font-medium">{new Date(trx.date_time).toLocaleDateString('id-ID')}</p>
                    </div>
                  </div>
                  <p className={`font-black text-lg ${trx.type === 'expense' ? 'text-red-600' : 'text-green-600'}`}>
                    {trx.type === 'expense' ? '-' : '+'}{formatRp(trx.amount)}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'history' && <HistoryView />}
      {activeTab === 'report' && <ReportView />}

      {/* BOTTOM NAVIGATION (FR-002 UX Requirement) */}
      <div className="absolute bottom-0 left-0 right-0 bg-white border-t border-gray-100 flex justify-around p-4 pb-6 shadow-[0_-10px_40px_rgba(0,0,0,0.05)]">
        <button onClick={() => setActiveTab('dashboard')} className={`flex flex-col items-center gap-1 ${activeTab === 'dashboard' ? 'text-blue-600' : 'text-gray-400 hover:text-gray-600'}`}>
          <span className="text-2xl">🏠</span>
          <span className="text-[10px] font-black uppercase">Beranda</span>
        </button>
        <button onClick={() => setActiveTab('history')} className={`flex flex-col items-center gap-1 ${activeTab === 'history' ? 'text-blue-600' : 'text-gray-400 hover:text-gray-600'}`}>
          <span className="text-2xl">📋</span>
          <span className="text-[10px] font-black uppercase">Riwayat</span>
        </button>
        <button onClick={() => setActiveTab('report')} className={`flex flex-col items-center gap-1 ${activeTab === 'report' ? 'text-blue-600' : 'text-gray-400 hover:text-gray-600'}`}>
          <span className="text-2xl">📊</span>
          <span className="text-[10px] font-black uppercase">Laporan</span>
        </button>
      </div>

      {isModalOpen && (
        <InputModal 
          type={txType} 
          onClose={() => setIsModalOpen(false)}
          onSuccess={() => {
            setIsModalOpen(false);
            if (activeTab === 'dashboard') loadDashboardData();
          }}
        />
      )}
    </div>
  );
}
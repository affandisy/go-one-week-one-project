import { useState, useEffect } from 'react';
import InputModal from './components/InputModal';
import ReportView from './components/ReportView';
import HistoryView from './components/HistoryView';
import AddWalletModal from './components/AddWalletModal'; // <--- Impor komponen baru

const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard'); 
  
  const [wallets, setWallets] = useState([]);
  const [activeWalletId, setActiveWalletId] = useState('');
  const [recentTrx, setRecentTrx] = useState([]);
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [txType, setTxType] = useState('expense');
  
  // State baru untuk modal dompet
  const [isAddWalletOpen, setIsAddWalletOpen] = useState(false);

  const loadWallets = async () => {
    try {
      const res = await fetch('/api/v1/wallets');
      const data = await res.json();
      const loadedWallets = data.data || [];
      setWallets(loadedWallets);
      
      // Update wallet aktif jika belum ada, ATAU jika dompet bertambah (opsional)
      if (loadedWallets.length > 0 && !activeWalletId) {
        setActiveWalletId(loadedWallets[0].id);
      }
    } catch (error) { console.error("Gagal memuat dompet", error); }
  };

  const loadTransactions = async () => {
    try {
      const res = await fetch(`/api/v1/transactions/recent`); 
      const data = await res.json();
      setRecentTrx(data.data || []);
    } catch (error) { console.error("Gagal memuat riwayat", error); }
  };

  useEffect(() => { loadWallets(); }, []);
  useEffect(() => { if (activeWalletId) loadTransactions(); }, [activeWalletId, activeTab]);

  const openInput = (type) => { setTxType(type); setIsModalOpen(true); };

  const activeWalletInfo = wallets.find(w => w.id === activeWalletId);

  return (
    <div className="max-w-md mx-auto h-screen bg-gray-50 shadow-2xl flex flex-col font-sans relative overflow-hidden">
      
      {activeTab === 'dashboard' && (
        <div className="flex-1 overflow-y-auto pb-24">
          <div className="bg-blue-800 text-white p-6 rounded-b-3xl shadow-md z-10">
            
            {/* PENGALIH DOMPET & TOMBOL TAMBAH */}
            <div className="flex justify-between items-center mb-6">
              <h1 className="text-xl font-bold tracking-wide">Keuangan</h1>
              <div className="flex items-center gap-2">
                <select 
                  value={activeWalletId} 
                  onChange={(e) => setActiveWalletId(e.target.value)}
                  className="bg-blue-900 border border-blue-700 text-white text-sm rounded-xl px-3 py-2 outline-none font-bold cursor-pointer max-w-[120px] truncate"
                >
                  {wallets.map(w => (
                    <option key={w.id} value={w.id}>{w.name}</option>
                  ))}
                </select>
                {/* Tombol "+" untuk membuat dompet baru */}
                <button 
                  onClick={() => setIsAddWalletOpen(true)}
                  className="w-9 h-9 bg-blue-600 hover:bg-blue-500 text-white font-black rounded-xl border border-blue-500 shadow-sm flex justify-center items-center"
                  title="Tambah Dompet"
                >
                  +
                </button>
              </div>
            </div>

            <h2 className="text-4xl font-black mt-1 mb-2">{formatRp(activeWalletInfo?.balance)}</h2>
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
              {recentTrx.length === 0 ? (
                <p className="text-center text-gray-400 font-medium py-8 border-2 border-dashed border-gray-200 rounded-2xl">
                  Belum ada transaksi.
                </p>
              ) : (
                recentTrx.slice(0, 3).map((trx) => (
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
                ))
              )}
            </div>
          </div>
        </div>
      )}

      {activeTab === 'history' && <HistoryView />}
      {activeTab === 'report' && <ReportView />}

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
          walletId={activeWalletId} 
          onClose={() => setIsModalOpen(false)}
          onSuccess={() => {
            setIsModalOpen(false);
            loadWallets(); 
            loadTransactions(); 
          }}
        />
      )}

      {/* Render Modal Tambah Dompet */}
      {isAddWalletOpen && (
        <AddWalletModal 
          onClose={() => setIsAddWalletOpen(false)}
          onSuccess={() => {
            setIsAddWalletOpen(false);
            loadWallets(); // Refresh daftar agar dompet baru muncul di dropdown
          }}
        />
      )}
    </div>
  );
}
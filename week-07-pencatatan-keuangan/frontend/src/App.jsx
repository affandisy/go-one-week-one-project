import { useState, useEffect } from 'react';
import InputModal from './components/InputModal';

// Utilitas format Rupiah
const formatRp = (num) => new Intl.NumberFormat('id-ID', { 
  style: 'currency', 
  currency: 'IDR', 
  minimumFractionDigits: 0 
}).format(num || 0);

export default function App() {
  const [wallet, setWallet] = useState(null);
  const [recentTrx, setRecentTrx] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [txType, setTxType] = useState('expense'); // 'expense' atau 'income'

  // Memuat data dari API Golang
  const loadDashboardData = async () => {
    try {
      // Endpoint dompet dari master_handler
      const walletRes = await fetch('http://localhost:3000/api/v1/wallet');
      const walletData = await walletRes.json();
      setWallet(walletData.data);

      // Endpoint riwayat dari transaction_handler
      const trxRes = await fetch('http://localhost:3000/api/v1/transactions/recent');
      const trxData = await trxRes.json();
      setRecentTrx(trxData.data || []);
    } catch (error) {
      console.error("Gagal memuat data", error);
    }
  };

  useEffect(() => {
    loadDashboardData();
  }, []);

  const openInput = (type) => {
    setTxType(type);
    setIsModalOpen(true);
  };

  return (
    <div className="max-w-md mx-auto min-h-screen bg-gray-50 shadow-xl overflow-hidden flex flex-col font-sans relative">
      
      {/* HEADER & SALDO UTAMA */}
      <div className="bg-blue-800 text-white p-6 rounded-b-3xl shadow-md z-10">
        <h1 className="text-xl font-bold tracking-wide mb-6">Keuangan Sederhana</h1>
        <p className="text-blue-200 text-sm font-medium">{wallet?.name || 'Memuat...'}</p>
        <h2 className="text-4xl font-black mt-1 mb-2">{formatRp(wallet?.balance)}</h2>
        <span className="bg-blue-700/50 text-xs px-3 py-1 rounded-full font-bold uppercase tracking-wider">
          Saldo Saat Ini
        </span>
      </div>

      {/* KONTEN UTAMA */}
      <div className="flex-1 p-6 overflow-y-auto">
        
        {/* Tombol Aksi Cepat (FR-002) */}
        <div className="flex gap-4 mb-8 mt-2">
          <button 
            onClick={() => openInput('expense')} 
            className="flex-1 bg-red-100 text-red-700 py-4 rounded-2xl font-black text-lg hover:bg-red-200 active:scale-95 transition-all shadow-sm border border-red-200"
          >
            - KELUAR
          </button>
          <button 
            onClick={() => openInput('income')} 
            className="flex-1 bg-green-100 text-green-700 py-4 rounded-2xl font-black text-lg hover:bg-green-200 active:scale-95 transition-all shadow-sm border border-green-200"
          >
            + MASUK
          </button>
        </div>

        {/* Riwayat Transaksi (FR-003) */}
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-gray-800 text-lg">Riwayat Terakhir</h3>
          <button className="text-blue-600 font-bold text-sm hover:underline">Semua Riwayat</button>
        </div>

        <div className="space-y-3">
          {recentTrx.length === 0 ? (
            <p className="text-center text-gray-400 font-medium py-8 border-2 border-dashed border-gray-200 rounded-2xl">
              Belum ada transaksi dicatat.
            </p>
          ) : (
            recentTrx.map((trx) => (
              <div key={trx.id} className="bg-white p-4 rounded-2xl shadow-sm border border-gray-100 flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-gray-50 rounded-full flex items-center justify-center text-2xl border border-gray-100">
                    {trx.category?.icon || '📝'}
                  </div>
                  <div>
                    <p className="font-bold text-gray-800">{trx.category?.name || 'Lainnya'}</p>
                    <p className="text-xs text-gray-400 font-medium">
                      {new Date(trx.date_time).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}
                    </p>
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

      {/* Modal Input Transaksi (Ditampilkan berdasar state) */}
      {isModalOpen && (
        <InputModal 
          type={txType} 
          onClose={() => setIsModalOpen(false)}
          onSuccess={() => {
            setIsModalOpen(false);
            loadDashboardData(); // Memperbarui saldo dan riwayat secara instan
          }}
        />
      )}
    </div>
  );
}
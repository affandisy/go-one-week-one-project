import { useState, useEffect } from 'react';

const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function HistoryView() {
  const [transactions, setTransactions] = useState([]);

  const loadHistory = async () => {
    const res = await fetch('http://localhost:3000/api/v1/transactions/recent');
    const data = await res.json();
    setTransactions(data.data || []);
  };

  useEffect(() => { loadHistory(); }, []);

  const handleDelete = async (id, name) => {
    // Sesuai prinsip kehati-hatian, beri konfirmasi sederhana
    if (!window.confirm(`Hapus transaksi ${name}? Saldo akan dihitung ulang otomatis.`)) return;

    try {
      await fetch(`http://localhost:3000/api/v1/transactions/${id}`, { method: 'DELETE' });
      loadHistory(); // Muat ulang setelah sukses dihapus
    } catch (err) {
      alert("Gagal menghapus: " + err.message);
    }
  };

  return (
    <div className="flex-1 overflow-y-auto bg-gray-50 p-6 pb-24">
      <h2 className="text-2xl font-black text-gray-800 tracking-wide uppercase mb-6">Riwayat Lengkap</h2>
      
      <div className="space-y-3">
        {transactions.map((trx) => (
          <div key={trx.id} className="bg-white p-4 rounded-2xl shadow-sm border border-gray-100">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-3">
                <span className="text-2xl">{trx.category?.icon || '📝'}</span>
                <div>
                  <p className="font-bold text-gray-800">{trx.category?.name || 'Lainnya'}</p>
                  <p className="text-xs text-gray-400 font-medium">{new Date(trx.date_time).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })}</p>
                </div>
              </div>
              <p className={`font-black text-lg ${trx.type === 'expense' ? 'text-red-600' : 'text-green-600'}`}>
                {trx.type === 'expense' ? '-' : '+'}{formatRp(trx.amount)}
              </p>
            </div>
            
            {/* Tombol Hapus (FR-003) */}
            <div className="flex justify-end border-t border-gray-50 pt-2">
              <button 
                onClick={() => handleDelete(trx.id, trx.category?.name)}
                className="text-xs font-bold text-red-500 hover:text-red-700 bg-red-50 px-3 py-1 rounded-lg"
              >
                Hapus
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
import { useState } from 'react';

export default function AddWalletModal({ onClose, onSuccess }) {
  const [name, setName] = useState('');
  const [initialBalance, setInitialBalance] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!name.trim()) return alert('Nama dompet wajib diisi');

    setIsLoading(true);
    try {
      const res = await fetch('/api/v1/wallets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: name.trim(),
          initial_balance: parseFloat(initialBalance) || 0
        })
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || 'Gagal menyimpan dompet');
      }
      
      onSuccess();
    } catch (err) {
      alert('Terjadi kesalahan: ' + err.message);
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-end sm:items-center justify-center z-50 p-0 sm:p-4 animate-fade-in">
      <div className="bg-white w-full max-w-md rounded-t-3xl sm:rounded-3xl shadow-2xl flex flex-col">
        
        {/* Header Modal */}
        <div className="p-6 border-b border-gray-100 flex justify-between items-center bg-blue-50 sm:rounded-t-3xl">
          <h3 className="text-xl font-black text-blue-900">Buat Dompet Baru</h3>
          <button onClick={onClose} className="w-8 h-8 bg-blue-200 rounded-full text-blue-800 font-bold hover:bg-blue-300">
            ✕
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6">
          {/* Input Nama Dompet */}
          <div className="mb-6">
            <label className="block text-sm font-bold text-gray-500 mb-2 uppercase tracking-wide">Nama Dompet *</label>
            <input 
              type="text" 
              autoFocus
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full text-2xl font-black text-gray-800 border-b-2 border-gray-200 focus:border-blue-600 outline-none py-2 bg-transparent"
              placeholder="Contoh: Tabungan Liburan"
            />
          </div>

          {/* Input Saldo Awal */}
          <div className="mb-8">
            <label className="block text-sm font-bold text-gray-500 mb-2 uppercase tracking-wide">Saldo Awal (Opsional)</label>
            <input 
              type="number" 
              value={initialBalance}
              onChange={(e) => setInitialBalance(e.target.value)}
              className="w-full text-2xl font-black text-blue-600 border-b-2 border-gray-200 focus:border-blue-600 outline-none py-2 bg-transparent"
              placeholder="0"
            />
          </div>

          {/* Tombol Simpan */}
          <button 
            type="submit" 
            disabled={isLoading || !name}
            className="w-full py-5 rounded-2xl font-black text-white text-lg tracking-wide transition-all shadow-lg active:scale-95 bg-blue-600 hover:bg-blue-700 shadow-blue-600/30 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? 'MENYIMPAN...' : 'SIMPAN DOMPET'}
          </button>
        </form>
      </div>
    </div>
  );
}
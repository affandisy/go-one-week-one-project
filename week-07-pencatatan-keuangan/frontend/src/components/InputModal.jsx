import { useState, useEffect } from 'react';

export default function InputModal({ type, onClose, onSuccess }) {
  const [categories, setCategories] = useState([]);
  const [selectedCat, setSelectedCat] = useState('');
  const [amount, setAmount] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const isExpense = type === 'expense';
  const themeColor = isExpense ? 'red' : 'green';

  useEffect(() => {
    // Menarik daftar kategori aktif dari Golang (FR-006)
    fetch(`http://localhost:3000/api/v1/categories?type=${type}`)
      .then(res => res.json())
      .then(data => {
        setCategories(data.data || []);
        if (data.data?.length > 0) {
          setSelectedCat(data.data[0].id); // Pilih kategori pertama secara bawaan
        }
      });
  }, [type]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!amount || parseFloat(amount) <= 0) return alert('Detail nominal harus lebih dari 0');
    
    setIsLoading(true);
    try {
      const res = await fetch('http://localhost:3000/api/v1/transactions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          type: type,
          category_id: selectedCat,
          amount: parseFloat(amount),
          note: '' // Dikosongkan untuk MVP agar sangat cepat
        })
      });

      if (!res.ok) throw new Error('Gagal menyimpan');
      onSuccess();
    } catch (err) {
      alert('Terjadi kesalahan: ' + err.message);
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-end sm:items-center justify-center z-50 p-0 sm:p-4 animate-fade-in">
      <div className="bg-white w-full max-w-md rounded-t-3xl sm:rounded-3xl shadow-2xl flex flex-col max-h-[90vh]">
        
        {/* Header Modal */}
        <div className="p-6 border-b border-gray-100 flex justify-between items-center">
          <h3 className="text-xl font-black text-gray-800">
            {isExpense ? 'Catat Pengeluaran' : 'Catat Pemasukan'}
          </h3>
          <button onClick={onClose} className="w-8 h-8 bg-gray-100 rounded-full text-gray-500 font-bold hover:bg-gray-200">
            ✕
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 flex-1 overflow-y-auto">
          {/* Input Nominal Cepat */}
          <div className="mb-8">
            <label className="block text-sm font-bold text-gray-400 mb-2 uppercase tracking-wide">Nominal (Rp)</label>
            <input 
              type="number" 
              autoFocus
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className={`w-full text-5xl font-black text-${themeColor}-600 border-b-2 border-gray-200 focus:border-${themeColor}-500 outline-none py-2 bg-transparent`}
              placeholder="0"
            />
          </div>

          {/* Pemilihan Kategori (Grid Ikon) */}
          <div className="mb-8">
            <label className="block text-sm font-bold text-gray-400 mb-4 uppercase tracking-wide">Pilih Kategori</label>
            <div className="grid grid-cols-4 gap-3">
              {categories.map((cat) => (
                <button
                  key={cat.id}
                  type="button"
                  onClick={() => setSelectedCat(cat.id)}
                  className={`p-3 rounded-2xl flex flex-col items-center justify-center transition-all border-2
                    ${selectedCat === cat.id 
                      ? `bg-${themeColor}-50 border-${themeColor}-500 shadow-sm scale-105` 
                      : 'bg-gray-50 border-transparent hover:bg-gray-100 text-gray-500'}`}
                >
                  <span className="text-3xl mb-2">{cat.icon}</span>
                  <span className={`text-[10px] font-black uppercase tracking-tight ${selectedCat === cat.id ? `text-${themeColor}-700` : ''}`}>
                    {cat.name}
                  </span>
                </button>
              ))}
            </div>
          </div>

          {/* Tombol Eksekusi Akhir */}
          <button 
            type="submit" 
            disabled={isLoading || !amount}
            className={`w-full py-5 rounded-2xl font-black text-white text-lg tracking-wide transition-all shadow-lg active:scale-95
              ${isExpense ? 'bg-red-600 hover:bg-red-700 shadow-red-600/30' : 'bg-green-600 hover:bg-green-700 shadow-green-600/30'}
              disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            {isLoading ? 'MENYIMPAN...' : 'SIMPAN TRANSAKSI'}
          </button>
        </form>
      </div>
    </div>
  );
}
import { useState } from 'react';
import { useBillingStore } from '../store/useBillingStore';
import { Search, Plus, Trash2, Receipt, User, CheckCircle2 } from 'lucide-react';

// --- DATA TIRUAN (MOCK DATA) UNTUK MVP UI ---
const MOCK_OWNERS = [
  { id: 'o1', name: 'Sinta', phone: '0812-3333-4444' },
  { id: 'o2', name: 'Budi Santoso', phone: '0856-7777-8888' }
];

const MOCK_PETS = [
  { id: 'p1', owner_id: 'o1', name: 'Milo', species: 'Kucing', breed: 'Domestik' },
  { id: 'p2', owner_id: 'o1', name: 'Luna', species: 'Kucing', breed: 'Persia' },
  { id: 'p3', owner_id: 'o2', name: 'Max', species: 'Anjing', breed: 'Golden Retriever' }
];

const MOCK_SERVICES = [
  { id: 's1', name: 'Premium Grooming', base_price: 150000 },
  { id: 's2', name: 'Terapi Nutrisi Bulu', base_price: 75000 },
  { id: 's3', name: 'Penitipan Harian (Boarding)', base_price: 100000 },
  { id: 's4', name: 'Cek Berat Badan & Gizi', base_price: 50000 }
];

// Utilitas Format Rupiah
const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function CashierView() {
  // Menghubungkan State Zustand
  const { 
    currentOwner, ownerPets, cartItems, 
    startTransaction, clearTransaction, addItem, removeItem, 
    getTotal, getItemsGroupedByPet 
  } = useBillingStore();

  const [searchQuery, setSearchQuery] = useState('');
  const [selectedPetId, setSelectedPetId] = useState(null);

  // Simulasi pencarian API
  const handleSearchOwner = (e) => {
    e.preventDefault();
    const owner = MOCK_OWNERS.find(o => o.name.toLowerCase().includes(searchQuery.toLowerCase()));
    
    if (owner) {
      const pets = MOCK_PETS.filter(p => p.owner_id === owner.id);
      startTransaction(owner, pets);
      if (pets.length > 0) setSelectedPetId(pets[0].id); // Otomatis pilih hewan pertama
    } else {
      alert('Pemilik tidak ditemukan!');
    }
  };

  const handleCheckout = async () => {
    if (cartItems.length === 0) return alert('Keranjang tagihan masih kosong!');
    
    // Di sini nantinya Anda akan melakukan POST /api/v1/invoices ke backend Golang
    alert(`Tagihan berhasil dibuat! Total: ${formatRp(getTotal())}`);
    clearTransaction();
    setSearchQuery('');
  };

  // Mengambil data keranjang yang sudah dikelompokkan per hewan
  const groupedCart = getItemsGroupedByPet();

  return (
    <div className="flex flex-col lg:flex-row h-full">
      
      {/* KOLOM KIRI: Interaksi & Pemilihan Layanan */}
      <div className="flex-1 p-8 flex flex-col h-full overflow-y-auto">
        <h2 className="text-2xl font-black text-slate-800 mb-6 flex items-center gap-2">
          <Receipt className="text-indigo-600" /> Modul Kasir Utama
        </h2>

        {/* 1. Pencarian Pemilik */}
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 mb-6">
          <form onSubmit={handleSearchOwner} className="flex gap-4">
            <div className="flex-1 relative">
              <Search className="absolute left-4 top-3 text-slate-400" size={20} />
              <input 
                type="text" 
                placeholder="Cari nama pemilik (Contoh: Sinta)..." 
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-12 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 font-medium"
              />
            </div>
            <button type="submit" className="px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl shadow-md transition-all active:scale-95">
              Cari Data
            </button>
          </form>
        </div>

        {currentOwner && (
          <div className="animate-fade-in flex-1 flex flex-col">
            {/* Info Pelanggan */}
            <div className="flex items-center gap-4 bg-indigo-50 p-4 rounded-xl border border-indigo-100 mb-6">
              <div className="w-12 h-12 bg-indigo-200 text-indigo-700 rounded-full flex items-center justify-center font-black text-xl">
                {currentOwner.name.charAt(0)}
              </div>
              <div>
                <p className="text-sm font-bold text-indigo-400 uppercase tracking-wider">Pelanggan Aktif</p>
                <p className="text-lg font-black text-indigo-900">{currentOwner.name} <span className="text-sm font-medium text-indigo-600 ml-2">({currentOwner.phone})</span></p>
              </div>
            </div>

            {/* 2. Pilihan Hewan Peliharaan */}
            <h3 className="text-sm font-bold text-slate-400 uppercase tracking-wider mb-3">Pilih Hewan Untuk Dikenakan Layanan</h3>
            <div className="flex gap-3 mb-8 overflow-x-auto pb-2">
              {ownerPets.map(pet => (
                <button
                  key={pet.id}
                  onClick={() => setSelectedPetId(pet.id)}
                  className={`min-w-[140px] p-4 rounded-2xl border-2 text-left transition-all flex flex-col gap-1
                    ${selectedPetId === pet.id 
                      ? 'border-indigo-600 bg-indigo-600 text-white shadow-lg shadow-indigo-200 scale-105' 
                      : 'border-slate-200 bg-white text-slate-600 hover:border-indigo-300'}`}
                >
                  <span className={`text-xs font-bold uppercase tracking-wider ${selectedPetId === pet.id ? 'text-indigo-200' : 'text-slate-400'}`}>
                    {pet.species}
                  </span>
                  <span className="text-lg font-black">{pet.name}</span>
                </button>
              ))}
            </div>

            {/* 3. Daftar Katalog Layanan */}
            {selectedPetId && (
              <div className="flex-1">
                <h3 className="text-sm font-bold text-slate-400 uppercase tracking-wider mb-3">
                  Tambahkan Layanan untuk {ownerPets.find(p => p.id === selectedPetId)?.name}
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  {MOCK_SERVICES.map(service => (
                    <div key={service.id} className="bg-white p-4 rounded-2xl border border-slate-100 shadow-sm flex flex-col justify-between group hover:border-indigo-200 transition-colors">
                      <div className="mb-4">
                        <h4 className="font-bold text-slate-800">{service.name}</h4>
                        <p className="text-indigo-600 font-black mt-1">{formatRp(service.base_price)}</p>
                      </div>
                      <button 
                        onClick={() => addItem(ownerPets.find(p => p.id === selectedPetId), service)}
                        className="w-full py-2 bg-slate-50 hover:bg-indigo-50 text-indigo-600 font-bold rounded-xl border border-slate-200 hover:border-indigo-200 transition-all flex items-center justify-center gap-2"
                      >
                        <Plus size={18} /> Tambah
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* KOLOM KANAN: Struk Tagihan Dinamis (FR-004) */}
      <div className="w-full lg:w-[400px] bg-white border-l border-slate-200 shadow-2xl flex flex-col h-full z-10">
        <div className="p-6 bg-slate-800 text-white border-b border-slate-700">
          <h2 className="text-xl font-black">Rincian Tagihan</h2>
          <p className="text-slate-400 text-sm mt-1">{currentOwner ? `Pelanggan: ${currentOwner.name}` : 'Belum ada pelanggan'}</p>
        </div>

        <div className="flex-1 p-6 overflow-y-auto bg-slate-50">
          {cartItems.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-slate-400">
              <Receipt size={64} className="mb-4 opacity-20" />
              <p className="font-bold">Belum ada layanan dipilih.</p>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Me-render item yang telah dikelompokkan berdasarkan hewan */}
              {Object.entries(groupedCart).map(([petId, data]) => (
                <div key={petId} className="bg-white p-4 rounded-xl shadow-sm border border-slate-200">
                  <div className="flex items-center gap-2 border-b border-slate-100 pb-3 mb-3">
                    <CheckCircle2 size={16} className="text-green-500" />
                    <h4 className="font-black text-slate-800 uppercase tracking-wide">Hewan: {data.pet_name}</h4>
                  </div>
                  
                  <div className="space-y-3">
                    {data.services.map(item => (
                      <div key={item.id} className="flex justify-between items-start group">
                        <div className="flex-1 pr-4">
                          <p className="text-sm font-bold text-slate-700">{item.service_name}</p>
                          <p className="text-xs font-black text-indigo-600 mt-0.5">{formatRp(item.price)}</p>
                        </div>
                        <button 
                          onClick={() => removeItem(item.id)}
                          className="p-1.5 text-slate-300 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors opacity-0 group-hover:opacity-100"
                          title="Hapus"
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    ))}
                  </div>
                  
                  <div className="mt-4 pt-3 border-t border-slate-100 flex justify-between items-center">
                    <span className="text-xs font-bold text-slate-400 uppercase">Subtotal {data.pet_name}</span>
                    <span className="text-sm font-black text-slate-800">{formatRp(data.subtotal)}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Panel Checkout Bawah */}
        <div className="p-6 bg-white border-t border-slate-200 shadow-[0_-10px_20px_-10px_rgba(0,0,0,0.05)]">
          <div className="flex justify-between items-end mb-6">
            <span className="font-bold text-slate-500 uppercase tracking-wider text-sm">Total Tagihan</span>
            <span className="text-3xl font-black text-slate-800">{formatRp(getTotal())}</span>
          </div>
          
          <button 
            onClick={handleCheckout}
            disabled={cartItems.length === 0}
            className="w-full py-4 bg-indigo-600 hover:bg-indigo-700 disabled:bg-slate-300 text-white font-black text-lg rounded-2xl shadow-xl shadow-indigo-200 transition-all active:scale-95 flex justify-center items-center gap-2"
          >
            Cetak Invoice & Bayar
          </button>
        </div>
      </div>
    </div>
  );
}
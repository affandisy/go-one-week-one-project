import { useState } from 'react';
import { Syringe, Plus, Activity, Calendar, Info } from 'lucide-react';

// --- DATA TIRUAN (MOCK DATA) ---
const MOCK_PETS = [
  { id: 'p1', owner: 'Sinta', name: 'Miko', species: 'Kucing', breed: 'Domestik / Kampung', weight: 4.2 },
  { id: 'p2', owner: 'Budi', name: 'Luna', species: 'Kucing', breed: 'Persia', weight: 3.8 },
  { id: 'p3', owner: 'Agus', name: 'Max', species: 'Anjing', breed: 'Golden Retriever', weight: 28.5 }
];

const MOCK_HISTORY = [
  {
    id: 'log1',
    pet_id: 'p1',
    date: '2026-05-20T10:00:00',
    food_brand: 'Orijen Original Cat',
    calories: 280,
    health_notes: 'Transisi ke makanan tinggi protein selesai. Bulu rontok berkurang drastis, mulai terlihat lebih lebat dan mengkilap.'
  },
  {
    id: 'log2',
    pet_id: 'p1',
    date: '2026-05-13T09:30:00',
    food_brand: 'Orijen (Campur Royal Canin Recovery)',
    calories: 310,
    health_notes: 'Fase awal perbaikan gizi. Massa otot mulai terbentuk, nafsu makan sangat baik.'
  }
];

export default function NutritionLog() {
  const [selectedPetId, setSelectedPetId] = useState(MOCK_PETS[0].id);
  const [formData, setFormData] = useState({
    food_brand: '',
    calories: '',
    health_notes: ''
  });

  const activePet = MOCK_PETS.find(p => p.id === selectedPetId);
  const activeHistory = MOCK_HISTORY.filter(h => h.pet_id === selectedPetId);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!formData.food_brand || !formData.calories) return alert('Merek makanan dan kalori wajib diisi!');
    
    // Di sini nantinya akan memanggil POST /api/v1/nutrition
    alert(`Rekam gizi baru untuk ${activePet.name} berhasil disimpan!`);
    setFormData({ food_brand: '', calories: '', health_notes: '' });
  };

  return (
    <div className="flex flex-col lg:flex-row h-full">
      
      {/* KOLOM KIRI: Daftar Hewan & Form Pencatatan */}
      <div className="w-full lg:w-1/3 bg-white border-r border-slate-200 shadow-sm flex flex-col h-full z-10">
        <div className="p-6 bg-orange-50 border-b border-orange-100">
          <h2 className="text-xl font-black text-orange-900 flex items-center gap-2">
            <Syringe className="text-orange-600" /> Rekam Gizi
          </h2>
          <p className="text-orange-700 text-sm mt-1">Pantau kualitas diet & transformasi fisik</p>
        </div>

        <div className="flex-1 overflow-y-auto">
          {/* Pilih Hewan */}
          <div className="p-6 border-b border-slate-100">
            <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3">Pilih Pasien / Hewan</h3>
            <select 
              value={selectedPetId}
              onChange={(e) => setSelectedPetId(e.target.value)}
              className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl font-bold text-slate-700 focus:outline-none focus:border-orange-500 cursor-pointer"
            >
              {MOCK_PETS.map(pet => (
                <option key={pet.id} value={pet.id}>
                  {pet.name} ({pet.species} {pet.breed}) - Pemilik: {pet.owner}
                </option>
              ))}
            </select>
            
            {activePet && (
              <div className="mt-4 p-4 bg-slate-50 rounded-xl border border-slate-100 flex items-center gap-3">
                <div className="w-10 h-10 bg-orange-100 text-orange-600 rounded-full flex items-center justify-center font-black">
                  {activePet.name.charAt(0)}
                </div>
                <div>
                  <p className="font-bold text-slate-800">{activePet.name}</p>
                  <p className="text-xs text-slate-500">{activePet.breed} • {activePet.weight} kg</p>
                </div>
              </div>
            )}
          </div>

          {/* Form Input Gizi Baru */}
          <form onSubmit={handleSubmit} className="p-6">
            <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-4 flex items-center gap-2">
              <Plus size={14} /> Tambah Catatan Baru
            </h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-bold text-slate-600 mb-1">Merek / Jenis Diet</label>
                <input 
                  type="text" 
                  value={formData.food_brand}
                  onChange={(e) => setFormData({...formData, food_brand: e.target.value})}
                  placeholder="Contoh: Raw Diet, Royal Canin..."
                  className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange-500 text-sm"
                />
              </div>
              
              <div>
                <label className="block text-sm font-bold text-slate-600 mb-1">Asupan Kalori (kkal)</label>
                <input 
                  type="number" 
                  value={formData.calories}
                  onChange={(e) => setFormData({...formData, calories: e.target.value})}
                  placeholder="0"
                  className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange-500 text-sm"
                />
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-600 mb-1">Catatan Transformasi / Detail Fisik</label>
                <textarea 
                  value={formData.health_notes}
                  onChange={(e) => setFormData({...formData, health_notes: e.target.value})}
                  placeholder="Catat perubahan pada bulu, pencernaan, atau aktivitas..."
                  rows="4"
                  className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange-500 text-sm resize-none"
                ></textarea>
              </div>

              <button 
                type="submit" 
                className="w-full py-3 mt-2 bg-orange-600 hover:bg-orange-700 text-white font-bold rounded-xl shadow-md transition-all active:scale-95"
              >
                Simpan Catatan Gizi
              </button>
            </div>
          </form>
        </div>
      </div>

      {/* KOLOM KANAN: Linimasa (Timeline) Transformasi */}
      <div className="flex-1 bg-slate-50 p-8 overflow-y-auto">
        <h2 className="text-2xl font-black text-slate-800 mb-6 flex items-center gap-2">
          <Activity className="text-orange-500" /> Riwayat Perkembangan: {activePet?.name}
        </h2>

        {activeHistory.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-64 text-slate-400">
            <Info size={48} className="mb-4 opacity-20" />
            <p className="font-bold">Belum ada riwayat pencatatan gizi.</p>
          </div>
        ) : (
          <div className="space-y-6 relative before:absolute before:inset-0 before:ml-5 before:-translate-x-px md:before:mx-auto md:before:translate-x-0 before:h-full before:w-0.5 before:bg-gradient-to-b before:from-transparent before:via-slate-200 before:to-transparent">
            
            {activeHistory.map((log, index) => (
              <div key={log.id} className="relative flex items-center justify-between md:justify-normal md:odd:flex-row-reverse group is-active">
                
                {/* Ikon Titik Tengah */}
                <div className="flex items-center justify-center w-10 h-10 rounded-full border-4 border-slate-50 bg-orange-100 text-orange-600 shadow shrink-0 md:order-1 md:group-odd:-translate-x-1/2 md:group-even:translate-x-1/2 z-10">
                  <Calendar size={16} />
                </div>
                
                {/* Kartu Konten */}
                <div className="w-[calc(100%-4rem)] md:w-[calc(50%-2.5rem)] bg-white p-5 rounded-2xl shadow-sm border border-slate-100 hover:border-orange-200 transition-colors">
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <h4 className="font-black text-slate-800 text-lg">{log.food_brand}</h4>
                      <p className="text-xs font-bold text-orange-600 bg-orange-50 inline-block px-2 py-1 rounded-md mt-1">
                        {log.calories} kkal / hari
                      </p>
                    </div>
                    <time className="text-xs font-bold text-slate-400">
                      {new Date(log.date).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
                    </time>
                  </div>
                  
                  <div className="text-sm text-slate-600 leading-relaxed bg-slate-50 p-3 rounded-xl border border-slate-100">
                    <p className="font-medium">Detail Transformasi:</p>
                    <p className="mt-1">{log.health_notes}</p>
                  </div>
                </div>

              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
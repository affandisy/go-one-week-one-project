import { Activity, DollarSign, Scissors, Home, CheckCircle2, Clock } from 'lucide-react';

// --- DATA TIRUAN (MOCK DATA) UNTUK DASBOR ---
const MOCK_STATS = {
  dailyRevenue: 850000,
  activeBoarding: 4,
  completedGrooming: 3,
  pendingGrooming: 2
};

const MOCK_QUEUE = [
  { id: 'q1', petName: 'Miko', species: 'Kucing Domestik', service: 'Premium Grooming & Spa', time: '09:00', status: 'Selesai' },
  { id: 'q2', petName: 'Luna', species: 'Kucing Persia', service: 'Konsultasi Nutrisi & Diet', time: '11:30', status: 'Proses' },
  { id: 'q3', petName: 'Max', species: 'Anjing Golden', service: 'Grooming Standar', time: '14:00', status: 'Menunggu' },
];

const MOCK_INPATIENTS = [
  { id: 'in1', petName: 'Oreo', species: 'Kucing Domestik', owner: 'Sinta', note: 'Observasi pencernaan pasca ganti diet' },
  { id: 'in2', petName: 'Belly', species: 'Kucing Anggora', owner: 'Budi', note: 'Penitipan harian (Boarding)' },
];

const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function Dashboard() {
  return (
    <div className="p-8 h-full overflow-y-auto bg-slate-50">
      
      {/* HEADER */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-black text-slate-800 tracking-tight">Dasbor Operasional</h1>
          <p className="text-slate-500 font-medium mt-1">Ringkasan aktivitas klinik hari ini: {new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })}</p>
        </div>
        <div className="flex items-center gap-3 bg-indigo-50 px-4 py-2 rounded-xl border border-indigo-100">
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-indigo-500"></span>
          </span>
          <span className="text-sm font-bold text-indigo-800">Sistem Online</span>
        </div>
      </div>

      {/* KARTU STATISTIK (METRICS CARDS) */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {/* Kartu Pendapatan */}
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition-shadow">
          <div className="w-14 h-14 bg-green-100 text-green-600 rounded-2xl flex items-center justify-center">
            <DollarSign size={28} strokeWidth={2.5} />
          </div>
          <div>
            <p className="text-sm font-bold text-slate-400 uppercase tracking-wider">Estimasi Kas</p>
            <h3 className="text-2xl font-black text-slate-800">{formatRp(MOCK_STATS.dailyRevenue)}</h3>
          </div>
        </div>

        {/* Kartu Rawat Inap */}
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition-shadow">
          <div className="w-14 h-14 bg-blue-100 text-blue-600 rounded-2xl flex items-center justify-center">
            <Home size={28} strokeWidth={2.5} />
          </div>
          <div>
            <p className="text-sm font-bold text-slate-400 uppercase tracking-wider">Rawat Inap</p>
            <h3 className="text-2xl font-black text-slate-800">{MOCK_STATS.activeBoarding} <span className="text-sm text-slate-500 font-medium">Hewan</span></h3>
          </div>
        </div>

        {/* Kartu Selesai */}
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition-shadow">
          <div className="w-14 h-14 bg-indigo-100 text-indigo-600 rounded-2xl flex items-center justify-center">
            <CheckCircle2 size={28} strokeWidth={2.5} />
          </div>
          <div>
            <p className="text-sm font-bold text-slate-400 uppercase tracking-wider">Selesai</p>
            <h3 className="text-2xl font-black text-slate-800">{MOCK_STATS.completedGrooming} <span className="text-sm text-slate-500 font-medium">Layanan</span></h3>
          </div>
        </div>

        {/* Kartu Antrean */}
        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 flex items-center gap-5 hover:shadow-md transition-shadow">
          <div className="w-14 h-14 bg-orange-100 text-orange-600 rounded-2xl flex items-center justify-center">
            <Clock size={28} strokeWidth={2.5} />
          </div>
          <div>
            <p className="text-sm font-bold text-slate-400 uppercase tracking-wider">Menunggu</p>
            <h3 className="text-2xl font-black text-slate-800">{MOCK_STATS.pendingGrooming} <span className="text-sm text-slate-500 font-medium">Layanan</span></h3>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* TABEL ANTREAN HARI INI */}
        <div className="lg:col-span-2 bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden flex flex-col">
          <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-slate-800">
            <h2 className="text-lg font-black text-white flex items-center gap-2">
              <Scissors className="text-indigo-400" size={20} /> Jadwal & Antrean Layanan
            </h2>
            <button className="text-sm font-bold text-indigo-300 hover:text-white transition-colors">Lihat Semua</button>
          </div>
          
          <div className="flex-1 p-0">
            <table className="w-full text-left">
              <thead className="bg-slate-50 border-b border-slate-100">
                <tr>
                  <th className="p-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Waktu</th>
                  <th className="p-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Pasien</th>
                  <th className="p-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Layanan</th>
                  <th className="p-4 text-xs font-bold text-slate-400 uppercase tracking-wider">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-50">
                {MOCK_QUEUE.map((item) => (
                  <tr key={item.id} className="hover:bg-slate-50 transition-colors">
                    <td className="p-4 font-black text-slate-700">{item.time}</td>
                    <td className="p-4">
                      <p className="font-bold text-slate-800">{item.petName}</p>
                      <p className="text-xs text-slate-500">{item.species}</p>
                    </td>
                    <td className="p-4 font-medium text-slate-600">{item.service}</td>
                    <td className="p-4">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider
                        ${item.status === 'Selesai' ? 'bg-green-100 text-green-700' : 
                          item.status === 'Proses' ? 'bg-blue-100 text-blue-700' : 
                          'bg-orange-100 text-orange-700'}`}>
                        {item.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* DAFTAR RAWAT INAP */}
        <div className="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden flex flex-col">
          <div className="p-6 border-b border-slate-100 flex justify-between items-center bg-indigo-50">
            <h2 className="text-lg font-black text-indigo-900 flex items-center gap-2">
              <Activity className="text-indigo-600" size={20} /> Fasilitas Rawat Inap
            </h2>
          </div>
          
          <div className="flex-1 p-6 space-y-4 overflow-y-auto">
            {MOCK_INPATIENTS.map((patient) => (
              <div key={patient.id} className="p-4 bg-slate-50 rounded-xl border border-slate-100 hover:border-indigo-200 transition-colors">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h4 className="font-black text-slate-800">{patient.petName}</h4>
                    <p className="text-xs font-bold text-slate-500">{patient.species}</p>
                  </div>
                  <span className="text-[10px] font-black uppercase tracking-wider bg-indigo-100 text-indigo-700 px-2 py-1 rounded-md">
                    Pemilik: {patient.owner}
                  </span>
                </div>
                <div className="text-sm text-slate-600 bg-white p-3 rounded-lg border border-slate-100 mt-3">
                  <span className="font-bold block mb-1">Catatan Staf:</span>
                  {patient.note}
                </div>
              </div>
            ))}
          </div>
        </div>

      </div>
    </div>
  );
}
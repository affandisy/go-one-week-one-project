import { useState, useEffect } from 'react';

const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function ReportView() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [report, setReport] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchReport = async (year, month) => {
    setIsLoading(true);
    try {
      const res = await fetch(`http://localhost:3000/api/v1/reports/monthly?year=${year}&month=${month}`);
      const data = await res.json();
      setReport(data.data);
    } catch (err) {
      console.error("Gagal memuat laporan", err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchReport(currentDate.getFullYear(), currentDate.getMonth() + 1);
  }, [currentDate]);

  const prevMonth = () => setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1));
  const nextMonth = () => setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1));

  const monthName = currentDate.toLocaleString('id-ID', { month: 'long', year: 'numeric' });

  return (
    <div className="flex-1 overflow-y-auto bg-gray-50 p-6 pb-24">
      {/* Navigasi Bulan (FR-004) */}
      <div className="flex justify-between items-center bg-white p-4 rounded-2xl shadow-sm mb-6 border border-gray-100">
        <button onClick={prevMonth} className="w-10 h-10 bg-gray-100 rounded-full font-black text-gray-600 hover:bg-gray-200">{"<"}</button>
        <h2 className="text-lg font-black text-gray-800 tracking-wide uppercase">{monthName}</h2>
        <button onClick={nextMonth} className="w-10 h-10 bg-gray-100 rounded-full font-black text-gray-600 hover:bg-gray-200">{">"}</button>
      </div>

      {isLoading ? (
        <p className="text-center text-gray-400 font-bold py-10">Memuat data laporan...</p>
      ) : report ? (
        <>
          {/* Ringkasan Saldo & Perbandingan */}
          <div className="bg-blue-800 text-white p-6 rounded-3xl shadow-lg mb-6 relative overflow-hidden">
            <p className="text-blue-200 text-sm font-bold uppercase tracking-wider mb-1">Bersih Bulan Ini</p>
            <h3 className="text-4xl font-black mb-4">{formatRp(report.net_balance)}</h3>
            
            <div className="flex gap-4 border-t border-blue-700 pt-4">
              <div className="flex-1">
                <p className="text-xs text-blue-300 font-bold uppercase">Masuk</p>
                <p className="font-bold text-green-400">{formatRp(report.total_income)}</p>
              </div>
              <div className="flex-1">
                <p className="text-xs text-blue-300 font-bold uppercase">Keluar</p>
                <p className="font-bold text-red-400">{formatRp(report.total_expense)}</p>
              </div>
            </div>

            {/* Indikator vs Bulan Lalu (FR-004) */}
            <div className={`absolute top-6 right-6 px-3 py-1 rounded-full text-xs font-black
              ${report.vs_last_month_percent > 0 ? 'bg-red-500 text-white' : 'bg-green-500 text-white'}`}>
              {report.vs_last_month_percent > 0 ? '↑' : '↓'} {Math.abs(report.vs_last_month_percent).toFixed(1)}% vs Lalu
            </div>
          </div>

          {/* Rincian Kategori Pengeluaran (Tabel Visual) */}
          <h3 className="font-black text-gray-800 text-lg mb-4">Rincian Pengeluaran</h3>
          <div className="bg-white rounded-3xl shadow-sm border border-gray-100 p-2">
            {report.expense_breakdown?.length === 0 ? (
              <p className="text-center text-gray-400 text-sm py-6 font-bold">Tidak ada pengeluaran.</p>
            ) : (
              report.expense_breakdown?.map((cat, idx) => (
                <div key={idx} className="p-4 border-b border-gray-50 last:border-0">
                  <div className="flex justify-between items-end mb-2">
                    <span className="font-bold text-gray-700">{cat.category_name}</span>
                    <span className="font-black text-gray-900">{formatRp(cat.amount)}</span>
                  </div>
                  {/* Progress Bar Sederhana sebagai pengganti Pie Chart berat */}
                  <div className="flex items-center gap-3">
                    <div className="flex-1 h-3 bg-gray-100 rounded-full overflow-hidden">
                      <div className="h-full rounded-full" style={{ width: `${cat.percentage}%`, backgroundColor: cat.color }}></div>
                    </div>
                    <span className="text-xs font-black text-gray-500 w-10 text-right">{cat.percentage.toFixed(0)}%</span>
                  </div>
                </div>
              ))
            )}
          </div>
        </>
      ) : null}
    </div>
  );
}
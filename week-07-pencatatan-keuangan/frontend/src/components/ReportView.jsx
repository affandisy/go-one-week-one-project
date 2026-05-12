import { useState, useEffect } from 'react';
import jsPDF from 'jspdf';
import 'jspdf-autotable';

const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num || 0);

export default function ReportView() {
  const [currentDate, setCurrentDate] = useState(new Date());
  const [report, setReport] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchReport = async (year, month) => {
    setIsLoading(true);
    try {
      const res = await fetch(`/api/v1/reports/monthly?year=${year}&month=${month}`);
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

  // --- FUNGSI EKSPOR PDF ---
  const handleExportPDF = () => {
    if (!report) return;

    // Inisiasi dokumen PDF (Orientasi Portrait, Satuan mm, Ukuran A4)
    const doc = new jsPDF('p', 'mm', 'a4');

    // 1. Header Dokumen
    doc.setFontSize(18);
    doc.setFont("helvetica", "bold");
    doc.text('Laporan Keuangan Bulanan', 14, 22);

    // 2. Ringkasan Laporan
    doc.setFontSize(11);
    doc.setFont("helvetica", "normal");
    doc.text(`Bulan: ${monthName}`, 14, 32);
    doc.text(`Total Pemasukan: ${formatRp(report.total_income)}`, 14, 38);
    doc.text(`Total Pengeluaran: ${formatRp(report.total_expense)}`, 14, 44);
    doc.text(`Saldo Bersih (Net): ${formatRp(report.net_balance)}`, 14, 50);

    // 3. Konfigurasi Tabel (jspdf-autotable)
    const tableColumn = ["Kategori", "Total Pengeluaran", "Persentase"];
    const tableRows = [];

    // Menyusun data baris tabel
    if (report.expense_breakdown && report.expense_breakdown.length > 0) {
      report.expense_breakdown.forEach(cat => {
        const rowData = [
          cat.category_name,
          formatRp(cat.amount),
          `${cat.percentage.toFixed(1)}%`
        ];
        tableRows.push(rowData);
      });
    } else {
      tableRows.push(["Tidak ada pengeluaran", "-", "-"]);
    }

    // Menggambar tabel di PDF
    doc.autoTable({
      startY: 58,
      head: [tableColumn],
      body: tableRows,
      theme: 'striped',
      headStyles: { fillColor: [30, 64, 175] }, // Menggunakan warna setara bg-blue-800
      styles: { fontSize: 10, cellPadding: 4 }
    });

    // 4. Memicu Unduhan (Download)
    const fileName = `Laporan_Keuangan_${currentDate.getFullYear()}_${currentDate.getMonth() + 1}.pdf`;
    doc.save(fileName);
  };

  return (
    <div className="flex-1 overflow-y-auto bg-gray-50 p-6 pb-24">
      
      {/* Navigasi Bulan & Tombol PDF */}
      <div className="flex flex-col gap-4 mb-6">
        <div className="flex justify-between items-center bg-white p-4 rounded-2xl shadow-sm border border-gray-100">
          <button onClick={prevMonth} className="w-10 h-10 bg-gray-100 rounded-full font-black text-gray-600 hover:bg-gray-200">{"<"}</button>
          <h2 className="text-lg font-black text-gray-800 tracking-wide uppercase">{monthName}</h2>
          <button onClick={nextMonth} className="w-10 h-10 bg-gray-100 rounded-full font-black text-gray-600 hover:bg-gray-200">{">"}</button>
        </div>

        {/* TOMBOL EKSPOR (Tampil jika data sudah dimuat) */}
        {!isLoading && report && (
          <button 
            onClick={handleExportPDF}
            className="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-bold rounded-xl shadow-md transition-all active:scale-95 flex justify-center items-center gap-2"
          >
            <span className="text-xl">📄</span> Unduh Laporan PDF
          </button>
        )}
      </div>

      {isLoading ? (
        <p className="text-center text-gray-400 font-bold py-10">Memuat data laporan...</p>
      ) : report ? (
        <>
          {/* Ringkasan Saldo & Perbandingan (Tidak Berubah) */}
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

            <div className={`absolute top-6 right-6 px-3 py-1 rounded-full text-xs font-black
              ${report.vs_last_month_percent > 0 ? 'bg-red-500 text-white' : 'bg-green-500 text-white'}`}>
              {report.vs_last_month_percent > 0 ? '↑' : '↓'} {Math.abs(report.vs_last_month_percent).toFixed(1)}% vs Lalu
            </div>
          </div>

          {/* Rincian Kategori Pengeluaran (Tidak Berubah) */}
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
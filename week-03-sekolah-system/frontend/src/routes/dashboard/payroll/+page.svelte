<script lang="ts">
    let teacherId = '';
    let month = new Date().getMonth() + 1; // Bulan saat ini
    let year = new Date().getFullYear();
    let message = '';
    let isError = false;
    let payslipData: any = null;

    // Dummy data untuk dropdown guru
    let teachers = [
        { id: 'uuid-guru-1', name: 'Bapak Budi (Matematika)' },
        { id: 'uuid-guru-2', name: 'Ibu Siti (Fisika)' }
    ];

    async function generatePayslip() {
        if (!teacherId) {
            isError = true;
            message = "Pilih guru terlebih dahulu!";
            return;
        }

        const token = localStorage.getItem('token');
        const res = await fetch('http://localhost:3000/api/v1/salary-payslips/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ teacher_id: teacherId, month: Number(month), year: Number(year) })
        });

        const data = await res.json();
        
        if (!res.ok) {
            isError = true;
            message = data.error;
            payslipData = null;
        } else {
            isError = false;
            message = data.message;
            payslipData = data.data; // Tampilkan rincian slip
        }
    }
</script>

<div class="p-8">
    <h2 class="text-3xl font-bold text-gray-800 mb-6">Sistem Penggajian Guru</h2>

    {#if message}
        <div class="p-4 mb-4 rounded {isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}">
            {message}
        </div>
    {/if}

    <div class="bg-white p-6 rounded shadow-sm border border-gray-200 mb-8 w-full max-w-2xl">
        <h3 class="text-lg font-semibold mb-4 text-gray-700">Kalkulasi Gaji Bulanan</h3>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div class="md:col-span-2">
                <label for="teacher-select" class="block text-sm font-medium text-gray-700 mb-1">Pilih Guru</label>
                <select id="teacher-select" bind:value={teacherId} class="w-full border-gray-300 rounded p-2 border">
                    <option value="">-- Pilih Guru --</option>
                    {#each teachers as t}
                        <option value={t.id}>{t.name}</option>
                    {/each}
                </select>
            </div>
            <div>
                <label for="period-month" class="block text-sm font-medium text-gray-700 mb-1">Periode</label>
                <div class="flex gap-2">
                    <input id="period-month" type="number" bind:value={month} min="1" max="12" class="w-16 border border-gray-300 rounded p-2" />
                    <input type="number" bind:value={year} class="w-24 border border-gray-300 rounded p-2" />
                </div>
            </div>
        </div>
        
        <button on:click={generatePayslip} class="bg-indigo-600 text-white px-6 py-2 rounded shadow hover:bg-indigo-700">
            Hitung & Generate Slip
        </button>
    </div>

    {#if payslipData}
        <div class="bg-white p-6 rounded shadow-lg border border-indigo-200 max-w-2xl">
            <div class="text-center mb-6 border-b pb-4">
                <h3 class="text-2xl font-bold text-gray-800">Slip Gaji Digital</h3>
                <p class="text-gray-500">Periode: {payslipData.month} / {payslipData.year}</p>
            </div>
            
            <table class="w-full text-left mb-6">
                <tbody>
                    {#each Object.entries(payslipData.details) as [komponen, nilai]}
                        <tr class="border-b">
                            <td class="py-2 text-gray-700">{komponen}</td>
                            <td class="py-2 text-right font-mono {Number(nilai) < 0 ? 'text-red-500' : 'text-gray-900'}">
                                Rp {Math.abs(Number(nilai)).toLocaleString('id-ID')}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            <div class="bg-gray-50 p-4 rounded flex justify-between items-center border border-gray-200">
                <span class="font-bold text-gray-700 text-lg">Total Diterima (Net)</span>
                <span class="font-black text-indigo-600 text-2xl font-mono">
                    Rp {Number(payslipData.net_salary).toLocaleString('id-ID')}
                </span>
            </div>
        </div>
    {/if}
</div>
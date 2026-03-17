<script lang="ts">
    import { onMount } from 'svelte';

    let reports: any[] = [];
    let students: any[] = [];
    let academicYears: any[] = [];
    
    let selectedStudent = '';
    let selectedSemester = 1;
    let selectedAcademicYear = '';
    
    let isLoading = false;
    let message = '';
    let isError = false;

    onMount(async () => {
        await fetchAllData();
    });

    async function fetchAllData() {
        const token = localStorage.getItem('token');
        const headers = { 'Authorization': `Bearer ${token}` };

        try {
            const [resReports, resStudents, resAY] = await Promise.all([
                fetch('http://localhost:3000/api/v1/report-cards', { headers }),
                fetch('http://localhost:3000/api/v1/students', { headers }),
                fetch('http://localhost:3000/api/v1/academic-years', { headers })
            ]);

            if (resReports.ok) { const data = await resReports.json(); reports = data.data || []; }
            if (resStudents.ok) { const data = await resStudents.json(); students = data.data || []; }
            if (resAY.ok) { const data = await resAY.json(); academicYears = data.data || []; }
        } catch (err) { console.error(err); }
    }

    async function generateReport() {
        if (!selectedStudent || !selectedAcademicYear) {
            isError = true; message = "Siswa dan Tahun Ajaran wajib dipilih!"; return;
        }

        isLoading = true; message = ''; isError = false;
        const token = localStorage.getItem('token');

        try {
            const res = await fetch('http://localhost:3000/api/v1/report-cards/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
                body: JSON.stringify({
                    student_id: selectedStudent,
                    semester: selectedSemester,
                    academic_year_id: selectedAcademicYear
                })
            });

            const data = await res.json();
            if (!res.ok) throw new Error(data.error);

            isError = false;
            // Pesan bahwa ini berjalan di background
            message = data.message + " Silakan refresh tabel dalam beberapa detik.";
            
            // Auto refresh setelah 3 detik untuk mengambil hasil PDF yang selesai digenerate
            setTimeout(fetchAllData, 3000);
        } catch (err: any) {
            isError = true; message = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="p-8 max-w-6xl mx-auto">
    <h2 class="text-3xl font-bold text-gray-800 mb-8">Pusat Cetak Rapor Akhir</h2>

    {#if message}
        <div class="p-4 mb-6 rounded-lg shadow-sm border-l-4 {isError ? 'border-red-500 bg-red-50 text-red-700' : 'border-blue-500 bg-blue-50 text-blue-700'}">
            {message}
        </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-1">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
                <h3 class="text-lg font-bold text-gray-800 mb-4">Minta Pembuatan Rapor</h3>
                <div class="space-y-4">
                    <div>
                        <label for="report-student" class="block text-sm font-semibold text-gray-700 mb-1">Pilih Siswa</label>
                        <select id="report-student" bind:value={selectedStudent} class="w-full border p-2 rounded bg-gray-50">
                            <option value="">-- Pilih Siswa --</option>
                            {#each students as s}
                                <option value={s.id}>{s.user?.full_name} ({s.nis})</option>
                            {/each}
                        </select>
                    </div>
                    <div>
                        <label for="report-semester" class="block text-sm font-semibold text-gray-700 mb-1">Semester</label>
                        <select id="report-semester" bind:value={selectedSemester} class="w-full border p-2 rounded bg-gray-50">
                            <option value={1}>1 (Ganjil)</option>
                            <option value={2}>2 (Genap)</option>
                        </select>
                    </div>
                    <div>
                        <label for="report-academic-year" class="block text-sm font-semibold text-gray-700 mb-1">Tahun Ajaran</label>
                        <select id="report-academic-year" bind:value={selectedAcademicYear} class="w-full border p-2 rounded bg-gray-50">
                            <option value="">-- Pilih Tahun --</option>
                            {#each academicYears as ay}
                                <option value={ay.id}>{ay.name}</option>
                            {/each}
                        </select>
                    </div>
                    <button on:click={generateReport} disabled={isLoading} class="w-full bg-blue-600 text-white font-bold p-2.5 rounded shadow hover:bg-blue-700 disabled:opacity-50">
                        {isLoading ? 'Memproses...' : 'Generate Rapor (PDF)'}
                    </button>
                </div>
            </div>
        </div>

        <div class="lg:col-span-2 bg-white rounded-xl shadow-sm border border-gray-100 p-6">
            <h3 class="font-bold text-gray-800 mb-4">Arsip Rapor Siswa</h3>
            <table class="w-full text-left">
                <thead class="bg-gray-50 text-gray-500 text-xs uppercase">
                    <tr>
                        <th class="p-3">Nama Siswa</th>
                        <th class="p-3">Semester</th>
                        <th class="p-3">Nilai Akhir</th>
                        <th class="p-3">File PDF</th>
                    </tr>
                </thead>
                <tbody class="divide-y">
                    {#each reports as r}
                        <tr>
                            <td class="p-3 font-semibold text-gray-800">{r.student?.user?.full_name}</td>
                            <td class="p-3">{r.semester}</td>
                            <td class="p-3 font-mono text-indigo-600 font-bold">{r.final_score}</td>
                            <td class="p-3">
                                {#if r.pdf_url}
                                    <a href={r.pdf_url} target="_blank" class="bg-red-500 hover:bg-red-600 text-white text-xs px-3 py-1.5 rounded font-bold transition flex items-center justify-center gap-1 w-max">
                                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M6 2a2 2 0 00-2 2v12a2 2 0 002 2h8a2 2 0 002-2V7.414A2 2 0 0015.414 6L12 2.586A2 2 0 0010.586 2H6zm5 6a1 1 0 10-2 0v3.586l-1.293-1.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 11.586V8z" clip-rule="evenodd"></path></svg>
                                        Download PDF
                                    </a>
                                {:else}
                                    <span class="text-gray-400 text-xs italic">Diproses...</span>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
</div>
<script lang="ts">
    let selectedClass = '';
    let semester = 1;
    let isLoading = false;
    let message = '';
    let isError = false;
    let rankedStudents: any[] = [];

    // Dummy data untuk form
    let classes = [
        { id: 'c1-uuid', name: 'XII IPA 1' },
        { id: 'c2-uuid', name: 'XII IPS 2' }
    ];

    async function calculateRank() {
        if (!selectedClass) {
            isError = true;
            message = "Pilih kelas terlebih dahulu!";
            return;
        }

        isLoading = true;
        message = '';
        const token = localStorage.getItem('token');
        
        // Di sistem nyata, academic_year_id didapat dari context tahun ajaran aktif
        const payload = {
            class_id: selectedClass,
            semester: semester,
            academic_year_id: "uuid-tahun-ajaran-aktif" 
        };

        const res = await fetch('http://localhost:3000/api/v1/student-recommendations/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        const data = await res.json();
        isLoading = false;

        if (!res.ok) {
            isError = true;
            message = data.error;
            rankedStudents = [];
        } else {
            isError = false;
            message = "Kalkulasi selesai! Berikut peringkat siswa terbaik.";
            rankedStudents = data.data; // Menampilkan hasil ranking
        }
    }
</script>

<div class="p-8">
    <h2 class="text-3xl font-bold text-gray-800 mb-6">Rekomendasi Siswa Berprestasi</h2>

    {#if message}
        <div class="p-4 mb-4 rounded {isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}">
            {message}
        </div>
    {/if}

    <div class="bg-white p-6 rounded shadow border border-gray-200 mb-8 flex gap-4 items-end">
        <div class="flex-1">
            <label for="class-select" class="block text-sm font-medium text-gray-700 mb-1">Pilih Kelas</label>
            <select id="class-select" bind:value={selectedClass} class="w-full border-gray-300 rounded p-2 border">
                <option value="">-- Pilih Kelas --</option>
                {#each classes as c}
                    <option value={c.id}>{c.name}</option>
                {/each}
            </select>
        </div>
        <div>
            <label for="semester-select" class="block text-sm font-medium text-gray-700 mb-1">Semester</label>
            <select id="semester-select" bind:value={semester} class="w-32 border-gray-300 rounded p-2 border">
                <option value={1}>Ganjil (1)</option>
                <option value={2}>Genap (2)</option>
            </select>
        </div>
        <button on:click={calculateRank} disabled={isLoading} class="bg-purple-600 text-white px-6 py-2 rounded shadow hover:bg-purple-700 disabled:bg-purple-300">
            {isLoading ? 'Menghitung...' : 'Mulai Kalkulasi'}
        </button>
    </div>

    {#if rankedStudents.length > 0}
        <div class="bg-white rounded shadow border border-gray-200 overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-purple-50">
                    <tr>
                        <th class="px-6 py-3 text-center text-xs font-bold text-purple-700 uppercase">Peringkat</th>
                        <th class="px-6 py-3 text-left text-xs font-bold text-purple-700 uppercase">Nama Siswa</th>
                        <th class="px-6 py-3 text-center text-xs font-bold text-purple-700 uppercase">Skor Akhir</th>
                        <th class="px-6 py-3 text-center text-xs font-bold text-purple-700 uppercase">Aksi Kepsek</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                    {#each rankedStudents as student, i}
                        <tr class={i < 3 ? 'bg-yellow-50' : 'bg-white'}>
                            <td class="px-6 py-4 text-center">
                                {#if i === 0} 🥇 
                                {:else if i === 1} 🥈 
                                {:else if i === 2} 🥉 
                                {:else} {student.rank} {/if}
                            </td>
                            <td class="px-6 py-4 font-medium text-gray-900">Siswa ID: {student.student_id}</td>
                            <td class="px-6 py-4 text-center font-mono font-bold text-purple-600">{student.score.toFixed(2)}</td>
                            <td class="px-6 py-4 text-center">
                                <button class="text-sm border border-purple-600 text-purple-600 px-3 py-1 rounded hover:bg-purple-600 hover:text-white transition">
                                    Pilih Kandidat
                                </button>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>
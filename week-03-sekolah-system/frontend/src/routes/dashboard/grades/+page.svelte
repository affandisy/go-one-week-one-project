<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';

    // State untuk form pilihan
    let selectedClass = '';
    let selectedSubject = '';
    let selectedComponent = '';
    
    // State untuk data dummy (Nantinya diganti dengan fetch API sungguhan)
    let classes = [
        { id: 'c1-uuid', name: 'XII IPA 1' },
        { id: 'c2-uuid', name: 'XII IPS 2' }
    ];
    let subjects = [
        { id: 's1-uuid', name: 'Matematika Lanjut' },
        { id: 's2-uuid', name: 'Fisika' }
    ];
    let components = [
        { id: 'comp1-uuid', name: 'Ulangan Harian 1', weight: 0.2 },
        { id: 'comp2-uuid', name: 'Ujian Tengah Semester', weight: 0.3 }
    ];

    // State untuk daftar siswa dan nilainya
    let students = [
        { id: 'stud1-uuid', nis: '10101', name: 'Budi Santoso', score: 0, notes: '' },
        { id: 'stud2-uuid', nis: '10102', name: 'Siti Aminah', score: 0, notes: '' },
        { id: 'stud3-uuid', nis: '10103', name: 'Agus Pratama', score: 0, notes: '' }
    ];

    let isLoading = false;
    let message = '';
    let isError = false;

    onMount(() => {
        const token = localStorage.getItem('token');
        if (!token) goto('/login');
        
        // Di aplikasi nyata, Anda akan fetch() daftar kelas, mapel, dan komponen di sini
    });

    // Fungsi untuk mengirim nilai ke Golang Backend
    async function saveGrades() {
        if (!selectedClass || !selectedSubject || !selectedComponent) {
            isError = true;
            message = "Harap pilih Kelas, Mata Pelajaran, dan Komponen Nilai!";
            return;
        }

        isLoading = true;
        message = '';
        isError = false;

        const token = localStorage.getItem('token');

        // Mengirim nilai satu per satu atau batch (Tergantung desain API Anda)
        // Di sini kita contohkan iterasi untuk mengirim semua nilai siswa
        try {
            for (const student of students) {
                const payload = {
                    class_id: selectedClass,
                    student_id: student.id,
                    subject_id: selectedSubject,
                    component_id: selectedComponent,
                    score: Number(student.score),
                    notes: student.notes
                };

                const res = await fetch('http://localhost:3000/api/v1/grades', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}` // JWT Auth Middleware!
                    },
                    body: JSON.stringify(payload)
                });

                if (!res.ok) {
                    const data = await res.json();
                    throw new Error(data.error || 'Gagal menyimpan nilai');
                }
            }
            
            isError = false;
            message = "Semua nilai berhasil disimpan!";
        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="p-8">
    <h2 class="text-3xl font-semibold text-gray-800 mb-6">Input Nilai Siswa</h2>

    {#if message}
        <div class="p-4 mb-4 rounded {isError ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}">
            {message}
        </div>
    {/if}

    <div class="bg-white p-6 rounded shadow-sm border border-gray-200 mb-8 grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
            <label for="class-select" class="block text-sm font-medium text-gray-700 mb-1">Pilih Kelas</label>
            <select id="class-select" bind:value={selectedClass} class="w-full border-gray-300 rounded-md shadow-sm p-2 border">
                <option value="">-- Pilih Kelas --</option>
                {#each classes as c}
                    <option value={c.id}>{c.name}</option>
                {/each}
            </select>
        </div>
        <div>
            <label for="subject-select" class="block text-sm font-medium text-gray-700 mb-1">Mata Pelajaran</label>
            <select id="subject-select" bind:value={selectedSubject} class="w-full border-gray-300 rounded-md shadow-sm p-2 border">
                <option value="">-- Pilih Mapel --</option>
                {#each subjects as s}
                    <option value={s.id}>{s.name}</option>
                {/each}
            </select>
        </div>
        <div>
            <label for="component-select" class="block text-sm font-medium text-gray-700 mb-1">Komponen Nilai</label>
            <select id="component-select" bind:value={selectedComponent} class="w-full border-gray-300 rounded-md shadow-sm p-2 border">
                <option value="">-- Pilih Komponen --</option>
                {#each components as comp}
                    <option value={comp.id}>{comp.name} (Bobot: {comp.weight * 100}%)</option>
                {/each}
            </select>
        </div>
    </div>

    {#if selectedClass && selectedSubject && selectedComponent}
        <div class="bg-white rounded shadow-sm border border-gray-200 overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">NIS</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nama Siswa</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nilai (0-100)</th>
                        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Catatan Guru</th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    {#each students as student}
                        <tr>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{student.nis}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{student.name}</td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <input type="number" min="0" max="100" bind:value={student.score} 
                                    class="border border-gray-300 rounded p-1 w-24 text-center focus:ring-blue-500 focus:border-blue-500" />
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <input type="text" placeholder="Cth: Sangat Aktif" bind:value={student.notes} 
                                    class="border border-gray-300 rounded p-1 w-full text-sm" />
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
            <div class="p-4 bg-gray-50 flex justify-end">
                <button on:click={saveGrades} disabled={isLoading} 
                    class="bg-blue-600 text-white px-6 py-2 rounded shadow hover:bg-blue-700 disabled:bg-blue-300 transition">
                    {isLoading ? 'Menyimpan...' : 'Simpan Semua Nilai'}
                </button>
            </div>
        </div>
    {:else}
        <div class="text-center py-12 bg-white rounded shadow-sm border border-gray-200 text-gray-500">
            Silakan pilih Kelas, Mata Pelajaran, dan Komponen Nilai terlebih dahulu untuk menampilkan daftar siswa.
        </div>
    {/if}
</div>
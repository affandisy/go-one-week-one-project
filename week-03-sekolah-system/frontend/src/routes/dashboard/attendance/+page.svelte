<script lang="ts">

    let selectedClass = '';
    let attendanceDate = new Date().toISOString().split('T')[0]; // Format YYYY-MM-DD
    let students: any[] = [];
    let isLoading = false;
    let message = '';
    let isError = false;

    // Dummy data kelas (idealnya ambil dari API /classes)
    let classes = [
        { id: 'c1-uuid', name: 'XII IPA 1' },
        { id: 'c2-uuid', name: 'XII IPS 2' }
    ];

    // Dummy data murid jika API belum berisi data seeder
    let dummyStudents = [
        { user_id: 'u1-uuid', id: 's1-uuid', nis: '10101', user: { full_name: 'Budi Santoso' }, status: 'hadir', notes: '' },
        { user_id: 'u2-uuid', id: 's2-uuid', nis: '10102', user: { full_name: 'Siti Aminah' }, status: 'hadir', notes: '' },
        { user_id: 'u3-uuid', id: 's3-uuid', nis: '10103', user: { full_name: 'Agus Pratama' }, status: 'hadir', notes: '' }
    ];

    // Fungsi fetch murid (Untuk MVP kita pakai dummy jika API kosong)
    async function fetchStudents() {
        if (!selectedClass) {
            students = [];
            return;
        }
        
        // Di sistem nyata:
        // const token = localStorage.getItem('token');
        // const res = await fetch(`http://localhost:3000/api/v1/attendance/class/${selectedClass}/students`, { ... });
        // const data = await res.json();
        // students = data.data.map(s => ({ ...s, status: 'hadir', notes: '' }));
        
        // Memakai dummy data yang sudah disiapkan formatnya:
        students = dummyStudents.map(s => ({ ...s }));
    }

    async function submitAttendance() {
        if (!selectedClass || students.length === 0) {
            isError = true;
            message = "Pilih kelas terlebih dahulu!";
            return;
        }

        isLoading = true;
        message = '';
        isError = false;
        const token = localStorage.getItem('token');

        const payload = {
            class_id: selectedClass,
            date: attendanceDate,
            records: students.map(s => ({
                student_id: s.id,
                user_id: s.user_id,
                status: s.status,
                notes: s.notes
            }))
        };

        try {
            const res = await fetch('http://localhost:3000/api/v1/attendance/student/batch', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });

            const data = await res.json();
            if (!res.ok) throw new Error(data.error);

            isError = false;
            message = "Absensi kelas berhasil disimpan!";
        } catch (err: any) {
            isError = true;
            message = err.message;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="p-8 max-w-6xl mx-auto">
    <h2 class="text-3xl font-bold text-gray-800 mb-2">Pencatatan Absensi Siswa</h2>
    <p class="text-gray-500 mb-8">Catat kehadiran siswa secara massal untuk satu hari penuh.</p>

    {#if message}
        <div class="p-4 mb-6 rounded shadow-sm border-l-4 {isError ? 'bg-red-50 text-red-700 border-red-500' : 'bg-green-50 text-green-700 border-green-500'}">
            {message}
        </div>
    {/if}

    <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 mb-8 flex flex-wrap gap-6 items-end">
        <div class="flex-1 min-w-[200px]">
            <label for="class-select" class="block text-sm font-semibold text-gray-700 mb-2">Pilih Kelas</label>
            <select id="class-select" bind:value={selectedClass} on:change={fetchStudents} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50">
                <option value="">-- Pilih Kelas --</option>
                {#each classes as c}
                    <option value={c.id}>{c.name}</option>
                {/each}
            </select>
        </div>
        <div class="flex-1 min-w-[200px]">
            <label for="attendance-date" class="block text-sm font-semibold text-gray-700 mb-2">Tanggal Absensi</label>
            <input id="attendance-date" type="date" bind:value={attendanceDate} class="w-full border border-gray-300 rounded-lg p-2.5 focus:ring-2 focus:ring-blue-500 bg-gray-50" />
        </div>
        <div>
            <button on:click={submitAttendance} disabled={isLoading || students.length === 0} 
                class="bg-blue-600 text-white font-semibold px-6 py-2.5 rounded-lg shadow-md hover:bg-blue-700 disabled:opacity-50 transition">
                {isLoading ? 'Menyimpan...' : 'Simpan Absensi'}
            </button>
        </div>
    </div>

    {#if students.length > 0}
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-6 py-4 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">NIS</th>
                        <th class="px-6 py-4 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Nama Siswa</th>
                        <th class="px-6 py-4 text-center text-xs font-bold text-gray-500 uppercase tracking-wider">Status Kehadiran</th>
                        <th class="px-6 py-4 text-left text-xs font-bold text-gray-500 uppercase tracking-wider">Keterangan / Surat</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-100">
                    {#each students as student}
                        <tr class="hover:bg-gray-50 transition">
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 font-mono">{student.nis}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm font-bold text-gray-800">{student.user.full_name}</td>
                            <td class="px-6 py-4 whitespace-nowrap text-center">
                                <div class="flex justify-center gap-3">
                                    <label class="flex flex-col items-center cursor-pointer group">
                                        <input type="radio" bind:group={student.status} value="hadir" class="hidden peer" />
                                        <span class="w-8 h-8 flex items-center justify-center rounded-full border-2 border-gray-200 peer-checked:bg-green-500 peer-checked:border-green-500 peer-checked:text-white text-gray-400 font-bold transition">H</span>
                                        <span class="text-[10px] mt-1 text-gray-500 peer-checked:text-green-600 peer-checked:font-bold">Hadir</span>
                                    </label>
                                    <label class="flex flex-col items-center cursor-pointer group">
                                        <input type="radio" bind:group={student.status} value="izin" class="hidden peer" />
                                        <span class="w-8 h-8 flex items-center justify-center rounded-full border-2 border-gray-200 peer-checked:bg-blue-500 peer-checked:border-blue-500 peer-checked:text-white text-gray-400 font-bold transition">I</span>
                                        <span class="text-[10px] mt-1 text-gray-500 peer-checked:text-blue-600 peer-checked:font-bold">Izin</span>
                                    </label>
                                    <label class="flex flex-col items-center cursor-pointer group">
                                        <input type="radio" bind:group={student.status} value="sakit" class="hidden peer" />
                                        <span class="w-8 h-8 flex items-center justify-center rounded-full border-2 border-gray-200 peer-checked:bg-yellow-500 peer-checked:border-yellow-500 peer-checked:text-white text-gray-400 font-bold transition">S</span>
                                        <span class="text-[10px] mt-1 text-gray-500 peer-checked:text-yellow-600 peer-checked:font-bold">Sakit</span>
                                    </label>
                                    <label class="flex flex-col items-center cursor-pointer group">
                                        <input type="radio" bind:group={student.status} value="alfa" class="hidden peer" />
                                        <span class="w-8 h-8 flex items-center justify-center rounded-full border-2 border-gray-200 peer-checked:bg-red-500 peer-checked:border-red-500 peer-checked:text-white text-gray-400 font-bold transition">A</span>
                                        <span class="text-[10px] mt-1 text-gray-500 peer-checked:text-red-600 peer-checked:font-bold">Alfa</span>
                                    </label>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <input type="text" bind:value={student.notes} placeholder="Catatan opsional..." 
                                    class="w-full border border-gray-200 rounded p-2 text-sm focus:ring-2 focus:ring-blue-500 bg-gray-50 focus:bg-white transition" />
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {:else if selectedClass}
        <div class="text-center py-12 bg-white rounded-xl shadow-sm border border-gray-100 text-gray-500">
            Tidak ada data siswa di kelas ini.
        </div>
    {/if}
</div>
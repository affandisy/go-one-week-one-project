<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // --- STATE MODUL ---
    let modules = $state<any[]>([]);
    let selectedModuleId = $state('');
    let isModulesLoading = $state(true);

    let newModule = $state({ title: '', description: '', level_order: 1 });

    // --- STATE MATERI ---
    let materials = $state<any[]>([]);
    let isMaterialsLoading = $state(false);

    let newMaterial = $state({
        content_type: 'learn_card', // 'learn_card', 'quiz_mcq', 'quiz_unscramble'
        question: '',
        correct_answer: '',
        optionsRaw: '', // Input bantuan (dipisah koma) untuk dikonversi ke JSON Array
        image_url: '',
        audio_url: '',
        display_order: 1
    });

    // File Upload State
    let imageFiles = $state<FileList | null>(null);
    let audioFiles = $state<FileList | null>(null);
    let isUploading = $state(false);
    let feedbackMsg = $state('');

    onMount(async () => {
        await loadModules();
    });

    async function loadModules() {
        try {
            const res = await apiFetch('/modules');
            modules = res.data || [];
            if (modules.length > 0 && !selectedModuleId) {
                selectModule(modules[0].id);
            }
        } catch (err: any) {
            feedbackMsg = 'Gagal memuat modul: ' + err.message;
        } finally {
            isModulesLoading = false;
        }
    }

    async function selectModule(id: string) {
        selectedModuleId = id;
        isMaterialsLoading = true;
        try {
            const res = await apiFetch(`/modules/${id}/materials`);
            materials = res.data || [];
            // Otomatis set urutan (display_order) selanjutnya
            newMaterial.display_order = materials.length + 1;
        } catch (err: any) {
            feedbackMsg = 'Gagal memuat materi: ' + err.message;
        } finally {
            isMaterialsLoading = false;
        }
    }

    // --- FUNGSI UPLOAD FILE ---
    // Kita gunakan fetch native karena apiFetch kita di-setting untuk application/json
    async function uploadFile(file: File): Promise<string> {
        const formData = new FormData();
        formData.append('file', file);

        const token = localStorage.getItem('token');
        const res = await fetch('http://localhost:3000/api/v1/upload', {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` },
            body: formData
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Gagal upload file');
        return data.file_url;
    }

    // --- FUNGSI CREATE ---
    async function handleCreateModule(e: Event) {
        e.preventDefault();
        try {
            // Karena rute POST /modules/seed saat ini tidak terproteksi di backend,
            // kita bisa menembaknya langsung (di masa depan, ganti ke rute khusus admin)
            await apiFetch('/modules/seed', {
                method: 'POST',
                data: newModule
            });
            feedbackMsg = 'Modul berhasil dibuat!';
            newModule = { title: '', description: '', level_order: modules.length + 2 };
            await loadModules();
        } catch (err: any) {
            feedbackMsg = err.message;
        }
    }

    async function handleCreateMaterial(e: Event) {
        e.preventDefault();
        isUploading = true;
        feedbackMsg = 'Memproses data dan mengunggah file...';

        try {
            // 1. Upload Gambar jika ada
            if (imageFiles && imageFiles.length > 0) {
                newMaterial.image_url = await uploadFile(imageFiles[0]);
            }
            // 2. Upload Audio jika ada
            if (audioFiles && audioFiles.length > 0) {
                newMaterial.audio_url = await uploadFile(audioFiles[0]);
            }

            // 3. Parsing Options (Pilihan Ganda/Unscramble)
            let optionsJSON = "";
            if (newMaterial.content_type !== 'learn_card' && newMaterial.optionsRaw) {
                // Memecah teks "A, B, C" menjadi array ["A", "B", "C"]
                const optArray = newMaterial.optionsRaw.split(',').map(s => s.trim());
                optionsJSON = JSON.stringify(optArray);
            }

            // 4. Kirim Data Materi
            await apiFetch(`/modules/${selectedModuleId}/materials`, {
                method: 'POST',
                data: {
                    module_id: selectedModuleId,
                    content_type: newMaterial.content_type,
                    question: newMaterial.question,
                    correct_answer: newMaterial.correct_answer,
                    options: optionsJSON,
                    image_url: newMaterial.image_url,
                    audio_url: newMaterial.audio_url,
                    display_order: newMaterial.display_order
                }
            });

            feedbackMsg = 'Materi berhasil ditambahkan!';
            
            // Reset form sebagian, biarkan tipe materi sama untuk memudahkan input beruntun
            newMaterial.question = '';
            newMaterial.correct_answer = '';
            newMaterial.optionsRaw = '';
            newMaterial.image_url = '';
            newMaterial.audio_url = '';
            imageFiles = null;
            audioFiles = null;
            
            await selectModule(selectedModuleId); // Refresh daftar
        } catch (err: any) {
            feedbackMsg = 'Gagal menambahkan materi: ' + err.message;
        } finally {
            isUploading = false;
        }
    }
</script>

<div class="min-h-screen bg-slate-50 p-6">
    <div class="max-w-7xl mx-auto">
        <div class="mb-8 flex justify-between items-center">
            <div>
                <h1 class="text-3xl font-black text-slate-800">LingoLearn Admin</h1>
                <p class="text-slate-500 font-medium">Manajemen Modul & Konten Pembelajaran</p>
            </div>
            <a href="/dashboard" class="text-indigo-600 font-bold hover:underline">Kembali ke App</a>
        </div>

        {#if feedbackMsg}
            <div class="p-4 mb-6 rounded-xl font-bold bg-indigo-50 text-indigo-700 border border-indigo-100">
                {feedbackMsg}
                <button onclick={() => feedbackMsg = ''} class="float-right text-indigo-400">✕</button>
            </div>
        {/if}

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
            
            <div class="lg:col-span-4 space-y-6">
                <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                    <h2 class="text-xl font-black text-slate-800 mb-4">Buat Level Baru</h2>
                    <form onsubmit={handleCreateModule} class="space-y-4">
                        <div>
                            <label for="module-title" class="block text-sm font-bold text-slate-700">Judul Level</label>
                            <input id="module-title" type="text" bind:value={newModule.title} required class="mt-1 w-full p-2 border rounded-lg bg-slate-50">
                        </div>
                        <div>
                            <label for="module-description" class="block text-sm font-bold text-slate-700">Deskripsi</label>
                            <input id="module-description" type="text" bind:value={newModule.description} class="mt-1 w-full p-2 border rounded-lg bg-slate-50">
                        </div>
                        <div>
                            <label for="module-level-order" class="block text-sm font-bold text-slate-700">Urutan Level</label>
                            <input id="module-level-order" type="number" bind:value={newModule.level_order} required class="mt-1 w-full p-2 border rounded-lg bg-slate-50">
                        </div>
                        <button type="submit" class="w-full py-2 bg-slate-800 text-white font-bold rounded-lg hover:bg-slate-900">
                            Tambah Level
                        </button>
                    </form>
                </div>

                <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                    <h2 class="text-xl font-black text-slate-800 mb-4">Daftar Level</h2>
                    {#if isModulesLoading}
                        <p class="text-slate-500">Memuat data...</p>
                    {:else}
                        <div class="space-y-2">
                            {#each modules as mod}
                                <button 
                                    onclick={() => selectModule(mod.id)}
                                    class="w-full text-left p-3 rounded-lg border-2 font-bold transition-all
                                    {selectedModuleId === mod.id ? 'border-indigo-500 bg-indigo-50 text-indigo-700' : 'border-slate-100 hover:border-slate-300'}">
                                    {mod.level_order}. {mod.title}
                                </button>
                            {/each}
                        </div>
                    {/if}
                </div>
            </div>

            <div class="lg:col-span-8 space-y-6">
                {#if selectedModuleId}
                    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                        <h2 class="text-xl font-black text-slate-800 mb-6 border-b pb-4">
                            Tambah Materi ke Level yang Dipilih
                        </h2>

                        <form onsubmit={handleCreateMaterial} class="space-y-6">
                            <div>
                                <p class="block text-sm font-bold text-slate-700 mb-2">Tipe Materi</p>
                                <div class="flex gap-4">
                                    <label class="flex items-center gap-2 cursor-pointer">
                                        <input type="radio" bind:group={newMaterial.content_type} value="learn_card" class="text-indigo-600 focus:ring-indigo-500">
                                        <span class="font-medium text-slate-700">Kartu Belajar (Flashcard)</span>
                                    </label>
                                    <label class="flex items-center gap-2 cursor-pointer">
                                        <input type="radio" bind:group={newMaterial.content_type} value="quiz_mcq" class="text-indigo-600 focus:ring-indigo-500">
                                        <span class="font-medium text-slate-700">Kuis (Pilihan Ganda)</span>
                                    </label>
                                    <label class="flex items-center gap-2 cursor-pointer">
                                        <input type="radio" bind:group={newMaterial.content_type} value="quiz_unscramble" class="text-indigo-600 focus:ring-indigo-500">
                                        <span class="font-medium text-slate-700">Kuis (Susun Kata)</span>
                                    </label>
                                </div>
                            </div>

                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div>
                                    <label for="material-question" class="block text-sm font-bold text-slate-700">
                                        {newMaterial.content_type === 'learn_card' ? 'Kata Asing (Pertanyaan)' : 'Teks Pertanyaan Kuis'}
                                    </label>
                                    <input id="material-question" type="text" bind:value={newMaterial.question} required class="mt-1 w-full p-2 border rounded-lg bg-slate-50" 
                                        placeholder={newMaterial.content_type === 'learn_card' ? 'Contoh: Apple' : 'Contoh: Apa bahasa inggrisnya Apel?'}>
                                </div>

                                <div>
                                    <label for="material-correct-answer" class="block text-sm font-bold text-slate-700">
                                        {newMaterial.content_type === 'learn_card' ? 'Arti / Terjemahan' : 'Kunci Jawaban Tepat'}
                                    </label>
                                    <input id="material-correct-answer" type="text" bind:value={newMaterial.correct_answer} required class="mt-1 w-full p-2 border rounded-lg bg-slate-50"
                                        placeholder={newMaterial.content_type === 'quiz_unscramble' ? 'Contoh: I eat an apple' : 'Jawaban benar'}>
                                </div>
                            </div>

                            {#if newMaterial.content_type !== 'learn_card'}
                                <div>
                                    <label for="material-options" class="block text-sm font-bold text-slate-700">Pilihan Opsi (Pisahkan dengan koma)</label>
                                    <input id="material-options" type="text" bind:value={newMaterial.optionsRaw} required class="mt-1 w-full p-2 border rounded-lg bg-slate-50"
                                        placeholder={newMaterial.content_type === 'quiz_mcq' ? 'A, B, C, D' : 'I, eat, an, apple, orange, drink'}>
                                    <p class="text-xs text-slate-500 mt-1">Pastikan kunci jawaban persis ada di dalam opsi ini.</p>
                                </div>
                            {/if}

                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 bg-slate-50 border border-slate-200 rounded-xl">
                                <div>
                                    <label for="material-image" class="block text-sm font-bold text-slate-700">Gambar (Opsional)</label>
                                    <input id="material-image" type="file" accept="image/*" bind:files={imageFiles} class="mt-1 w-full text-sm">
                                </div>
                                <div>
                                    <label for="material-audio" class="block text-sm font-bold text-slate-700">Audio (Opsional)</label>
                                    <input id="material-audio" type="file" accept="audio/*" bind:files={audioFiles} class="mt-1 w-full text-sm">
                                    <p class="text-xs text-slate-500 mt-1">Jika kosong, sistem akan menggunakan fitur Text-to-Speech browser.</p>
                                </div>
                            </div>

                            <button type="submit" disabled={isUploading} class="w-full py-3 bg-indigo-600 text-white font-bold rounded-xl shadow hover:bg-indigo-700 disabled:opacity-70 transition-all">
                                {isUploading ? 'Menyimpan & Mengunggah File...' : 'Simpan Materi Ini'}
                            </button>
                        </form>
                    </div>

                    <div class="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                        <h2 class="text-xl font-black text-slate-800 mb-4">Isi Level Saat Ini</h2>
                        {#if isMaterialsLoading}
                            <p class="text-slate-500">Memuat materi...</p>
                        {:else if materials.length === 0}
                            <p class="text-slate-500 italic">Level ini masih kosong.</p>
                        {:else}
                            <div class="overflow-x-auto">
                                <table class="w-full text-left border-collapse">
                                    <thead>
                                        <tr class="bg-slate-50 text-sm font-bold text-slate-600 border-b">
                                            <th class="p-3">#</th>
                                            <th class="p-3">Tipe</th>
                                            <th class="p-3">Pertanyaan</th>
                                            <th class="p-3">Kunci Jawaban</th>
                                            <th class="p-3">Media</th>
                                        </tr>
                                    </thead>
                                    <tbody class="text-sm">
                                        {#each materials as mat, i}
                                            <tr class="border-b border-slate-100 hover:bg-slate-50">
                                                <td class="p-3 font-medium text-slate-500">{i + 1}</td>
                                                <td class="p-3">
                                                    <span class="px-2 py-1 rounded text-xs font-bold
                                                        {mat.content_type === 'learn_card' ? 'bg-blue-100 text-blue-700' : 'bg-orange-100 text-orange-700'}">
                                                        {mat.content_type.replace('quiz_', '')}
                                                    </span>
                                                </td>
                                                <td class="p-3 font-bold text-slate-800 truncate max-w-xs">{mat.question}</td>
                                                <td class="p-3 text-slate-600">{mat.correct_answer}</td>
                                                <td class="p-3 flex gap-2">
                                                    {#if mat.image_url} <span title="Ada Gambar">🖼️</span> {/if}
                                                    {#if mat.audio_url} <span title="Ada Audio">🔊</span> {/if}
                                                </td>
                                            </tr>
                                        {/each}
                                    </tbody>
                                </table>
                            </div>
                        {/if}
                    </div>
                {:else}
                    <div class="flex items-center justify-center h-64 bg-slate-100 rounded-2xl border-2 border-dashed border-slate-300">
                        <p class="text-slate-400 font-bold">Pilih level di sebelah kiri untuk mengelola materinya.</p>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>
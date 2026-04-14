<script lang="ts">
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    let moduleId = $page.params.moduleId;
    
    let questions = $state<any[]>([]);
    let currentIndex = $state(0);
    let currentSelection = $state(''); 
    let collectedAnswers = $state<{material_id: string, selected: string}[]>([]);
    
    let availableWords = $state<string[]>([]);
    let selectedWords = $state<string[]>([]);
    
    let isLoading = $state(true);
    let isSubmitting = $state(false);
    let errorMessage = $state('');
    let result = $state<any>(null);

    onMount(async () => {
        try {
            const res = await apiFetch(`/modules/${moduleId}/quiz`);
            questions = res.data.questions || [];
            if (questions.length === 0) errorMessage = "Kuis belum tersedia untuk level ini.";
        } catch (err: any) {
            errorMessage = err.message || "Gagal memuat soal kuis.";
        } finally {
            isLoading = false;
        }
    });

    function parseOptions(optionsStr: string): string[] {
        try {
            return JSON.parse(optionsStr);
        } catch {
            return [];
        }
    }

    $effect(() => {
        const currentQ = questions[currentIndex];
        if (currentQ && currentQ.content_type === 'quiz_unscramble') {
            // Memuat kata-kata acak dari database (dari kolom options)
            availableWords = parseOptions(currentQ.options);
            selectedWords = [];
            currentSelection = '';
        } else {
            // Reset jika kembali ke soal MCQ
            currentSelection = '';
        }
    });

    function selectOption(opt: string) {
        currentSelection = opt;
    }

    function selectUnscrambleWord(word: string, index: number) {
        // Pindahkan kata dari available ke selected
        selectedWords.push(word);
        availableWords.splice(index, 1);
        // Gabungkan array menjadi string jawaban (contoh: "I eat apple")
        currentSelection = selectedWords.join(' ');
    }

    function removeUnscrambleWord(word: string, index: number) {
        // Kembalikan kata dari selected ke available
        availableWords.push(word);
        selectedWords.splice(index, 1);
        // Update string jawaban
        currentSelection = selectedWords.join(' ');
    }


    function handleNext() {
            collectedAnswers.push({
                material_id: questions[currentIndex].id,
                selected: currentSelection
            });

            if (currentIndex < questions.length - 1) {
                currentIndex++;
                // Note: $effect akan otomatis mereset state unscramble untuk soal berikutnya
            } else {
                submitQuiz();
            }
        }

    async function submitQuiz() {
        isSubmitting = true;
        const payload = { answers: collectedAnswers };

        try {
            const res = await apiFetch(`/modules/${moduleId}/quiz/submit`, {
                method: 'POST',
                data: payload
            });
            
            result = res.data;
        } catch (err: any) {
            // Cek apakah error disebabkan oleh koneksi terputus (Network Error)
            if (err.message === 'Terjadi kesalahan server' || err.message.includes('Failed to fetch')) {
                
                // Simpan ke Local Storage (Mekanisme Offline PRD Phase 4)
                const offlineData = {
                    moduleId: moduleId,
                    payload: payload,
                    timestamp: new Date().toISOString()
                };
                localStorage.setItem(`offline_quiz_${moduleId}`, JSON.stringify(offlineData));

                errorMessage = "Koneksi internet terputus. Jawaban Anda telah disimpan di perangkat dan akan dikirim saat Anda kembali online.";
            } else {
                errorMessage = err.message || "Gagal mengirimkan jawaban kuis.";
            }
        } finally {
            isSubmitting = false;
        }
    }

    
</script>

<div class="min-h-screen bg-slate-50 flex items-center justify-center py-10 px-4">
    {#if isLoading}
        <div class="w-full max-w-xl h-64 bg-slate-200 rounded-3xl animate-pulse"></div>
        
    {:else if errorMessage}
        <div class="w-full max-w-xl p-6 bg-red-50 text-red-700 rounded-2xl text-center font-bold border border-red-100">
            {errorMessage}
            <button onclick={() => goto('/dashboard')} class="block w-full mt-4 underline">Kembali ke Peta Belajar</button>
        </div>

    {:else if result}
        <div class="w-full max-w-xl bg-white rounded-3xl shadow-xl border border-slate-100 p-8 text-center transition-all">
            <h2 class="text-2xl font-black text-slate-800 mb-6">Sesi Selesai!</h2>
            
            <div class="relative w-40 h-40 mx-auto mb-6 flex items-center justify-center rounded-full border-8 {result.score >= 80 ? 'border-green-500 bg-green-50' : 'border-orange-400 bg-orange-50'}">
                <span class="text-5xl font-black {result.score >= 80 ? 'text-green-600' : 'text-orange-500'}">
                    {result.score}%
                </span>
            </div>

            {#if result.completed}
                <h3 class="text-2xl font-black text-green-600 mb-2">Luar Biasa! 🎉</h3>
                <p class="text-slate-600 font-medium mb-6">
                    Anda berhasil menyelesaikan level ini.
                    {#if result.next_module_unlocked}
                        <br><span class="text-indigo-600 font-bold mt-2 block">Level '{result.next_module_unlocked.title}' telah terbuka!</span>
                    {/if}
                </p>
            {:else}
                <h3 class="text-2xl font-black text-orange-600 mb-2">Hampir Saja! 💪</h3>
                <p class="text-slate-600 font-medium mb-6">
                    Anda butuh minimal skor 80% untuk lulus. Terus berlatih, Anda pasti bisa!
                </p>
            {/if}

            <button onclick={() => goto('/dashboard')} class="w-full py-4 bg-slate-800 text-white font-bold rounded-2xl hover:bg-slate-900 transition-colors shadow-md">
                Kembali ke Peta Belajar
            </button>
        </div>

    {:else}
        <div class="w-full max-w-xl bg-white rounded-3xl shadow-lg border border-slate-100 overflow-hidden flex flex-col min-h-[500px]">
            
            <div class="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
                <button onclick={() => goto('/dashboard')} class="text-slate-400 hover:text-red-500 font-bold transition-colors">
                    Tutup ✕
                </button>
                <div class="font-black text-slate-500">
                    Soal {currentIndex + 1} <span class="text-slate-300">/</span> {questions.length}
                </div>
            </div>

            <div class="w-full h-1 bg-slate-100">
                <div class="h-full bg-indigo-500 transition-all duration-300" style="width: {((currentIndex) / questions.length) * 100}%"></div>
            </div>

            <div class="flex-1 p-6 flex flex-col">
                <h2 class="text-2xl sm:text-3xl font-black text-slate-800 mb-8 text-center mt-4">
                    {questions[currentIndex].question}
                </h2>

                <div class="space-y-3 mt-auto mb-8">
                    {#each parseOptions(questions[currentIndex].options) as opt}
                        <button 
                            onclick={() => selectOption(opt)}
                            class="w-full p-4 rounded-2xl border-2 text-left font-bold transition-all duration-200
                            {currentSelection === opt 
                                ? 'border-indigo-500 bg-indigo-50 text-indigo-700 shadow-sm' 
                                : 'border-slate-200 text-slate-600 hover:border-indigo-300 hover:bg-slate-50'}">
                            {opt}
                        </button>
                    {/each}
                </div>

                <button 
                    disabled={!currentSelection || isSubmitting}
                    onclick={handleNext}
                    class="w-full py-4 rounded-2xl font-black text-lg transition-all
                    {!currentSelection 
                        ? 'bg-slate-100 text-slate-400 cursor-not-allowed' 
                        : 'bg-indigo-600 text-white shadow-md hover:bg-indigo-700 hover:shadow-lg'}">
                    
                    {#if isSubmitting}
                        Memeriksa...
                    {:else if currentIndex === questions.length - 1}
                        Kirim Jawaban
                    {:else}
                        Selanjutnya
                    {/if}
                </button>
            </div>
        </div>
    {/if}
    {#if questions[currentIndex].content_type === 'quiz_unscramble'}
                        <div class="w-full flex flex-col gap-6">
                            
                            <div class="min-h-[80px] p-4 border-2 border-dashed border-slate-300 bg-slate-50 rounded-2xl flex flex-wrap content-start gap-2 transition-all">
                                {#each selectedWords as word, i}
                                    <button 
                                        onclick={() => removeUnscrambleWord(word, i)}
                                        class="px-4 py-2.5 bg-indigo-600 text-white font-bold rounded-xl shadow hover:bg-indigo-700 transition-transform hover:-translate-y-0.5">
                                        {word}
                                    </button>
                                {/each}
                                
                                {#if selectedWords.length === 0}
                                    <span class="text-slate-400 font-medium m-auto text-sm">Susun kata-kata di sini...</span>
                                {/if}
                            </div>

                            <div class="flex flex-wrap justify-center gap-2 min-h-[60px]">
                                {#each availableWords as word, i}
                                    <button 
                                        onclick={() => selectUnscrambleWord(word, i)}
                                        class="px-4 py-2.5 bg-white text-slate-700 border-2 border-slate-200 font-bold rounded-xl shadow-sm hover:border-indigo-300 hover:text-indigo-600 transition-transform hover:-translate-y-0.5">
                                        {word}
                                    </button>
                                {/each}
                            </div>
                        </div>

                    {:else}
                        <div class="space-y-3">
                            {#each parseOptions(questions[currentIndex].options) as opt}
                                <button 
                                    onclick={() => selectOption(opt)}
                                    class="w-full p-4 rounded-2xl border-2 text-left font-bold transition-all duration-200
                                    {currentSelection === opt 
                                        ? 'border-indigo-500 bg-indigo-50 text-indigo-700 shadow-sm' 
                                        : 'border-slate-200 text-slate-600 hover:border-indigo-300 hover:bg-slate-50'}">
                                    {opt}
                                </button>
                            {/each}
                        </div>
                    {/if}
</div>
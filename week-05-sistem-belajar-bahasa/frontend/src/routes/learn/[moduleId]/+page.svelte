<script lang="ts">
    import { page } from '$app/stores';
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    // Mengambil ID modul dari URL
    let moduleId = $page.params.moduleId;
    
    // State Svelte 5 Runes
    let cards = $state<any[]>([]);
    let currentIndex = $state(0);
    let isLoading = $state(true);
    let errorMessage = $state('');

    onMount(async () => {
        try {
            // Mengambil semua materi untuk modul ini
            const res = await apiFetch(`/modules/${moduleId}/materials`);
            
            // Filter HANYA materi yang bertipe 'learn_card' (Kartu Belajar)
            cards = (res.data || []).filter((m: any) => m.content_type === 'learn_card');
            
            if (cards.length === 0) {
                errorMessage = "Belum ada materi belajar untuk level ini. Admin perlu menambahkannya.";
            }
        } catch (err: any) {
            errorMessage = err.message || "Gagal memuat materi pembelajaran.";
        } finally {
            isLoading = false;
        }
    });

    function nextCard() {
        if (currentIndex < cards.length - 1) currentIndex++;
    }

    function prevCard() {
        if (currentIndex > 0) currentIndex--;
    }

    function goToQuiz() {
        goto(`/quiz/${moduleId}`);
    }

    function playAudio(url: string) {
        if (!url) return;
        const audio = new Audio(url);
        audio.play().catch(e => console.error("Gagal memutar audio", e));
    }
</script>

<div class="min-h-screen bg-slate-50 flex flex-col items-center justify-center py-10 px-4 sm:px-6">
    <div class="w-full max-w-lg flex justify-between items-center mb-6">
        <button onclick={() => goto('/dashboard')} class="text-slate-500 hover:text-slate-800 flex items-center gap-1 font-bold transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
            Kembali
        </button>
        
        {#if cards.length > 0}
            <div class="text-sm font-black text-slate-400 bg-white px-3 py-1 rounded-full border border-slate-200 shadow-sm">
                Kartu {currentIndex + 1} dari {cards.length}
            </div>
        {/if}
    </div>

    {#if isLoading}
        <div class="w-full max-w-lg h-96 bg-slate-200 rounded-3xl animate-pulse"></div>
    {:else if errorMessage}
        <div class="w-full max-w-lg p-6 bg-red-50 text-red-700 rounded-2xl text-center font-medium border border-red-100">
            {errorMessage}
        </div>
    {:else}
        <div class="w-full max-w-lg bg-white rounded-3xl shadow-lg border border-slate-100 overflow-hidden relative min-h-[400px] flex flex-col">
            
            {#if cards[currentIndex].image_url}
                <div class="w-full h-48 bg-slate-100 flex items-center justify-center border-b border-slate-50 p-4">
                    <img src={cards[currentIndex].image_url} alt="Ilustrasi" class="max-h-full max-w-full object-contain drop-shadow-sm" />
                </div>
            {/if}

            <div class="flex-1 p-8 flex flex-col items-center justify-center text-center">
                <h2 class="text-5xl font-black text-slate-800 mb-2 tracking-tight">
                    {cards[currentIndex].question}
                </h2>
                
                <p class="text-xl font-medium text-slate-500 mb-6">
                    {cards[currentIndex].correct_answer}
                </p>

                {#if cards[currentIndex].audio_url}
                    <button 
                        onclick={() => playAudio(cards[currentIndex].audio_url)}
                        aria-label="Dengarkan audio"
                        class="w-14 h-14 bg-indigo-100 text-indigo-600 rounded-full flex items-center justify-center hover:bg-indigo-600 hover:text-white transition-all shadow-sm hover:shadow-md cursor-pointer">
                        <svg class="w-6 h-6 ml-1" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                    </button>
                    <span class="text-xs font-bold text-slate-400 mt-2 uppercase tracking-wider">Dengarkan</span>
                {/if}
            </div>
            
            <div class="w-full h-1.5 bg-slate-100">
                <div 
                    class="h-full bg-indigo-500 transition-all duration-300" 
                    style="width: {((currentIndex + 1) / cards.length) * 100}%">
                </div>
            </div>
        </div>

        <div class="w-full max-w-lg mt-8 flex justify-between items-center gap-4">
            <button 
                onclick={prevCard} 
                disabled={currentIndex === 0}
                class="flex-1 py-4 px-6 bg-white border border-slate-200 text-slate-700 font-bold rounded-2xl hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm">
                Sebelumnya
            </button>

            {#if currentIndex === cards.length - 1}
                <button 
                    onclick={goToQuiz}
                    class="flex-1 py-4 px-6 bg-indigo-600 border border-indigo-600 text-white font-bold rounded-2xl hover:bg-indigo-700 transition-all shadow-md hover:shadow-lg">
                    Mulai Kuis 🚀
                </button>
            {:else}
                <button 
                    onclick={nextCard}
                    class="flex-1 py-4 px-6 bg-slate-800 border border-slate-800 text-white font-bold rounded-2xl hover:bg-slate-900 transition-all shadow-md">
                    Selanjutnya
                </button>
            {/if}
        </div>
    {/if}
</div>
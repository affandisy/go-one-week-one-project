<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    let products = $state<any[]>([]);
    let suppliers = $state<any[]>([]);
    let isLoading = $state(true);
    let isSubmitting = $state(false);
    let feedbackMsg = $state('');

    // Form State
    let formData = $state({
        product_id: '',
        supplier_id: '',
        quantity: 1,
        note: ''
    });

    onMount(async () => {
        try {
            // Muat produk dan pemasok secara bersamaan
            const [prodRes, suppRes] = await Promise.all([
                apiFetch('/products'),
                apiFetch('/suppliers')
            ]);
            products = prodRes.data || [];
            suppliers = suppRes.data || [];
        } catch (err: any) {
            feedbackMsg = "Gagal memuat data master: " + err.message;
        } finally {
            isLoading = false;
        }
    });

    async function handleSubmit(e: Event) {
        e.preventDefault();
        isSubmitting = true;
        feedbackMsg = '';

        try {
            // Siapkan payload sesuai kebutuhan Backend
            const payload = {
                product_id: formData.product_id,
                supplier_id: formData.supplier_id || undefined, // Opsional
                quantity: formData.quantity,
                note: formData.note
            };

            await apiFetch('/stocks/in', {
                method: 'POST',
                data: payload
            });

            feedbackMsg = "✅ Stok berhasil ditambahkan!";
            
            // Reset form
            formData.product_id = '';
            formData.supplier_id = '';
            formData.quantity = 1;
            formData.note = '';

        } catch (err: any) {
            feedbackMsg = "❌ Gagal menambah stok: " + err.message;
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="max-w-2xl bg-white p-8 rounded-3xl shadow-sm border border-gray-100">
    <div class="mb-6 border-b border-gray-100 pb-4">
        <h2 class="text-2xl font-black text-gray-800 uppercase tracking-tight">Barang Masuk</h2>
        <p class="text-gray-500 font-medium">Catat penambahan stok dari pemasok.</p>
    </div>

    {#if feedbackMsg}
        <div class="mb-6 p-4 {feedbackMsg.startsWith('✅') ? 'bg-green-50 text-green-700 border-green-200' : 'bg-red-50 text-red-700 border-red-200'} font-bold rounded-xl border">
            {feedbackMsg}
        </div>
    {/if}

    {#if isLoading}
        <div class="animate-pulse flex space-x-4">
            <div class="flex-1 space-y-4 py-1">
                <div class="h-10 bg-gray-200 rounded"></div>
                <div class="h-10 bg-gray-200 rounded"></div>
            </div>
        </div>
    {:else}
        <form onsubmit={handleSubmit} class="space-y-5">
            <!-- Pilih Produk -->
            <div>
                <label for="product_id" class="block text-sm font-bold text-gray-700 mb-1">Pilih Produk *</label>
                <select id="product_id" required bind:value={formData.product_id} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 font-medium">
                    <option value="" disabled>-- Pilih Produk --</option>
                    {#each products as product}
                        <option value={product.id}>{product.name} (Sisa Stok: {product.stock})</option>
                    {/each}
                </select>
            </div>

            <!-- Pilih Pemasok (Opsional) -->
            <div>
                <label for="supplier_id" class="block text-sm font-bold text-gray-700 mb-1">Pemasok / Supplier (Opsional)</label>
                <select id="supplier_id" bind:value={formData.supplier_id} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 font-medium">
                    <option value="">-- Tanpa Pemasok --</option>
                    {#each suppliers as supplier}
                        <option value={supplier.id}>{supplier.name}</option>
                    {/each}
                </select>
            </div>

            <div class="grid grid-cols-2 gap-4">
                <!-- Kuantitas -->
                <div>
                    <label for="quantity" class="block text-sm font-bold text-gray-700 mb-1">Jumlah Masuk *</label>
                    <input id="quantity" type="number" required min="1" bind:value={formData.quantity} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 font-black text-xl">
                </div>
                
                <!-- Catatan -->
                <div>
                    <label for="note" class="block text-sm font-bold text-gray-700 mb-1">No. Faktur / Catatan</label>
                    <input id="note" type="text" bind:value={formData.note} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 font-medium" placeholder="Contoh: INV-001">
                </div>
            </div>

            <div class="pt-4">
                <button type="submit" disabled={isSubmitting || !formData.product_id} class="w-full py-4 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white font-black text-lg rounded-xl shadow-lg transition-transform active:scale-95">
                    {isSubmitting ? 'MENYIMPAN...' : 'SIMPAN STOK MASUK'}
                </button>
            </div>
        </form>
    {/if}
</div>
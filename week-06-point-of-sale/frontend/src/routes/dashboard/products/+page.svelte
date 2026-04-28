<script lang="ts">
    import { apiFetch } from '$lib/api';
    
    let products = $state<any[]>([]);
    let showModal = $state(false);
    let formData = $state({ barcode: '', name: '', purchase_price: 0, selling_price: 0, unit: 'pcs', min_stock: 5 });

    async function handleSubmit() {
        await apiFetch('/products', { method: 'POST', data: formData });
        showModal = false;
        // Reset & Refresh...
    }
</script>

<div class="flex justify-between items-center mb-6">
    <div class="relative w-72">
        <input type="text" placeholder="Cari barcode/nama..." class="w-full pl-10 pr-4 py-2 rounded-xl border-gray-200 focus:ring-blue-500">
        <span class="absolute left-3 top-2.5">🔍</span>
    </div>
    <button onclick={() => showModal = true} class="px-6 py-3 bg-blue-600 text-white font-black rounded-xl shadow-lg hover:bg-blue-700">
        + TAMBAH PRODUK
    </button>
</div>
<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';
    import { goto } from '$app/navigation';

    // --- STATE DATA ---
    let products = $state<any[]>([]);
    let isLoadingProducts = $state(true);
    let searchQuery = $state('');

    // --- STATE KERANJANG (CART) ---
    type CartItem = { product_id: string; name: string; price: number; quantity: number; max_stock: number };
    let cart = $state<CartItem[]>([]);
    
    // --- STATE PEMBAYARAN ---
    let discount = $state<number>(0);
    let paymentMethod = $state<'cash' | 'qris'>('cash');
    let cashGiven = $state<number>(0);
    let isProcessing = $state(false);
    let checkoutSuccess = $state(false);
    let lastReceipt = $state('');

    // --- RUNES: KALKULASI INSTAN ($derived) ---
    let filteredProducts = $derived(
        products.filter(p => 
            p.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
            p.barcode.includes(searchQuery)
        )
    );

    let subtotal = $derived(cart.reduce((sum, item) => sum + (item.price * item.quantity), 0));
    let finalAmount = $derived(Math.max(0, subtotal - discount));
    
    // Kembalian hanya dihitung jika bayar tunai dan uang yang diberikan cukup
    let changeAmount = $derived(paymentMethod === 'cash' ? Math.max(0, cashGiven - finalAmount) : 0);
    let isCashEnough = $derived(paymentMethod === 'qris' || cashGiven >= finalAmount);
    let isCartEmpty = $derived(cart.length === 0);

    // --- LIFECYCLE ---
    onMount(async () => {
        await loadProducts();
    });

    async function loadProducts() {
        try {
            const res = await apiFetch('/products');
            products = res.data || [];
        } catch (err) {
            console.error("Gagal memuat produk", err);
        } finally {
            isLoadingProducts = false;
        }
    }

    // --- FUNGSI KERANJANG ---
    function addToCart(product: any) {
        if (product.stock <= 0) {
            alert(`Stok ${product.name} habis!`);
            return;
        }

        const existingIndex = cart.findIndex(item => item.product_id === product.id);
        
        if (existingIndex >= 0) {
            if (cart[existingIndex].quantity < product.stock) {
                cart[existingIndex].quantity += 1;
            } else {
                alert(`Maksimal stok ${product.name} tercapai.`);
            }
        } else {
            cart.push({
                product_id: product.id,
                name: product.name,
                price: product.selling_price,
                quantity: 1,
                max_stock: product.stock
            });
        }
    }

    function updateQuantity(index: number, delta: number) {
        const item = cart[index];
        const newQty = item.quantity + delta;

        if (newQty <= 0) {
            cart.splice(index, 1);
        } else if (newQty <= item.max_stock) {
            cart[index].quantity = newQty;
        }
    }

    function clearCart() {
        cart = [];
        discount = 0;
        cashGiven = 0;
        paymentMethod = 'cash';
        checkoutSuccess = false;
    }

    // --- FUNGSI CHECKOUT ---
    async function handleCheckout() {
        if (isCartEmpty || !isCashEnough) return;

        isProcessing = true;
        try {
            const payload = {
                items: cart.map(c => ({ product_id: c.product_id, quantity: c.quantity })),
                discount: discount,
                payment_method: paymentMethod,
                cash_given: paymentMethod === 'cash' ? cashGiven : finalAmount
            };

            const res = await apiFetch('/checkout', {
                method: 'POST',
                data: payload
            });

            lastReceipt = res.data.receipt_number;
            checkoutSuccess = true;
            await loadProducts();

        } catch (err: any) {
            alert("Gagal memproses transaksi: " + err.message);
        } finally {
            isProcessing = false;
        }
    }

    // --- UTILITY & PENCETAKAN STRUK ---
    const formatRp = (num: number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num);

    function handleLogout() {
        localStorage.clear();
        goto('/login');
    }

    function cetakStruk() {
        window.print();
    }

    const formatWaktu = () => new Date().toLocaleString('id-ID', { 
        year: 'numeric', month: 'short', day: 'numeric', 
        hour: '2-digit', minute: '2-digit' 
    });
</script>

<div class="h-screen flex flex-col bg-gray-100 overflow-hidden font-sans">
    
    <!-- NAVBAR KASIR (Akan disembunyikan saat dicetak) -->
    <header class="bg-blue-800 text-white px-6 py-4 flex justify-between items-center shadow-md z-10 print:hidden">
        <div class="flex items-center gap-3">
            <div class="bg-blue-600 p-2 rounded-lg">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
            </div>
            <h1 class="text-xl font-black tracking-wide">POS Cepat</h1>
        </div>
        <button onclick={handleLogout} class="bg-blue-900 hover:bg-red-600 px-4 py-2 rounded-xl text-sm font-bold transition-colors border border-blue-700">
            Tutup Kasir
        </button>
    </header>

    <!-- LAYOUT UTAMA -->
    <main class="flex-1 flex flex-col lg:flex-row overflow-hidden">
        
        <!-- KOLOM KIRI: KATALOG PRODUK (Disembunyikan saat dicetak) -->
        <section class="flex-1 flex flex-col bg-gray-50 border-r border-gray-200 print:hidden">
            <div class="p-4 bg-white shadow-sm z-0">
                <input type="text" bind:value={searchQuery} placeholder="Cari nama produk atau scan barcode..." 
                    class="w-full bg-gray-100 px-5 py-4 rounded-2xl border-none focus:ring-4 focus:ring-blue-100 text-lg font-medium transition-all">
            </div>

            <div class="flex-1 overflow-y-auto p-4">
                {#if isLoadingProducts}
                    <div class="flex justify-center items-center h-full">
                        <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-blue-600"></div>
                    </div>
                {:else if filteredProducts.length === 0}
                    <div class="text-center text-gray-400 mt-10 font-bold text-lg">Produk tidak ditemukan.</div>
                {:else}
                    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                        {#each filteredProducts as p}
                            <button 
                                onclick={() => addToCart(p)}
                                disabled={p.stock <= 0}
                                class="flex flex-col text-left bg-white p-4 rounded-2xl shadow-sm border border-gray-100 hover:shadow-md hover:border-blue-300 active:scale-95 transition-all disabled:opacity-50 disabled:active:scale-100">
                                <div class="text-sm font-bold text-gray-400 mb-1">{p.barcode || '-'}</div>
                                <h3 class="font-black text-gray-800 leading-tight mb-2 flex-1">{p.name}</h3>
                                <div class="flex justify-between items-end w-full mt-2">
                                    <span class="text-blue-700 font-black">{formatRp(p.selling_price)}</span>
                                    <span class="text-xs font-bold px-2 py-1 rounded-md {p.stock > 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">
                                        Stok: {p.stock}
                                    </span>
                                </div>
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
        </section>

        <!-- KOLOM KANAN: PANEL KERANJANG (Disembunyikan saat dicetak) -->
        <section class="w-full lg:w-[400px] xl:w-[480px] flex flex-col bg-white shadow-[-4px_0_15px_-3px_rgba(0,0,0,0.05)] z-20 print:hidden">
            
            {#if checkoutSuccess}
                <!-- LAYAR SUKSES -->
                <div class="flex-1 flex flex-col items-center justify-center p-8 text-center">
                    <div class="w-24 h-24 bg-green-100 text-green-600 rounded-full flex items-center justify-center mb-6">
                        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg>
                    </div>
                    <h2 class="text-3xl font-black text-gray-800 mb-2">Transaksi Sukses!</h2>
                    <p class="text-gray-500 font-medium mb-6">No. Struk: <span class="font-bold text-gray-800">{lastReceipt}</span></p>
                    
                    <div class="w-full bg-gray-50 p-6 rounded-2xl mb-8 border border-gray-200">
                        <div class="text-sm text-gray-500 font-bold mb-1">Kembalian</div>
                        <div class="text-4xl font-black text-blue-600">{formatRp(changeAmount)}</div>
                    </div>

                    <div class="flex flex-col sm:flex-row gap-4 w-full">
                        <button onclick={cetakStruk} class="flex-1 py-4 bg-blue-100 hover:bg-blue-200 text-blue-700 font-black rounded-2xl transition-colors shadow-sm">
                            🖨️ Cetak Struk
                        </button>
                        <button onclick={clearCart} class="flex-1 py-4 bg-gray-800 hover:bg-gray-900 text-white font-black rounded-2xl transition-colors shadow-sm">
                            Transaksi Baru
                        </button>
                    </div>
                </div>
            {:else}
                <!-- LAYAR KERANJANG -->
                <div class="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                    <h2 class="text-lg font-black text-gray-800">Keranjang ({cart.length})</h2>
                    <button onclick={clearCart} class="text-red-500 hover:text-red-700 text-sm font-bold">Kosongkan</button>
                </div>

                <div class="flex-1 overflow-y-auto p-4 space-y-3">
                    {#if isCartEmpty}
                        <div class="h-full flex items-center justify-center text-gray-400 font-bold text-center px-8">
                            Belum ada barang di keranjang. Silakan pilih produk di sebelah kiri.
                        </div>
                    {/if}

                    {#each cart as item, i}
                        <div class="flex justify-between items-center p-3 bg-white border border-gray-100 rounded-xl shadow-sm">
                            <div class="flex-1 pr-4">
                                <h4 class="font-bold text-gray-800 text-sm">{item.name}</h4>
                                <div class="text-blue-600 font-bold text-sm mt-1">{formatRp(item.price * item.quantity)}</div>
                            </div>
                            <div class="flex items-center gap-3 bg-gray-100 rounded-lg p-1">
                                <button onclick={() => updateQuantity(i, -1)} class="w-8 h-8 bg-white rounded flex items-center justify-center font-black shadow-sm text-gray-600 hover:text-red-500">-</button>
                                <span class="w-6 text-center font-black text-gray-800">{item.quantity}</span>
                                <button onclick={() => updateQuantity(i, 1)} class="w-8 h-8 bg-white rounded flex items-center justify-center font-black shadow-sm text-gray-600 hover:text-blue-600">+</button>
                            </div>
                        </div>
                    {/each}
                </div>

                <div class="bg-gray-50 p-5 border-t border-gray-200">
                    <div class="flex justify-between mb-2 text-gray-500 font-bold">
                        <span>Subtotal</span>
                        <span>{formatRp(subtotal)}</span>
                    </div>
                    
                    <div class="flex justify-between items-center mb-4 text-gray-500 font-bold">
                        <span>Diskon (Rp)</span>
                        <input type="number" bind:value={discount} min="0" class="w-32 text-right p-2 rounded-lg border border-gray-300 text-gray-800 font-bold focus:ring-2 focus:ring-blue-500">
                    </div>

                    <div class="flex justify-between items-center mb-6">
                        <span class="text-xl font-black text-gray-800">Total</span>
                        <span class="text-3xl font-black text-blue-700">{formatRp(finalAmount)}</span>
                    </div>

                    <div class="grid grid-cols-2 gap-3 mb-4">
                        <button onclick={() => paymentMethod = 'cash'} class="py-3 rounded-xl font-bold border-2 transition-all {paymentMethod === 'cash' ? 'border-blue-600 bg-blue-50 text-blue-700' : 'border-gray-200 text-gray-500 bg-white hover:border-blue-300'}">
                            💵 TUNAI
                        </button>
                        <button onclick={() => paymentMethod = 'qris'} class="py-3 rounded-xl font-bold border-2 transition-all {paymentMethod === 'qris' ? 'border-blue-600 bg-blue-50 text-blue-700' : 'border-gray-200 text-gray-500 bg-white hover:border-blue-300'}">
                            📱 QRIS
                        </button>
                    </div>

                    {#if paymentMethod === 'cash'}
                        <div class="mb-6">
                            <label for="cashInput" class="block text-sm font-bold text-gray-600 mb-2">Uang Diterima (Rp)</label>
                            <input id="cashInput" type="number" bind:value={cashGiven} class="w-full text-2xl font-black p-4 rounded-xl border-2 border-gray-300 focus:border-blue-600 focus:ring-0 text-right bg-white shadow-inner">
                        </div>
                    {/if}

                    <button 
                        onclick={handleCheckout}
                        disabled={isCartEmpty || !isCashEnough || isProcessing}
                        class="w-full py-5 rounded-2xl font-black text-xl text-white transition-all shadow-lg
                        {isCartEmpty || !isCashEnough 
                            ? 'bg-gray-300 cursor-not-allowed shadow-none' 
                            : 'bg-blue-600 hover:bg-blue-700 active:scale-95 hover:shadow-xl'}">
                        {isProcessing ? 'Memproses...' : (paymentMethod === 'cash' ? 'TERIMA PEMBAYARAN' : 'BUAT KODE QRIS')}
                    </button>
                </div>
            {/if}
        </section>
    </main>

    <!-- AREA CETAK STRUK (Hanya muncul di kertas printer) -->
    <!-- Ukuran w-[58mm] menyesuaikan standar printer kasir thermal -->
    {#if checkoutSuccess}
        <div class="hidden print:block text-black bg-white w-[58mm] text-[12px] font-mono mx-auto p-2">
            <div class="text-center mb-4">
                <h2 class="font-bold text-sm uppercase">Toko Kelontong</h2>
                <p>Jl. Raya Indramayu No. 1</p>
                <p>Telp: 08123456789</p>
            </div>
            
            <div class="border-b border-dashed border-black pb-1 mb-1">
                <p>No   : {lastReceipt}</p>
                <p>Waktu: {formatWaktu()}</p>
                <p>Kasir: Kasir Utama</p>
            </div>
            
            <table class="w-full text-left mb-2">
                <tbody>
                    {#each cart as item}
                    <tr>
                        <td colspan="3" class="font-bold">{item.name}</td>
                    </tr>
                    <tr>
                        <td>{item.quantity}x</td>
                        <td class="text-right">{formatRp(item.price)}</td>
                        <td class="text-right">{formatRp(item.price * item.quantity)}</td>
                    </tr>
                    {/each}
                </tbody>
            </table>
            
            <div class="border-t border-dashed border-black pt-1">
                <div class="flex justify-between"><p>Subtotal:</p><p>{formatRp(subtotal)}</p></div>
                {#if discount > 0}
                    <div class="flex justify-between"><p>Diskon:</p><p>-{formatRp(discount)}</p></div>
                {/if}
                <div class="flex justify-between font-bold text-sm mt-1"><p>TOTAL:</p><p>{formatRp(finalAmount)}</p></div>
                <div class="flex justify-between mt-1"><p>Tunai/QRIS:</p><p>{formatRp(paymentMethod === 'cash' ? cashGiven : finalAmount)}</p></div>
                <div class="flex justify-between"><p>Kembali:</p><p>{formatRp(changeAmount)}</p></div>
            </div>
            
            <div class="text-center mt-4 border-t border-dashed border-black pt-2">
                <p>Terima Kasih</p>
                <p class="text-[10px] mt-1">Barang yang sudah dibeli tidak dapat ditukar/dikembalikan.</p>
            </div>
        </div>
    {/if}
</div>
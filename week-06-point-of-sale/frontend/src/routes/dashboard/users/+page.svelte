<script lang="ts">
    import { onMount } from 'svelte';
    import { apiFetch } from '$lib/api';

    // State Data
    let users = $state<any[]>([]);
    let isLoading = $state(true);
    let feedbackMsg = $state('');

    // State Form Modal
    let showModal = $state(false);
    let isEditing = $state(false);
    let editId = $state('');
    let formData = $state({ username: '', pin: '', role: 'cashier' });

    onMount(async () => {
        await loadUsers();
    });

    async function loadUsers() {
        isLoading = true;
        try {
            const res = await apiFetch('/users');
            users = res.data || [];
        } catch (err: any) {
            feedbackMsg = "Gagal memuat pengguna: " + err.message;
        } finally {
            isLoading = false;
        }
    }

    function openAddModal() {
        isEditing = false;
        formData = { username: '', pin: '', role: 'cashier' };
        showModal = true;
    }

    function openEditModal(user: any) {
        isEditing = true;
        editId = user.id;
        formData = { username: user.username, pin: '', role: user.role };
        showModal = true;
    }

    async function handleSubmit(e: Event) {
        e.preventDefault();
        try {
            if (isEditing) {
                // Jika PIN kosong saat edit, kita tidak mengirimkannya (Backend akan mengabaikannya)
                const payload: Partial<typeof formData> = { ...formData };
                if (!payload.pin) delete payload.pin;
                
                await apiFetch(`/users/${editId}`, { method: 'PUT', data: payload });
                feedbackMsg = "Data pengguna berhasil diperbarui.";
            } else {
                await apiFetch('/users', { method: 'POST', data: formData });
                feedbackMsg = "Pengguna baru berhasil ditambahkan.";
            }
            showModal = false;
            await loadUsers();
        } catch (err: any) {
            alert(err.message);
        }
    }

    async function handleDelete(id: string, username: string) {
        if (!confirm(`Hapus akun '${username}' permanen? Kasir ini tidak akan bisa login lagi.`)) return;
        try {
            await apiFetch(`/users/${id}`, { method: 'DELETE' });
            feedbackMsg = "Akun berhasil dihapus.";
            await loadUsers();
        } catch (err: any) {
            alert(err.message);
        }
    }

    // Utility untuk warna Badge Role
    const roleColors: Record<string, string> = {
        'admin': 'bg-purple-100 text-purple-700',
        'owner': 'bg-blue-100 text-blue-700',
        'cashier': 'bg-green-100 text-green-700'
    };
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <div>
            <h2 class="text-2xl font-black text-gray-800 uppercase tracking-tight">Manajemen Pegawai</h2>
            <p class="text-gray-500 font-medium">Kelola akses login untuk Kasir dan Admin.</p>
        </div>
        <button onclick={openAddModal} class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-black rounded-xl shadow-lg transition-transform active:scale-95 flex items-center gap-2">
            <span>+</span> TAMBAH PEGAWAI
        </button>
    </div>

    {#if feedbackMsg}
        <div class="p-4 bg-green-50 border-l-4 border-green-500 text-green-700 font-bold rounded-r-xl flex justify-between">
            {feedbackMsg}
            <button onclick={() => feedbackMsg = ''} class="text-green-500 hover:text-green-800">✕</button>
        </div>
    {/if}

    <div class="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
        <div class="overflow-x-auto">
            <table class="w-full text-left">
                <thead class="bg-gray-50 text-gray-400 text-xs font-black uppercase">
                    <tr>
                        <th class="px-6 py-4">Username</th>
                        <th class="px-6 py-4">Hak Akses (Role)</th>
                        <th class="px-6 py-4">Tanggal Dibuat</th>
                        <th class="px-6 py-4 text-right">Aksi</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-50">
                    {#if isLoading}
                        <tr><td colspan="4" class="text-center py-8 text-gray-400 font-bold">Memuat data...</td></tr>
                    {:else}
                        {#each users as user}
                            <tr class="hover:bg-gray-50/50">
                                <td class="px-6 py-4 font-black text-gray-800">{user.username}</td>
                                <td class="px-6 py-4">
                                    <span class="px-3 py-1 rounded-full text-xs font-black uppercase {roleColors[user.role] || 'bg-gray-100 text-gray-700'}">
                                        {user.role}
                                    </span>
                                </td>
                                <td class="px-6 py-4 text-gray-500 font-medium">
                                    {new Date(user.created_at).toLocaleDateString('id-ID')}
                                </td>
                                <td class="px-6 py-4 text-right space-x-2">
                                    <button onclick={() => openEditModal(user)} class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 font-bold rounded-lg text-sm transition-colors">Edit</button>
                                    <button onclick={() => handleDelete(user.id, user.username)} class="px-4 py-2 bg-red-50 hover:bg-red-100 text-red-600 font-bold rounded-lg text-sm transition-colors">Hapus</button>
                                </td>
                            </tr>
                        {/each}
                    {/if}
                </tbody>
            </table>
        </div>
    </div>
</div>

{#if showModal}
    <div class="fixed inset-0 bg-gray-900/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden transform transition-all">
            <div class="p-6 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                <h3 class="text-xl font-black text-gray-800">{isEditing ? 'Edit Pegawai' : 'Tambah Pegawai Baru'}</h3>
                <button onclick={() => showModal = false} class="text-gray-400 hover:text-red-500 font-bold">✕</button>
            </div>
            
            <form onsubmit={handleSubmit} class="p-6 space-y-5">
                <div>
                    <label for="username" class="block text-sm font-bold text-gray-700 mb-1">Username</label>
                    <input id="username" type="text" required bind:value={formData.username} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 focus:ring-0 font-medium">
                </div>
                <div>
                    <label for="pin" class="block text-sm font-bold text-gray-700 mb-1">
                        PIN Login {isEditing ? '(Kosongkan jika tidak ingin diubah)' : ''}
                    </label>
                    <input id="pin" type="password" required={!isEditing} bind:value={formData.pin} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 focus:ring-0 font-medium tracking-widest" placeholder="••••••">
                </div>
                <div>
                    <label for="role" class="block text-sm font-bold text-gray-700 mb-1">Pilih Hak Akses</label>
                    <select id="role" bind:value={formData.role} class="w-full p-3 bg-gray-50 border border-gray-200 rounded-xl focus:border-blue-500 focus:ring-0 font-medium uppercase">
                        <option value="cashier">CASHIER (Kasir - Transaksi Saja)</option>
                        <option value="owner">OWNER (Pemilik - Laporan & Master)</option>
                        <option value="admin">ADMIN (Teknisi Sistem)</option>
                    </select>
                </div>
                <div class="pt-4 flex gap-3">
                    <button type="button" onclick={() => showModal = false} class="flex-1 py-3 bg-gray-100 hover:bg-gray-200 text-gray-700 font-bold rounded-xl transition-colors">Batal</button>
                    <button type="submit" class="flex-1 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl shadow-md transition-colors">Simpan</button>
                </div>
            </form>
        </div>
    </div>
{/if}
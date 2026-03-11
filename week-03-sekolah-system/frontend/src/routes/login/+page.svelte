<script lang="ts">
	import { goto } from '$app/navigation';
	let email = '';
	let password = '';
	let errorMsg = '';

	async function handleLogin() {
		const res = await fetch('http://localhost:3000/api/v1/auth/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email, password })
		});
		
		const data = await res.json();
		if (!res.ok) {
			errorMsg = data.error;
		} else {
			localStorage.setItem('token', data.data.token);
			goto('/dashboard');
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100">
	<div class="bg-white p-8 rounded shadow-md w-96">
		<h2 class="text-2xl font-bold mb-6 text-center text-blue-600">Sistem Sekolah SMA</h2>
		{#if errorMsg}
			<div class="bg-red-100 text-red-700 p-2 rounded mb-4 text-sm">{errorMsg}</div>
		{/if}
		<form on:submit|preventDefault={handleLogin} class="space-y-4">
			<div>
				<label class="block text-sm font-medium text-gray-700">Email</label>
				<input type="email" bind:value={email} class="mt-1 block w-full border border-gray-300 rounded p-2" required />
			</div>
			<div>
				<label class="block text-sm font-medium text-gray-700">Password</label>
				<input type="password" bind:value={password} class="mt-1 block w-full border border-gray-300 rounded p-2" required />
			</div>
			<button type="submit" class="w-full bg-blue-600 text-white p-2 rounded hover:bg-blue-700">Login</button>
		</form>
	</div>
</div>
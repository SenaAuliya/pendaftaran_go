<script lang="ts">
	import { goto } from '$app/navigation';
	import axios from 'axios';

	let username = '';
	let password = '';
	let message = '';
	let loading = false;

	// Define the login function
	async function login(username: string, password: string) {
		try {
			const response = await axios.post('http://localhost:3030/login', { username, password });
			const token = response.data.token;
			sessionStorage.setItem('token', token);
			sessionStorage.setItem('username', username);
			return true; // Indicate successful login
		} catch (error) {
			console.error('Login failed:', error);
			return false;
		}
	}

	const handleLogin = async () => {
		if (!username || !password) {
			message = 'Username and password are required.';
			return;
		}

		loading = true; // Show loading indicator

		const success = await login(username, password);

		if (success) {
			goto("/");
		} else {
			message = 'Login failed. Please check your username and password.';
		}

		loading = false; // Hide loading indicator
	};
</script>

<div class="flex flex-col items-center pt-[3rem] w-full">
	<main class="flex justify-center items-center h-[100vh]">
		<div class="flex flex-col justify-between items-center bg-slate-800 p-5 h-[20rem] rounded-md shadow-lg">
			<h1 class="text-2xl font-semibold text-white mb-4">Login</h1>
			<input class="rounded-md h-10 w-[20rem] p-3 mb-3" type="text" placeholder="Username" bind:value={username} />
			<input class="rounded-md h-10 w-[20rem] p-3 mb-3" type="password" placeholder="Password" bind:value={password} />
			<button 
				on:click={handleLogin} 
				class="bg-slate-500 py-2 px-5 rounded-md text-lg font-medium text-gray-200"
				disabled={loading}
			>
				{#if loading}
					<span>Loading...</span>
				{:else}
					<span>Login</span>
				{/if}
			</button>
			{#if message}
				<p class="text-red-400 mt-2">{message}</p>
			{/if}
		</div>
	</main>
</div>

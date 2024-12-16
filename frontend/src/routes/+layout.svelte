<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import Accordion from '../components/accordion.svelte';
	import { logout } from '$lib/auth';
	import { Url } from '../function/var';
	import axios from 'axios';
	import { onMount } from 'svelte';
	import type { User } from '../function/interface';

	let signOut = false;
	let userRole: string | null = null;
	let username: string | null = null;

	// Check if running in the browser
	if (typeof window !== 'undefined') {
		username = sessionStorage.getItem('username'); // Get username from localStorage
	}

	// Function to fetch user data
	async function fetchUser() {
		try {
			// Ensure this code only runs in the browser
			if (typeof window === 'undefined') return;

			const token = sessionStorage.getItem('token');
			if (!token) {
				console.error('No token found in localStorage');
				return;
			}

			const res = await axios.get('http://localhost:3030/user', {
				headers: { Authorization: `Bearer ${token}` }
			});

			console.log('API response:', res.data); // Log the full response

			const users = res.data.data;
			console.log('User data:', users); // Log the user data array

			const currentUsername = sessionStorage.getItem('username');
			if (!currentUsername) {
				console.error('No username found in localStorage');
				return;
			}

			const currentUser = users.find((user: User) => user.Username === currentUsername);

			if (currentUser) {
				console.log('Current User:', currentUser); // Log the current user details
				if (currentUser.Role) {
					userRole = currentUser.Role;
				} else {
					console.error('Role is not available for the current user');
				}
			} else {
				console.error('Current user not found in the user data array');
			}

			console.log('User Role:', userRole); // Log the role
		} catch (err) {
			console.error('Error fetching user data:', err);
		}
	}

	onMount(() => {
		if (typeof window !== 'undefined') {
			fetchUser(); // Fetch user data when component mounts
		}
	});

	$: currentPath = $page.url.pathname;

	// Function to handle user logout
	function handleLogOut() {
		logout();
	}

	// Function to toggle the visibility of the logout button
	function showLogOut() {
		signOut = !signOut;

		// Disable scrolling when the modal is open
		if (signOut) {
			document.body.style.overflow = 'hidden';
			document.documentElement.style.overflow = 'hidden';
		} else {
			document.body.style.overflow = '';
			document.documentElement.style.overflow = '';
		}
	}
</script>

{#if currentPath !== '/auth/login'}
	<!-- If not on the login page, show this layout -->
	<header class="py-3 px-5 bg-slate-500 w-full flex justify-between items-center">
		<h1 class="text-xl font-medium text-white" > Welcome {username}</h1>

		<button class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 font-medium"on:click={showLogOut}>Logout</button>
	</header>

	{#if signOut}
    <!-- Modal overlay for logout confirmation -->
    <div class="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50 z-50">
        <div class="bg-white p-6 rounded-md shadow-lg z-60">
            <h2 class="text-lg font-semibold mb-4">Are you sure you want to logout?</h2>
            <div class="flex justify-end gap-3">
                <button class="bg-gray-300 px-4 py-2 rounded-md" on:click={showLogOut}>Cancel</button>
                <button class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700" on:click={handleLogOut}>Logout</button>
            </div>
        </div>
    </div>
{/if}


	<section class="flex flex-row gap-5 overflow-y-hidden">
		<section class="flex h-[56.755rem] w-[20%] flex-col items-center gap-5 bg-gray-600 py-5">
			<h1 class="text-3xl font-semibold text-white">Dashboard</h1>
			<a href="/" class="text-white text-lg font-medium">Home</a>
			<Accordion title="Master Data">
				<ul>
					{#each Url.filter(item => userRole === 'admin' || item.showForUser) as item}
						<a
							class="py-2 text-white cursor-pointer flex flex-col justify-center items-center hover:text-gray-100 font-medium"
							href={item.link}
						>
							{item.nama}
						</a>
					{/each}
				</ul>
			</Accordion>
			<a href="/laporan" class="text-white text-lg font-medium">Laporan</a>
		</section>

		<section class="flex w-full flex-col items-center pt-[3rem]">
			<slot></slot>
		</section>
	</section>
{:else}
	<!-- If on the login page, only show the content of the page -->
	<slot></slot>
{/if}

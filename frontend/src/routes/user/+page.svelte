<script lang="ts">
	import { onMount } from 'svelte';
	import axios from 'axios';
	import type { UserData, ApiResponse } from '../../function/interface';
	import { goto } from '$app/navigation';
  
	let admin: UserData[] = [];
	let searchQuery = '';
	let page = 1;
	let limit = 5;
	let totalPages = 1;
	let showModalCreate = false;
	let Username = '';
	let Password = '';
	let showError = false;
	let errorMessage = '';
	let showModalEdit = false;
	let editId: number | null = null;
	let showPopUp = false;
	let deleteId: number | null = null;
  
	async function fetchData(page: number, limit: number, query: string) {
	  try {
		const token = sessionStorage.getItem('token');
		if (!token) {
		  throw new Error('Token not found');
		}
		const response = await axios.get<ApiResponse<UserData>>(`http://localhost:3030/user`, {
		  params: { page, limit, query },
		  headers: { Authorization: `Bearer ${token}` }
		});
		admin = response.data.data;
		totalPages = response.data.totalPages;
	  } catch (err) {
		console.error('Error:', err);
		goto('/auth/login');
		admin = [];
	  }
	}
  
	onMount(() => {
	  fetchData(page, limit, searchQuery);
	});
  
	function prevPage() {
	  if (page > 1) {
		page--;
		fetchData(page, limit, searchQuery);
	  }
	}
  
	function nextPage() {
	  if (page < totalPages) {
		page++;
		fetchData(page, limit, searchQuery);
	  }
	}
  
	function handleSearch(event: Event) {
	  event.preventDefault();
	  const input = (event.target as HTMLFormElement).querySelector('input') as HTMLInputElement;
	  searchQuery = input.value;
	  page = 1;
	  fetchData(page, limit, searchQuery);
	}
  
	function openModalCreate() {
	  showModalCreate = true;
	}
  
	async function onSubmitCreate(e: SubmitEvent) {
	  e.preventDefault();
	  try {
		const token = sessionStorage.getItem('token');
		if (!token) {
		  throw new Error('Token not found');
		}
		await axios.post('http://localhost:3030/user', { Username, Password }, {
		  headers: { Authorization: `Bearer ${token}` }
		});
		closeModalCreate();
		fetchData(page, limit, searchQuery);
	  } catch (err: any) {
		showError = true;
		errorMessage = err.response?.data.message || 'Failed to send data';
	  }
	}
  
	function closeError() {
	  showError = false;
	  errorMessage = '';
	}
  
	function closeModalCreate() {
	  showModalCreate = false;
	}
  
	function openModalEdit(id: number) {
	  editId = id;
	  showModalEdit = true;
	  fetchEditData(id);
	}
  
	async function fetchEditData(id: number) {
	  try {
		const token = sessionStorage.getItem('token');
		if (!token) {
		  throw new Error('Token not found');
		}
		const response = await axios.get(`http://localhost:3030/user/${id}`, {
		  headers: { Authorization: `Bearer ${token}` }
		});
		const data: UserData = response.data;
		Username = data.Username;
		Password = data.Password;
	  } catch (err) {
		console.error(err);
	  }
	}
  
	async function onUpdate(e: SubmitEvent) {
	  e.preventDefault();
	  try {
		const token = sessionStorage.getItem('token');
		if (!token) {
		  throw new Error('Token not found');
		}
		await axios.put(`http://localhost:3030/user/${editId}`, { Username, Password }, {
		  headers: { Authorization: `Bearer ${token}` }
		});
		closeModalEdit();
		fetchData(page, limit, searchQuery);
	  } catch (err) {
		console.error(err);
	  }
	}
  
	function closeModalEdit() {
	  showModalEdit = false;
	  editId = null;
	}
  
	function openPopUpDelete(id: number) {
	  deleteId = id;
	  showPopUp = true;
	}
  
	async function deleteUser() {
	  if (deleteId !== null) {
		try {
		  const token = sessionStorage.getItem('token');
		  if (!token) {
			throw new Error('Token not found');
		  }
		  await axios.delete(`http://localhost:3030/user/${deleteId}`, {
			headers: { Authorization: `Bearer ${token}` }
		  });
		  fetchData(page, limit, searchQuery);
		  closePopUp();
		} catch (err) {
		  console.error(err);
		}
	  }
	}
  
	function closePopUp() {
	  showPopUp = false;
	  deleteId = null;
	}
  </script>
  
  <h1 class="mb-2 text-3xl font-medium">Data Admin</h1>
  <div class="mb-4 flex flex-row items-center">
	<button
	  on:click={openModalCreate}
	  class="mr-4 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
	>
	  Tambah User
	</button>
	<form on:submit={handleSearch} class="flex flex-row items-center">
	  <input
		type="text"
		placeholder="Cari Admin..."
		class="rounded border px-4 py-2"
		value={searchQuery}
	  />
	  <button
		type="submit"
		class="ml-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
	  >
		Cari
	  </button>
	</form>
  </div>
  <div class="relative overflow-x-auto">
	<table class="w-full text-left text-sm text-gray-500 rtl:text-right dark:text-gray-400">
	  <thead class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-gray-700 dark:text-gray-400">
		<tr>
		  <th scope="col" class="px-6 py-3">ID</th>
		  <th scope="col" class="px-6 py-3">User</th>
		  <th scope="col" class="px-6 py-3">Password</th>
		  <th scope="col" class="px-6 py-3">Aksi</th>
		</tr>
	  </thead>
	  <tbody>
		{#if admin.length > 0}
		  {#each admin as item}
			<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
			  <th
				scope="row"
				class="whitespace-nowrap px-6 py-4 font-medium text-gray-900 dark:text-white"
			  >
				{item.ID}
			  </th>
			  <td class="px-6 py-4">{item.Username}</td>
			  <td class="px-6 py-4">*********</td>
			  <td class="flex items-center justify-center gap-4 px-6 py-4">
				<button
				  class="flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
				  on:click={() => openModalEdit(item.ID)}
				>
				  Edit
				</button>
				<button
				  class="flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
				  on:click={() => openPopUpDelete(item.ID)}
				>
				  Hapus
				</button>
			  </td>
			</tr>
		  {/each}
		{:else}
		  <tr>
			<td colspan="4" class="px-6 py-4 text-center">Tidak ada data</td>
		  </tr>
		{/if}
	  </tbody>
	</table>
  </div>
  <div class="mt-4 flex justify-center">
	<button
	  class="focus:shadow-outline mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
	  on:click={prevPage}
	  disabled={page === 1}
	>
	  Prev
	</button>
	<button
	  class="focus:shadow-outline mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
	  on:click={nextPage}
	  disabled={page === totalPages}
	>
	  Next
	</button>
  </div>
  
  {#if showModalCreate}
	<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
	  <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-md">
		<h2 class="mb-4 text-xl font-bold">Tambah Admin</h2>
		{#if showError}
		  <div class="mb-4 rounded bg-red-500 p-2 text-white">{errorMessage}</div>
		{/if}
		<form on:submit={onSubmitCreate}>
		  <div class="mb-4">
			<label class="mb-2 block font-bold text-gray-700" for="username">Username</label>
			<input
			  class="w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 focus:outline-none"
			  id="username"
			  type="text"
			  bind:value={Username}
			/>
		  </div>
		  <div class="mb-4">
			<label class="mb-2 block font-bold text-gray-700" for="password">Password</label>
			<input
			  class="w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 focus:outline-none"
			  id="password"
			  type="password"
			  bind:value={Password}
			/>
		  </div>
		  <div class="flex justify-end">
			<button
			  class="focus:shadow-outline rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
			  type="submit"
			>
			  Simpan
			</button>
			<button
			  class="focus:shadow-outline ml-2 rounded bg-red-500 px-4 py-2 font-bold text-white hover:bg-red-700 focus:outline-none"
			  on:click={closeModalCreate}
			  type="button"
			>
			  Batal
			</button>
		  </div>
		</form>
	  </div>
	</div>
  {/if}
  
  {#if showModalEdit}
	<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
	  <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-md">
		<h2 class="mb-4 text-xl font-bold">Edit Admin</h2>
		<form on:submit={onUpdate}>
		  <div class="mb-4">
			<label class="mb-2 block font-bold text-gray-700" for="username">Username</label>
			<input
			  class="w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 focus:outline-none"
			  id="username"
			  type="text"
			  bind:value={Username}
			/>
		  </div>
		  <div class="mb-4">
			<label class="mb-2 block font-bold text-gray-700" for="password">Password</label>
			<input
			  class="w-full appearance-none rounded border px-3 py-2 leading-tight text-gray-700 focus:outline-none"
			  id="password"
			  type="password"
			  bind:value={Password}
			/>
		  </div>
		  <div class="flex justify-end">
			<button
			  class="focus:shadow-outline rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
			  type="submit"
			>
			  Simpan
			</button>
			<button
			  class="focus:shadow-outline ml-2 rounded bg-red-500 px-4 py-2 font-bold text-white hover:bg-red-700 focus:outline-none"
			  on:click={closeModalEdit}
			  type="button"
			>
			  Batal
			</button>
		  </div>
		</form>
	  </div>
	</div>
  {/if}
  
  {#if showPopUp}
	<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
	  <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-md">
		<h2 class="mb-4 text-xl font-bold">Konfirmasi Hapus</h2>
		<p>Apakah Anda yakin ingin menghapus admin ini?</p>
		<div class="mt-4 flex justify-end">
		  <button
			class="focus:shadow-outline rounded bg-red-500 px-4 py-2 font-bold text-white hover:bg-red-700 focus:outline-none"
			on:click={deleteUser}
		  >
			Hapus
		  </button>
		  <button
			class="focus:shadow-outline ml-2 rounded bg-gray-500 px-4 py-2 font-bold text-white hover:bg-gray-700 focus:outline-none"
			on:click={closePopUp}
		  >
			Batal
		  </button>
		</div>
	  </div>
	</div>
  {/if}
  
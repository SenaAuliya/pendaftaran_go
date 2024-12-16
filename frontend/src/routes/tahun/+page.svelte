<script lang="ts">
	import { onMount } from 'svelte';
	import axios from 'axios';
	import type { TahunAjaran } from '../../function/interface';
	import { goto } from '$app/navigation';
	import DeletePopUp from '../../components/deletePopUp.svelte';

	let TahunAjaran: TahunAjaran[] = [];
	let searchQuery = '';
	let page = 1;
	let limit = 5;
	let totalPages = 1;
	let showModalCreate = false;
	let Tahun = '';
	let showError = false;
	let errorMessage = '';
	let showModalEdit = false;
	let editId: number | null = null;
	let showPopUp = false;
	let deleteId: number | null = null;
	let isLoading = false;

	async function fetchData(page: number, limit: number, query: string) {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get(`http://localhost:3030/tahun`, {
				params: {
					page,
					limit,
					query
				},
				headers: { Authorization: `Bearer ${token}` }
			});
			TahunAjaran = response.data.data;
			totalPages = response.data.totalPages;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			TahunAjaran = [];
		} finally {
			isLoading = false;
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

	function modalCreate() {
		showModalCreate = true;
		Tahun = '';
	}

	async function onSubmitCreate(e: SubmitEvent) {
		e.preventDefault();
		try {
			const token = sessionStorage.getItem('token');
			const response = await axios.post('http://localhost:3030/tahun', {
				Tahun,
			}, {
				headers: { Authorization: `Bearer ${token}` }
			});
			closeModalCreate();
			fetchData(page, limit, searchQuery);
		} catch (err: any) {
			showError = true;
			errorMessage = err.response?.data.message || 'Gagal Mengirim Data';
		}
	}

	function closeError() {
		showError = false;
		errorMessage = '';
	}

	function closeModalCreate() {
		showModalCreate = false;
	}

	function modalEdit(id: number) {
		editId = id;
		showModalEdit = true;
		fetchEditData(id);
	}

	async function fetchEditData(id: number) {
		try {
			const token = sessionStorage.getItem('token');
			const response = await axios.get(`http://localhost:3030/tahun/${id}`, {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data: TahunAjaran = response.data;
			Tahun = data.Tahun;
		} catch (err) {
			console.log(err);
		}
	}

	async function onUpdate(e: SubmitEvent) {
		e.preventDefault();
		try {
			const token = sessionStorage.getItem('token');
			await axios.put(`http://localhost:3030/tahun/${editId}`, {
				Tahun,
			}, {
				headers: { Authorization: `Bearer ${token}` }
			});
			closeModalEdit();
			fetchData(page, limit, searchQuery);
		} catch (err) {
			console.log(err);
		}
	}

	function closeModalEdit() {
		showModalEdit = false;
		editId = null;
	}

	function popUpDelete(id: number) {
		deleteId = id;
		showPopUp = true;
	}

	async function deleteTahun() {
		if (deleteId !== null) {
			try {
				const token = sessionStorage.getItem('token');
				await axios.delete(`http://localhost:3030/smp/${deleteId}`, {
					headers: { Authorization: `Bearer ${token}` }
				});
				fetchData(page, limit, searchQuery);
				closePopUp();
			} catch (err) {
				console.log(err);
			}
		}
	}

	function closePopUp() {
		showPopUp = false;
		deleteId = null;
	}
</script>

<h1 class="mb-2 text-3xl font-medium">Data Tahun Ajaran</h1>
<div class="mb-4 flex flex-row items-center">
	<button
		on:click={modalCreate}
		class="mr-4 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
		>Tambah Tahun</button
	>
	<form on:submit={handleSearch} class="flex flex-row items-center">
		<input
			type="text"
			placeholder="Cari Tahun..."
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
{#if isLoading}
	<p>Loading...</p>
{:else}
	<div class="relative overflow-x-auto">
		<table class="w-[120vh] text-left text-sm text-gray-500 rtl:text-right dark:text-gray-400">
			<thead class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-gray-700 dark:text-gray-400">
				<tr>
					<th scope="col" class="px-6 py-3">ID</th>
					<th scope="col" class="px-6 py-3">Tahun</th>
					<th scope="col" class="px-6 py-3">Aksi</th>
				</tr>
			</thead>
			<tbody>
				{#if TahunAjaran.length > 0}
					{#each TahunAjaran as item}
						<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
							<th
								scope="row"
								class="whitespace-nowrap px-6 py-4 font-medium text-gray-900 dark:text-white"
							>
								{item.ID}
							</th>
							<td class="px-6 py-4">{item.Tahun}</td>
							<td class="flex items-center justify-center gap-4 px-6 py-4">
								<button
									class="flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
									on:click={() => modalEdit(item.ID)}>Edit</button
								>
								<button
									class="flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
									on:click={() => popUpDelete(item.ID)}>Hapus</button
								>
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
			disabled={page === 1}>Prev</button
		>
		<button
			class="focus:shadow-outline mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
			on:click={nextPage}
			disabled={page === totalPages}>Next</button
		>
	</div>
{/if}

{#if showPopUp}
	<DeletePopUp onClose={closePopUp} onDelete={deleteTahun}/>
{/if}

{#if showModalCreate || showModalEdit}
	<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
		<div class="rounded bg-white p-4 shadow-lg">
			<form on:submit={showModalCreate ? onSubmitCreate : onUpdate}>
				<label for="Nama" class="block text-sm font-bold text-gray-700">Nama:</label>
				<input
					type="text"
					id="Nama"
					bind:value={Tahun}
					class="mb-2 w-full rounded border px-3 py-2"
				/>
				<div class="flex justify-end">
					<button
						type="submit"
						class="focus:shadow-outline rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
					>{showModalCreate ? 'Kirim' : 'Update'}</button
					>
					<button
						type="button"
						class="focus:shadow-outline ml-2 rounded bg-gray-500x-4 py-2 font-bold text-white hover:bg-gray-700 focus:outline-none"
						on:click={showModalCreate ? closeModalCreate : closeModalEdit}
					>Batal</button>
				</div>
			</form>
			{#if showError}
				<div class="mt-4 text-red-500">{errorMessage}</div>
				<button
					class="mt-2 rounded bg-red-500 px-4 py-2 text-white hover:bg-red-700"
					on:click={closeError}>Tutup</button>
			{/if}
		</div>
	</div>
{/if}
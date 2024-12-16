<script lang="ts">
	import { onMount } from 'svelte';
	import axios from 'axios';
	import type { Guru, Jurusan, Pendaftaran, Siswa, TahunAjaran } from '../../function/interface';
	import { goto } from '$app/navigation';
	import DeletePopUp from '../../components/deletePopUp.svelte';

	let pendaftaran: Pendaftaran[] = [];
	let DataSiswa: Siswa[] = [];
	let DataGuru: Guru[] = [];
	let DataTahun: TahunAjaran[] = [];
	let DataJurusan: Jurusan[] = [];
	let searchQuery = '';
	let page = 1;
	let limit = 5;
	let totalPages = 1;
	let showModalCreate = false;
	let GuruID: number | null = null;
	let SiswaID: number | null = null;
	let TahunID: number | null = null;
	let JurusanID: number | null = null;
	let TanggalPendaftaran = '';
	let TipePendaftaran = '';
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
			const response = await axios.get('http://localhost:3030/pendaftaran', {
				params: { page, limit, query },
				headers: { Authorization: `Bearer ${token}` }
			});
			pendaftaran = response.data.data;
			totalPages = response.data.totalPages;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			pendaftaran = [];
		} finally {
			isLoading = false;
		}
	}

	onMount(async () => {
		fetchData(page, limit, searchQuery);
		await fetchDataSiswa();
		await fetchDataGuru();
		await fetchDataJurusan();
		await fetchDataTahun();
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
		resetForm();
	}

	async function onSubmitCreate(e: SubmitEvent) {
		e.preventDefault();

		if (!SiswaID) {
			showError = true;
			errorMessage = 'Please fill Siswa fields.';
			return;
		} else if (!GuruID) {
			showError = true;
			errorMessage = 'Please fill GuruID fields.';
			return;
		} else if (!TahunID) {
			showError = true;
			errorMessage = 'Please fill TahunID fields.';
			return;
		} else if (!JurusanID) {
			showError = true;
			errorMessage = 'Please fill AsalSekolah fields.';
			return;
		} else if (!TanggalPendaftaran) {
			showError = true;
			errorMessage = 'Please fill TanggalPendaftaran fields.';
			return;
		} else if (!TipePendaftaran) {
			showError = true;
			errorMessage = 'Please fill TipePendaftaran fields.';
			return;
		}

		const payload = {
			SiswaID,
			GuruID,
			TahunID,
			JurusanID,
			TipePendaftaran,
			TanggalPendaftaran: new Date(TanggalPendaftaran).toISOString()
		};

		console.log('Payload:', payload); // Log the payload

		const formData = new FormData();
		formData.append('SiswaID', SiswaID.toString());
		formData.append('GuruID', GuruID.toString());
		formData.append('TahunID', TahunID.toString());
		formData.append('JurusanID', JurusanID.toString());
		formData.append('TipePendaftaran', TipePendaftaran);
		if (TanggalPendaftaran) {
			const formattedDate = new Date(TanggalPendaftaran).toISOString();
			formData.append('TanggalPendaftaran', formattedDate);
		}

		try {
			const token = sessionStorage.getItem('token');
			await axios.post('http://localhost:3030/pendaftaran', formData, {
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'multipart/form-data'
				}
			});
			closeModalCreate();
			fetchData(page, limit, searchQuery);
		} catch (err: any) {
			showError = true;
			errorMessage = err;
			console.error('Error during POST request:', err.response?.data);
		}
	}

	function closeError() {
		showError = false;
		errorMessage = '';
	}

	function closeModalCreate() {
		showModalCreate = false;
	}

	function resetForm() {
		SiswaID = null;
		GuruID = null;
		TahunID = null;
		JurusanID = null;
		TanggalPendaftaran = '';
		TipePendaftaran = '';
	}

	function modalEdit(id: number) {
		editId = id;
		showModalEdit = true;
		fetchEditData(id);
	}

	async function fetchEditData(id: number) {
		try {
			const token = sessionStorage.getItem('token');
			const response = await axios.get(`http://localhost:3030/pendaftaran/${id}`, {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data: Pendaftaran = response.data;
			SiswaID = data.SiswaID;
			GuruID = data.GuruID;
			TahunID = data.TahunID;
			JurusanID = data.JurusanID;
			TipePendaftaran = data.TipePendaftaran;
			TanggalPendaftaran = data.TanggalPendaftaran; // Adjust if necessary
		} catch (err) {
			console.log(err);
		}
	}

	async function onUpdate(e: SubmitEvent) {
		e.preventDefault();
		const formData = new FormData();
		if (SiswaID !== null) {
			formData.append('SiswaID', SiswaID.toString());
		}
		if (GuruID !== null) {
			formData.append('GuruID', GuruID.toString());
		}
		if (TahunID !== null) {
			formData.append('TahunID', TahunID.toString());
		}
		if (JurusanID !== null) {
			formData.append('JurusanID', JurusanID.toString());
		}
		formData.append('TipePendaftaran', TipePendaftaran);
		if (TanggalPendaftaran) {
			formData.append('TanggalPendaftaran', new Date(TanggalPendaftaran).toISOString());
		}
		try {
			const token = sessionStorage.getItem('token');
			await axios.put(`http://localhost:3030/pendaftaran/${editId}`, formData, {
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'multipart/form-data'
				}
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

	async function deletePendaftaran() {
		if (deleteId !== null) {
			try {
				const token = sessionStorage.getItem('token');
				await axios.delete(`http://localhost:3030/pendaftaran/${deleteId}`, {
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

	async function fetchDataSiswa() {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get('http://localhost:3030/siswa', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataSiswa = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			DataSiswa = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchDataGuru() {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get('http://localhost:3030/guru', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataGuru = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			DataGuru = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchDataTahun() {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get('http://localhost:3030/tahun', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataTahun = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			DataTahun = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchDataJurusan() {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get('http://localhost:3030/jurusan', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataJurusan = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			DataJurusan = [];
		} finally {
			isLoading = false;
		}
	}
</script>

<h1 class="mb-2 text-3xl font-medium">Data Siswa Ajaran</h1>
<div class="mb-4 flex flex-row items-center">
	<button
		on:click={modalCreate}
		class="mr-4 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
		>Tambah Siswa</button
	>
	<form on:submit={handleSearch} class="flex flex-row items-center">
		<input
			type="text"
			placeholder="Cari Siswa..."
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
		<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
			<thead class="bg-gray-50 dark:bg-gray-700">
				<tr>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>ID</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Siswa</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Guru</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Jurusan</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Tahun Ajaran</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white">
						Tipe Pendaftaran</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Tanggal Pendaftaran</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Actions</th
					>
				</tr>
			</thead>
			<tbody>
				{#each pendaftaran as item}
					<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
						<td class="px-6 py-4 text-white">{item.ID}</td>
						<td class="px-6 py-4 text-white">
							{#if item.SiswaID}
								{DataSiswa.find((murid) => murid.ID === item.SiswaID)?.Nama || 'unknown'}
							{:else}
								not name
							{/if}
						</td>
						<td class="px-6 py-4 text-white">
							{#if item.GuruID}
								{DataGuru.find((g) => g.ID === item.GuruID)?.Nama || 'unknown'}
							{:else}
								not name
							{/if}
						</td>
						<td class="px-6 py-4 text-white">
							{#if item.JurusanID}
								{DataJurusan.find((j) => j.ID === item.JurusanID)?.Nama || 'unknown'}
							{:else}
								not name
							{/if}
						</td>
						<td class="px-6 py-4 text-white">
							{#if item.TahunID}
								{DataTahun.find((t) => t.ID === item.TahunID)?.Tahun || 'Unknown'}
							{:else}
								No tahun
							{/if}
						</td>
						<td class="px-6 py-4 text-white">{item.TipePendaftaran}</td>
						<td class="px-6 py-4 text-white">
							{#if item.TanggalPendaftaran}
								{new Date(item.TanggalPendaftaran).toLocaleDateString()}
							{:else}
								No Date
							{/if}
						</td>
						<td class="flex flex-row gap-2 px-6 py-4 text-white">
							<button
								on:click={() => modalEdit(item.ID)}
								class="flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
								>Edit</button
							>
							<button
								on:click={() => popUpDelete(item.ID)}
								class="ml-2 flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
								>Hapus</button
							>
						</td>
					</tr>
				{/each}
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

<!-- Modal Create -->
{#if showModalCreate}
	<div class="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50">
		<div class="rounded bg-white p-4 shadow-lg">
			<h2 class="mb-4 text-2xl font-bold">Tambah Pendaftaran</h2>
			<form on:submit={onSubmitCreate}>
				<!-- Form fields -->
				<label for="Siswa">Siswa : </label>
				<select id="Siswa" bind:value={SiswaID}>
					{#each DataSiswa as siswa}
						<option value={siswa.ID}>{siswa.Nama}</option>
					{/each}
				</select>

				<label for="Guru">Guru : </label>
				<select id="Guru" bind:value={GuruID}>
					{#each DataGuru as guru}
						<option value={guru.ID}>{guru.Nama}</option>
					{/each}
				</select>

				<label for="Tahun">Tahun Ajaran : </label>
				<select id="Tahun" bind:value={TahunID}>
					{#each DataTahun as tahun}
						<option value={tahun.ID}>{tahun.Tahun}</option>
					{/each}
				</select>

				<label for="Jurusan">Jurusan : </label>
				<select id="Jurusan" bind:value={JurusanID}>
					{#each DataJurusan as jurusan}
						<option value={jurusan.ID}>{jurusan.Nama}</option>
					{/each}
				</select>

				<label for="TanggalPendaftaran">Tanggal Pendaftaran : </label>
				<input
					type="date"
					id="TanggalPendaftaran"
					bind:value={TanggalPendaftaran}
					class="mb-2 border p-2"
				/>

				<label for="TipePendaftaran">Tipe Pendaftaran : </label>
				<select name="TipePendaftaran" id="TipePendaftaran" bind:value={TipePendaftaran}>
					<option value="zonasi">Zonasi</option>
					<option value="prestasi">Prestasi</option>
				</select>

				<button type="submit" class="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-700"
					>Simpan</button
				>
				<button
					type="button"
					on:click={closeModalCreate}
					class="ml-2 rounded bg-white px-4 py-2 text-white hover:bg-gray-700">Batal</button
				>
			</form>
			{#if showError}
				<p class="mt-2 text-red-500">{errorMessage}</p>
				<button
					on:click={closeError}
					class="mt-2 rounded bg-white px-4 py-2 text-white hover:bg-gray-700">Tutup</button
				>
			{/if}
		</div>
	</div>
{/if}

<!-- Modal Edit -->
{#if showModalEdit}
	<div class="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50">
		<div class="rounded bg-white p-4 shadow-lg">
			<h2 class="mb-4 text-2xl font-bold">Edit Pendaftaran</h2>
			<form on:submit={onUpdate}>
				<!-- Form fields -->
				<label for="Siswa">Siswa : </label>
				<select id="Siswa" bind:value={SiswaID}>
					{#each DataSiswa as siswa}
						<option value={siswa.ID}>{siswa.Nama}</option>
					{/each}
				</select>

				<label for="Guru">Guru : </label>
				<select id="Guru" bind:value={GuruID}>
					{#each DataGuru as guru}
						<option value={guru.ID}>{guru.Nama}</option>
					{/each}
				</select>

				<label for="Tahun">Tahun Ajaran : </label>
				<select id="Tahun" bind:value={TahunID}>
					{#each DataTahun as tahun}
						<option value={tahun.ID}>{tahun.Tahun}</option>
					{/each}
				</select>

				<label for="Jurusan">Jurusan : </label>
				<select id="Jurusan" bind:value={JurusanID}>
					{#each DataJurusan as jurusan}
						<option value={jurusan.ID}>{jurusan.Nama}</option>
					{/each}
				</select>

				<label for="TanggalPendaftaran">Tanggal Pendaftaran : </label>
				<input
					type="date"
					id="TanggalPendaftaran"
					bind:value={TanggalPendaftaran}
					class="mb-2 border p-2"
				/>

				<label for="TipePendaftaran">Tipe Pendaftaran : </label>
				<input
					type="text"
					id="TipePendaftaran"
					bind:value={TipePendaftaran}
					class="mb-2 border p-2"
				/>

				<button type="submit" class="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-700"
					>Simpan</button
				>
				<button
					type="button"
					on:click={closeModalEdit}
					class="ml-2 rounded bg-white px-4 py-2 text-white hover:bg-gray-700">Batal</button
				>
			</form>
		</div>
	</div>
{/if}

<!-- Confirmation Popup for Delete -->
{#if showPopUp}
<DeletePopUp onClose={closePopUp} onDelete={deletePendaftaran}/>

{/if}

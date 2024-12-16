<script lang="ts">
	import { onMount } from 'svelte';
	import axios from 'axios';
	import type { AsalSekolah, Siswa } from '../../function/interface';
	import { goto } from '$app/navigation';
	import DeletePopUp from '../../components/deletePopUp.svelte';

	let siswa: Siswa[] = [];
	let SMP: AsalSekolah[] = [];
	let searchQuery = '';
	let page = 1;
	let limit = 5;
	let totalPages = 1;
	let showModalCreate = false;
	let showModalEdit = false;
	let showPopUp = false;
	let errorMessage = '';
	let showError = false;
	let editId: number | null = null;
	let deleteId: number | null = null;
	let isLoading = false;

	let Nama = '';
	let NIK = 0;
	let NISN = 0;
	let AsalSekolahID:number | null = null;
	let Nilai = 0;
	let JarakTempuh = 0;
	let Alamat = '';
	let TempatLahir = '';
	let TanggalLahir: Date | null = null;
	let Gender = '';
	let FotoURL: File | null = null;

	async function fetchData(page: number, limit: number, query: string) {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get('http://localhost:3030/siswa', {
				params: { page, limit, query },
				headers: { Authorization: `Bearer ${token}` }
			});
			siswa = response.data.data;
			totalPages = response.data.totalPages;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			siswa = [];
		} finally {
			isLoading = false;
		}
	}

	onMount(async () => {
		fetchData(page, limit, searchQuery);
		fetchDataAsalSekolah();
		
		const asalSekolahIDs = siswa.map(student => student.AsalSekolahID);
		if (asalSekolahIDs.length > 0) {
			for (const id of asalSekolahIDs) {
				await fetchEditDataSmp(id);
			}
		}
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

		if (!Nama) {
			showError = true;
			errorMessage = 'Please fill nama fields.';
			return;
		} else if (NIK === 0) {
			showError = true;
			errorMessage = 'Please fill NIK fields.';
			return;
		} else if (NISN === 0) {
			showError = true;
			errorMessage = 'Please fill NISN fields.';
			return;
		} else if (!AsalSekolahID) {
			showError = true;
			errorMessage = 'Please fill AsalSekolah fields.';
			return;
		} else if (Nilai === 0) {
			showError = true;
			errorMessage = 'Please fill Nilai fields.';
			return;
		} else if (JarakTempuh === 0) {
			showError = true;
			errorMessage = 'Please fill Jarak temppuh fields.';
			return;
		} else if (!Alamat) {
			showError = true;
			errorMessage = 'Please fill alamat fields.';
			return;
		} else if (!TempatLahir) {
			showError = true;
			errorMessage = 'Please fill TempatLahir fields.';
			return;
		} else if (!TanggalLahir) {
			showError = true;
			errorMessage = 'Please fill tanggal fields.';
			return;
		} else if (!Gender) {
			showError = true;
			errorMessage = 'Please fill gender fields.';
			return;
		}

		// Create FormData object
		const formData = new FormData();
		formData.append('Nama', Nama);
		formData.append('NIK', NIK.toString());
		formData.append('NISN', NISN.toString());
		formData.append('AsalSekolahID', AsalSekolahID.toString());
		formData.append('Nilai', Nilai.toString());
		formData.append('JarakTempuh', JarakTempuh.toString());
		formData.append('Alamat', Alamat);
		formData.append('TempatLahir', TempatLahir);
		formData.append('Gender', Gender);

		if (TanggalLahir) {
			const formattedDate = new Date(TanggalLahir).toISOString(); // Format to ISO 8601
			formData.append('TanggalLahir', formattedDate);
		}

		if (FotoURL) formData.append('foto', FotoURL);

		try {
			const token = sessionStorage.getItem('token');
			await axios.post('http://localhost:3030/siswa', formData, {
				headers: {
					Authorization: `Bearer ${token}`,
					'Content-Type': 'multipart/form-data'
				}
			});
			closeModalCreate();
			console.log({
				Nama,
				NIK,
				NISN,
				AsalSekolahID,
				Nilai,
				JarakTempuh,
				Alamat,
				TempatLahir,
				TanggalLahir: formData.get('TanggalLahir'), // Log formatted date
				Gender,
				FotoURL
			});

			fetchData(page, limit, searchQuery);
		} catch (err: any) {
			showError = true;
			errorMessage = err.response?.data.message || 'Failed to submit data';
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
		Nama = '';
		NIK = 0;
		NISN = 0;
		AsalSekolahID = null;
		Nilai = 0;
		JarakTempuh = 0;
		Alamat = '';
		TempatLahir = '';
		TanggalLahir = null;
		Gender = '';
		FotoURL = null;
	}

	function modalEdit(id: number) {
		editId = id;
		showModalEdit = true;
		fetchEditData(id);
	}

	async function fetchEditData(id: number) {
		try {
			const token = sessionStorage.getItem('token');
			const response = await axios.get(`http://localhost:3030/siswa/${id}`, {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data: Siswa = response.data;
			Nama = data.Nama;
			NIK = data.NIK;
			NISN = data.NISN;
			AsalSekolahID = data.AsalSekolahID;
			Nilai = data.Nilai;
			JarakTempuh = data.JarakTempuh;
			Alamat = data.Alamat;
			TempatLahir = data.TempatLahir;
			TanggalLahir = new Date(data.TanggalLahir);
			Gender = data.Gender;
		} catch (err) {
			console.log(err);
		}
	}

	async function onUpdate(e: SubmitEvent) {
		e.preventDefault();
		const formData = new FormData();
		formData.append('Nama', Nama);
		formData.append('NIK', NIK.toString());
		formData.append('NISN', NISN.toString());
		if (AsalSekolahID !== null) formData.append('AsalSekolahID', AsalSekolahID.toString());
		formData.append('Nilai', Nilai.toString());
		formData.append('JarakTempuh', JarakTempuh.toString());
		formData.append('Alamat', Alamat);
		formData.append('TempatLahir', TempatLahir);
		if (TanggalLahir) {
			formData.append('TanggalLahir', TanggalLahir.toISOString());
		}
		formData.append('Gender', Gender);
		if (FotoURL) formData.append('foto', FotoURL);

		try {
			const token = sessionStorage.getItem('token');
			await axios.put(`http://localhost:3030/siswa/${editId}`, formData, {
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

	async function deleteSiswa() {
		if (deleteId !== null) {
			try {
				const token = sessionStorage.getItem('token');
				await axios.delete(`http://localhost:3030/siswa/${deleteId}`, {
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

	function handleFileUpload(event: Event) {
		const input = event.target as HTMLInputElement;
		if (input.files?.length) {
			FotoURL = input.files[0];
		}
	}

	async function fetchDataAsalSekolah() {
		try {
			const token = sessionStorage.getItem('token');
			isLoading = true;
			const response = await axios.get('http://localhost:3030/smp', {
				headers: { Authorization: `Bearer ${token}` }
			});
			SMP = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			goto('/auth/login');
			SMP = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchEditDataSmp(id: number) {
		try {
			const token = sessionStorage.getItem('token');
			const response = await axios.get(`http://localhost:3030/smp/${id}`, {
				headers: { Authorization: `Bearer ${token}` }
			});
			const data: AsalSekolah = response.data;
			Nama = data.Nama;
			console.log(Nama);
		} catch (err) {
			console.log(err);
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
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">ID</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nama</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">NISN</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">NIK</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Asal Sekolah</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Nilai</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Jarak Tempuh</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Alamat</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Tempat Lahir</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Tanggal Lahir</th>
					<th class="px-6 py-3 text-left text-xs font-medium text-white uppercase tracking-wider">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each siswa as item}
					<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
						<td class="px-6 py-4 text-white">{item.ID}</td>
						<td class="px-6 py-4 text-white">{item.Nama}</td>
						<td class="px-6 py-4 text-white">{item.NISN}</td>
						<td class="px-6 py-4 text-white">{item.NIK}</td>
						<td class="px-6 py-4 text-white">
							{#if item.AsalSekolahID}
								{SMP.find(smp => smp.ID === item.AsalSekolahID)?.Nama || 'Unknown'}
							{:else}
								No Asal Sekolah
							{/if}
						</td>
						<td class="px-6 text-white py-4">{item.Nilai}</td>
						<td class="px-6 text-white py-4">{item.JarakTempuh}</td>
						<td class="px-6 text-white py-4">{item.Alamat}</td>
						<td class="px-6 text-white py-4">{item.TempatLahir}</td>
						<td class="px-6 text-white py-4">
							{#if item.TanggalLahir}
								{new Date(item.TanggalLahir).toLocaleDateString()}
							{:else}
								No Date
							{/if}
						</td>
						<td class="px-6 text-white py-4 gap-2 flex flex-row">
							<button
								on:click={() => modalEdit(item.ID)}
								class="flex items-center justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
								>Edit</button>
							<button	
								on:click={() => popUpDelete(item.ID)}
								class="flex items-center ml-2 justify-center rounded-md bg-white px-4 py-2 font-medium text-slate-500"
								>Hapus</button>
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
			<h2 class="mb-4 text-2xl font-bold">Tambah Siswa</h2>
			<form on:submit={onSubmitCreate}>
				<!-- Form fields -->
				<label for="Nama">Nama : </label>
				<input type="text" placeholder="Nama" bind:value={Nama} class="mb-2 border p-2" />
				<label for="NIK">NIK :</label>
				<input type="number" placeholder="NIK" bind:value={NIK} class="mb-2 border p-2" />
				<label for="NISN">NISN :</label>
				<input type="number" placeholder="NISN" bind:value={NISN} class="mb-2 border p-2" />
				<label for="AsalSekolahID">Asal Sekolah</label>
				<select name="AsalSekolahID" id="AsalSekolahID" bind:value={AsalSekolahID}>
					{#each SMP as sekolah}
						<option value={sekolah.ID}>{sekolah.Nama}</option>
					{/each}
				</select>
				<label for="Nilai">Nilai : </label>
				<input type="number" placeholder="Nilai" bind:value={Nilai} class="mb-2 border p-2" />
				<label for="JarakTempuh">Jarak Tempuh : </label>
				<input
					type="number"
					placeholder="Jarak Tempuh"
					bind:value={JarakTempuh}
					class="mb-2 border p-2"
				/>
				<label for="Alamat">Alamat : </label>
				<input type="text" placeholder="Alamat" bind:value={Alamat} class="mb-2 border p-2" />
				<label for="TempatLahir">Tempat Lahir:</label>
				<input
					type="text"
					placeholder="Tempat Lahir"
					bind:value={TempatLahir}
					class="mb-2 border p-2"
				/>
				<label for="TanggalLahir">Tanggal Lahir</label>
				<input
					type="date"
					placeholder="Tanggal Lahir"
					bind:value={TanggalLahir}
					class="mb-2 border p-2"
				/>
				<label for="Gender">Gender : </label>
				<select name="Gender" id="Gender" bind:value={Gender}>
					<option value="laki-laki">Laki-Laki</option>
					<option value="perempuan">Perempuan</option>
				</select>
				<label for="FotoUrl">Foto : </label>
				<input type="file" accept="image/*" on:change={handleFileUpload} class="mb-2 border p-2" />
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
			<h2 class="mb-4 text-2xl font-bold">Edit Siswa</h2>
			<form on:submit={onUpdate}>
				<!-- Form fields -->
				<input type="text" placeholder="Nama" bind:value={Nama} class="mb-2 border p-2" />
				<input type="number" placeholder="NIK" bind:value={NIK} class="mb-2 border p-2" />
				<input type="number" placeholder="NISN" bind:value={NISN} class="mb-2 border p-2" />
				<select name="AsalSekolah" id="AsalSekolah">
					{#each SMP as sekolah}
						<option value={sekolah.ID}>{sekolah.Nama}</option>
					{/each}
				</select>
				<input type="number" placeholder="Nilai" bind:value={Nilai} class="mb-2 border p-2" />
				<input
					type="number"
					placeholder="Jarak Tempuh"
					bind:value={JarakTempuh}
					class="mb-2 border p-2"
				/>
				<input type="text" placeholder="Alamat" bind:value={Alamat} class="mb-2 border p-2" />
				<input
					type="text"
					placeholder="Tempat Lahir"
					bind:value={TempatLahir}
					class="mb-2 border p-2"
				/>
				<input
					type="date"
					placeholder="Tanggal Lahir"
					bind:value={TanggalLahir}
					class="mb-2 border p-2"
				/>
				<input type="text" placeholder="Gender" bind:value={Gender} class="mb-2 border p-2" />
				<input type="file" accept="image/*" on:change={handleFileUpload} class="mb-2 border p-2" />
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
	<DeletePopUp onClose={closePopUp} onDelete={deleteSiswa}/>
{/if}

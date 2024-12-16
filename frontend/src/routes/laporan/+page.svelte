<script lang="ts">
	import { onMount } from 'svelte';
	import axios from 'axios';
	import type { Guru, Jurusan, Pendaftaran, Siswa, TahunAjaran } from '../../function/interface';
	import { goto } from '$app/navigation';
	import * as XLSX from 'xlsx';
	import jsPDF from 'jspdf';
	import 'jspdf-autotable';

	let pendaftaran: Pendaftaran[] = [];
	let DataSiswa: Siswa[] = [];
	let DataGuru: Guru[] = [];
	let DataTahun: TahunAjaran[] = [];
	let DataJurusan: Jurusan[] = [];
	let searchQuery = '';
	let selectedTipe = '';
	let startDate: string = '';
	let endDate: string = '';
	let page = 1;
	let limit = 10;
	let totalPages = 1;
	let totalData = 0;
	let isLoading = false;

	async function fetchData(
		page: number,
		limit: number,
		query: string,
		tipe: string,
		startDate: string,
		endDate: string
	) {
		try {
			const token = sessionStorage.getItem('token');
			if (!token) {
				goto('/auth/login');
			}
			isLoading = true;
			const response = await axios.get('http://localhost:3030/pendaftaran', {
				params: { page, limit, query, tipe, startDate, endDate },
				headers: { Authorization: `Bearer ${token}` }
			});
			pendaftaran = response.data.data;
			totalPages = response.data.totalPages;
			totalData = response.data.total
		} catch (err) {
			console.log('Error:', err);
			pendaftaran = [];
		} finally {
			isLoading = false;
		}
	}

	onMount(async () => {
		fetchData(page, limit, searchQuery, selectedTipe, startDate, endDate);
		await fetchDataSiswa();
		await fetchDataGuru();
		await fetchDataJurusan();
		await fetchDataTahun();
	});

	function handleFilterChange(event: Event) {
		selectedTipe = (event.target as HTMLSelectElement).value;
		page = 1;
		fetchData(page, limit, searchQuery, selectedTipe, startDate, endDate);
	}

	function prevPage() {
		if (page > 1) {
			page--;
			fetchData(page, limit, searchQuery, selectedTipe, startDate, endDate);
		}
	}

	function nextPage() {
		if (page < totalPages) {
			page++;
			fetchData(page, limit, searchQuery, selectedTipe, startDate, endDate);
		}
	}

	function handleSearch(event: Event) {
		event.preventDefault();
		const input = (event.target as HTMLFormElement).querySelector(
			'input[type="text"]'
		) as HTMLInputElement;
		searchQuery = input.value;
		page = 1;
		fetchData(page, limit, searchQuery, selectedTipe, startDate, endDate);
	}

	async function fetchDataSiswa() {
		try {
			const token = sessionStorage.getItem('token');
			if (!token) {
				goto('/auth/login');
			}
			isLoading = true;
			const response = await axios.get('http://localhost:3030/siswa', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataSiswa = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			DataSiswa = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchDataGuru() {
		try {
			const token = sessionStorage.getItem('token');
			if (!token) {
				goto('/auth/login');
			}
			isLoading = true;
			const response = await axios.get('http://localhost:3030/guru', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataGuru = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			DataGuru = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchDataTahun() {
		try {
			const token = sessionStorage.getItem('token');
			if (!token) {
				goto('/auth/login');
			}
			isLoading = true;
			const response = await axios.get('http://localhost:3030/tahun', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataTahun = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			DataTahun = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchDataJurusan() {
		try {
			const token = sessionStorage.getItem('token');
			if (!token) {
				goto('/auth/login');
			}
			isLoading = true;
			const response = await axios.get('http://localhost:3030/jurusan', {
				headers: { Authorization: `Bearer ${token}` }
			});
			DataJurusan = response.data.data;
		} catch (err) {
			console.log('Error:', err);
			DataJurusan = [];
		} finally {
			isLoading = false;
		}
	}

	function exportToExcel() {
		const data = pendaftaran.map((item) => ({
			ID: item.ID,
			Siswa: DataSiswa.find((s) => s.ID === item.SiswaID)?.Nama || 'unknown',
			Guru: DataGuru.find((g) => g.ID === item.GuruID)?.Nama || 'unknown',
			Jurusan: DataJurusan.find((j) => j.ID === item.JurusanID)?.Nama || 'unknown',
			Nilai: DataSiswa.find((s) => s.ID === item.SiswaID)?.Nilai || 'unknown',
			JarakTempuh: DataSiswa.find((s) => s.ID === item.SiswaID)?.JarakTempuh + 'km' || 'unknown',
			Tahun: DataTahun.find((t) => t.ID === item.TahunID)?.Tahun || 'Unknown',
			TipePendaftaran: item.TipePendaftaran,
			TanggalPendaftaran: new Date(item.TanggalPendaftaran).toLocaleDateString()
		}));

		const worksheet = XLSX.utils.json_to_sheet(data);
		const workbook = XLSX.utils.book_new();
		XLSX.utils.book_append_sheet(workbook, worksheet, 'Pendaftaran');

		XLSX.writeFile(workbook, 'pendaftaran.xlsx');
	}

	function exportToExcelByID(id: number) {
		const item = pendaftaran.find((p) => p.ID === id);
		if (!item) {
			console.error('Data not found');
			return;
		}

		const data = [
			{
				ID: item.ID,
				Siswa: DataSiswa.find((s) => s.ID === item.SiswaID)?.Nama || 'unknown',
				Guru: DataGuru.find((g) => g.ID === item.GuruID)?.Nama || 'unknown',
				Jurusan: DataJurusan.find((j) => j.ID === item.JurusanID)?.Nama || 'unknown',
				Nilai: DataSiswa.find((s) => s.ID === item.SiswaID)?.Nilai || 'unknown',
				JarakTempuh: DataSiswa.find((s) => s.ID === item.SiswaID)?.JarakTempuh + 'km' || 'unknown',
				Tahun: DataTahun.find((t) => t.ID === item.TahunID)?.Tahun || 'Unknown',
				TipePendaftaran: item.TipePendaftaran,
				TanggalPendaftaran: new Date(item.TanggalPendaftaran).toLocaleDateString()
			}
		];

		const worksheet = XLSX.utils.json_to_sheet(data);
		const workbook = XLSX.utils.book_new();
		XLSX.utils.book_append_sheet(workbook, worksheet, 'Pendaftaran');

		XLSX.writeFile(workbook, 'pendaftaran.xlsx');
	}

	function exportToPDF() {
		const doc = new jsPDF();
		const headers = [
			'ID',
			'Siswa',
			'Guru',
			'Jurusan',
			'Nilai',
			'Jarak Tempuh',
			'Tahun',
			'Tipe Pendaftaran',
			'Tanggal Pendaftaran'
		];

		const data = pendaftaran.map((item) => [
			item.ID,
			DataSiswa.find((s) => s.ID === item.SiswaID)?.Nama || 'unknown',
			DataGuru.find((g) => g.ID === item.GuruID)?.Nama || 'unknown',
			DataJurusan.find((j) => j.ID === item.JurusanID)?.Nama || 'unknown',
			DataSiswa.find((s) => s.ID === item.SiswaID)?.Nilai || 'unknown',
			DataSiswa.find((s) => s.ID === item.SiswaID)?.JarakTempuh || 'unknown',
			DataTahun.find((t) => t.ID === item.TahunID)?.Tahun || 'Unknown',
			item.TipePendaftaran,
			new Date(item.TanggalPendaftaran).toLocaleDateString()
		]);

		doc.autoTable({
			head: [headers],
			body: data
		});

		doc.save('pendaftaran.pdf');
	}
</script>

<h1 class="mb-2 text-3xl font-medium">Data Siswa Ajaran</h1>
<div class="mb-4 flex flex-row items-center">
	<button
		on:click={exportToExcel}
		class="ml-2 rounded bg-green-500 px-4 py-2 text-white hover:bg-green-700"
	>
		Ekspor ke Excel
	</button>
	<button
		on:click={exportToPDF}
		class="ml-2 rounded bg-red-500 px-4 py-2 text-white hover:bg-red-700"
	>
		Ekspor ke PDF
	</button>
	<form on:submit={handleSearch} class="flex flex-row items-center">
		<input
			type="text"
			placeholder="Cari Siswa..."
			class="rounded border px-4 py-2"
			value={searchQuery}
		/>
		<input
		type="date"
		placeholder="Start Date"
		class="ml-2 rounded border px-4 py-2"
		bind:value={startDate}
		/>
		<input
		type="date"
		placeholder="End Date"
		class="ml-2 rounded border px-4 py-2"
		bind:value={endDate}
		/>
		<button
			type="submit"
			class="ml-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
		>
			Cari
		</button>
	</form>

	<select on:change={handleFilterChange} class="ml-2 rounded border px-4 py-2">
		<option value="">Semua</option>
		<option value="zonasi">Zonasi</option>
		<option value="prestasi">Prestasi</option>
	</select>
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
					{#if selectedTipe !== 'zonasi'}
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
							>Nilai</th
						>
					{/if}
					{#if selectedTipe !== 'prestasi'}
						<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
							>Jarak Tempuh</th
						>
					{/if}
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Tahun Ajaran</th
					>
					<th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-white"
						>Tipe Pendaftaran</th
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
						{#if selectedTipe !== 'zonasi'}
							<td class="px-6 py-4 text-white">
								{#if item.SiswaID}
									{DataSiswa.find((murid) => murid.ID === item.SiswaID)?.Nilai || 'unknown'}
								{:else}
									not name
								{/if}
							</td>
						{/if}
						{#if selectedTipe !== 'prestasi'}
							<td class="px-6 py-4 text-white">
								{#if item.SiswaID}
									{DataSiswa.find((murid) => murid.ID === item.SiswaID)?.JarakTempuh || 'unknown'} km
								{:else}
									not name
								{/if}
							</td>
						{/if}
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
								on:click={() => exportToExcelByID(item.ID)}
								class="ml-2 rounded bg-green-500 px-4 py-2 text-white hover:bg-green-700"
							>
								Ekspor ke Excel
							</button>
							<button
								on:click={() => goto(`/laporan/${item.ID}`)}
								class="ml-2 rounded bg-green-500 px-4 py-2 text-white hover:bg-green-700"
							>
								Detail
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<div class="mt-4 flex justify-center">
		{#if totalData > limit}
		<button
		class="focus:shadow-outline mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
		on:click={prevPage}
		disabled={page === 1}>Prev</button
	>
		{/if}
		{#if totalData > limit}
		<button
			class="focus:shadow-outline mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700 focus:outline-none"
			on:click={nextPage}
			disabled={page === totalPages}>Next</button
		>
		{/if}
	</div>
{/if}

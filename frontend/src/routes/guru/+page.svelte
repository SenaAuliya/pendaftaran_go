<script lang="ts">
    import { onMount } from 'svelte';
    import axios from 'axios';
    import type { Guru } from '../../function/interface';
    import { goto } from '$app/navigation';
    import DeletePopUp from '../../components/deletePopUp.svelte';

    let guru: Guru[] = [];
    let searchQuery = '';
    let page = 1;
    let limit = 5;
    let totalPages = 1;
    let showModal = false;
    let showModalMode: 'create' | 'edit' = 'create';
    let Nama = '';
    let KodeGuru = '';
    let showError = false;
    let errorMessage = '';
    let editId: number | null = null;
    let showPopUp = false;
    let deleteId: number | null = null;
    let isLoading = false;

    async function fetchData(page: number, limit: number, query: string) {
        try {
            const token = sessionStorage.getItem('token');
            isLoading = true;
            const response = await axios.get('http://localhost:3030/guru', {
                params: { page, limit, query },
                headers: { Authorization: `Bearer ${token}` }
            });
            guru = response.data.data;
            totalPages = response.data.totalPages;
        } catch (err) {
            console.log('Error:', err);
            goto('/auth/login');
            guru = [];
        } finally {
            isLoading = false;
        }
    }

    onMount(() => fetchData(page, limit, searchQuery));

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

    function openModal(mode: 'create' | 'edit', id?: number) {
        showModalMode = mode;
        if (mode === 'edit' && id) {
            editId = id;
            fetchEditData(id);
        } else {
            Nama = '';
            KodeGuru = '';
            editId = null;
        }
        showModal = true;
    }

    async function fetchEditData(id: number) {
        try {
            const token = sessionStorage.getItem('token');
            const response = await axios.get(`http://localhost:3030/guru/${id}`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            const data: Guru = response.data;
            Nama = data.Nama;
            KodeGuru = data.KodeGuru;
        } catch (err) {
            console.log(err);
        }
    }

    async function submitForm(e: SubmitEvent) {
        e.preventDefault();
        try {
            const token = sessionStorage.getItem('token');
            if (showModalMode === 'create') {
                await axios.post('http://localhost:3030/guru', { Nama, KodeGuru }, {
                    headers: { Authorization: `Bearer ${token}` }
                });
            } else if (showModalMode === 'edit' && editId) {
                await axios.put(`http://localhost:3030/guru/${editId}`, { Nama, KodeGuru }, {
                    headers: { Authorization: `Bearer ${token}` }
                });
            }
            closeModal();
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

    function closeModal() {
        showModal = false;
    }

    function openDeletePopUp(id: number) {
        deleteId = id;
        showPopUp = true;
    }

    async function deleteGuru() {
        if (deleteId !== null) {
            try {
                const token = sessionStorage.getItem('token');
                await axios.delete(`http://localhost:3030/guru/${deleteId}`, {
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

<h1 class="mb-2 text-3xl font-medium">Data Guru</h1>
<div class="mb-4 flex items-center">
    <button
        on:click={() => openModal('create')}
        class="mr-4 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
    >
        Tambah Guru
    </button>
    <form on:submit={handleSearch} class="flex items-center">
        <input
            type="text"
            placeholder="Cari guru..."
            class="rounded border px-4 py-2"
            bind:value={searchQuery}
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
        <table class="w-full text-left text-sm text-gray-500 dark:text-gray-400">
            <thead class="bg-gray-50 text-xs uppercase text-gray-700 dark:bg-gray-700 dark:text-gray-400">
                <tr>
                    <th class="px-6 py-3">ID</th>
                    <th class="px-6 py-3">Nama</th>
                    <th class="px-6 py-3">Kode Guru</th>
                    <th class="px-6 py-3">Aksi</th>
                </tr>
            </thead>
            <tbody>
                {#if guru.length > 0}
                    {#each guru as item}
                        <tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
                            <th
                                scope="row"
                                class="whitespace-nowrap px-6 py-4 font-medium text-gray-900 dark:text-white"
                            >
                                {item.ID}
                            </th>
                            <td class="px-6 py-4">{item.Nama}</td>
                            <td class="px-6 py-4">{item.KodeGuru}</td>
                            <td class="flex gap-4 px-6 py-4">
                                <button
                                    class="rounded bg-white px-4 py-2 text-slate-500 hover:bg-gray-200"
                                    on:click={() => openModal('edit', item.ID)}
                                >
                                    Edit
                                </button>
                                <button
                                    class="rounded bg-white px-4 py-2 text-slate-500 hover:bg-gray-200"
                                    on:click={() => openDeletePopUp(item.ID)}
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
            class="mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
            on:click={prevPage}
            disabled={page === 1}
        >
            Prev
        </button>
        <button
            class="mx-2 rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
            on:click={nextPage}
            disabled={page === totalPages}
        >
            Next
        </button>
    </div>
{/if}

{#if showPopUp}
    <DeletePopUp onClose={closePopUp} onDelete={deleteGuru} />
{/if}

{#if showModal}
    <div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
        <div class="bg-white p-4 rounded shadow-lg">
            <form on:submit={submitForm}>
                <label for="Nama" class="block text-sm font-bold text-gray-700">Nama:</label>
                <input
                    type="text"
                    id="Nama"
                    bind:value={Nama}
                    class="mb-2 w-full rounded border px-3 py-2"
                />
                <label for="KodeGuru" class="block text-sm font-bold text-gray-700">Kode Guru:</label>
                <input
                    type="text"
                    id="KodeGuru"
                    bind:value={KodeGuru}
                    class="mb-2 w-full rounded border px-3 py-2"
                />
                <div class="flex justify-end mt-4">
                    <button
                        type="submit"
                        class="rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
                    >
                        {showModalMode === 'create' ? 'Kirim' : 'Update'}
                    </button>
                    <button
                        type="button"
                        class="ml-2 rounded bg-gray-500 px-4 py-2 font-bold text-white hover:bg-gray-700"
                        on:click={closeModal}
                    >
                        Batal
                    </button>
                </div>
            </form>
            {#if showError}
                <div class="mt-4 text-red-500">{errorMessage}</div>
                <button
                    class="mt-2 rounded bg-red-500 px-4 py-2 text-white hover:bg-red-700"
                    on:click={closeError}
                >
                    Tutup
                </button>
            {/if}
        </div>
    </div>
{/if}

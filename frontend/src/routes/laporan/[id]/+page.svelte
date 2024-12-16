<script lang="ts">
    import { onMount } from 'svelte';
    import axios from 'axios';
    import type { Guru, Jurusan, Pendaftaran, Siswa, TahunAjaran } from '../../../function/interface';
    import { goto } from '$app/navigation';
    import { page } from '$app/stores';
  
    let pendaftaran: Pendaftaran | null = null;
    let DataSiswa: Siswa[] = [];
    let DataGuru: Guru[] = [];
    let DataJurusan: Jurusan[] = [];
    let DataTahun: TahunAjaran[] = [];
  
    $: id = $page.params.id;
    let isLoading = false;
  
    async function fetchData() {
      try {
        const token = sessionStorage.getItem('token');
        if(!token){
            goto('/auth/login');
        }
        isLoading = true;
        const response = await axios.get(`http://localhost:3030/pendaftaran/${id}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        console.log('Pendaftaran Data:', response.data);
        pendaftaran = response.data;
      } catch (err) {
        console.log('Error:', err);
        pendaftaran = null;
      } finally {
        isLoading = false;
      }
    }
  
    onMount(async () => {
      await Promise.all([fetchData(),
    fetchDataSiswa(),
    fetchDataGuru(),
    fetchDataJurusan(),
    fetchDataTahun(),]) 
    });
  
    async function fetchDataSiswa() {
      try {
        const token = sessionStorage.getItem('token');
        if(!token){
            goto('/auth/login');
        }
        isLoading = true;
        const response = await axios.get('http://localhost:3030/siswa', {
          headers: { Authorization: `Bearer ${token}` }
        });
        console.log('Data Siswa:', response.data);
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
        if(!token){
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
        if(!token){
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
        if(!token){
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
    $: siswa =  DataSiswa.find((s) => s.ID === pendaftaran.SiswaID)
    $: guru = DataGuru.find((g) => g.ID === pendaftaran.GuruID)
    $: jurusan = DataJurusan.find((j) => j.ID === pendaftaran.JurusanID)
    $: tahun = DataTahun.find((t) => t.ID === pendaftaran.TahunID)

  </script>
  
  {#if pendaftaran}
    <h1 class="text-2xl font-bold">Detail Pendaftaran</h1>
    <p><strong>ID:</strong> {pendaftaran.ID}</p>
    <p><strong>Nama Siswa:</strong> 
      {#if DataSiswa.length > 0}
        {#if siswa}
          <img src={`http://localhost:3030/${siswa?.FotoURL}`} alt="Foto Siswa" class="w-32 h-32 object-cover" />
        {:else}
          No Foto Available
        {/if}
      {:else}
        Loading Siswa Data...
      {/if}
    </p>
    <p><strong>Nama Guru:</strong> 
        {guru?.Nama || 'unknown'}
    </p>
    <p><strong>Nama Jurusan:</strong> 
        {jurusan?.Nama || 'unknown'}
    </p>
    <p><strong>Tahun Ajaran:</strong> 
      {#if DataTahun.length > 0}
        {tahun?.Tahun || 'Unknown'}
      {:else}
        Loading Tahun Data...
      {/if}
    </p>
  {:else}
    <p>Loading...</p>
  {/if}
  
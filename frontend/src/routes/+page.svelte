<script lang="ts">
	import { goto } from "$app/navigation";
	import axios from "axios";
	import { onMount } from "svelte";
	import Card from "../components/card.svelte";
    import { cardItem } from "../function/var";

	async function checkAuthorizes() {
		try {
			const token = sessionStorage.getItem('token');
			const response = await axios.get("http://localhost:3030/", {
				headers: { Authorization: `Bearer ${token}` }
			});
		} catch (err) {
			goto("/auth/login");
		}
	}

	onMount(() => {
		checkAuthorizes();
	});
</script>

<div class="grid grid-cols-3 items-center gap-4">
	{#each cardItem as item}
		<Card color={item.color} title={item.title} link={item.link} desc={item.desc} />
	{/each}
</div>

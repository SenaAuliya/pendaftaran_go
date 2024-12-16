<script lang="ts">
    import Chevron from '$lib/assets/chevron-right-solid.svg'; // Chevron Right icon
    import ChevronDown from '$lib/assets/chevron-down-solid.svg'; // Chevron Down icon

    export let title: string;
    let isOpen = false;

    function toggle() {
        isOpen = !isOpen;
    }
</script>

<div class="w-full flex flex-col items-center">
    <button class="p-4 text-left mx-10 font-medium text-lg flex items-center text-white" on:click={toggle}>
        {title}
        {#if isOpen}
        <img src={ChevronDown} class="w-[1rem] ml-3" alt="Chevron Icon">
        {:else}
        <img src={Chevron} class="w-[0.59rem] ml-3" alt="Chevron Icon">
        {/if}
    </button>
    {#if isOpen}
        <div class="p-4">
            <slot></slot>
        </div>
    {/if}
</div>

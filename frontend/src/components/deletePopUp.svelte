<script lang="ts">
    export let onClose: () => void;
    export let onDelete: () => void;
</script>

<div class="fixed left-0 right-0 top-0 z-50 flex items-center justify-center overflow-y-auto overflow-x-hidden h-full w-full bg-gray-900 bg-opacity-50">
    <div class="relative w-full max-w-md p-4 bg-white rounded-lg shadow dark:bg-gray-700">
        <button
            type="button"
            class="absolute top-3 right-3 text-gray-400 hover:bg-gray-200 rounded-lg p-2"
            on:click={onClose}
        >
            <svg class="h-5 w-5" fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            <span class="sr-only">Close modal</span>
        </button>
        <div class="text-center">
            <svg class="mx-auto mb-4 h-12 w-12 text-gray-400" fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                <path stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M10 11V6m0 8h.01M19 10a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>
            <h3 class="mb-5 text-lg font-normal text-gray-500">Are you sure you want to delete this item?</h3>
            <button
                type="button"
                class="inline-flex items-center rounded-lg bg-red-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-red-800"
                on:click={() => {
                    onDelete();
                    onClose();
                }}
            >
                Yes, I'm sure
            </button>
            <button
                type="button"
                class="ms-3 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100"
                on:click={onClose}
            >
                No, cancel
            </button>
        </div>
    </div>
</div>

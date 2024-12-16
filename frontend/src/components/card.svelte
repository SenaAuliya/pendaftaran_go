<script lang="ts">
    import Arrow from "$lib/assets/arrow-right-solid.svg";
    export let color: string;
    export let title: string;
    export let desc: string;
    export let link: string;

    $: bgClass = color === 'red' ? 'bg-red-600' 
                 : color === 'blue' ? 'bg-blue-600' 
                 : color === 'green' ? 'bg-green-600' 
                 : color === 'gray' ? 'bg-gray-400'
                 : color === 'violet' ? 'bg-violet-400'               
                 : 'bg-purple-400'; // Default fallback
</script>

<div class={`${bgClass} text-white flex flex-col justify-center items-center p-4 w-[20rem] h-[10rem] gap-2`}>
    <h1 class="text-xl font-semibold text-center m-3">{title}</h1> 
    <p>{desc}</p>
    <a class="hover:text-gray-200 flex items-center justify-center mt-2" href={link}>Ke Data {title} <img src={Arrow} class="w-[1rem] ml-2" alt=""></a>
</div>

<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen flex-col overflow-hidden bg-gray-100">
	<header>
		<nav class="flex h-24 shrink-0 items-center px-6 shadow-md"></nav>
	</header>
	<main class="min-h-0 flex-1 pt-8">
		{@render children()}
	</main>

	<footer class="h-20"></footer>
</div>

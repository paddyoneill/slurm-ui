<script lang="ts">
	import type { Snippet } from 'svelte';
	let dialog = $state() as HTMLDialogElement;
	let {
		showModal = $bindable(),
		title,
		children
	}: {
		showModal: boolean;
		title: string;
		children: Snippet;
	} = $props();

	$effect(() => {
		if (showModal) dialog.showModal();
	});
</script>

<dialog
	class="m-auto w-full max-w-md rounded-lg bg-gray-100 p-6 shadow-lg backdrop:bg-black/60"
	bind:this={dialog}
	onclose={() => (showModal = false)}
	onclick={(e) => {
		if (e.target === dialog) dialog.close();
	}}
>
	<button
		type="button"
		onclick={() => dialog.close()}
		class="absolute top-2 right-2 p-2 text-xl font-semibold text-gray-400 hover:text-gray-600"
		>x</button
	>
	<h2 class="mb-4 p-6 pb-0 text-lg font-semibold">{title}</h2>
	{@render children()}
</dialog>

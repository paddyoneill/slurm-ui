<script lang="ts">
	import PageIntro from '$lib/components/PageIntro.svelte';
	import OverviewCard from '$lib/components/OverviewCard.svelte';
</script>

<div class="page">
	<PageIntro
		title="Overview"
		description="High level overview of current Slurm jobs and Jupyter notebooks"
	/>
	<section class="summary-grid" aria-label="Overviews summaries">
		<OverviewCard href="/jobs" label="Jobs" value="..." copy="placeholder" linkLabel="Open Jobs" />
		<OverviewCard
			href="/notebooks"
			label="Notebooks"
			value="..."
			copy="placeholder"
			linkLabel="Open Notebooks"
		/>
	</section>
</div>

<style>
	.page {
		display: grid;
		gap: var(--space-6);
		width: min(100%, 68rem);
		margin: 0 auto;
	}

	.summary-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--space-6);
	}
</style>

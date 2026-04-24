<script lang="ts">
	let {
		title,
		description
	}: {
		title: string;
		description: string;
	} = $props();
</script>

<div class="intro">
	<div class="copy">
		<h1>{title}</h1>
		<p>{description}</p>
	</div>
</div>

<style>
	.intro {
		position: relative;
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-4);
		padding: var(--space-6);
		border: 1px solid var(--colour-border);
		border-radius: var(--radius-lg);
		background: var(--colour-surface);
		box-shadow: var(--shadow-md);
		overflow: hidden;
	}

	.intro::before {
		content: '';
		position: absolute;
		inset: 0;
		background:
			linear-gradient(135deg, rgb(var(--colour-primary-rgb) / 0.08), transparent 42%),
			linear-gradient(180deg, rgb(255 255 255 / 0.75), transparent 38%);
		pointer-events: none;
	}

	.copy {
		position: relative;
		display: grid;
		gap: var(--space-2);
		z-index: 1;
	}

	h1,
	p {
		margin: 0;
	}

	h1 {
		font-size: clamp(1.75rem, 3vw, 2.35rem);
		line-height: 1.1;
	}

	p {
		max-width: 42rem;
		color: var(--colour-text-muted);
	}
</style>

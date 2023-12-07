<script>
	export let icon = null
	export let stickyactions = false

	import { Icon } from "$lib/c"
</script>

<style>
	.stickyactions {
		position: sticky;
		top: 0;
		z-index: 200;
	}
</style>

<div class="is-pulled-right" class:stickyactions>
	<slot name="actions" />
</div>
<h1 class="title">
	{#if icon}
		<Icon {icon} class="is-size-4 {$$props.iclass??""}"><slot /></Icon>
	{:else}
		<slot />
	{/if}
</h1>

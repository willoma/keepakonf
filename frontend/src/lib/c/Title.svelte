<script>
	export let icon = null

	import { Icon } from "$lib/c"
</script>

<div class="is-pulled-right">
	<slot name="actions" />
</div>
<h1 class="title">
	{#if icon}
		<Icon {icon} class="is-size-4 {$$props.iclass??""}"><slot /></Icon>
	{:else}
		<slot />
	{/if}
</h1>

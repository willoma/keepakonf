<script>
	export let href
	export let icon = null

	import { Icon } from "$lib/c"
</script>

<a class="button {$$props.class??""}" {href}>
	{#if icon}
		<Icon {icon} />
		{#if $$slots.default}
			<span><slot /></span>
		{/if}
	{:else}
		<slot />
	{/if}
</a>

<script>
	export let icon = null
	export let disabled = false
	export let submit = false

	import { Icon } from "$lib/c"

	$: type = submit ? "submit" : "button"
</script>

<button
	class="button {$$props.class??""}"
	{type}
	{disabled}
	on:mousedown
	on:click|stopPropagation
>
	{#if icon}
		<Icon {icon} />
		{#if $$slots.default}
			<span><slot /></span>
		{/if}
	{:else}
		<slot />
	{/if}
</button>

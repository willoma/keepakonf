<!-- 
	props:

	- class: class on the icon span
	- tclass: class on the icon-text span (if slot)
 -->
<script>
	export let icon
	export let baseclass = "icon"

	import { icon as iconFn } from "$lib/icons"
</script>

<style>
	.icon > :global(svg) {
		height: 1em;
		width: 1em;
	}
</style>

{#if $$slots.default}
	<span class="icon-text {$$props.tclass??""}">
		<span class="{baseclass} {$$props.class??""}">
			{@html iconFn(icon)}
		</span>
		<span><slot /></span>
	</span>
{:else}
	<span class="{baseclass} {$$props.class??""}">
		{@html iconFn(icon)}
	</span>
{/if}

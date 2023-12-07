<!-- 
	props:

	- class: class on the icon span
	- tclass: class on the icon-text span (if slot)
 -->
<script>
	export let icon
	export let baseclass = "icon"
	export let ellipsis = false

	import { icon as iconFn } from "$lib/icons"
</script>

<style>
	.icon > :global(svg) {
		height: 1em;
		width: 1em;
	}

	.ellipsis {
		text-overflow: ellipsis;
		overflow: hidden;
	}
</style>

{#if $$slots.default}
	<span class="icon-text is-flex-wrap-nowrap {$$props.tclass??""}">
		<span class="{baseclass} {$$props.class??""}">
			{@html iconFn(icon)}
		</span>
		<span class="text" class:ellipsis><slot /></span>
	</span>
{:else}
	<span class="{baseclass} {$$props.class??""}">
		{@html iconFn(icon)}
	</span>
{/if}

<script>
	export let field
	export let label

	import { Field, Icon } from "$lib/c"
	import { randomID } from "$lib/random"

	const id = randomID()
</script>

<Field {field} {id} {label} horizontal>
	<button type="button" class="button" class:is-info={$field.value} on:click={() => $field.value = !$field.value}>
		{#if $field.value}
			<Icon icon="check" />
		{:else}
			<span class="icon" />
		{/if}
	</button>
</Field>

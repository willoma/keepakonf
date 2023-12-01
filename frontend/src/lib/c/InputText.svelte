<script>
	export let field
	export let label = null
	export let placeholder = label
	export let icon = null

	import { Icon } from "$lib/c"
	import { randomID } from "$lib/random"

	const id = randomID()
</script>

<div class="field">
	{#if label}
		<label class="label" for={id}>{label}</label>
	{/if}
	<div class="control" class:has-icons-left={icon}>
		<input
			class="input"
			class:is-danger={!$field.valid}
			{id}
			type="text"
			{placeholder}
			required
			bind:value={$field.value}
		/>
		{#if icon}<Icon {icon} class="is-left" />{/if}
	</div>
</div>

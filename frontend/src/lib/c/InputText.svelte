<script>
	export let field
	export let label = null
	export let placeholder = label
	export let icon = null

	import { Icon } from "$lib/c"
	import { errors } from "$lib/fielderrors"
	import { randomID } from "$lib/random"

	const id = randomID()
</script>

<style>
	.field {
		position: relative;
	}
	.error-message {
		position: absolute;
		z-index: 100;
		bottom: -0.75em;
		right: 1em;
		pointer-events: none;
	}
</style>

<div class="field">
	{#if label}
		<label class="label" for={id}>{label}</label>
	{/if}
	{#if !$field.valid}
		<div class="error-message tag is-danger">{errors($field.errors)}</div>
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

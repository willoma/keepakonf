<script>
	export let field
	export let id = ""
	export let label = null
	export let horizontal = false
	export let addons = false
	export let grouped = false

	import { errors } from "$lib/fielderrors"
</script>

<style>
	.error-message {
		position: absolute;
		z-index: 100;
		bottom: -1em;
		left: 1em;
		pointer-events: none;
	}
</style>

{#if horizontal}
	<div class="field is-horizontal">
		<div class="field-label is-normal">
			{#if label}
				<label class="label" for={id}>{label}</label>
			{/if}
		</div>
		<div class="field-body is-relative">
			{#if !$field.valid}
				<div class="error-message tag is-danger">{errors($field.errors)}</div>
			{/if}
			<div class="field" class:has-addons={addons} class:is-grouped={grouped} class:is-grouped-multiline={grouped}>
				<slot />
			</div>
		</div>
	</div>
{:else}
	<div class="field is-relative" class:has-addons={addons} class:is-grouped={grouped} class:is-grouped-multiline={grouped}>
		{#if label}
			<label class="label" for={id}>{label}</label>
		{/if}
		{#if !$field.valid}
			<div class="error-message tag is-danger">{errors($field.errors)}</div>
		{/if}
		<slot />
	</div>
{/if}

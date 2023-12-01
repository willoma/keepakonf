<script>
	export let field
	export let label = null

	import { icon } from "$lib/icons"
	import { randomID } from "$lib/random"

	const id = randomID()
</script>

<style lang="scss">
	@import "bulma/sass/utilities/derived-variables";

	input[type="checkbox"] {
		display: none;
	}
	.icon-text {
		align-items: center;
	}
	.icon {
		height: 2em;
		width: 2em;
		margin-right: 0.5em !important;
		cursor: pointer;

		:global(svg) {
			fill: $grey-light;
		}

		&.checked :global(svg) {
			fill: $success;
		}
	}
</style>

<div class="field is-grouped">
	<label class="control">
		<input type="checkbox" {id} bind:checked={$field.value} />
		<span class="icon-text">
			<span class="icon" class:checked={$field.value}>
				{@html icon($field.value?"checkboxyes":"checkboxno")}
			</span>
			{#if label}<span>{label}</span>{/if}
		</span>
	</label>
</div>

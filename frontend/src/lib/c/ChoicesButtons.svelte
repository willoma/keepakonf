<script>
	export let value
	export let options
</script>

<div class="field has-addons">
	{#each options as option}
		<div class="control is-expanded">
			<button
				class="button is-fullwidth is-small"
				class:is-primary={value===option[0]}
				type="button"
				on:click|stopPropagation={() => value=option[0]}
			>
				{option[1]}
			</button>
		</div>
	{/each}
</div>

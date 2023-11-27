<script>
	export let value
	export let label = null

	import { Button } from "$lib/c"

	function doRemove(i) {
		value.splice(i, 1)
		value = value
	}
</script>

{#if label}
	<div class="field">
		<!-- svelte-ignore a11y-label-has-associated-control -->
		<label class="label">{label}</label>
	</div>
{/if}

<div class="field is-grouped is-grouped-multiline">
	{#each value as row, i}
		<div class="control">
			<div class="field has-addons">
				<div class="control">
					<input class="input" type="text" bind:value={row} />
				</div>
				<div class="control">
					<Button class="is-info" icon="remove" on:click={() => doRemove(i)} />
				</div>
			</div>
		</div>
	{/each}
	<div class="control">
		<Button icon="add" on:click={() => value = [...value, ""]}>Add</Button>
	</div>
</div>

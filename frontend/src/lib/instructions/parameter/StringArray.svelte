<script>
	export let field
	export let label

	import { Button, Field } from "$lib/c"
	import { randomID } from "$lib/random"

	const id = randomID()
	
	function doRemove(i) {
		$field.value.splice(i, 1)
		$field.value = $field.value
	}
</script>

<Field {field} {id} {label} horizontal grouped>
	{#each $field.value as row, i}
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
		<Button icon="add" on:click={() => $field.value = [...$field.value, ""]}>Add</Button>
	</div>
</Field>
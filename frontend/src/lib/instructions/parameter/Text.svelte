<script>
	export let field
	export let label

	import { Field } from "$lib/c"
	import { randomID } from "$lib/random"

	const id = randomID()
</script>

<Field {field} {id} {label} horizontal>
	<div class="control is-expanded">
		<textarea
			class="textarea"
			class:is-danger={!$field.valid}
			{id}
			required
			bind:value={$field.value}
		/>
	</div>
</Field>
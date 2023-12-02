<script>
	export let group = null
	export let valid

	import { field } from 'svelte-forms'
	import { required } from 'svelte-forms/validators'

	import { Field, Icon } from "$lib/c"

	import EditInstructions from "./EditInstructions.svelte"

	$: name = field('name', group?.name ?? "", [required()], { "checkOnInit": true})
	
	let instructions
	let instructionsValid

	export function makeData() {
		const data = {
			name: $name.value,
			instructions: instructions.makeData(),
		}
		if (group?.id) {
			data.id = group.id
		}
		return data
	}
	$: valid = $name.valid && instructionsValid
</script>

<Field field={name}>
	<div class="control has-icons-left">
		<input
			type="text"
			class="input"
			class:is-danger={!$name.valid}
			placeholder="Group name"
			required
			bind:value={$name.value}
		/>
		<Icon icon={group?.icon??"group"} class="is-left" />
	</div>
</Field>

<EditInstructions bind:this={instructions} bind:valid={instructionsValid} initial={group?.instructions ?? []} />

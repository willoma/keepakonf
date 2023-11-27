<script>
	export let group = null
	export let valid

	import { field } from 'svelte-forms'
	import { required } from 'svelte-forms/validators'

	import { InputText } from "$lib/c"

	import EditInstructions from "./EditInstructions.svelte"

	$: id = field('id', group?.id)
	$: name = field('name', group?.name ?? "", [required()], { "checkOnInit": true})
	
	let instructions
	let instructionsValid

	export function makeData() {
		return {
		id: group?.id,
		name: $name.value,
		instructions: instructions.makeData(),
		}
	}
	$: valid = $name.valid && instructionsValid
</script>

<InputText icon={group?.icon??"group"} placeholder="Group name" bind:value={$name.value} />

<EditInstructions bind:this={instructions} bind:valid={instructionsValid} initial={group?.instructions ?? []} />

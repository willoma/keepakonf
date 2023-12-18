<script>
	export let initial
	export let valid

	import EditInstruction from "./EditInstruction.svelte"
    import InstructionAdder from "./InstructionAdder.svelte"
	import VariableInserter from "./VariableInserter.svelte"

	let instructions = []

	let allComponents = []
	let allValid = []

	function initialLoad() {
		instructions = [...initial]
	}
	$: initialLoad(initial)

	function addInstruction(data) {
		switch (data.type) {
			case "if":
				instructions = [...instructions, {"type": "if"}]
				break;
			case "command":
			default:
				instructions = [...instructions, {"type": "command", "command": data.params.name}]
				break;
		}
	}

	function removeInstruction(index) {
		instructions.splice(index, 1)
		allComponents.splice(index, 1)
		allValid.splice(index, 1)
		instructions = instructions
		allComponents = allComponents
		allValid = allValid
	}
	export function makeData() {
		return allComponents.map(comp => comp?.makeData()).filter(Boolean)
	}

	$: valid = !allValid.includes(false)
</script>

<VariableInserter>
	{#each instructions as instruction, i}
		<EditInstruction
			initial={instruction}
			bind:this={allComponents[i]}
			bind:valid={allValid[i]}
			on:remove={() => removeInstruction(i)}
		/>
	{/each}
</VariableInserter>

<InstructionAdder on:add={(evt) => addInstruction(evt.detail)} />

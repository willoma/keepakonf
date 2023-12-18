<script>
	export let initial
	export let valid


	import EditInstructionCommand from "./EditInstructionCommand.svelte"
	import EditInstructionIf from "./EditInstructionIf.svelte";

	let component

	export function makeData() {
		return component.makeData()
	}
</script>

{#if initial.type === "command" || !initial.type}
	<EditInstructionCommand bind:this={component} {initial} {valid} />
{:else if initial.type === "if"}
	<EditInstructionIf bind:this={component} {initial} {valid} />
{/if}

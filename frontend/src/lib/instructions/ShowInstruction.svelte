<script>
	export let instruction
	export let highlight

	import ShowInstructionCommand from "./ShowInstructionCommand.svelte";

</script>

{#if instruction?.type == "command" || !instruction?.type}
	<ShowInstructionCommand {instruction} {highlight} />
{:else}
	<div>
		{JSON.stringify(instruction)}
	</div>
{/if}

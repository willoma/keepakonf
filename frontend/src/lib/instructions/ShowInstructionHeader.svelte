<script>
	export let instruction
	export let cmd

	import { Button, Icon } from "$lib/c"
	import { applyInstruction } from "$lib/instructions"
</script>

<div class="is-flex is-justify-content-space-between is-align-items-center" style:width="100%">
	<div class="icon-text py-1">
		<Icon icon={cmd?.icon??"command"}>
			<b>{instruction.command}</b>: {instruction.info ?? "Unknown"}
		</Icon>
	</div>
	{#if instruction.status === "todo"}
		<Button class="is-small is-info" icon="run" on:click={() => applyInstruction(instruction.id)}>Apply</Button>
	{:else if instruction.status === "running"}
		<span class="icon loader mr-1"></span>
	{:else if instruction.status === "failed"}
		<Button class="is-small" icon="run" on:click={() => applyInstruction(instruction.id)}>Try again</Button>
	{/if}
</div>

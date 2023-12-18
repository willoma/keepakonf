<script>
	export let initial
	export let valid

	import { createEventDispatcher } from "svelte"

	import { Button, Icon } from "$lib/c"

    import EditInstructions from "./EditInstructions.svelte"
	import { Equal } from "./condition"

	const dispatch = createEventDispatcher()

	const operations = [
		["=", "Equality", "equal"],
	]

	let operation = initial?.condition?.op ?? "="
	let condition
	let instructions
	let instructionsValid

	$: opData = operations.find(d => d[0] === operation) ?? ["", "None", "unknown"]

	let showOperationDropdown = false

	export function makeData() {
		return {
			"type": "if",
			"cond": condition.makeData(),
			"instructions": instructions.makeData(),
		}
	}
	$: valid = instructionsValid
</script>

<div class="box p-2">
	<div class="is-flex mb-3">
		<Icon icon="condition" tclass="is-flex-grow-1">
			<b>if</b>: Execute commands conditionally
		</Icon>
		<Button class="is-warning is-small is-flex-grow-0" icon="remove" on:click={() => dispatch("remove")}>Remove</Button>
	</div>

	<div class="block is-flex">
		<div class="dropdown is-flex-grow-0 mr-3" class:is-active={showOperationDropdown}>
			<div class="dropdown-trigger">
				<button class="button" type="button" on:click={() => showOperationDropdown = !showOperationDropdown}>
					<Icon icon={opData[2]} />
					<span>{opData[1]}</span>
				</button>
			</div>
			<div class="dropdown-menu">
				<div class="dropdown-content">
					{#each operations as op}
						<a
							href="#top"
							class="dropdown-item"
							on:click|preventDefault={() => {
								operation=op[0]
								showOperationDropdown = false
							}}
						>
							<Icon icon={op[2]} />
							{op[1]}
						</a>
					{/each}
				</div>
			</div>
		</div>
		<div class="is-flex-grow-0 mr-3">
			<input
				type="text"
				class="input has-text-centered is-static"
				value="Execute instructions if..."
				readonly
			/>
		</div>
		<div class="is-flex-grow-1">
			{#if operation === "="}
				<Equal bind:this={condition} initial={initial?.cond ?? {}} />
			{/if}
		</div>
	</div>

	<EditInstructions bind:this={instructions} bind:valid={instructionsValid} initial={initial?.instructions ?? []} />
</div>
<script>
	export let initial
	export let valid

	// import { field, combined } from 'svelte-forms'

	import { Button, Buttons, Icon } from "$lib/c"
	import { commands } from "$lib/store"

	import EditInstruction from "./EditInstruction.svelte"

	let instructions = []

	let allComponents = []
	let allValid = []

	function initialLoad() {
		instructions = [...initial]
	}
	$: initialLoad(initial)

	let showCommandChoice = false

	function addInstruction(name) {
		instructions = [...instructions, {"command": name}]
		showCommandChoice = false
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

{#each instructions as instruction, i}
	<EditInstruction
		initial={instruction}
		bind:this={allComponents[i]}
		bind:valid={allValid[i]}
		on:remove={() => removeInstruction(i)}
	/>
{/each}

<Buttons right>
	<Button class="is-primary" icon="add" on:click={() => showCommandChoice = true}>
		Add instruction
	</Button>
</Buttons>

{#if showCommandChoice}
	<div class="modal is-active">
		<div class="modal-background" role="none" on:click={() => showCommandChoice = false}></div>
		<div class="modal-content">
			<div class="box p-0">
				<div class="menu">
					<p class="menu-label pt-3 pl-5">
						Select command
					</p>
					<ul class="menu-list">
						{#each $commands as cmd}
							<li>
								<a href="#top" on:click|preventDefault={addInstruction(cmd.name)}>
									<Icon icon={cmd?.icon??"command"}>
										<b>{cmd?.name ?? "Unknown"}</b>: {cmd?.description ?? "Unknown"}
									</Icon>
								</a>
							</li>
						{/each}
					</ul>
					<p class="menu-label has-text-right mt-1 pb-1 pr-1">
						<Button icon="cancel" class="is-small" on:click={() => showCommandChoice = false}>
							Cancel
						</Button>
					</p>
				</div>
			</div>
		</div>
	</div>
{/if}

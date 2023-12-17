<script>
	export let group = null
	export let valid

	import { field } from 'svelte-forms'
	import { required } from 'svelte-forms/validators'

	import { Field, Icon } from "$lib/c"
	import { icons } from "$lib/icons"

	import EditInstructions from "./EditInstructions.svelte"

	$: name = field('name', group?.name ?? "", [required()], { "checkOnInit": true})
	$: icon = field('icon', group?.icon ?? "group", [required()], { "checkOnInit": true})

	let instructions
	let instructionsValid

	let showIconDropdown

	export function makeData() {
		const data = {
			name: $name.value,
			icon: $icon.value,
			instructions: instructions.makeData(),
		}
		if (group?.id) {
			data.id = group.id
		}
		return data
	}
	$: valid = $name.valid && instructionsValid
</script>

<Field field={name} addons>
	<div class="control">
		<div class="dropdown" class:is-active={showIconDropdown}>
			<div class="dropdown-trigger">
				<button class="button" type="button" on:click={() => showIconDropdown = !showIconDropdown}>
					<Icon icon={$icon.value} />
					<span>Icon</span>
				</button>
			</div>
			<div class="dropdown-menu">
				<div class="dropdown-content">
					{#each Object.entries(icons) as ic}
						<a
							href="#top"
							class="dropdown-item"
							on:click|preventDefault={() => {
								$icon.value = ic[0]
								showIconDropdown = false
							}}
						>
							<Icon icon={ic[0]} />
							{ic[0]}
					</a>
					{/each}
				</div>
			</div>
		</div>
	</div>
	<div class="control">
		<input
			type="text"
			class="input"
			class:is-danger={!$name.valid}
			placeholder="Group name"
			required
			bind:value={$name.value}
		/>
	</div>
</Field>

<EditInstructions bind:this={instructions} bind:valid={instructionsValid} initial={group?.instructions ?? []} />

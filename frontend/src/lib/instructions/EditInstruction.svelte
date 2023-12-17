<script>
	export let initial
	export let valid

	import { combined, field } from "svelte-forms"
	import { createEventDispatcher } from "svelte"

	import { Button, Icon } from "$lib/c"
	import { command } from "$lib/store"
	import { makefield } from "$lib/instructions"

	import EditParameter from './EditParameter.svelte'

	const dispatch = createEventDispatcher()

	$: cmd = $command(initial.command)

	let parametersFields = []
	let parametersCombined = field("parameters", "")
	let ready = false

	function makeParametersFields() {
		if (parametersFields.length || !cmd?.parameters) {
			return
		}
		parametersFields = cmd.parameters.map((param) => makefield(param, initial?.parameters?.[param.id]))
		parametersCombined = combined("parameters", parametersFields, ($fields) => Object.fromEntries($fields.map((f) => [f.name, f.value])))
		ready = true
	}

	$: makeParametersFields(cmd, initial)

	export function makeData() {
		return {
			"command": initial.command,
			"parameters": $parametersCombined.value,
		}
	}

	$:valid = $parametersCombined.valid
</script>

<div class="box p-2">
	<Button class="is-warning is-small is-pulled-right" icon="remove" on:click={() => dispatch("remove")}>Remove</Button>
	<Icon icon={cmd?.icon??"command"} tclass="block mb-1">
		<b>{initial.command}</b>: {cmd?.description ?? "Unknown"}
	</Icon>
	{#if ready}
		{#each cmd?.parameters ?? [] as param, i}
			<EditParameter {param} field={parametersFields[i]} />
		{/each}
	{/if}
</div>

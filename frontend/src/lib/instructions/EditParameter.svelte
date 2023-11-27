<script>
	export let param
	export let field

	import { InputText, InputTextArray, Toggle } from "$lib/c"

	$: label = param?.title ?? "Unknown"
</script>

{#if param.type === "bool"}
	<Toggle bind:value={$field.value} {label} />
{:else if param.type === "string"}
	<InputText bind:value={$field.value} {label} />
{:else if param.type === "[string]"}
	<InputTextArray bind:value={$field.value} {label} />
{:else}
	<pre>{JSON.stringify(param)}</pre>
{/if}

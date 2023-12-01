<script>
	export let param
	export let field

	import { InputText, InputTextArray, Toggle } from "$lib/c"

	$: label = param?.title ?? "Unknown"
</script>

{#if param.type === "bool"}
	<Toggle {field} {label} />
{:else if param.type === "filepath"}
	<InputText {field} {label} />
{:else if param.type === "string"}
	<InputText {field} {label} />
{:else if param.type === "[string]"}
	<InputTextArray {field} {label} />
{:else if param.type === "username"}
	<InputText {field} {label} />
{:else}
	<pre>{JSON.stringify(param)}</pre>
{/if}

<script>
	export let param
	export let field

	import { Bool, Filepath, String, StringArray, Text, Username } from "./parameter"

	$: label = param?.title ?? "Unknown"
</script>

{#if param.type === "bool"}
	<Bool {field} {label} />
{:else if param.type === "filepath"}
	<Filepath {field} {label} />
{:else if param.type === "string"}
	<String {field} {label} />
{:else if param.type === "[string]"}
	<StringArray {field} {label} />
{:else if param.type === "text"}
	<Text {field} {label} />
{:else if param.type === "username"}
	<Username {field} {label} />
{:else}
	<pre>{JSON.stringify(param)}</pre>
{/if}

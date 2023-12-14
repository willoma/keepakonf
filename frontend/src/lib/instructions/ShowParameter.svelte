<script>
	export let param
	export let value

	$: label = param?.title ?? "Unknown"
</script>

{#if param.type === "bool"}
	<b>{label}</b> : {value?"yes":"no"}
{:else if param.type === "filepath"}
	<b>{label}</b>: {value}
{:else if param.type === "string"}
	<b>{label}</b>: {value}
{:else if param.type === "[string]"}
<b>{label}</b>: {value.join(", ")}
{:else if param.type === "text"}
	<b>{label}</b>: {#if value.length > 33}{value.substring(0,30)}[...]{:else}{value}{/if}
{:else if param.type === "username"}
	<b>{label}</b>: {value}
{:else}
	{JSON.stringify(param)}
{/if}

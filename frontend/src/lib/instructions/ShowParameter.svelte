<script>
	export let param
	export let value

	$: label = param?.title ?? "Unknown"
</script>

{#if param.type === "string"}
	<b>{label}</b>: {value}
{:else if param.type === "[string]"}
	<b>{label}</b>: {value.join(", ")}
{:else if param.type === "bool"}
	<b>{label}</b> : {value?"yes":"no"}
{:else}
	{JSON.stringify(param)}
{/if}

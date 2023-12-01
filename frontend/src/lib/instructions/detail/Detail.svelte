<script>
	export let data

	import Error from "./Error.svelte"
	import Raw from "./Raw.svelte"
	import Table from "./Table.svelte"
	import Terminal from "./Terminal.svelte"
</script>

{#if data.t === "error"}
	<Error data={data.d} />
{:else if data.t === "table"}
	<Table data={data.d} />
{:else if data.t === "terminal"}
	<Terminal data={data.d} />
{:else}
	<Raw data={data.d} />
{/if}

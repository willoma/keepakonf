<script>
	export let data

	import Raw from "./Raw.svelte"
	import Table from "./Table.svelte"
	import Terminal from "./Terminal.svelte"
</script>

{#if data.t === "table"}
	<Table data={data.d} />
{:else if data.t === "terminal"}
	<Terminal data={data.d} />
{:else}
	<Raw data={data.d} />
{/if}

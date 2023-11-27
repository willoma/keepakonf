<script>
	export let data

	import { statuscolor } from "$lib/color"
</script>

<table>
	{#if data.h}
		<tr>
			{#each data.h as cell}
				<th>{cell}</th>
			{/each}
		</tr>
	{/if}
	{#each data.r as row}
		<tr>
			{#each row as cell}
				<td class="has-text-{statuscolor[cell.s]}-dark">{cell.c}</td>
			{/each}
		</tr>
	{/each}
</table>

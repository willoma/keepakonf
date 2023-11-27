<script>
	import { page } from '$app/stores'

	import { Button, ButtonA, Buttons, Title } from "$lib/c"
	import { applyGroup } from "$lib/instructions"
	import { group } from "$lib/store"
	import { statuscolordark } from "$lib/color"

	import ShowInstruction from "$lib/instructions/ShowInstruction.svelte"

	$: grp = $group($page.params.id)

	$: highlightedInstruction = $page.url.hash.startsWith("#") ? $page.url.hash.substring(1) : null
</script>

<svelte:head>
	<title>Keepakonf - {grp?.name ?? "Group not found"}</title>
</svelte:head>

{#if grp}
	<Title icon={grp?.icon??"group"} iclass="has-text-{statuscolordark[grp?.status]}">
		{grp?.name}
		<Buttons slot="actions">
			{#if grp?.status === "todo"}
				<Button class="is-info" icon="run" on:click={() => applyGroup($page.params.id)}>Apply</Button>
			{:else if grp?.status === "running"}
				<span class="icon loader mr-1"></span>
			{:else if grp?.status === "failed"}
				<Button icon="run" on:click={() => applyGroup($page.params.id)}>Try again</Button>
			{/if}

			<ButtonA class="is-primary" icon="edit" href="/{$page.params.id}/modify">Modify</ButtonA>
		</Buttons>
	</Title>

	{#each grp?.instructions ?? [] as instruction}
		<ShowInstruction {instruction} highlight={highlightedInstruction === instruction.id} />
	{/each}
{:else}
	<Title icon="group">
		Group not found
		<Buttons slot="actions">
			<ButtonA class="is-primary" icon="cancel" href="/">Back to summary</ButtonA>
		</Buttons>
	</Title>
{/if}
<script>
	import { derived } from "svelte/store"

	import { page } from "$app/stores"

	import { Button, ButtonA, Icon } from "$lib/c"
	import { applyGroup } from "$lib/instructions"
	import { groups } from "$lib/store"
	import { statuscolor, statuscolordark } from "$lib/color"

	const todo = derived(groups, ($gs) => $gs?.filter(g => g.status === "todo" || g.status === "->todo"))
	const failed = derived(groups, ($gs) => $gs?.filter(g => g.status === "failed"))
	const running = derived(groups, ($gs) => $gs?.filter(g => g.status === "running"))
	const applied = derived(groups, ($gs) => $gs?.filter(g => g.status === "applied" || g.status === "->applied"))

	const current = derived(
		[page, groups, todo, failed, running, applied],
		([$page, $groups, $todo, $failed, $running, applied]) => {
			const currentHash = $page.url.hash.startsWith("#") ? $page.url.hash.substring(1) : ""
			switch (currentHash) {
				case "todo":
					return ["todo", $todo]
				case "failed":
					return ["failed", $failed]
				case "running":
					return ["running", $running]
				case "applied":
					return ["applied", applied]
				case "all":
					return ["all", $groups]
				default:
					return ["todo", $todo]
			}
		}
	)
</script>

<style lang="scss">
	@import "bulma/sass/components/panel";

	.panel-block:hover {
			background-color: $panel-block-hover-background-color;
	}
</style>

<div class="panel is-{statuscolor[$current[0]]}">
	<p class="panel-heading">
		<Icon icon="group">
			Groups
		</Icon>
	</p>
	<p class="panel-tabs">
		<a class:is-active={$current[0]==="todo"} href="/#todo">Todo ({$todo.length})</a>
		<a class:is-active={$current[0]==="failed"} href="/#failed">Failed ({$failed.length})</a>
		<a class:is-active={$current[0]==="running"} href="/#running">Running ({$running.length})</a>
		<a class:is-active={$current[0]==="applied"} href="/#applied">Applied ({$applied.length})</a>
		<a class:is-active={$current[0]==="all"} href="/#all">All ({$groups.length})</a>
	</p>
	{#each $current[1] as group}
		<div class="panel-block p-0 is-clickable is-flex is-align-items-center">
			<a href="/{group.id}" class="is-flex-grow-1 has-text-black pl-3 py-3">
				<Icon
					icon={group?.icon??"group"}
					class="has-text-{statuscolordark[group.status]}"
				>
					{group.name}
				</Icon>
			</a>
			<div class="buttons is-flex-grow-0 is-flex-shrink-0 py-2 pr-3">
				{#if group.status === "todo"}
					<Button class="is-small is-info" icon="run" on:click={() => applyGroup(group.id)}>Apply</Button>
				{:else if group.status === "failed"}
					<Button class="is-small" icon="run" on:click={() => applyGroup(group.id)}>Try again</Button>
				{/if}
				<ButtonA class="is-small is-primary" href="/{group.id}/modify?back={$current[0]}" icon="edit">Modify</ButtonA>
			</div>
		</div>
	{:else}
		<div class="panel-block py-3">
			<Icon icon="group">
				No group
			</Icon>
		</div>
	{/each}
</div>
<script>
	import { ButtonA, Title } from "$lib/c"

    import GroupsPanel from "./GroupsPanel.svelte"
    import LogsPanel from "./LogsPanel.svelte"
</script>

<style lang="scss">
	@import "bulma/sass/components/panel";

	.panel-block:hover {
			background-color: $panel-block-hover-background-color;
	}
</style>

<svelte:head>
	<title>Keepakonf - Summary</title>
</svelte:head>


<Title icon="search">
	Summary
	<ButtonA slot="actions" class="is-info" icon="add" href="/new">
		Create a new group
	</ButtonA>
</Title>

<div class="columns">
	<div class="column is-half">
		<GroupsPanel />
	</div>
	<div class="column is-half">
		<LogsPanel />
	</div>
</div>

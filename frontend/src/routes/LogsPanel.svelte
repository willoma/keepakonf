<script>
	import { Icon } from "$lib/c"
	import { logs, logsReachedTheEnd, socket } from "$lib/store"

	import LogsPanelEntry from "./LogsPanelEntry.svelte"

	function loadMore() {
		socket.emit("logs", {"offset":$logs.length}, (response) => {
			$logs = [...$logs, ...response.logs?.toReversed()]
			if (response.reached_the_end) {
				$logsReachedTheEnd = true
			}
		})
	}
</script>

<div class="panel">
	<p class="panel-heading">
		<Icon icon="logs">
			Logs
		</Icon>
	</p>
	{#each $logs as log}
		<LogsPanelEntry {log} />
	{/each}
	{#if !$logsReachedTheEnd}
		<div class="panel-block is-clickable" role="button" tabindex="0" on:click={loadMore} on:keyup={(e) => e.key === "Enter" && loadMore}>
			Load more...
		</div>
	{/if}
</div>
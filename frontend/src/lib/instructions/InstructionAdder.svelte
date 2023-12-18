<script>
	import { createEventDispatcher } from "svelte"

	import { Button, Icon } from "$lib/c"
	import { commands } from "$lib/store"

	const dispatch = createEventDispatcher()

	function add(name) {
		dispatch("add", name)
		show = false
	}

	let show = false
	let section = "command"
	let filter = ""
</script>

<style lang="scss">
	@import "bulma/sass/components/panel";

	.panel-block a:hover {
			background-color: $panel-block-hover-background-color;
	}
</style>

<Button class="is-primary" icon="add" on:click={() => show = true}>
	Add instruction
</Button>

{#if show}
	<div class="modal is-active">
		<div class="modal-background" role="none" on:click={() => show = false}></div>
		<div class="modal-content">
			<div class="box p-0">
				<div class="panel">
					<div class="panel-heading is-flex">
						<div class="is-flex-grow-1">
							Select instruction
						</div>
						<div class="is-flex-grow-0">
							<Button icon="cancel" class="is-small" on:click={() => show = false}>
								Cancel
							</Button>
						</div>
					</div>
					<div class="panel-block">
					  <p class="control has-icons-left">
						<input type="text" class="input" placeholder="Search" bind:value={filter} />
						<Icon icon="search" class="is-left" />
					  </p>
					</div>
					<p class="panel-tabs">
						<a class:is-active={section === "command"} href="#top" on:click|preventDefault={() => section="command"}>Command</a>
					</p>
					{#each $commands as cmd}
						{#if !filter.length || cmd?.name?.includes(filter) || cmd?.description?.includes(filter)}
							<a href="#top" class="panel-block has-text-black" on:click|preventDefault={() => add(cmd.name)}>
								<div class="columns w-100">
									<div class="column is-one-third">
										<Icon icon={cmd?.icon??"command"}>
											<b>{cmd?.name ?? "Unknown"}</b>
										</Icon>
									</div>
									<div class="column is-two-thirds">
										{cmd?.description ?? "Unknown"}
									</div>
								</div>
							</a>
						{/if}
					{/each}
				</div>
			</div>
		</div>
	</div>
{/if}

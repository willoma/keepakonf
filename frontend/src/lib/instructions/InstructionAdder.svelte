<script>
	import { createEventDispatcher } from "svelte"
	import { fade } from 'svelte/transition';

	import { Button, Icon } from "$lib/c"
	import { commands } from "$lib/store"

	const dispatch = createEventDispatcher()

	function add(type, params) {
		dispatch("add", {type, params})
		show = false
	}

	let show = false
	let section = "command"
	let filter = ""
</script>

<style lang="scss">
	@import "bulma/sass/components/panel";

	.modal-content {
		height: 80vh;
	}
	.panel-block a:hover {
			background-color: $panel-block-hover-background-color;
	}
</style>

<Button class="is-primary" icon="add" on:click={() => show = true}>
	Add instruction
</Button>

{#if show}
	<div class="modal is-active" transition:fade={{duration:50}}>
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
						<a class:is-active={section === "control"} href="#top" on:click|preventDefault={() => section="control"}>Control</a>
					</p>
					{#if section==="command"}
						{#each $commands as cmd}
							{#if !filter.length || cmd?.name?.includes(filter) || cmd?.description?.includes(filter)}
								<a href="#top" class="panel-block has-text-black" on:click|preventDefault={() => add("command", {"name": cmd.name})}>
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
					{:else if section === "control"}
						<a href="#top" class="panel-block has-text-black" on:click|preventDefault={() => add("if")}>
							<div class="columns w-100">
								<div class="column is-one-third">
									<Icon icon="condition">
										<b>if</b>
									</Icon>
								</div>
								<div class="column is-two-thirds">
									Execute commands conditionally
								</div>
							</div>
							
						</a>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}

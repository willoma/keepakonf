<script>
	export let instruction
	export let highlight

	import { command } from "$lib/store"
	import {  statuscolor } from "$lib/color"

	import Detail from "./detail/Detail.svelte"
	import ShowInstructionHeader from "./ShowInstructionHeader.svelte"
	import ShowParameter from "./ShowParameter.svelte"

	$: cmd = $command(instruction.command)

	let showDetail = highlight ?? false

	$: color = statuscolor[instruction.status]
</script>

<style lang="scss">
	@import "bulma/sass/utilities/derived-variables";

	.highlight {
		outline: 2px solid $info;
		outline-offset: 2px;
	}
</style>

{#if showDetail}
	<article class="message is-{color}" class:highlight>
		<div class="message-header is-clickable" role="button" tabindex="0" on:click={() => showDetail = !showDetail} on:keyup={(e) => e.key === "Enter" ? showDetail = !showDetail : null}>
			<ShowInstructionHeader {instruction} {cmd} />
		</div>
		<div class="message-body">
			<div class="columns">
				<div class="column content is-one-third is-size-7">
					<p class="is-size-6">{cmd?.description}</p>
					{#if cmd?.parameters?.length}
						<p>Parameters:</p>
						<ul>
							{#each cmd.parameters as param}
								<li>
									<ShowParameter {param} value={instruction.parameters[param.id]} />
								</li>
							{/each}
						</ul>
					{/if}
				</div>
				<div class="column content is-two-thirds">
					{#if instruction.detail}
						<Detail data={instruction.detail} />
					{/if}
				</div>
			</div>
		</div>
	</article>
{:else}
	<div
		role="button" tabindex="0"
		class="message is-clickable is-{color}"
		class:highlight
		on:click={() => showDetail = !showDetail} on:keyup={(e) => e.key === "Enter" ? showDetail = !showDetail : null}
	>
		<div class="message-body">
			<ShowInstructionHeader {instruction} {cmd} />
		</div>
	</div>
{/if}

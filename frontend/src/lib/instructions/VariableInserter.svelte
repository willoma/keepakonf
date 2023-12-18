<script>
	import { fade, fly } from 'svelte/transition';

	import { Button } from "$lib/c"
	import { globalVariables } from "$lib/store"

	let focusOnInput = false
	let show = false

	function writeVariable(name) {
		const el = document.activeElement
		const before = el.value.substring(0, el.selectionStart)
		const after = el.value.substring(el.selectionEnd)
		const newPos = before.length + name.length + 2
		el.value = `${before}<${name}>${after}`
		el.selectionStart = newPos
		el.selectionEnd = newPos
		show = false
	}
</script>

<style lang="scss">
	@import "bulma/sass/components/panel";

	.variableinserterbutton {
		position: absolute;
		bottom: 0;
		right: 0;
		z-index: 10;
	}

	a.panel-block:hover {
			background-color: $panel-block-hover-background-color;
	}
</style>

<div class="block variableinserter" on:focusin={(e) => focusOnInput = e.target.nodeName === "INPUT"} on:focusout={() => focusOnInput = false}>
	<slot />
</div>

{#if focusOnInput}
	<div class="variableinserterbutton pb-5 pr-5" role="none" on:mousedown|preventDefault in:fly={{x:"-100vw"}} out:fly={{x:"100%"}}>
		<Button class="is-primary" icon="variable" on:click={() => show = true}>
			Variable
		</Button>
	</div>
{/if}

{#if show}
	<div class="modal is-active" role="none" on:mousedown|preventDefault transition:fade={{duration:50}}>
		<div class="modal-background" role="none" on:click={() => show = false}></div>
		<div class="modal-content">
			<div class="box p-0">
				<div class="panel">
					<div class="panel-heading is-flex">
						<div class="is-flex-grow-1">
							Select variable
						</div>
						<div class="is-flex-grow-0">
							<Button icon="cancel" class="is-small" on:click={() => show = false}>
								Cancel
							</Button>
						</div>
					</div>
					{#each $globalVariables as v}
						<a href="#top" class="panel-block has-text-black w-100" on:click|preventDefault={() => writeVariable(v.name)}>
							<b>{v?.name ?? "Unknown"}</b>: {v?.description ?? "Unknown"} (<code>{v?.value ?? "Unknown"}</code>)
						</a>
					{/each}
				</div>
			</div>
		</div>
	</div>
{/if}
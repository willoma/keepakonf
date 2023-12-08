<script>
	import { Button } from "$lib/c"
	import { globalVariables } from "$lib/store"

	let focusOnInput = false
	let showSelector = false

	function writeVariable(name) {
		const el = document.activeElement
		const before = el.value.substring(0, el.selectionStart)
		const after = el.value.substring(el.selectionEnd)
		const newPos = before.length + name.length + 2
		el.value = `${before}<${name}>${after}`
		el.selectionStart = newPos
		el.selectionEnd = newPos
		showSelector = false
	}
</script>

<style>
	.variableinserterbutton {
		position: absolute;
		bottom: 0;
		right: 0;
		z-index: 200;
	}
</style>

<div class="variableinserter" on:focusin={(e) => focusOnInput = e.target.nodeName === "INPUT"} on:focusout={() => focusOnInput = false}>
	<slot />
</div>

{#if focusOnInput}
	<div class="variableinserterbutton pb-5 pr-5" role="none" on:mousedown|preventDefault>
		<Button class="is-primary" icon="variable" on:click={() => showSelector = true}>
			Variable
		</Button>
	</div>
{/if}

{#if showSelector}
	<div class="modal is-active" role="none" on:mousedown|preventDefault>
		<div class="modal-background" role="none" on:click={() => showSelector = false}></div>
		<div class="modal-content">
			<div class="box p-0">
				<div class="menu">
					<p class="menu-label pt-3 pl-5">
						Select variable
					</p>
					<ul class="menu-list">
						{#each $globalVariables as v}
							<li>
								<a href="#top" on:click|preventDefault={() => writeVariable(v.name)}>
									<b>{v?.name ?? "Unknown"}</b>: {v?.description ?? "Unknown"} (<code>{v?.value ?? "Unknown"}</code>)
								</a>
							</li>
						{/each}
					</ul>
					<p class="menu-label has-text-right mt-1 pb-1 pr-1">
						<Button icon="cancel" class="is-small" on:click={() => showSelector = false}>
							Cancel
						</Button>
					</p>
				</div>
			</div>
		</div>
	</div>
{/if}
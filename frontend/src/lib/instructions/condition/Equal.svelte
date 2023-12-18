<script>
	export let initial

	let left = initial?.left ?? ""
	let right = initial?.right ?? ""

	export function makeData() {
		return {
			"op": "=",
			"left": left,
			"right": right,
		}
	}
</script>

<div class="field has-addons">
	<div class="control is-expanded">
		<input
			type="text"
			class="input"
			bind:value={left}
		/>
	</div>
	<p class="control">
		<input
			tabindex="-1"
			type="text"
			class="input has-text-centered"
			value="="
			readonly
		/>
	</p>
	<div class="control is-expanded">
		<input
			type="text"
			class="input"
			bind:value={right}
		/>
	</div>
</div>

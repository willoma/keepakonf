<script>
	export let data

	import { ChoicesButtons } from "../../c";
	import DiffMatchPatch, { DIFF_DELETE, DIFF_EQUAL, DIFF_INSERT } from "diff-match-patch"

	let mode = "char"

	const options = [
		["char", "Character mode"],
		["line", "Line mode"],
	]

	let diff
	$: {
		const dmp = new DiffMatchPatch()
		if (mode === "line") {
			const asChars = dmp.diff_linesToChars_(data.b, data.a)
			diff = dmp.diff_main(asChars.chars1, asChars.chars2, false)
			dmp.diff_charsToLines_(diff, asChars.lineArray)

		} else {
			diff = dmp.diff_main(data.b, data.a)
		}
		dmp.diff_cleanupSemantic(diff)
	}
</script>

<ChoicesButtons bind:value={mode} {options} />

<pre class="has-text-white has-background-dark">
{#each diff as token}
	{#if token[0] === DIFF_INSERT}
		<ins class="has-background-success-dark">{token[1]}</ins>
	{:else if token[0] === DIFF_DELETE}
		<del class="has-background-danger-dark">{token[1]}</del>
	{:else if token[0] === DIFF_EQUAL}
		<span>{token[1]}</span>
	{/if}
{/each}
</pre>

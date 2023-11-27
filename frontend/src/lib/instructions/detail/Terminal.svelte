<script>
	export let data
</script>

{#if data.cmd}
	<pre class="has-text-white has-background-dark">{data.cmd}</pre>
{/if}
<pre class="has-text-white has-background-dark">{data.out}</pre>
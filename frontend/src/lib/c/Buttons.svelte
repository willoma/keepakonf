<script>
	export let right = false
</script>

<style>
	.is-justify-content-end {
		width: 100%;
	}
</style>

<div class="field">
	<div class="buttons" class:is-justify-content-end={right}>
		<slot />
	</div>
</div>

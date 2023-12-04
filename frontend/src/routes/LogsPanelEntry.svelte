<script>
	export let log

	import Detail from "$lib/instructions/detail/Detail.svelte"
	import { Icon } from "$lib/c"

	function twoDigits(nb) {
		if (nb < 10) {
			return`0${nb}`
		}
			return `${nb}` 
	}

	$: time = log.ts ? new Date(log.ts) : null

	$: applied = log?.st === "applied"
	$: failed = log?.st === "failed"
	$: todo = log?.st === "todo"

	$: href = log.gid ?
		log.iid ?
			`/${log.gid}#${log.iid}` :
			`/${log.gid}` :
		null

	$: group = log.grp

	let showdetail = false
</script>

<style>
	.showdetail {
		border-bottom: 0;
	}

	.time {
		margin: -0.75em 0.5em -0.75em -0.5em;
		font-size: 0.8em;
		line-height: 1.2;
		text-align: center;
	}

	.no-wrap {
		white-space: nowrap;
	}

	.no-overflow {
		overflow: hidden;
		text-overflow: ellipsis;
	}
</style>

<div
	class="panel-block"
	class:has-background-success-light={applied}
	class:has-background-danger-light={failed}
	class:has-background-warning-light={todo}
	class:showdetail
>
	<div class="is-flex-grow-0 is-flex-shrink-0 time">
		{#if time}
			{twoDigits(time.getDate())}/{twoDigits(time.getMonth())}<br/>
			{twoDigits(time.getHours())}:{twoDigits(time.getMinutes())}
		{:else}
			00/00<br/>
			00/00
		{/if}
	</div>
	<div class="is-flex-grow-1 is-flex-shrink-1 no-wrap no-overflow">
		<Icon tclass="w-100" icon={log.ico}>{log.msg??"-"}</Icon>
	</div>
	<div class="is-flex-grow-0 is-flex-shrink-0 field is-grouped">
		{#if group}
			<div class="control">
				<span class="tag is-info is-light">
					<Icon icon="group">{group}</Icon>
				</span>
			</div>
		{/if}
		{#if href}
			<div class="control">
				<a class="tag is-primary" {href}>
					<Icon icon="link" />
				</a>
			</div>
		{/if}
		{#if log.dtl}
			<div class="control">
				<div role="button" tabindex="0" class="tag is-link is-clickable" on:click={() => showdetail = !showdetail} on:keyup={(e) => e.key === "Enter" ? showdetail = !showdetail : null}>
					<Icon icon={showdetail?"less":"more"} />
				</div>
			</div>
		{/if}
	</div>
</div>

{#if showdetail}
	<div
		class="panel-block"
		class:has-background-success-light={applied}
		class:has-background-danger-light={failed}
		class:has-background-warning-light={todo}
	>
		<div class="content no-overflow">
			<Detail data={log.dtl} />
		</div>
	</div>
{/if}

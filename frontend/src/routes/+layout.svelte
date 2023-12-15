<script>
	import "../app.scss"

	import { page } from '$app/stores'

	import { groups, reconnecting } from "$lib/store"
	import { Icon } from "$lib/c"
	import { statuscolordark } from "$lib/color"

	// Only show the "reconnecting" error after 5 seconds
	let showReconnecting = false
	$: if ($reconnecting) {
		setTimeout(() => {
			if ($reconnecting) {
				showReconnecting = true
			}
		}, 5000)
	}

	let filter = ""
</script>

<style>
	nav {
		box-shadow: 0 0.5em 1em -0.125em rgba(10, 10, 10, 0.1), 0 0px 0 1px rgba(10, 10, 10, 0.02)
	}
</style>

<div class="is-flex">

	<!-- Left navbar doesn't grow nor shrink horizontally -->
	<nav class="is-flex-grow-0 is-flex-shrink-0">
		<!-- Navbar contains 2 parts (in order to allow scrolling only on groups list) -->
		<div class="is-flex is-flex-direction-column" style:height="100vh">
			<!-- Top part of navbar -->
			<div class="is-flex-grow-0 is-flex-shrink-0">
				<!-- App title -->
				<h1 class="title is-4 m-3">Keepakonf</h1>
				<!-- Search box -->
				<div class="m-2">
					<input class="input is-small" type="text" placeholder="Filter groups" bind:value="{filter}">
				</div>
			</div>
			<!-- Scrollable menu -->
			<div class="is-flex-grow-1 is-flex-shrink-1 p-1" style:overflow-y="scroll">
				<div class="menu">
					<ul class="menu-list">
						<!-- Summary -->
						<li>
							<a class:is-active={$page.route.id === "/"} href="/">
								<Icon icon="search">
									Summary
								</Icon>
							</a>
						</li>
						<!-- List of groups -->
						{#each $groups as group}
							{#if group.name.includes(filter)}
								{@const active = $page.params.id === group.id}
								<li>
									<a class:is-active={active} href="/{group.id}">
										<Icon icon={group?.icon??"group"} class="{active ? "": "has-text-"+statuscolordark[group.status]}">
											{group.name}
										</Icon>
									</a>
								</li>
							{/if}
						{/each}
					</ul>
				</div>
			</div>
		</div>
	</nav>

	<!-- Main content grows and shrinks, and is scrollable independently -->
	<main class="is-flex-grow-1 is-flex-shrink-1 p-5" style:overflow-y="scroll" style:height="100vh">
		<!-- Alert message when server is not reachable -->
		{#if showReconnecting}
			<div class="notification is-danger is-light has-text-centered">
				<span class="icon-text">
					<span class="icon loader mr-4"></span>
					<span>Connection lost, trying again...</span>
				</span>
			</div>
		{/if}
		<!-- Content itself -->
		<slot></slot>
	</main>

</div>

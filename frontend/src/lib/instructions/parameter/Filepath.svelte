<script>
	export let field
	export let label

	import { Field, Icon } from "$lib/c"
	import { randomID } from "$lib/random"
	import { users } from "$lib/store"

	const id = randomID()

	let input

	let showDropdown = false
</script>

<Field {field} {id} {label} horizontal addons>
	<div class="control is-expanded">
		<input
			type="text"
			class="input"
			class:is-danger={!$field.valid}
			{id}
			required
			bind:value={$field.value}
			bind:this={input}
		/>
	</div>
	{#if $users.length}
		<div class="control">
			<button
				class="button"
				type="button"
				on:click={() => {
					$field.value = $users[0].home
					input.focus()
				}}
			>
				{$users[0].home}
			</button>
		</div>
		<div class="control">
			<button
				class="button"
				type="button"
				on:click={() => {
					$field.value = $users[1].home
					input.focus()
				}}
			>
				{$users[1].home}
			</button>
		</div>
		<div class="control">
			<div class="dropdown is-right" class:is-active={showDropdown}>
				<div class="dropdown-trigger">
					<button class="button" type="button" on:click={() => showDropdown = !showDropdown}>
						<Icon icon="more" />
					</button>
				</div>
				<div class="dropdown-menu">
					<div class="dropdown-content">
						{#each $users as user}
							<a
								href="#top"
								class="dropdown-item"
								on:click|preventDefault={() => {
									$field.value = user.home
									showDropdown = false
									input.focus()
								}}
							>{user.name}: {user.home}</a>
						{/each}
					</div>
				</div>
			</div>
		</div>
	{/if}
</Field>
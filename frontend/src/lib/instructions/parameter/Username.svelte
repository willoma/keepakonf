<script>
	export let field
	export let label

	import { Field, Icon } from "$lib/c"
	import { randomID } from "$lib/random"
	import { users } from "$lib/store"

	const id = randomID()

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
		/>
	</div>
	{#if $users.length}
		<div class="control">
			<button
				class="button"
				type="button"
				on:click={() => $field.value = $users[0].name}
			>
				{$users[0].name}
			</button>
		</div>
		<div class="control">
			<button
				class="button"
				type="button"
				on:click={() => $field.value = $users[1].name}
			>
				{$users[1].name}
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
									$field.value = user.name
									showDropdown = false
								}}
							>{user.name}</a>
						{/each}
					</div>
				</div>
			</div>
		</div>
	{/if}
</Field>

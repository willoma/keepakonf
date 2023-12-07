<script>
	import { goto } from "$app/navigation"
	import { page } from '$app/stores'

	import { Button, ButtonA, Buttons, Title } from "$lib/c"
	import { group, groups, socket } from "$lib/store"
	import EditGroup from "$lib/instructions/EditGroup.svelte"

	$: grp = $group($page.params.id)

	let editor
	let valid
	
	let confirmRemove = false

	async function doModify() {
		socket.emit("modify group", editor.makeData(), (response) => {
			$groups = $groups.map((someGroup) => someGroup.id === response.id ? response : someGroup)
			goto(`/${response.id}`)
		})
	}

	async function doRemove() {
		socket.emit("remove group", $page.params.id, (response) => {
			if (response === $page.params.id) {
				$groups = $groups.filter(
					(group) => group.id !== response
				)
				goto("/")
			}
		})
	}
	
	$: cancelTarget = $page.url.searchParams.get("back") ?
		`/#${$page.url.searchParams.get("back")}` :
		`/${$page.params.id}`
</script>

<svelte:head>
	<title>Keepakonf - Modify {grp?.name}</title>
</svelte:head>

<form on:submit={doModify}>
	<Title stickyactions icon="edit">
		{grp?.name}
		<Buttons slot="actions">
			{#if confirmRemove}
				<p class="is-align-self-center mr-2 mb-2">Do you really want to remove <i>{grp?.name}</i>?</p>
				<Button class="is-info" icon="cancel" on:click={() => confirmRemove = false}>
					No
				</Button>
				<Button class="is-danger" icon="remove" on:click={doRemove}>
					Yes
				</Button>
			{:else}
				<ButtonA icon="cancel" href={cancelTarget}>
					Cancel
				</ButtonA>
				<Button class="is-warning" icon="remove" on:click={() => confirmRemove = true}>
					Remove
				</Button>
				<Button submit class="is-info" icon="save" disabled={!valid}>
					Save
				</Button>
			{/if}
			</Buttons>
	</Title>

	<EditGroup bind:this={editor} bind:valid group={grp} />
</form>

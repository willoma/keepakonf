<script>
	import { goto } from "$app/navigation"

	import { Button, ButtonA, Buttons, Title } from "$lib/c"
	import EditGroup from "$lib/instructions/EditGroup.svelte"
	import { groups, socket } from "$lib/store"

	let makeData
	let valid

	async function doAdd() {
		socket.emit("add group", makeData(), (response) => {
			$groups = [...$groups, response]
			goto(`/${response.id}`)
		})
	}
</script>

<svelte:head>
	<title>Keepakonf - Create a new group</title>
</svelte:head>

<form on:submit={doAdd}>
	<Title icon="add">
		Create a new group
		<Buttons slot="actions">
			<ButtonA icon="cancel" href="/">
				Cancel
			</ButtonA>
			<Button submit class="is-info" icon="add" disabled={!valid}>
				Create
			</Button>
		</Buttons>
	</Title>

	<EditGroup bind:makeData bind:valid />
</form>

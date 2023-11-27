import { derived, writable } from 'svelte/store'
import { io } from "socket.io-client"

export const socket = io()
// export const socket = io("http://127.0.0.1:35653")

export const reconnecting = writable(false)

export const groups = writable([])

export const commands = writable([])

export const logs = writable([])

export const logsReachedTheEnd = writable(false)

const groupsMap = derived(
	groups,
	($groups) => Object.fromEntries($groups.map(grp => [grp.id, grp]))
)

export const group = derived(
	groupsMap,
	($groupsMap) => (id) => $groupsMap[id]
)

const commandsMap = derived(
	commands,
	($commands) => Object.fromEntries($commands.map(cmd => [cmd.name, cmd]))
)

export const command = derived(
	commandsMap,
	($commandsMap) => (name) => $commandsMap[name]
)

socket.on("connect", () => {
	reconnecting.set(false)
	socket.emit("groups", (response) => groups.set(response))
	socket.emit("commands", (response) => commands.set(response))
	socket.emit("logs", (response) => {
		logs.set(response.logs?.toReversed())
		if (response.reached_the_end) {
			logsReachedTheEnd.set(true)
		}
	})
})

socket.on("disconnect", (reason) => {
	console.log(reason)
	if (reason !== "io server disconnect" && reason !== "io client disconnect") {
		reconnecting.set(true)
	}
})

socket.on("add group", (data) => {
	groups.update((groupz) => [...groupz, data])
})

socket.on("modify group", (data) => {
	groups.update((groupz) => groupz.map(
		(group) => group.id === data.id ? data : group
	))
})

socket.on("remove group", (data) => {
	groups.update((groupz) => groupz.filter(
		(group) => group.id !== data
	))
})

socket.on("status", (data) => {
	groups.update((groupz) => {
		for (const group of groupz) {
			for (const instruction of group.instructions) {
				if (instruction.id === data.instruction) {
					instruction.status = data.status
					instruction.info = data.info
					instruction.detail_type = data.detail_type
					instruction.detail = data.detail
				}
			}
		}
		return groupz
	})
})

socket.on("group status", (data) => {
	groups.update((groupz) => {
		for (const group of groupz) {
			if (group.id === data.group) {
				group.status = data.status
			}
		}
		return groupz
	})
})

socket.on("log", (log) => {
	logs.update(logs => [log, ...logs])
})

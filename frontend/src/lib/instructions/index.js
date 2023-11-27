import { field } from 'svelte-forms'
import { required } from 'svelte-forms/validators'
import { socket } from "$lib/store"

export function makefield(param, initial) {
	let value
	switch (param.type) {
	case "bool":
		value = initial??false
	case "string":
		value = initial??""
		break
	case "[string]":
		value =  initial?[...initial]:[""]
		break
	default:
		value = null
	}
	return field(param.id, value, [required()])
}

export function applyGroup(id) {
	socket.emit("apply group", id)
}

export function applyInstruction(id) {
	socket.emit("apply instruction", id)
}

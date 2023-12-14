import { field } from 'svelte-forms'
import { required } from 'svelte-forms/validators'
import { socket } from "$lib/store"

function filepathValidator() {
	return (value) => ({
		"valid": value.startsWith("/") && !value.endsWith("/"),
		"name": "filepath",
	})
}

function usernameValidator() {
	return (value) => {
		const firstIsAlphabetic = (value[0] >= "a" && value[0] <= "z") || (value[0] >= "A" && value[0] <= "Z")
		const tooLong = value.length > 30

		const allCharactersValid = true
		for (let i = 1; i < value.length; i++) {
			// Valid Unix Characters: letters, numbers, "-"", ""."", and "_" ​
			const c = value.charAt(i)
			const cValid = (c === "-" || c === "." || c === "_" || (c >= "0" && c <= "9") || (c >= "a" && c <= "z") || (c >= "A" && c <= "Z"))
			if (!cValid) {
				allCharactersValid = false
				break
			}
		}

		return {
			"valid": firstIsAlphabetic && !tooLong && allCharactersValid,
			"name": "username",
		}
	}
}

export function makefield(param, initial) {
	let value
	let validators
	switch (param.type) {
	case "bool":
		value = initial??false
		validators = []
		break
	case "filepath":
		value = initial??""
		validators = [required(), filepathValidator()]
		break
	case "string":
		value = initial??""
		validators = [required()]
		break
	case "[string]":
		value =  initial?[...initial]:[""]
		validators = [required()]
	break
	case "text":
		value = initial??""
		validators = [required()]
		break
	case "username":
		value = initial??""
		validators = [required(), usernameValidator()]
		break
	default:
		value = null
		validators = []
	}
	return field(param.id, value, validators, {"stopAtFirstError": true})
}

export function applyGroup(id) {
	socket.emit("apply group", id)
}

export function applyInstruction(id) {
	socket.emit("apply instruction", id)
}

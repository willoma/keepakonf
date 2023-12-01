import Raw from "./Raw.svelte"
import Error from "./Error.svelte"
import Table from "./Table.svelte"
import Terminal from "./Terminal.svelte"

const detailComponents = {
	"error": Error,
	"table": Table,
	"terminal": Terminal,
}

export function detailComponent(detail) {
	return detailComponents[detail.type] ?? Raw
}
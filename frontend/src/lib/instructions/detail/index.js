import Raw from "./Raw.svelte"
import Table from "./Table.svelte"
import Terminal from "./Terminal.svelte"

const detailComponents = {
	"table": Table,
	"terminal": Terminal,
}

export function detailComponent(detail) {
	return detailComponents[detail.type] ?? Raw
}
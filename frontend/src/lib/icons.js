import add from "@fortawesome/fontawesome-free/svgs/solid/plus.svg?raw"
import cancel from "@fortawesome/fontawesome-free/svgs/solid/rotate-left.svg?raw"
import check from "@fortawesome/fontawesome-free/svgs/solid/check.svg?raw"
import command from "@fortawesome/fontawesome-free/svgs/solid/terminal.svg?raw"
import database from "@fortawesome/fontawesome-free/svgs/solid/database.svg?raw"
import edit from "@fortawesome/fontawesome-free/svgs/solid/pen-to-square.svg?raw"
import error from "@fortawesome/fontawesome-free/svgs/solid/triangle-exclamation.svg?raw"
import file from "@fortawesome/fontawesome-free/svgs/solid/file.svg?raw"
import folder from "@fortawesome/fontawesome-free/svgs/regular/folder.svg?raw"
import group from "@fortawesome/fontawesome-free/svgs/solid/layer-group.svg?raw"
import less from "@fortawesome/fontawesome-free/svgs/solid/caret-up.svg?raw"
import link from "@fortawesome/fontawesome-free/svgs/solid/link.svg?raw"
import logs from "@fortawesome/fontawesome-free/svgs/solid/bars-staggered.svg?raw"
import more from "@fortawesome/fontawesome-free/svgs/solid/caret-down.svg?raw"
import packages from "@fortawesome/fontawesome-free/svgs/solid/cubes.svg?raw"
import remove from "@fortawesome/fontawesome-free/svgs/solid/trash-can.svg?raw"
import run from "@fortawesome/fontawesome-free/svgs/solid/gears.svg?raw"
import save from "@fortawesome/fontawesome-free/svgs/solid/floppy-disk.svg?raw"
import search from "@fortawesome/fontawesome-free/svgs/solid/magnifying-glass.svg?raw"
import ubuntu from "@fortawesome/fontawesome-free/svgs/brands/ubuntu.svg?raw"
import unknown from "@fortawesome/fontawesome-free/svgs/solid/question.svg?raw"
import variable from "@fortawesome/fontawesome-free/svgs/solid/code.svg?raw"

export const icons = {
	add,
	cancel,
	check,
	command,
	database,
	edit,
	error,
	file,
	folder,
	group,
	less,
	link,
	logs,
	more,
	packages,
	remove,
	run,
	save,
	search,
	ubuntu,
	unknown,
	variable,
}

export function icon(i) {
	return icons[i] ?? unknown
}

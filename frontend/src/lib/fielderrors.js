const errorMessages = {
	"required": "Cannot be empty",
}

export function errors(errs) {
	return errs.map(e => errorMessages[e]).join(", ")
}
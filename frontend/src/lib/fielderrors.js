const errorMessages = {
	"filepath": "Must be an absolute path: /.../...",
	"required": "Cannot be empty",
	"username": "Invalid username (0-9, a-z, A-Z, \".\", \"-\", \"_\")",
}

export function errors(errs) {
	return errs.map(e => errorMessages[e]).join(", ")
}
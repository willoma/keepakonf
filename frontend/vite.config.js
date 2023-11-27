import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
	plugins: [sveltekit()],

	css: {
		preprocessorOptions: {
			scss: {
				additionalData: '@use "src/variables.scss" as *;',
			},
		},
	},

	server: {
		proxy: {
			"/socket.io": {
				target: "http://localhost:35653",
				ws: true
			}
		}
	}
});

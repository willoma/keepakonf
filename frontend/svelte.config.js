import { vitePreprocess } from "@sveltejs/kit/vite";
import adapter from "@sveltejs/adapter-static";

const config = {
	kit: {
		adapter: adapter({
			fallback: '200.html'
		}),
		prerender: {
			handleMissingId: "ignore",
		}
	},

	preprocess: [vitePreprocess({})],
};

export default config;

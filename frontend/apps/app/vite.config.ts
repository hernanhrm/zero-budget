import tailwindcss from "@tailwindcss/vite"
import { devtools } from "@tanstack/devtools-vite"
import { tanstackRouter } from "@tanstack/router-plugin/vite"
import viteReact from "@vitejs/plugin-react"
import { defineConfig } from "vite"
import tsconfigPaths from "vite-tsconfig-paths"

const config = defineConfig(({ command }) => ({
	plugins: [
		// Devtools plugin + UI are dev-only; skipping on `vite build` cuts CI/Pages memory and time.
		...(command === "serve" ? [devtools()] : []),
		tsconfigPaths({ projects: ["./tsconfig.json"] }),
		tailwindcss(),
		tanstackRouter({ target: "react", autoCodeSplitting: true }),
		viteReact(),
	],
}))

export default config

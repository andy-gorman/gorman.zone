import js from "@eslint/js";
import globals from "globals";
import markdown from "@eslint/markdown";
import { defineConfig, globalIgnores } from "eslint/config";


export default defineConfig([
	globalIgnores(["public/"]),
  { files: ["**/*.{js,mjs,cjs}"], plugins: { js }, extends: ["js/recommended"] },
  { files: ["**/*.{js,mjs,cjs}"], languageOptions: { globals: globals.browser } },
  { files: ["**/*.md"], plugins: { markdown }, language: "markdown/commonmark", extends: ["markdown/recommended"] },
]);

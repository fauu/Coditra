import { defineConfig, globalIgnores } from "eslint/config";
import js from "@eslint/js";
import prettierConfig from "eslint-config-prettier";
import sveltePlugin from "eslint-plugin-svelte3";
import globals from "globals";

export default defineConfig([
  globalIgnores(["webpack.config.js"]),

  js.configs.recommended,
  prettierConfig,
  {
    languageOptions: {
      ecmaVersion: 2019,
      sourceType: "module",
      globals: {
        ...globals.browser,
        ...globals.es2019,
      },
    },
    rules: {
      quotes: ["warn", "double"],
      semi: ["error", "always"],
      "no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
          caughtErrorsIgnorePattern: "^_",
        },
      ],
    },
  },
  {
    files: ["**/*.svelte"],
    plugins: {
      svelte3: sveltePlugin,
    },
    processor: "svelte3/svelte3",
  },
]);

// eslint-disable-next-line no-undef
module.exports = {
  extends: ["eslint:recommended", "prettier"],
  parserOptions: {
    ecmaVersion: 2019,
    sourceType: "module",
  },
  env: {
    es6: true,
    browser: true,
  },
  plugins: ["svelte3"],
  ignorePatterns: ["webpack.config.js"],
  overrides: [
    {
      files: ["**/*.svelte"],
      processor: "svelte3/svelte3",
    },
  ],
  rules: {
    quotes: ["warn", "double"],
    semi: ["error", "always"],
  },
};

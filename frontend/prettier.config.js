export default {
  plugins: ["prettier-plugin-go-template"],
  overrides: [
    {
      files: ["*.html"],
      options: {
        parser: "go-template",
      },
    },
  ],
};

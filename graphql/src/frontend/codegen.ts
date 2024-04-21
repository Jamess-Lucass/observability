import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  overwrite: true,
  schema: "http://localhost:4000/graphql",
  documents: ["app/**/*.tsx", "!gql/**/*"],
  generates: {
    "./gql/": {
      preset: "client",
      plugins: [],
      config: {
        scalars: {
          UUID: {
            input: "string",
            output: "string",
          },
          Decimal: {
            input: "number",
            output: "number",
          },
        },
      },
    },
  },
};

export default config;

import { initGraphQLTada } from "gql.tada";
import type { introspection } from "./graphql-env.d.ts";

export const graphql = initGraphQLTada<{
  introspection: introspection;
  scalars: {
    UUID: string;
    Decimal: number;
  };
}>();

export type { FragmentOf, ResultOf, VariablesOf } from "gql.tada";
export { readFragment } from "gql.tada";

export type ExtractNodeType<T> = T extends { nodes: infer U }
  ? U extends (infer V)[]
    ? V
    : never
  : never;

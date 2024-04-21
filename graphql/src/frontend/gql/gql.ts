/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n  mutation CreateProduct($input: CreateProductInput!) {\n    createProduct(input: $input) {\n      response {\n        __typename\n        ... on Product {\n          id\n          name\n          description\n          price\n        }\n        ... on ErrorPayload {\n          errors {\n            message\n            path\n          }\n        }\n      }\n    }\n  }\n": types.CreateProductDocument,
    "\n  query Product($productId: UUID!) {\n    product(id: $productId) {\n      price\n      name\n      id\n      description\n    }\n  }\n": types.ProductDocument,
    "\n  query Products($where: ProductFilterInput) {\n    products(where: $where) {\n      nodes {\n        id\n        name\n        description\n        price\n      }\n    }\n  }\n": types.ProductsDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  mutation CreateProduct($input: CreateProductInput!) {\n    createProduct(input: $input) {\n      response {\n        __typename\n        ... on Product {\n          id\n          name\n          description\n          price\n        }\n        ... on ErrorPayload {\n          errors {\n            message\n            path\n          }\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  mutation CreateProduct($input: CreateProductInput!) {\n    createProduct(input: $input) {\n      response {\n        __typename\n        ... on Product {\n          id\n          name\n          description\n          price\n        }\n        ... on ErrorPayload {\n          errors {\n            message\n            path\n          }\n        }\n      }\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query Product($productId: UUID!) {\n    product(id: $productId) {\n      price\n      name\n      id\n      description\n    }\n  }\n"): (typeof documents)["\n  query Product($productId: UUID!) {\n    product(id: $productId) {\n      price\n      name\n      id\n      description\n    }\n  }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n  query Products($where: ProductFilterInput) {\n    products(where: $where) {\n      nodes {\n        id\n        name\n        description\n        price\n      }\n    }\n  }\n"): (typeof documents)["\n  query Products($where: ProductFilterInput) {\n    products(where: $where) {\n      nodes {\n        id\n        name\n        description\n        price\n      }\n    }\n  }\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;
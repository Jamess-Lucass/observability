"use client";

import { graphql } from "@/gql/gql";
import { useQuery } from "@tanstack/react-query";
import request from "graphql-request";

const productQuery = graphql(`
  query Product($productId: UUID!) {
    product(id: $productId) {
      price
      name
      id
      description
    }
  }
`);

type Props = {
  params: Params;
};

type Params = {
  id: string;
};

export default function ProductDetailsPage({ params }: Props) {
  const { data } = useQuery({
    queryKey: ["products", params.id],
    queryFn: async () =>
      request("http://localhost:4000/graphql", productQuery, {
        productId: params.id,
      }),
  });

  if (!data?.product) {
    return <h1>No product found.</h1>;
  }

  return (
    <div>
      <p>id: {data.product.id}</p>
      <p>name: {data.product.name}</p>
      <p>description: {data.product.description}</p>
      <p>
        price:{" "}
        {new Intl.NumberFormat("en-US", {
          style: "currency",
          currency: "GBP",
        }).format(data.product.price)}
      </p>
    </div>
  );
}

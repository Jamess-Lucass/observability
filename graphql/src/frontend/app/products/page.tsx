"use client";

import { useQuery } from "@tanstack/react-query";
import { graphql } from "@/graphql";
import request from "graphql-request";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ChangeEvent, useState } from "react";
import { Input } from "@/components/ui/input";
import Link from "next/link";

const productsQuery = graphql(`
  query Products($where: ProductFilterInput) {
    products(where: $where) {
      nodes {
        id
        name
        description
        price
      }
    }
  }
`);

export default function Products() {
  const [search, setSearch] = useState<string>("");

  const { data } = useQuery({
    queryKey: ["products", search],
    queryFn: async () =>
      request("http://localhost:4000/graphql", productsQuery, {
        where: { name: { contains: search } },
      }),
  });

  const handleOnChange = (e: ChangeEvent<HTMLInputElement>) => {
    setTimeout(() => {
      setSearch(e.target.value);
    }, 500);
  };

  if (!data) {
    return <h1>Could not load products</h1>;
  }

  return (
    <div className="flex flex-col gap-4">
      <Input type="text" placeholder="T-Shirt" onChange={handleOnChange} />

      <div className="flex flex gap-4">
        {data.products?.nodes?.map((product) => (
          <Link key={product.id} href={`/products/${product.id}`}>
            <Card className="min-w-64">
              <CardHeader>
                <CardTitle>{product.name}</CardTitle>
                <CardDescription>{product.description}</CardDescription>
              </CardHeader>
              <CardContent>
                <p>
                  {new Intl.NumberFormat("en-US", {
                    style: "currency",
                    currency: "GBP",
                  }).format(product.price)}
                </p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}

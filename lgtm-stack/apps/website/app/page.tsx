import { Button } from "@/components/ui/button";
import { Product } from "@/lib/types";
import Link from "next/link";

async function getProducts() {
  const res = await fetch("http://host.docker.internal:1000/api/products");

  if (!res.ok) {
    throw new Error("Failed to fetch data");
  }

  return res.json() as Promise<Product[]>;
}

export default async function Home() {
  const products = await getProducts();

  return (
    <>
      <h1>Products</h1>

      {products.map((product) => (
        <li key={product.id}>
          <Button variant="link">
            <Link href={`/products/${product.id}`}>{product.name}</Link>
          </Button>
        </li>
      ))}
    </>
  );
}

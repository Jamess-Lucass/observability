import { Product } from "@/lib/types";

async function getProduct(id: string) {
  const res = await fetch(
    `http://host.docker.internal:1000/api/products/${id}`
  );

  if (!res.ok) {
    throw new Error("Failed to fetch data");
  }

  return res.json() as Promise<Product>;
}

type Props = {
  params: Params;
};

type Params = {
  id: string;
};

export default async function ProductsPage({ params }: Props) {
  const product = await getProduct(params.id);

  return (
    <>
      <h1>Product: {product.name}</h1>
    </>
  );
}

import { graphql } from "@/gql/gql";
import request from "graphql-request";
import { QueryClient } from "@tanstack/react-query";

const basketQuery = graphql(`
  query ($id: String) {
    basket(id: $id) {
      id
    }
  }
`);

type Props = {
  params: Params;
};

type Params = {
  id: string;
};

export default async function Basket({ params }: Props) {
  const queryClient = new QueryClient();

  const basket = await queryClient.fetchQuery({
    queryKey: ["basket", params.id],
    queryFn: async () => request("http://localhost:4000/graphql", basketQuery),
  });

  return <div className="flex flex-col gap-4">{/* {basket.} */}</div>;
}

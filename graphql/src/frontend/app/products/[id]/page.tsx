"use client";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { VariablesOf, graphql } from "@/graphql";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery } from "@tanstack/react-query";
import request from "graphql-request";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

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

const createBasketMutation = graphql(`
  mutation CreateBasket($input: CreateBasketRequest!) {
    createBasket(input: $input) {
      response {
        __typename
        ... on ErrorPayload {
          errors {
            path
            message
          }
        }
        ... on Basket {
          id
        }
      }
    }
  }
`);

type Props = {
  params: Params;
};

type Params = {
  id: string;
};

const schema = z.object({
  items: z.array(
    z.object({
      quantity: z.number().gt(0),
      productId: z.string().uuid(),
    })
  ),
}) satisfies z.ZodType<VariablesOf<typeof createBasketMutation>["input"]>;

export default function ProductDetailsPage({ params }: Props) {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
  });

  const { data } = useQuery({
    queryKey: ["products", params.id],
    queryFn: async () =>
      request("http://localhost:4000/graphql", productQuery, {
        productId: params.id,
      }),
  });

  const createBasket = useMutation({
    mutationFn: (body: VariablesOf<typeof createBasketMutation>) =>
      request("http://localhost:4000/graphql", createBasketMutation, body),
    onSuccess: (data) => {
      const response = data.createBasket.response;

      if (response?.__typename === "ErrorPayload" && response.errors) {
        for (const error of response.errors) {
          // TODO: make this fully type-safe by adding an enum for the properties
          // to the graphql error endpoint
          form.setError(error.path as keyof z.infer<typeof schema>, {
            type: "custom",
            message: error.message,
          });
        }

        toast.error("Failed to create product");

        return;
      }

      toast.success("Item added to your basket");
    },
  });

  function onSubmit(data: z.infer<typeof schema>) {
    createBasket.mutate({ input: data });
  }

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

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="grid gap-2 space-y-2 mt-8"
        >
          <FormField
            control={form.control}
            name={`items.0.quantity`}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Quantity</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    placeholder="1"
                    defaultValue={0}
                    step={1}
                    {...field}
                    onChange={(x) => field.onChange(Number(x.target.value))}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <Input
            className="hidden"
            value={data.product.id}
            {...form.register("items.0.productId")}
          />

          <Button type="submit">
            {createBasket.isPending && <Spinner className="mr-2" />}
            Add to basket
          </Button>
        </form>
      </Form>
    </div>
  );
}

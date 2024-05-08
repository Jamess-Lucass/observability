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
import { useMutation, useQueryClient } from "@tanstack/react-query";
import request from "graphql-request";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

const createProductMutation = graphql(`
  mutation CreateProduct($input: CreateProductInput!) {
    createProduct(input: $input) {
      response {
        __typename
        ... on Product {
          id
          name
          description
          price
        }
        ... on ErrorPayload {
          errors {
            message
            path
          }
        }
      }
    }
  }
`);

const schema = z.object({
  name: z.string().min(2).max(128),
  description: z.string().min(2).max(128),
  price: z.coerce.number(),
}) satisfies z.ZodType<VariablesOf<typeof createProductMutation>["input"]>;

export default function AdminCreateProductForm() {
  const queryClient = useQueryClient();

  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
      description: "",
      price: 0,
    },
  });

  const createProduct = useMutation({
    mutationFn: (body: VariablesOf<typeof createProductMutation>) =>
      request("http://localhost:5000/graphql", createProductMutation, body),
    onSuccess: (data) => {
      const response = data.createProduct.response;

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

      queryClient.invalidateQueries({ queryKey: ["products"] });

      toast.success("Product has been created");
    },
  });

  function onSubmit(data: z.infer<typeof schema>) {
    createProduct.mutate({ input: data });
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="grid gap-2 space-y-2"
      >
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="T-Shirt" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Input placeholder="A very cool T-Shirt" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="price"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Price</FormLabel>
              <FormControl>
                <Input
                  placeholder="10.00"
                  step="0.01"
                  type="number"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit">
          {createProduct.isPending && <Spinner className="mr-2" />}
          Create Product
        </Button>
      </form>
    </Form>
  );
}

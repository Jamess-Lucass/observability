import { ApolloGateway, IntrospectAndCompose } from "@apollo/gateway";
import { ApolloServer } from "@apollo/server";
import { startStandaloneServer } from "@apollo/server/standalone";

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: "baskets", url: "http://localhost:8080/graphql" },
      { name: "products", url: "http://localhost:5000/graphql" },
    ],
  }),
});

const server = new ApolloServer({
  gateway,
});

const { url } = await startStandaloneServer(server);
console.log(`🚀 Server ready at ${url}`);

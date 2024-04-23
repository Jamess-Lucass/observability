import { ApolloGateway, IntrospectAndCompose } from "@apollo/gateway";
import { ApolloServer } from "@apollo/server";
import { startStandaloneServer } from "@apollo/server/standalone";

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: "products", url: process.env.PRODUCTS_SUBGRAPH_URL },
      // { name: "baskets", url: process.env.BASKETS_SUBGRAPH_URL },
    ],
  }),
});

const server = new ApolloServer({
  gateway,
});

const { url } = await startStandaloneServer(server);
console.log(`🚀 Server ready at ${url}`);

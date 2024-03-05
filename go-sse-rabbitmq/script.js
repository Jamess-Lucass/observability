import http from "k6/http";

export const options = {
  scenarios: {
    consumer: {
      executor: "constant-vus",
      exec: "consumers",
      vus: 50,
      duration: "30s",
    },
    producer: {
      executor: "per-vu-iterations",
      exec: "producers",
      vus: 50,
      iterations: 100,
      startTime: "30s",
      maxDuration: "1m",
    },
  },
};

export function consumers() {
  http.get("http://localhost:8080/sse", {
    tags: { my_custom_tag: "consumer" },
  });
}

export function producers() {
  http.post("http://localhost:8080/producer", {
    tags: { my_custom_tag: "producer" },
  });
}

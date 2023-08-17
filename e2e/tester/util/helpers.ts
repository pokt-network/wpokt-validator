export const sleep = (ms: number) =>
  new Promise((resolve) => setTimeout(resolve, ms));

export const debug = process.env.DEBUG === "true" ? console.log : () => {};

import type { ApiError } from "~/types";

export const getApiErrorMessage = (
  err: unknown,
  fallback: string = "Something went wrong"
) => {
  const error = err as ApiError;
  if (
    error &&
    error.response &&
    error.response.data &&
    error.response.data.message
  ) {
    return error.response.data.message;
  }

  if (err instanceof Error) {
    return err.message;
  }

  return fallback;
};

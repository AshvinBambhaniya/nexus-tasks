import type { UseFetchOptions } from "#app";
import type { JSendResponse } from "~/types";

export const useApi = <T>(
  url: string | (() => string),
  options: UseFetchOptions<JSendResponse<T>, T> = {}
) => {
  const config = useRuntimeConfig();

  return useFetch(url, {
    baseURL: config.public.apiUrl,
    // Include credentials for session cookies
    onResponse() {
      // Common JSend unwrapping could also happen here,
      // but 'transform' is more specific to the returned 'data' property
    },
    transform: (response) => {
      if (response && response.status === "success") {
        return response.data as T;
      }
      return response as unknown as T;
    },
    ...options,
  });
};

/**
 * Custom $fetch wrapper for mutations (POST, PATCH, DELETE)
 */
export const api$fetch = <T>(
  url: string,
  options: Record<string, unknown> = {}
) => {
  const config = useRuntimeConfig();

  return $fetch<JSendResponse<T>>(url, {
    baseURL: config.public.apiUrl,
    ...options,
    onResponse({ response }) {
      if (response._data && response._data.status === "success") {
        response._data = response._data.data;
      }
    },
  });
};

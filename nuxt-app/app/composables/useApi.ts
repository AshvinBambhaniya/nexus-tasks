import type { UseFetchOptions } from "#app";
import type { JSendResponse } from "~/types";

export const useApi = <T>(
  url: string | (() => string),
  options: UseFetchOptions<JSendResponse<T>, T> = {}
) => {
  const config = useRuntimeConfig();
  const headers = useRequestHeaders(["cookie"]);

  // Use the internal API URL for SSR and the public API URL for the browser
  const baseURL = import.meta.server ? config.apiUrl : config.public.apiUrl;

  return useFetch(url, {
    baseURL: baseURL,
    credentials: "include",
    headers: headers,
    transform: (response) => {
      const res = response as JSendResponse<T>;
      if (res && res.status === "success") {
        return res.data as T;
      }
      return res as unknown as T;
    },
    ...options,
  });
};

/**
 * Custom $fetch wrapper for mutations (POST, PATCH, DELETE)
 */
export const useMutation = <T>(
  url: string,
  options: Record<string, unknown> = {}
) => {
  const config = useRuntimeConfig();

  // Use the internal API URL for SSR and the public API URL for the browser
  const baseURL = import.meta.server ? config.apiUrl : config.public.apiUrl;

  return $fetch<JSendResponse<T>>(url, {
    baseURL: baseURL,
    credentials: "include",
    ...options,
    onResponse({ response }) {
      const data = response._data as JSendResponse<T>;
      if (data && data.status === "success") {
        // We cast to T here to unwrap the data for the consumer.
        response._data = data.data as T;
      }
    },
  });
};

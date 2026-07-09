import type { ApiKey, CreateApiKeyResponse } from "~/types";

export const useApiKeys = () => {
  const {
    data: keys,
    pending: isLoading,
    error,
    refresh,
  } = useApi<ApiKey[]>("/api/v2/auth/api-keys", {
    key: "api-keys-list",
  });

  const createKey = async (name: string): Promise<CreateApiKeyResponse> => {
    try {
      const result = await useMutation<CreateApiKeyResponse>(
        "/api/v2/auth/api-keys",
        {
          method: "POST",
          body: { name },
        }
      );
      await refresh();
      return result;
    } catch (err) {
      console.error("Failed to create API key", err);
      throw err;
    }
  };

  const revokeKey = async (keyId: string) => {
    try {
      await useMutation(`/api/v2/auth/api-keys/${keyId}`, {
        method: "DELETE",
      });
      await refresh();
    } catch (err) {
      console.error("Failed to revoke API key", err);
      throw err;
    }
  };

  return {
    keys: computed(() => keys.value || []),
    isLoading,
    isError: !!error.value,
    createKey,
    revokeKey,
    refresh,
  };
};

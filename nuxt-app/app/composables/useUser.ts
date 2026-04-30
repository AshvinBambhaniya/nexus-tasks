export const useUser = () => {
  const config = useRuntimeConfig();
  
  const { data: user, pending: isLoading, error, refresh } = useFetch(`${config.public.apiUrl}/api/v1/auth/me`, {
    key: 'current-user',
    pick: ['id', 'email', 'full_name'], // Adjust based on your API response
    // JSend unwrap
    transform: (response: any) => {
      if (response && response.status === 'success') {
        return response.data;
      }
      return response;
    },
    // Include credentials for session cookie
    onResponseError({ response }) {
      if (response.status === 401) {
        // Handle unauthorized
      }
    }
  });

  return {
    user,
    isLoading,
    isError: !!error.value,
    refresh
  };
};

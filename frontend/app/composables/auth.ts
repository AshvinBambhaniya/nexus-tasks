import { useUsersStore } from "~/stores/user";
import type { User, JSendResponse } from "~/types";

export const setUserDataStore = async () => {
  const { apiUrl } = useRuntimeConfig().public;
  const { setUserData } = useUsersStore();

  try {
    const response = await fetch(apiUrl + "/api/v2/auth/me", {
      method: "GET",
      credentials: "include",
      mode: "cors",
    });

    if (response.status !== 200) {
      throw new Error(response.status.toString());
    } else if (response.status === 200) {
      const data = (await response.json()) as JSendResponse<User>;
      // Mapping the project's User type
      setUserData(data?.data || null);
    }
  } catch (error: unknown) {
    if (error instanceof Error && error.message === "401") {
      console.log(error.message);
      setUserData(null);
      return;
    }
    if (error instanceof Error) {
      console.log(error.message);
    }
  }
};

import { defineStore } from "pinia";
import type { User } from "~/types";

export const useUsersStore = defineStore(
  "users-store",
  () => {
    const userData = ref<User | null>(null);

    const setUserData = (data: User | null) => {
      userData.value = data;
    };

    const getUserData = () => {
      return userData.value;
    };

    return { userData, setUserData, getUserData };
  },
  {
    persist: true,
  }
);

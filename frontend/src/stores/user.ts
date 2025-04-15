// src/stores/auth.ts
import {defineStore} from "pinia";
import {ref} from "vue";
import {useAuthStore} from "./auth";

interface User {
  id: string;
  email: string;
  username: string;
  password: string | null;
  firstname: string;
  surname: string;
  created_at: Date;
}


export const useUserStore = defineStore("user", () => {

  const state = ref<User>({
    user: null,
    isLoading: true,
  });
  async function updateUser(user: User) {
    try {
      const response = await fetch(`http://localhost:8001/v1/user/${user.id}`, {
        method: "PUT",
        headers: { 'Content-Type': 'application/json' },
        credentials: "include",
        body: JSON.stringify({
          firstname: user.firstname,
          surname: user.surname,
          password: user.password,
        }),
      });
      console.log(response);

      if (!response.ok) {
        throw new Error('Failed to update user');
      }
      useAuthStore().state.user = user
    } catch (error) {
      console.error('Error updating user:', error);
      throw error;
    }
  }
  async function fetchUser(id) {
    try {
      state.value.isLoading = true;
      const response = await fetch(`http://localhost:8001/v1/user/${id}`, {
        credentials: "include",
      });
      if (!response.ok) throw new Error("Failed to fetch User");
      state.value.user = await response.json();
      console.log("fetch ended")
      state.value.isLoading = false;
      return response;
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  }

  return { state, updateUser, fetchUser };
});

// src/stores/auth.ts
import { defineStore } from "pinia";
import { ref, computed, watch } from "vue";
import Cookies from "js-cookie";

interface User {
  id: number;
  email: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  exp: number | null;
  isLoading: boolean;
}

const STORAGE_KEY = "auth_data";

export const useAuthStore = defineStore("auth", () => {

  const state = ref<AuthState>({
    user: null,
    token: null,
    exp: null,
    isLoading: true,
  });
  async function fetchUser() {
    try {
      const response = await fetch("http://localhost:8001/auth/me", {
        method: "POST",
        credentials: "include",
      });
      console.log(response);

      if (!response.ok) {
        clearAuthData();
        return;
      }

      const data = await response.json();
      console.log(data);

      state.value.user = data?.user;
      state.value.exp = data?.exp;
      state.value.isLoading = false;
    } catch (error) {
      console.error("Session validation error:", error);
      clearAuthData();
    }
 }



  const isAuthenticated = computed(() => !!state.value.user && !isJwtExpired());
  const isLoading = computed(() => state.value.isLoading);

  function isJwtExpired(): boolean {
    if (!state.value.exp) return true;
    const currentTime = Math.floor(Date.now() / 1000);
    return currentTime >= state.value.exp;
  }

  async function signIn(email: string, password: string) {
    try {
      const response = await fetch("http://localhost:8001/auth/signin", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        throw new Error((await response.text()) || "Something went wrong");
      }

      const data = await response.json();
      state.value.user = data.user;
      state.value.token = data.token;
      state.value.exp = data.exp;

      saveAuthData();
    } catch (error) {
      console.error("Sign in error:", error);
      throw error;
    }
  }

  async function signOut() {
    try {
      const response = await fetch("http://localhost:8001/auth/signout", {
        method: "POST",
      });

      if (!response.ok) {
        throw new Error((await response.text()) || "Something went wrong");
      }

      clearAuthData();
    } catch (error) {
      console.error("Sign out error:", error);
      throw error;
    }
  }

  function saveAuthData() {
    Cookies.set("jwt", state.value.token);
  }

  function clearAuthData() {
    state.value.user = null;
    state.value.token = null;
    state.value.exp = null;
    state.value.isLoading = false;
    Cookies.remove("jwt");
  }

  watch(state, () => {
    if (state.value.user) {
      saveAuthData();
    } else {
      clearAuthData();
    }
  });
  fetchUser()
  return { state, isAuthenticated, isLoading, signIn, signOut };
});

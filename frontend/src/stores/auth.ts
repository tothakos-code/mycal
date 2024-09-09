// src/stores/auth.ts
import { defineStore } from "pinia";
import { ref, computed } from "vue";

interface User {
  id: number;
  email: string;
}

interface AuthState {
  user: User | null;
  exp: number | null;
}

export const useAuthStore = defineStore("auth", () => {
  const state = ref<AuthState>({
    user: null,
    exp: null,
  });

  const isAuthenticated = computed(() => !!state.value.user && !isJwtExpired());

  function isJwtExpired(): boolean {
    if (!state.value.exp) return true;
    const currentTime = Math.floor(Date.now() / 1000);
    return currentTime >= state.value.exp;
  }

  async function signIn(email: string, password: string) {
    try {
      const response = await fetch("http://localhost:8080/auth/signin", {
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
      state.value.exp = data.exp;
    } catch (error) {
      console.error("Sign in error:", error);
      throw error;
    }
  }

  async function signOut() {
    try {
      const response = await fetch("http://localhost:8080/auth/signout", {
        method: "POST",
      });

      if (!response.ok) {
        throw new Error((await response.text()) || "Something went wrong");
      }

      state.value.user = null;
      state.value.exp = null;
    } catch (error) {
      console.error("Sign out error:", error);
      throw error;
    }
  }

  return { state, isAuthenticated, signIn, signOut };
});

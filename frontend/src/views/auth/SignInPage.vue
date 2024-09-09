<template>
  <main>
    <router-link class="home-link" to="/">◄ Home</router-link>
    <form class="main-container" @submit.prevent="handleSubmit">
      <h1 class="header-text">Sign In</h1>
      <input
        v-model="formValues.email"
        name="email"
        type="email"
        placeholder="Email"
      />
      <input
        v-model="formValues.password"
        name="password"
        type="password"
        placeholder="Password"
      />
      <button type="submit">Login</button>
      <router-link class="auth-link" to="/auth/sign-up">
        Don't have an account? Sign Up
      </router-link>
      <p v-if="error" style="color: red">{{ error }}</p>
    </form>
  </main>
</template>

<script lang="ts" setup>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../../stores/auth";

const router = useRouter();
const authStore = useAuthStore();

// Redirect if user is already logged in
if (authStore.isAuthenticated) {
  router.push("/");
}

const error = ref("");
const formValues = reactive({
  email: "",
  password: "",
});

const handleSubmit = async () => {
  try {
    await authStore.signIn(formValues.email, formValues.password);
    // Will automatically redirect to authenticated page when successful
    router.push("/");
  } catch (err) {
    console.error(err);
    error.value = err instanceof Error ? err.message : "An error occurred";
    // Show a Toast message or render an error message or something
  }
};
</script>

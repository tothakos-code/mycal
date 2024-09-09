<template>
  <main>
    <router-link class="home-link" to="/">◄ Home</router-link>
    <form class="main-container" @submit.prevent="handleSubmit">
      <h1 class="header-text">Sign Up</h1>
      <p class="demo-text">
        Demo app, please don't use your real email or password
      </p>
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
      <button type="submit">Create Account</button>
      <router-link class="auth-link" to="/auth/sign-in">
        Already have an account? Sign In
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
    const response = await fetch("http://localhost:8080/auth/signup", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formValues),
    });

    if (!response.ok) {
      const errorMessage = (await response.text()) || "Something went wrong";
      throw new Error(errorMessage);
    }

    // If successful registered, navigate to sign-in page
    router.push("/auth/sign-in");
  } catch (err) {
    console.error(err);
    error.value = err instanceof Error ? err.message : "An error occurred";
    // Show a Toast message or render an error message or something
  }
};
</script>

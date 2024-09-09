<template>
  <router-view></router-view>
</template>

<script lang="ts" setup>
import { onMounted, watch } from "vue";
import { useAuthStore } from "./stores/auth";
import { useRouter } from "vue-router";

const authStore = useAuthStore();
const router = useRouter();

onMounted(() => {
  if (authStore.isAuthenticated) {
    // Check if JWT is expired on app load
    if (
      authStore.state.exp &&
      authStore.state.exp <= Math.floor(Date.now() / 1000)
    ) {
      authStore.signOut();
      router.push("/auth/sign-in");
    }
  }
});

watch(
  () => authStore.isAuthenticated,
  (newValue) => {
    if (!newValue) {
      router.push("/auth/sign-in");
    }
  }
);
</script>

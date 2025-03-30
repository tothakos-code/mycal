<template>
  <main>
    <router-link class="home-link" to="/">◄ Home</router-link>
    <section class="main-container">
      <h1 class="header-text">User profile</h1>
      <p>Email : {{ state.user?.email || "None" }}</p>
      <p>Username : {{ state.user?.username || "None" }}</p>
      <p>Firstname : {{ state.user?.firstname || "None" }}</p>
      <p>Surname : {{ state.user?.surname || "None" }}</p>
    </section>
  </main>
</template>

<script lang="ts" setup>
import { useAuthStore } from "../stores/auth";

const authStore = useAuthStore();
const { state } = authStore;
</script>

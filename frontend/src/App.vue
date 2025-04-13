<template>
  <v-app>
  <v-app-bar
  class="px-3"
  density="compact"
  flat
>
  <v-avatar
    class="hidden-md-and-up"
    color="grey-darken-1"
    size="32"
  ></v-avatar>

  <v-spacer></v-spacer>

  <v-tabs
    align-tabs="center"
    color="grey-darken-2"
  >
    <v-tab
      v-if="state.user"
      text="Sign Out"
      @click="signOut"
    ></v-tab>
    <v-tab
      v-else
      tag="router-link"
      to="/auth/sign-in"
      text="Sign In"
    ></v-tab>
    <v-tab
      tag="router-link"
      to="/"
      text="Home"
    ></v-tab>
    <v-tab
      tag="router-link"
      to="/profile"
      text="Profile Page 🛡️"
    ></v-tab>
  </v-tabs>
  <v-spacer></v-spacer>

  <v-avatar
    class="hidden-sm-and-down"
    color="grey-darken-1"
    size="32"
  ></v-avatar>
</v-app-bar>
<v-main class="bg-grey-lighten-3">
  <router-view></router-view>
  </v-main>
</v-app>

</template>

<script lang="ts" setup>
import { onMounted, watch } from "vue";
import { useAuthStore } from "./stores/auth";
import { useRouter } from "vue-router";

const authStore = useAuthStore();
const router = useRouter();
const { state, signOut } = authStore;

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

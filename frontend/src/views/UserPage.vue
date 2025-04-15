<template>
  <div class="container mx-auto p-4">
    <h2 class="text-2xl font-bold mb-4">User Profile</h2>

    <v-card v-if="!userStore.isLoading" class="p-6" elevation="2">

         <v-text-field
            v-model="userStore.user.username"
            label="Username"
            readonly
            required
        />
        <v-text-field
            v-model="userStore.user.firstname"
            label="First Name"
            :readonly="!isEditing"
        />
        <v-text-field
            v-model="userStore.user.surname"
            label="Surname"
            :readonly="!isEditing"
        />

      <div class="flex justify-between mt-6">
        <div class="flex gap-2">
          <v-btn color="primary" to="/" >Back</v-btn>


        </div>
      </div>
    </v-card>
  </div>
</template>

<script lang="ts" setup>
import { useAuthStore } from "../stores/auth";
import { useUserStore } from "../stores/user";
import { ref, reactive, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const userId = route.params.id as string

console.log(userId)

const { state: authStore } = useAuthStore();
const { state: userStore, fetchUser } = useUserStore();
fetchUser(userId);
</script>
<style scoped>
.container {
  max-width: 600px;
}
</style>
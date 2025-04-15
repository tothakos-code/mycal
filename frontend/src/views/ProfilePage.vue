<template>
  <div class="container mx-auto p-4">
    <h2 class="text-2xl font-bold mb-4">User Profile</h2>

    <v-card class="p-6" elevation="2">
      <v-form ref="formRef" v-model="formValid" lazy-validation>
        <v-text-field
            v-model="editedUser.email"
            label="Email"
            readonly
            required
        />
        <v-text-field
            v-model="editedUser.username"
            label="Username"
            readonly
            required
        />
        <v-text-field
            v-model="editedUser.firstname"
            label="First Name"
            :readonly="!isEditing"
        />
        <v-text-field
            v-model="editedUser.surname"
            label="Surname"
            :readonly="!isEditing"
        />

        <div v-if="isEditing" class="mt-6">
          <h3 class="text-lg font-semibold mb-2">Change Password</h3>
          <v-text-field
              v-model="password1"
              label="New Password"
              type="password"
          />
          <v-text-field
              v-model="password2"
              label="Confirm New Password"
              type="password"
              :error="password1 !== password2"
              :error-messages="password1 !== password2 ? 'Passwords do not match' : ''"
          />
        </div>
      </v-form>

      <div class="flex justify-between mt-6">
        <div>
          <v-btn v-if="!isEditing" @click="isEditing = true" color="primary">
            Edit
          </v-btn>
          <v-btn v-else @click="cancelEdit" color="info" variant="outlined">
            Cancel
          </v-btn>
        </div>
        <div class="flex gap-2">
          <v-btn
              v-if="isEditing"
              :disabled="!formValid || (password1 !== password2 && password1)"
              color="success"
              @click="saveUser"
          >
            Save
          </v-btn>

          <v-dialog v-model="deleteDialog" max-width="400">
            <template #activator="{ props }">
              <v-btn color="error" v-bind="props">Delete</v-btn>
            </template>
            <v-card>
              <v-card-title class="text-h6">Confirm Deletion</v-card-title>
              <v-card-text>Are you sure you want to delete this account? This action is irreversible.</v-card-text>
              <v-card-actions>
                <v-spacer />
                <v-btn text @click="deleteDialog = false">Cancel</v-btn>
                <v-btn color="error" @click="confirmDelete">Delete</v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </div>
      </div>
    </v-card>
  </div>
</template>

<script lang="ts" setup>
import { useAuthStore } from "../stores/auth";
import { useUserStore } from "../stores/user";
import { ref, reactive, watch } from 'vue'


const { state: authStore } = useAuthStore();
const { state: userStore } = useUserStore();

const isEditing = ref(false)
const deleteDialog = ref(false)
const formValid = ref(true)

const editedUser = { ...authStore.user }
const password1 = ref('')
const password2 = ref('')
const formRef = ref()


function cancelEdit() {
  Object.assign(editedUser, authStore.user)
  password1.value = ''
  password2.value = ''
  isEditing.value = false
}

function saveUser() {
  // Example: update the user info in your backend
  console.log('Saving user:', editedUser);
  editedUser.password = password1.value;
  useUserStore().updateUser(editedUser);
  if (password1.value) {
    console.log('New password:', password1.value);
  }

  // Call your API here
  isEditing.value = false;
}

function confirmDelete() {
  console.log('Deleting user...')
  deleteDialog.value = false
  // Call your delete API here
}

function formattedDate(date: Date | string): string {
  return new Date(date).toLocaleString('HU')
}

</script>
<style scoped>
.container {
  max-width: 600px;
}
</style>
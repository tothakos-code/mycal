<template>
      <div class="container mx-auto p-4">
          <h2 class="text-2xl font-bold mb-4">Events & Invitations</h2>

          <!-- Create Event Form -->
          <div v-if="showCreateForm" class="bg-gray-100 p-4 rounded mb-4">
            <h3 class="text-lg font-semibold mb-2">Create New Event</h3>
            <form @submit.prevent="createEvent">
              <v-text-field v-model="newEvent.title" label="Title" required />

              <v-textarea v-model="newEvent.description" label="Description" required />

              <v-text-field v-model="newEvent.location" label="Location" required />

              <v-date-picker
                v-model="newEvent.start"
                show-adjacent-months
                title="Start Time"
              ></v-date-picker>

              <v-date-picker
                v-model="newEvent.finish"
                show-adjacent-months
                title="End Time"
              ></v-date-picker>

              <v-text-field label="Notify Before (minutes)" v-model="newEvent.notify_before" type="number" required />
              <v-switch label="Public" v-model="newEvent.is_public" />

              <v-btn color="light-blue" @click="showCreateForm = false">
                Back
              </v-btn>
              <v-btn type="submit" color="light-blue">
                Create Event
              </v-btn>
            </form>
          </div>

          <!-- Event List -->
          <div v-if="eventStore.isLoading" class="text-gray-500">Loading...</div>

          <v-list v-else>
            <v-list-item
              v-for="event in sortedEvents"
              :key="event.id"
              :value="event.id"
              class="py-3"
              @click="openEvent(event)"
            >
              <v-list-item-title>{{ event.title || "No Title" }}</v-list-item-title>

              <v-list-item-subtitle class="mb-1 text-high-emphasis opacity-100">{{ new Date(event.start).toLocaleString('HU') }}</v-list-item-subtitle>

              <v-list-item-subtitle class="text-high-emphasis">{{ event.description }}</v-list-item-subtitle>

              <span v-if="event.is_public" class="text-green-600">Public</span>
              <span v-else class="text-green-600">Private</span>
              <br>
              <span class="text-green-600">{{ event.user.firstname }} {{ event.user.surname }}</span>

              <template v-slot:append="{  }">
                <v-list-item-action class="flex-column align-end">
                  <small class="mb-4 text-high-emphasis opacity-60">{{ event.action }}</small>

                  <v-spacer></v-spacer>

                  <v-icon color="green-darken-3">mdi-check-circle-outline</v-icon>
                  <v-icon color="yellow-darken-3">mdi-help-circle-outline</v-icon>
                  <v-icon color="red-darken-3">mdi-close-circle-outline</v-icon>
                </v-list-item-action>
              </template>
            </v-list-item>
          </v-list>

          <div v-if="!eventStore.isLoading && sortedEvents.length === 0" class="text-gray-500">
            No upcoming events.
          </div>
          <v-fab
            v-if="authStore.isAuthenticated"
            @click="showCreateForm = !showCreateForm" icon="mdi-plus"
            location="right bottom"
            app
            size="70"
          ></v-fab>
        </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { useAuthStore } from "../stores/auth";
import { useEventStore, Event } from "../stores/event";

const authStore = useAuthStore();
const eventStore = useEventStore();

const invitations = ref([]);
const showCreateForm = ref(false);

// Form data for creating an event
const newEvent: Event = ref({
  title: "",
  description: "",
  location: "",
  start: "",
  finish: "",
  notify_before: 10,
  is_public: false,
});

function openEvent(event: Event) {
  eventStore.setEvent(event)
}

// Computed: Merge events & pending invitations and sort by date
const sortedEvents = computed(() => {
  return [...eventStore.eventsPub, ...eventStore.eventsPriv || []].sort(
    (a, b) => new Date(a.start) - new Date(b.start)
  );
});

// Create a new event
async function createEvent() {
  try {
    const response = await eventStore.createEvent(newEvent.value);
    if (!response.ok) {
      throw new Error("Failed to create event");
    }

    newEvent.value = {
      user_id: "",
      title: "",
      description: "",
      location: "",
      start: "",
      finish: "",
      is_public: false,
      notify_before: 0,
    };

    showCreateForm.value = false;
    await eventStore.fetchEvents();
    await eventStore.fetchUserEvents();
  } catch (error) {
    console.error("Error creating event:", error);
  }
}

// Load data on mount
eventStore.fetchEvents();
if (authStore.user) eventStore.fetchUserEvents();
</script>

<style scoped>
.container {
  max-width: 600px;
}
</style>

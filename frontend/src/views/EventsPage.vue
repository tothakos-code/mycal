<template>
  <main>
    <router-link class="home-link" to="/">◄ Home</router-link>
    <section class="main-container">
      <div class="container mx-auto p-4">
          <h2 class="text-2xl font-bold mb-4">Events & Invitations</h2>

          <button
            v-if="authStore.isAuthenticated"
            @click="showCreateForm = !showCreateForm"
            class="bg-blue-500 text-white px-4 py-2 rounded mb-4"
          >
            {{ showCreateForm ? "Cancel" : "Create Event" }}
          </button>

          <!-- Create Event Form -->
          <div v-if="showCreateForm" class="bg-gray-100 p-4 rounded mb-4">
            <h3 class="text-lg font-semibold mb-2">Create New Event</h3>
            <form @submit.prevent="createEvent">
              <label class="block">Title:</label>
              <input v-model="newEvent.title" class="border p-2 w-full mb-2" required />

              <label class="block">Description:</label>
              <input v-model="newEvent.description" class="border p-2 w-full mb-2" required />

              <label class="block">Location:</label>
              <input v-model="newEvent.location" class="border p-2 w-full mb-2" required />

              <label class="block">Start Time:</label>
              <input type="datetime-local" v-model="newEvent.start" class="border p-2 w-full mb-2" required />

              <label class="block">End Time:</label>
              <input type="datetime-local" v-model="newEvent.finish" class="border p-2 w-full mb-2" required />

              <label class="block">Notify Before (minutes):</label>
              <input type="number" v-model="newEvent.notify_before" class="border p-2 w-full mb-2" required />

              <label class="block">Public</label>
              <input type="checkbox" v-model="newEvent.is_public" class="border p-2 w-full mb-2" />

              <button type="submit" class="bg-green-500 text-white px-4 py-2 rounded">
                Create Event
              </button>
            </form>
          </div>

          <!-- Event List -->
          <div v-if="isLoading" class="text-gray-500">Loading...</div>

          <ul v-else class="space-y-4">
            <li v-for="event in sortedEvents" :key="event.id" class="p-4 bg-white shadow rounded-lg">
              <h3 class="text-lg font-semibold">
                {{ event.title || "No Title" }}
              </h3>
              <p class="text-gray-600">{{ new Date(event.start).toLocaleString('HU') }}</p>
              <span v-if="event.is_public" class="text-green-600">Public</span>
              <span v-else class="text-green-600">Private</span>
            </li>
          </ul>

          <div v-if="!isLoading && sortedEvents.length === 0" class="text-gray-500">
            No upcoming events.
          </div>
        </div>
    </section>
  </main>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "../stores/auth";

const authStore = useAuthStore();

const events = ref([]);
const invitations = ref([]);
const isLoading = ref(true);
const showCreateForm = ref(false);

// Form data for creating an event
const newEvent = ref({
  title: "",
  description: "",
  location: "",
  start: "",
  finish: "",
  notify_before: 0,
  is_public: false,
});

// Fetch events and invitations from the API
async function fetchEventsAndInvitations() {
  try {
    isLoading.value = true;

    // Fetch last 30 events
    const eventsRes = await fetch("http://localhost:8001/v1/event/last-30", {
      credentials: "include",
    });
    if (!eventsRes.ok) throw new Error("Failed to fetch events");
    events.value = await eventsRes.json();

    // // Fetch pending invitations if authenticated
    // if (authStore.isAuthenticated) {
    //   const invitationsRes = await fetch("http://localhost:8001/invitations/pending", {
    //     credentials: "include",
    //   });
    //   if (!invitationsRes.ok) throw new Error("Failed to fetch invitations");
    //   invitations.value = await invitationsRes.json();
    // }
  } catch (error) {
    console.error("Error fetching data:", error);
  } finally {
    isLoading.value = false;
  }
}

// Computed: Merge events & pending invitations and sort by date
const sortedEvents = computed(() => {
  return [...events.value, ...invitations.value].sort(
    (a, b) => new Date(a.start) - new Date(b.start)
  );
});

// Create a new event
async function createEvent() {
  try {
    const response = await fetch("http://localhost:8001/v1/event", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        title: newEvent.value.title,
        description: newEvent.value.description,
        location: newEvent.value.location,
        start: new Date(newEvent.value.start).toISOString(),
        finish: new Date(newEvent.value.finish).toISOString(),
        is_public: newEvent.value.is_public,
        notify_before: parseInt(newEvent.value.notify_before),
      }),
    });

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
      is_public: False,
      notify_before: 0,
    };

    showCreateForm.value = false;
    fetchEventsAndInvitations(); // Refresh events
  } catch (error) {
    console.error("Error creating event:", error);
  }
}

// Load data on mount
onMounted(fetchEventsAndInvitations);
</script>

<style scoped>
.container {
  max-width: 600px;
}
</style>

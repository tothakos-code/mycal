<template>
  <div class="container mx-auto p-4">
    <v-row>
      <!-- Title -->
      <v-col cols="12">
        <h2 class="text-2xl font-bold mb-4">Event</h2>
      </v-col>

      <!-- Event Details -->
      <v-col cols="12" v-if="!isEditing">
        <v-list>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title><strong>Title:</strong> {{ eventStore.event.title }}</v-list-item-title>
              <v-list-item-subtitle>{{ eventStore.event.description }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title><strong>Location:</strong> {{ eventStore.event.location }}</v-list-item-title>
              <v-list-item-subtitle><strong>Start:</strong> {{ new Date(eventStore.event.start).toLocaleString('HU') }}</v-list-item-subtitle>
              <v-list-item-subtitle><strong>Finish:</strong> {{ new Date(eventStore.event.finish).toLocaleString('HU') }}</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title><strong>Notify Before:</strong> {{ eventStore.event.notify_before }} minutes</v-list-item-title>
              <v-list-item-title><strong>Public:</strong> {{ eventStore.event.is_public ? 'Yes' : 'No' }}</v-list-item-title>
              <v-list-item-title><strong>Created At:</strong> {{ new Date(eventStore.event.created_at).toLocaleString('HU') }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-col>

      <!-- Editable Fields -->
      <v-col cols="12" v-if="isEditing">
        <v-form ref="eventForm" v-model="formValid">
          <v-text-field label="Title" v-model="editedEvent.title" :rules="[rules.required]" />
          <v-textarea label="Description" v-model="editedEvent.description" :rules="[rules.required]" />
          <v-text-field label="Location" v-model="editedEvent.location" :rules="[rules.required]" />
          <v-text-field label="Start" v-model="editedEvent.start" :rules="[rules.required]" />
          <v-text-field label="Finish" v-model="editedEvent.finish" :rules="[rules.required]" />
          <v-text-field label="Notify Before" v-model="editedEvent.notify_before" type="number" :rules="[rules.required]" />
          <v-switch label="Public" v-model="editedEvent.is_public" />
        </v-form>
      </v-col>

      <!-- Buttons -->
      <v-col cols="12" class="d-flex justify-start">
        <v-btn color="primary" @click="goBack">Back</v-btn>
        <v-btn v-if="!isEditing" color="blue" @click="editEvent">Edit</v-btn>
        <v-btn v-if="isEditing" color="red" @click="deleteEvent">Delete</v-btn>
        <v-btn v-if="isEditing" color="green" @click="saveEvent">Save</v-btn>
        <v-btn v-if="!isEditing" color="teal" @click="invitePeople">Invite People</v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {useEventStore, Event} from '../stores/event'

// Store and data initialization
const eventStore = useEventStore()
const isEditing = ref(false)
const formValid = ref(false)
const editedEvent: Event = { ...eventStore.event }

// Validation rules for form fields
const rules = {
  required: (value: any) => !!value || 'This field is required.',
}

// Methods
const editEvent = () => {
  isEditing.value = true
  const {event} = eventStore;
  editedEvent.value = { ...event }
}

const saveEvent = async () => {
  // Call the API to save the event (empty function body for now)
  // API call logic here to update the event

  console.log(editedEvent)
  editedEvent.start = new Date(editedEvent.start)
  editedEvent.finish = new Date(editedEvent.finish)
  await eventStore.editEvent(eventStore.event.id, editedEvent)
  eventStore.event = {...editedEvent}
  isEditing.value = false
}

const deleteEvent = async () => {
  // Call the API to delete the event (empty function body for now)
  // API call logic here to delete the event
  eventStore.deleteEvent(eventStore.event.id)
  await eventStore.fetchEvents();
  eventStore.event = null // Simulate event deletion
}

const invitePeople = () => {
  // Placeholder for invite people logic (empty function body for now)
  console.log('Invite people functionality not implemented yet.')
}

const goBack = () => {
  if (isEditing.value) {
    isEditing.value = false
    editedEvent.value = { ...eventStore.event } // Reset the form to the current event data
  } else {
    eventStore.event = null
  }
}
</script>

<style scoped>
.container {
  max-width: 600px;
}

.v-btn {
  margin-left: 8px;
}
</style>

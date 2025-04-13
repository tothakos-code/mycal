<template>
  <v-container>
    <v-row>
      <v-col
        cols="12"
        md="2"
      >
        <v-sheet
          min-height="268"
          rounded="lg"
        >
          <!--  -->
        </v-sheet>
      </v-col>

      <v-col
        cols="12"
        md="8"
      >
        <v-sheet
          min-height="70vh"
          rounded="lg"
        >
          <Events v-if="!eventStore.event"/>
          <Event v-else />
        </v-sheet>
      </v-col>

      <v-col
        cols="12"
        md="2"
      >
        <v-sheet
          min-height="268"
          rounded="lg"
        >
          <!--  -->
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import { useEventStore } from "../stores/event";
import Events from './EventsPage.vue'
import Event from './EventPage.vue'

const eventStore = useEventStore();
const isLoading = ref(true);

</script>

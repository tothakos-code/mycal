import { defineStore } from "pinia";

export interface Event {
  title: string,
  user: User,
  description: string,
  location: string,
  start: Date,
  finish: Date,
  is_public: boolean,
  notify_before: number
}

export const useEventStore = defineStore('event', {
  state: () => ({
    event: null,
    eventsPub: [],
    eventsPriv: [],
    isLoading: true
  }),
  actions: {
    setEvent(event: any) {
      this.event = event
    },
    async fetchEvents() {
      try {
        this.isLoading = true;

        // Fetch last 30 events
        const eventsRes = await fetch("http://localhost:8001/v1/event/last-30", {
          credentials: "include",
        });
        if (!eventsRes.ok) throw new Error("Failed to fetch events");
        this.eventsPub = await eventsRes.json();
        this.isLoading = true;
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        this.isLoading = false;
      }
    },
    async fetchUserEvents() {
      try {
        this.isLoading = true;

        // Fetch last 30 events
        const eventsRes = await fetch("http://localhost:8001/v1/event/last-30-private", {
          credentials: "include",
        });
        if (!eventsRes.ok) throw new Error("Failed to fetch events");
        this.eventsPriv = await eventsRes.json();
        this.isLoading = true;
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        this.isLoading = false;
      }
    },
    async createEvent(newEvent: Event) {
      try {
        const response = await fetch('http://localhost:8001/v1/event', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({
            title: newEvent.title,
            description: newEvent.description,
            location: newEvent.location,
            start: newEvent.start.toISOString(),
            finish: newEvent.finish.toISOString(),
            is_public: newEvent.is_public,
            notify_before: newEvent.notify_before,
          }),
        });

        if (!response.ok) {
          throw new Error('Failed to create event');
        }

        return response;
      } catch (error) {
        console.error('Error creating event:', error);
        throw error;
      }
    },
    async editEvent(eventId: string, updatedEventData: Event) {
      try {
        console.log(updatedEventData)
        const response = await fetch(`http://localhost:8001/v1/event/${eventId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({
            title: updatedEventData.title,
            description: updatedEventData.description,
            location: updatedEventData.location,
            start: updatedEventData.start.toISOString(),
            finish: updatedEventData.finish.toISOString(),
            is_public: updatedEventData.is_public,
            notify_before: updatedEventData.notify_before,
          }),
        });
        if (!response.ok) {
          throw new Error('Failed to edit event');
        }
        return;
      } catch (error) {
        console.error('Error editing event:', error);
        throw error;
      }
    },

    async deleteEvent(eventId: string) {
      try {
        const response = await fetch(`http://localhost:8001/v1/event/${eventId}`, {
          method: 'DELETE',
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Failed to delete event');
        }

        return response;
      } catch (error) {
        console.error('Error deleting event:', error);
        throw error;
      }
    },
  }
})

import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import HomePage from "../views/HomePage.vue";
import SignInPage from "../views/auth/SignInPage.vue";
import SignUpPage from "../views/auth/SignUpPage.vue";
import ProfilePage from "../views/ProfilePage.vue";
import UserPage from "../views/UserPage.vue";
import EventPage from "../views/EventPage.vue";
import NotFoundPage from "../views/NotFoundPage.vue";
import { useAuthStore } from "../stores/auth";
import { watch } from "vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    component: HomePage,
  },
  {
    path: "/auth/sign-in",
    component: SignInPage,
  },
  {
    path: "/auth/sign-up",
    component: SignUpPage,
  },
  {
    path: "/event",
    component: EventPage,
    name: 'EventPage',
    meta: { requiresAuth: true },
  },
  {
    path: "/profile",
    component: ProfilePage,
    meta: { requiresAuth: true },
  },  {
    path: "/user/:id",
    component: UserPage,
  },
  {
    path: "/:pathMatch(.*)*",
    component: NotFoundPage,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _, next) => {
  const authStore = useAuthStore();
  if (authStore.isLoading) {
      const stopWatching = watch(
        () => authStore.isLoading,
        (isLoading) => {
          if (!isLoading) {
            stopWatching(); // Stop watching after it's done loading
            proceedWithAuthCheck();
          }
        },
        { immediate: true }
      );
    } else {
      proceedWithAuthCheck();
    }

    function proceedWithAuthCheck() {
      if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        next("/auth/sign-in");
      } else {
        next();
      }
    }
});

export default router;

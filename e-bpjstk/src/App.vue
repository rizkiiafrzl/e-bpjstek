<template>
  <v-app>
    <!-- App Bar untuk halaman yang memerlukan header -->
    <v-app-bar v-if="showHeader" app color="primary" dark elevation="2" class="header-bar">
      <v-app-bar-nav-icon></v-app-bar-nav-icon>

      <v-toolbar-title class="header-title">
        <div class="d-flex align-center">
          <v-img
            src="@/assets/logo.svg"
            alt="Vue logo"
            class="logo"
            max-height="32"
            max-width="32"
          />
          <span class="title-text">E-BPJS Ketenagakerjaan</span>
        </div>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <div class="nav-buttons">
        <v-btn
          v-for="item in menuItems"
          :key="item.name"
          :to="item.to"
          variant="text"
          class="nav-btn"
        >
          <v-icon left size="20">{{ item.icon }}</v-icon>
          {{ item.name }}
        </v-btn>
      </div>
    </v-app-bar>

    <!-- Main Content -->
    <v-main>
      <RouterView />
    </v-main>

    <!-- Footer -->
    <v-footer v-if="showHeader" app color="grey-lighten-1" class="text-center">
      <v-col cols="12">
        {{ new Date().getFullYear() }} — <strong>E-BPJS Ketenagakerjaan</strong>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script setup>
import { RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()

// Hide header on auth pages
const showHeader = computed(() => {
  const authRoutes = ['/login', '/register', '/dashboard']
  return !authRoutes.includes(route.path)
})

// Menu items untuk navigation
const menuItems = [
  { name: 'Home', to: '/', icon: 'mdi-home' },
  { name: 'About', to: '/about', icon: 'mdi-information' },
  { name: 'Dashboard', to: '/dashboard', icon: 'mdi-view-dashboard' },
]
</script>

<style scoped>
.logo {
  height: 32px;
  width: 32px;
}

.header-bar {
  padding: 0 16px;
}

.header-title {
  font-weight: 600;
  font-size: 1.1rem;
}

.header-title .d-flex {
  align-items: center;
  gap: 12px;
}

.title-text {
  font-weight: 600;
  font-size: 1.1rem;
  letter-spacing: 0.5px;
  line-height: 1.4;
  vertical-align: middle;
}

.nav-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-btn {
  font-weight: 500;
  font-size: 0.9rem;
  text-transform: none;
  letter-spacing: 0.5px;
  min-width: auto;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
  transform: translateY(-1px);
}

/* Mobile responsive */
@media (max-width: 768px) {
  .nav-buttons {
    gap: 4px;
  }

  .nav-btn {
    padding: 6px 12px;
    font-size: 0.8rem;
  }

  .nav-btn .v-icon {
    display: none;
  }
}

@media (max-width: 480px) {
  .header-title {
    font-size: 1rem;
  }

  .nav-btn {
    padding: 4px 8px;
    font-size: 0.75rem;
  }
}
</style>

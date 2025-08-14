<template>
  <div id="app-container" class="bg-gray-50 min-h-screen font-sans">
    <header class="bg-white shadow-md sticky top-0 z-50">
      <nav class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <router-link to="/" class="flex items-center space-x-2">
            <svg
              class="h-8 w-8 text-blue-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a2 2 0 01-2-2V7a2 2 0 012-2h4M9 12h6"
              />
            </svg>
            <h1 class="text-2xl font-bold text-gray-800">LingoLift</h1>
          </router-link>

          <div class="flex items-center space-x-4">
            <router-link
              to="/"
              class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              exact-active-class="text-blue-600 bg-blue-50"
              >Transcriber</router-link
            >
            <router-link
              to="/history"
              class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              active-class="text-blue-600 bg-blue-50"
              >History</router-link
            >
            <router-link
              to="/settings"
              class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              active-class="text-blue-600 bg-blue-50"
              >Settings</router-link
            >
          </div>
        </div>
      </nav>
    </header>
    <main class="container mx-auto p-4 sm:p-6 lg:p-8">
      <router-view />
    </main>
  </div>
</template>

<script setup>
// No script changes needed for this layout update.
</script>

<style>
/* Global styles are now primarily managed by Tailwind in index.css */
body {
  background-color: #f9fafb; /* bg-gray-50 */
}
</style>

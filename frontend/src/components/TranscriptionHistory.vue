<template>
  <div class="container mx-auto p-4 md:p-6">
    <div class="bg-white shadow-xl rounded-lg">
      <header
        class="border-b pb-4 p-4 md:p-6 flex justify-between items-center"
      >
        <h2 class="text-2xl font-bold text-gray-800 tracking-tight">
          Transcription History
        </h2>
        <router-link
          to="/"
          class="px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 active:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-300"
        >
          Back to Transcriber
        </router-link>
      </header>

      <div class="p-4 md:p-6">
        <div v-if="transcriptions.length > 0" class="space-y-4">
          <div
            v-for="item in transcriptions"
            :key="item.id"
            class="p-4 bg-gray-50 rounded-lg border border-gray-200 hover:shadow-md transition-shadow duration-200"
          >
            <div
              class="flex flex-col sm:flex-row justify-between items-start gap-4"
            >
              <div class="flex-grow overflow-hidden">
                <h3 class="font-semibold text-lg text-gray-800 mb-1 truncate">
                  {{ item.fileName || "Transcription" }}
                </h3>
                <p class="text-sm text-gray-500">
                  {{ new Date(item.timestamp).toLocaleString() }}
                </p>
                <p class="mt-2 text-gray-700 leading-relaxed line-clamp-2">
                  {{ item.transcript }}
                </p>
              </div>
              <div
                class="mt-2 sm:mt-0 flex-shrink-0 flex flex-row sm:flex-col gap-2 items-stretch"
              >
                <button
                  @click="viewTranscription(item.id)"
                  class="px-3 py-1.5 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 text-center"
                >
                  View
                </button>
                <div class="relative">
                  <button
                    @click="toggleExportMenu(item.id)"
                    class="w-full px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 text-center"
                  >
                    Export
                  </button>
                  <div
                    v-if="activeExportMenu === item.id"
                    class="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-md shadow-lg z-10"
                  >
                    <a
                      @click.prevent="handleExport(item, 'txt')"
                      href="#"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >Export as Text (.txt)</a
                    >
                    <a
                      @click.prevent="handleExport(item, 'srt')"
                      href="#"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                      >Export as Timed Text (.txt)</a
                    >
                  </div>
                </div>
                <button
                  @click="deleteTranscription(item.id)"
                  class="px-3 py-1.5 text-sm font-medium text-red-600 bg-red-50 rounded-lg hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-red-300 text-center"
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-12 text-gray-500">
          <p>No transcription history found.</p>
          <p class="mt-2">
            Go to the
            <router-link to="/" class="text-blue-600 hover:underline"
              >transcriber</router-link
            >
            to create one.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { exportTranscription } from "../utils.js";

const transcriptions = ref([]);
const router = useRouter();
const activeExportMenu = ref(null);

const loadTranscriptions = () => {
  const allKeys = Object.keys(localStorage);
  const transcriptionKeys = allKeys.filter((key) =>
    key.startsWith("transcription-")
  );
  const loaded = transcriptionKeys
    .map((key) => {
      try {
        return JSON.parse(localStorage.getItem(key));
      } catch (e) {
        console.error(
          `Failed to parse transcription from localStorage for key: ${key}`,
          e
        );
        return null;
      }
    })
    .filter((item) => item !== null);

  transcriptions.value = loaded.sort(
    (a, b) => new Date(b.timestamp) - new Date(a.timestamp)
  );
};

const viewTranscription = (id) => {
  router.push({ name: "AudioTranscriber", query: { id: id } });
};

const deleteTranscription = (id) => {
  if (confirm("Are you sure you want to delete this transcription?")) {
    localStorage.removeItem(`transcription-${id}`);
    loadTranscriptions();
  }
};

const toggleExportMenu = (id) => {
  if (activeExportMenu.value === id) {
    activeExportMenu.value = null;
  } else {
    activeExportMenu.value = id;
  }
};

const handleExport = (item, format) => {
  exportTranscription(item.segments, item.fileName, format);
  activeExportMenu.value = null; // Close menu after selection
};

onMounted(() => {
  loadTranscriptions();
});
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

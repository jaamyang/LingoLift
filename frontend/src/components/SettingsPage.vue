<template>
  <div class="max-w-2xl mx-auto p-4 sm:p-6 lg:p-8">
    <div class="bg-white shadow-lg rounded-lg">
      <header class="border-b border-gray-200 px-6 py-4">
        <h1 class="text-2xl font-bold text-gray-800 tracking-tight">
          LLM & Proxy Settings
        </h1>
        <p class="text-sm text-gray-500 mt-1">
          Configure the Language Model and network settings for word
          explanations and translations.
        </p>
      </header>

      <form @submit.prevent="saveSettings" class="p-6 space-y-6">
        <!-- LLM Settings Section -->
        <div>
          <h2 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-4">
            Language Model (LLM)
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Base URL -->
            <div class="col-span-1">
              <label
                for="baseUrl"
                class="block text-sm font-medium text-gray-700"
                >Base URL</label
              >
              <input
                type="url"
                id="baseUrl"
                v-model.trim="settings.baseUrl"
                placeholder="https://api.openai.com/v1"
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <!-- API Key -->
            <div class="col-span-1">
              <label
                for="apiKey"
                class="block text-sm font-medium text-gray-700"
                >API Key</label
              >
              <input
                type="password"
                id="apiKey"
                v-model.trim="settings.apiKey"
                placeholder="sk-..."
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <!-- Model Name -->
            <div class="col-span-2">
              <label for="model" class="block text-sm font-medium text-gray-700"
                >Model Name</label
              >
              <input
                type="text"
                id="model"
                v-model.trim="settings.model"
                placeholder="gpt-4, qwen/qwen3-30b-a3b:free, etc."
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <!-- Custom Prompt -->
            <div class="col-span-2">
              <label
                for="customPrompt"
                class="block text-sm font-medium text-gray-700"
              >
                Custom Prompt Template
                <span class="text-gray-500 font-normal">(optional)</span>
              </label>
              <textarea
                id="customPrompt"
                v-model="settings.customPrompt"
                rows="5"
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-mono"
                placeholder="Use {word}, {context}, {language} for explanations, or {sentence}, {language} for translations."
              ></textarea>
              <p class="mt-2 text-xs text-gray-500">
                If you provide a custom prompt, the default system prompt will
                be ignored. Use the available placeholders.
              </p>
            </div>

            <!-- Additional Params -->
            <div class="col-span-2">
              <label
                for="additionalParams"
                class="block text-sm font-medium text-gray-700"
              >
                Additional Request Body Parameters (JSON format)
              </label>
              <textarea
                id="additionalParams"
                v-model="additionalParamsString"
                @blur="validateAndFormatJson"
                class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm font-mono"
                :class="{ 'border-red-500': jsonError }"
                rows="4"
              ></textarea>
              <p v-if="jsonError" class="mt-1 text-xs text-red-600">
                {{ jsonError }}
              </p>
              <p v-else class="mt-1 text-xs text-gray-500">
                Override or add parameters like `temperature`, `top_p`, etc.
              </p>
            </div>
          </div>
        </div>

        <!-- Proxy Settings Section -->
        <div>
          <h2 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-4">
            Network
          </h2>
          <div>
            <label for="proxy" class="block text-sm font-medium text-gray-700">
              Proxy URL
              <span class="text-gray-500 font-normal">(optional)</span>
            </label>
            <input
              type="url"
              id="proxy"
              v-model.trim="settings.proxy"
              placeholder="http://user:pass@host:port"
              class="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
            />
            <p class="mt-2 text-xs text-gray-500">
              If set, all backend requests to the LLM will be routed through
              this proxy.
            </p>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-between pt-4 border-t">
          <button
            type="button"
            @click="resetToDefaults"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Reset to Defaults
          </button>
          <button
            type="submit"
            :disabled="!!jsonError"
            class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
          >
            Save Settings
          </button>
        </div>
        <div
          v-if="saveStatus"
          class="mt-4 text-center text-sm font-medium"
          :class="
            saveStatus.type === 'success' ? 'text-green-600' : 'text-red-600'
          "
        >
          {{ saveStatus.message }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";

const defaultSettings = {
  baseUrl: "https://openrouter.ai/api/v1/chat/completions",
  apiKey: "",
  model: "qwen/qwen3-30b-a3b:free",
  customPrompt: "",
  proxy: "",
  // additionalParams is handled separately as an object
  additionalParams: {},
};

const settings = ref(JSON.parse(JSON.stringify(defaultSettings)));
// For the textarea binding
const additionalParamsString = ref("");
const jsonError = ref("");
const saveStatus = ref(null);

// Load settings from localStorage when the component mounts
onMounted(() => {
  const savedSettings = localStorage.getItem("llmSettings");
  if (savedSettings) {
    try {
      const parsed = JSON.parse(savedSettings);
      // Merge saved settings into the ref, ensuring all keys from default are present
      Object.assign(settings.value, parsed);
      // Handle the JSON object separately for the textarea
      if (parsed.additionalParams) {
        additionalParamsString.value = JSON.stringify(
          parsed.additionalParams,
          null,
          2
        );
      } else {
        additionalParamsString.value = "{}";
      }
    } catch (e) {
      console.error("Failed to parse settings from localStorage", e);
      resetToDefaults(); // Reset if settings are corrupt
    }
  } else {
    additionalParamsString.value = "{}";
  }
});

function validateAndFormatJson() {
  try {
    if (additionalParamsString.value.trim() === "") {
      settings.value.additionalParams = {};
      additionalParamsString.value = ""; // or '{}' if you prefer
      jsonError.value = "";
      return;
    }
    const parsed = JSON.parse(additionalParamsString.value);
    settings.value.additionalParams = parsed;
    additionalParamsString.value = JSON.stringify(parsed, null, 2); // Auto-format
    jsonError.value = "";
  } catch (e) {
    jsonError.value = "Invalid JSON format. Please correct it.";
  }
}

// Watch for changes in the text area and validate
watch(additionalParamsString, (newValue) => {
  try {
    if (newValue.trim() === "") {
      jsonError.value = "";
      return;
    }
    JSON.parse(newValue);
    jsonError.value = "";
  } catch (e) {
    jsonError.value = "Invalid JSON format.";
  }
});

function saveSettings() {
  validateAndFormatJson(); // Final validation before saving
  if (jsonError.value) {
    showSaveStatus("Cannot save, please fix the JSON error.", "error");
    return;
  }

  try {
    localStorage.setItem("llmSettings", JSON.stringify(settings.value));
    showSaveStatus("Settings saved successfully!", "success");
  } catch (e) {
    showSaveStatus("Failed to save settings.", "error");
    console.error("Error saving to localStorage:", e);
  }
}

function resetToDefaults() {
  if (
    window.confirm(
      "Are you sure you want to reset all settings to their defaults?"
    )
  ) {
    settings.value = JSON.parse(JSON.stringify(defaultSettings));
    additionalParamsString.value = "{}";
    jsonError.value = "";
    showSaveStatus("Settings have been reset to default values.", "success");
  }
}

function showSaveStatus(message, type) {
  saveStatus.value = { message, type };
  setTimeout(() => {
    saveStatus.value = null;
  }, 3000);
}
</script>

<style scoped>
/* Scoped styles can be added here if needed */
</style>

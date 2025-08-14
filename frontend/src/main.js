console.debug = console.log;

const originalFetch = window.fetch;
window.fetch = async (input, init) => {
  console.log('[🏗 fetch]', input);
  const res = await originalFetch(input, init);
  console.log(`[🏗 fetch] status:`, res.status, 'for', input);
  return res;
};

import { createApp } from 'vue';
import App from './App.vue';
import './index.css'; // Assuming your Tailwind styles are imported via index.css
import { createRouter, createWebHistory } from 'vue-router';
import SettingsPage from './components/SettingsPage.vue';
import AudioTranscriber from './components/AudioTranscriber.vue';
import TranscriptionHistory from './components/TranscriptionHistory.vue';

const routes = [
  { path: '/', name: 'AudioTranscriber', component: AudioTranscriber },
  { path: '/settings', name: 'Settings', component: SettingsPage },
  { path: '/history', name: 'TranscriptionHistory', component: TranscriptionHistory },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const app = createApp(App);
app.use(router);
app.mount('#app'); 
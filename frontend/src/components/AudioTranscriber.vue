<template>
  <div class="flex flex-col md:flex-row gap-6 max-w-full mx-auto p-4 md:p-6">
    <!-- Left Column: Audio Transcriber -->
    <div
      class="flex-shrink-0 transition-all duration-300 ease-in-out"
      :class="isExplanationPanelVisible ? 'w-full md:w-2/3' : 'w-full'"
    >
      <div class="bg-white shadow-xl rounded-lg">
        <header
          class="mb-6 border-b pb-4 p-4 md:p-6 flex justify-between items-center"
        >
          <h2 class="text-2xl font-bold text-gray-800 tracking-tight">
            Enhanced Audio Transcription
          </h2>
          <button
            @click="toggleExplanationPanel"
            class="px-3 py-1.5 text-sm font-medium text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 active:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-300"
            title="Toggle explanation panel"
          >
            <svg
              v-if="isExplanationPanelVisible"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12.75 15l3-3m0 0l-3-3m3 3h-7.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <svg
              v-else
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="w-5 h-5"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M11.25 9l-3 3m0 0l3 3m-3-3h7.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span class="ml-1"
              >{{ isExplanationPanelVisible ? "Hide" : "Show" }} Notes</span
            >
          </button>
        </header>

        <div class="p-4 md:p-6">
          <!-- Audio Upload -->
          <div class="mb-8 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <div v-if="!isHistoryView">
              <!-- Tabs for upload method -->
              <div class="border-b border-gray-200">
                <nav class="-mb-px flex space-x-8" aria-label="Tabs">
                  <button
                    @click="inputMethod = 'upload'"
                    :class="[
                      'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm',
                      inputMethod === 'upload'
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
                    ]"
                  >
                    Upload File
                  </button>
                  <button
                    @click="inputMethod = 'url'"
                    :class="[
                      'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm',
                      inputMethod === 'url'
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
                    ]"
                  >
                    From URL
                  </button>
                  <button
                    @click="inputMethod = 'rss'"
                    :class="[
                      'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm',
                      inputMethod === 'rss'
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
                    ]"
                  >
                    RSS Feed
                  </button>
                </nav>
              </div>

              <!-- Content for tabs -->
              <div class="pt-6">
                <!-- File Upload View -->
                <div v-show="inputMethod === 'upload'">
                  <label
                    for="audioFile"
                    class="block text-sm font-medium text-gray-700 mb-2"
                    >Select Audio File</label
                  >
                  <div class="mb-4">
                    <label
                      for="audioLanguage"
                      class="block text-sm font-medium text-gray-700 mb-1"
                      >Audio Language (for transcription)</label
                    >
                    <select
                      id="audioLanguage"
                      v-model="selectedAudioLanguage"
                      :disabled="isLoading"
                      class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md shadow-sm"
                    >
                      <option
                        v-for="lang in audioLanguages"
                        :key="lang.code"
                        :value="lang.code"
                      >
                        {{ lang.name }}
                      </option>
                    </select>
                  </div>
                  <input
                    id="audioFile"
                    type="file"
                    accept="audio/mp3, audio/mpeg,audio/mp4,audio/wav,audio/ogg,audio/flac,audio/x-m4a"
                    @change="handleFileUpload"
                    class="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors duration-200"
                    :disabled="isLoading"
                  />
                  <p
                    v-if="currentAudioFile && !isHistoryView"
                    class="mt-2 text-xs text-gray-500"
                  >
                    Selected: {{ currentAudioFile.name }}
                  </p>
                </div>

                <!-- URL Input View -->
                <div v-show="inputMethod === 'url'">
                  <label
                    for="audioUrl"
                    class="block text-sm font-medium text-gray-700 mb-2"
                    >Audio File URL</label
                  >
                  <div class="flex items-center gap-2">
                    <input
                      id="audioUrl"
                      type="url"
                      v-model="audioUrlInput"
                      placeholder="https://example.com/audio.mp3"
                      class="flex-grow block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md shadow-sm"
                      :disabled="isLoading"
                      @keyup.enter="handleUrlTranscription"
                    />
                    <button
                      @click="handleUrlTranscription"
                      :disabled="isLoading || !audioUrlInput"
                      class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      Transcribe
                    </button>
                  </div>
                  <p class="mt-2 text-xs text-gray-500">
                    Note: Transcription from a URL may fail due to web security
                    (CORS) restrictions on the audio server.
                  </p>
                </div>

                <!-- RSS Feed View -->
                <div v-show="inputMethod === 'rss'">
                  <label
                    for="rssUrl"
                    class="block text-sm font-medium text-gray-700 mb-2"
                    >Podcast RSS Feed URL</label
                  >
                  <div class="flex items-center gap-2 mb-4">
                    <input
                      id="rssUrl"
                      type="url"
                      v-model="rssUrlInput"
                      placeholder="https://feeds.npr.org/510289/podcast.xml"
                      class="flex-grow block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md shadow-sm"
                      :disabled="isLoading"
                      @keyup.enter="handleFetchRssFeed"
                    />
                    <button
                      @click="handleFetchRssFeed"
                      :disabled="isLoading || !rssUrlInput"
                      class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                      Fetch Feed
                    </button>
                  </div>

                  <!-- Episode List -->
                  <div v-if="podcastEpisodes.length > 0" class="mt-4">
                    <h4 class="font-semibold mb-2">{{ podcastTitle }}</h4>
                    <div
                      class="space-y-2 max-h-60 overflow-y-auto border p-2 rounded-md"
                    >
                      <div
                        v-for="episode in podcastEpisodes"
                        :key="episode.guid"
                        class="flex items-center p-2 bg-gray-100 rounded-md"
                      >
                        <input
                          type="checkbox"
                          :id="episode.guid"
                          :value="episode"
                          v-model="selectedEpisodes"
                          class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                        />
                        <label
                          :for="episode.guid"
                          class="ml-3 text-sm text-gray-800"
                        >
                          <strong class="font-semibold">{{
                            episode.title
                          }}</strong>
                          - {{ new Date(episode.pubDate).toLocaleDateString() }}
                        </label>
                      </div>
                    </div>
                    <button
                      @click="handleBatchTranscription"
                      :disabled="isLoading || selectedEpisodes.length === 0"
                      class="mt-4 w-full px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-lg hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                    >
                      Transcribe {{ selectedEpisodes.length }} Selected
                      Episode(s)
                    </button>
                  </div>
                  <p v-if="rssError" class="mt-2 text-sm text-red-600">
                    {{ rssError }}
                  </p>
                </div>
              </div>
            </div>
            <div v-else>
              <h3 class="text-lg font-semibold text-gray-800 mb-2">
                Viewing Saved Transcription
              </h3>
              <p class="text-sm text-gray-600 mb-4">
                <span class="font-medium">Original file:</span>
                {{ currentAudioFile.name }}
              </p>
              <label
                for="playbackFile"
                class="inline-block px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 cursor-pointer transition-colors"
                >Load Audio for Playback</label
              >
              <input
                id="playbackFile"
                type="file"
                accept="audio/*"
                @change="handlePlaybackFileLoad"
                class="hidden"
                :disabled="isLoading"
              />
              <p class="mt-2 text-xs text-gray-500">
                Load the original audio file to enable the player.
              </p>
            </div>
          </div>

          <!-- Audio Player & Controls -->
          <div
            v-if="audioUrl"
            class="mb-6 p-4 bg-gray-50 rounded-lg border border-gray-200"
          >
            <audio
              ref="audioPlayer"
              :src="audioUrl"
              class="w-full rounded-lg shadow-sm mb-4"
              controls
              @loadedmetadata="handleAudioMetadataLoaded"
              @timeupdate="handleTimeUpdate"
              @play="isPlaying = true"
              @pause="isPlaying = false"
              @error="handleAudioError"
            ></audio>
            <div
              v-if="audioError"
              class="text-red-600 text-sm mt-2 mb-2 bg-red-50 p-2 rounded-md"
            >
              {{ audioError }}
            </div>

            <div
              class="flex flex-col sm:flex-row items-center justify-between gap-4"
            >
              <!-- Playback Speed Control -->
              <div class="flex items-center gap-2">
                <label
                  for="playbackSpeed"
                  class="text-sm font-medium text-gray-700"
                  >Speed:</label
                >
                <select
                  id="playbackSpeed"
                  v-model="currentPlaybackRate"
                  @change="setPlaybackSpeed"
                  class="bg-white border border-gray-300 text-gray-700 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-1.5 transition-colors duration-200"
                >
                  <option value="0.5">0.5x</option>
                  <option value="0.75">0.75x</option>
                  <option value="1">1x (Normal)</option>
                  <option value="1.25">1.25x</option>
                  <option value="1.5">1.5x</option>
                  <option value="1.75">1.75x</option>
                  <option value="2">2x</option>
                </select>
              </div>

              <!-- Segment Navigation -->
              <div class="flex items-center gap-2">
                <button
                  @click="navigateToSegment('previous')"
                  :disabled="!canNavigateSegment('previous') || isLoading"
                  class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  &larr; Previous
                </button>
                <button
                  @click="navigateToSegment('next')"
                  :disabled="!canNavigateSegment('next') || isLoading"
                  class="px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  Next &rarr;
                </button>
              </div>
            </div>
          </div>

          <!-- Transcription Display -->
          <div
            v-if="
              transcript &&
              transcript.segments &&
              transcript.segments.length > 0
            "
            class="mt-6"
          >
            <div class="flex justify-between items-center mb-3">
              <h3 class="text-xl font-semibold text-gray-700">Transcription</h3>
              <!-- Export Dropdown -->
              <div class="relative">
                <button
                  @click="showExportOptions = !showExportOptions"
                  class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  Export
                  <svg
                    class="w-4 h-4 inline-block ml-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 9l-7 7-7-7"
                    ></path>
                  </svg>
                </button>
                <div
                  v-if="showExportOptions"
                  class="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded-md shadow-lg z-10"
                >
                  <a
                    @click.prevent="
                      exportTranscription(
                        transcript.segments,
                        currentAudioFile.name,
                        'txt'
                      )
                    "
                    href="#"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >Export as Text (.txt)</a
                  >
                  <a
                    @click.prevent="
                      exportTranscription(
                        transcript.segments,
                        currentAudioFile.name,
                        'srt'
                      )
                    "
                    href="#"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >Export as Timed Text (.txt)</a
                  >
                </div>
              </div>
            </div>
            <div class="relative">
              <div
                ref="transcriptContainerRef"
                class="space-y-2 max-h-[300px] md:max-h-[400px] overflow-y-auto p-4 bg-gray-50 rounded-lg border border-gray-200 shadow-inner scroll-smooth"
                @scroll="handleTranscriptScroll"
              >
                <div
                  v-for="(segment, index) in transcript.segments"
                  :key="segment.start + '-' + index + '-' + segment.text.length"
                  :ref="(el) => assignSegmentRef(el, segment)"
                  :class="[
                    'p-3 rounded-lg transition-all duration-200 ease-in-out transform hover:scale-[1.01] hover:shadow-md',
                    isCurrentSegment(segment)
                      ? 'bg-blue-100 scale-[1.02] shadow-lg border-l-4 border-blue-600'
                      : 'bg-white border border-gray-200 hover:bg-gray-100',
                  ]"
                  @mouseenter="hoveredSegmentId = segment.id"
                  @mouseleave="hoveredSegmentId = null"
                >
                  <div class="flex items-center justify-between">
                    <div class="flex-grow mr-4">
                      <!-- NEW: Flex container for sentence and translate button -->
                      <div class="flex items-center justify-between gap-4">
                        <p
                          class="flex-grow text-gray-800 text-base leading-relaxed"
                          :class="{
                            'font-semibold text-blue-700':
                              isCurrentSegment(segment),
                          }"
                          @dblclick="handleWordDoubleClick"
                        >
                          <!-- Word rendering logic -->
                          <template
                            v-if="segment.words && segment.words.length > 0"
                          >
                            <span
                              v-for="(word, wordIndex) in segment.words"
                              :key="word.word + '-' + wordIndex"
                              :data-start="word.start"
                              :data-end="word.end"
                              class="cursor-pointer hover:bg-yellow-200 rounded px-0.5 py-px transition-colors duration-150"
                            >
                              {{
                                word.word +
                                (wordIndex < segment.words.length - 1
                                  ? " "
                                  : "")
                              }}
                            </span>
                          </template>
                          <template v-else>
                            {{ segment.text }}
                          </template>
                        </p>

                        <!-- MODIFIED: Translate button is now smaller and here -->
                        <button
                          v-if="
                            hoveredSegmentId === segment.id ||
                            translatingSegmentId === segment.id ||
                            segment.translationError
                          "
                          @click="translateSegment(segment, index)"
                          :disabled="
                            translatingSegmentId === segment.id ||
                            (segment.translation && !segment.translationError)
                          "
                          class="flex-shrink-0 px-2 py-1 text-xs font-medium text-indigo-600 bg-indigo-50 rounded-md hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-indigo-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-150 ease-in-out"
                        >
                          <span v-if="translatingSegmentId === segment.id"
                            >...</span
                          >
                          <span
                            v-else-if="
                              segment.translation && !segment.translationError
                            "
                            >✓</span
                          >
                          <span v-else-if="segment.translationError"
                            >Retry</span
                          >
                          <span v-else>Translate</span>
                        </button>
                      </div>

                      <!-- MODIFIED: Translation result section -->
                      <div class="mt-2 min-h-[20px]">
                        <div
                          v-if="
                            segment.translationLoading &&
                            hoveredSegmentId === segment.id
                          "
                          class="mt-1 text-xs text-gray-500"
                        >
                          Loading translation...
                        </div>
                        <p
                          v-if="
                            segment.translation &&
                            !segment.translationError &&
                            (hoveredSegmentId === segment.id ||
                              translatingSegmentId !== segment.id)
                          "
                          class="mt-1 text-sm text-gray-600 bg-gray-100 p-2 rounded-md shadow-sm"
                        >
                          {{ segment.translation }}
                        </p>
                        <p
                          v-if="
                            segment.translationError &&
                            (hoveredSegmentId === segment.id ||
                              translatingSegmentId !== segment.id)
                          "
                          class="mt-1 text-sm text-red-600 bg-red-50 p-2 rounded-md"
                        >
                          {{ segment.translationError }}
                        </p>
                      </div>
                    </div>
                    <button
                      @click="seekToSegment(segment)"
                      class="p-1.5 text-blue-600 hover:text-blue-800 focus:outline-none rounded-full hover:bg-blue-100 active:bg-blue-200 transition-colors duration-150 flex-shrink-0"
                      title="Play this segment"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                        class="w-6 h-6"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm14.024-.983a1.125 1.125 0 0 1 0 1.966l-5.625 3.125A1.125 1.125 0 0 1 9 15.125V8.875c0-.87.986-1.406 1.65-.983l5.625 3.125Z"
                          clip-rule="evenodd"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
              <button
                v-if="showReturnToPlaybackButton"
                @click="scrollToCurrentSegment"
                title="Scroll to current playing segment"
                class="absolute bottom-4 right-4 z-10 bg-blue-600 text-white p-3 rounded-full shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 transition-opacity duration-300"
                :class="{
                  'opacity-100': showReturnToPlaybackButton,
                  'opacity-0 pointer-events-none': !showReturnToPlaybackButton,
                }"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                  class="w-6 h-6"
                >
                  <path
                    d="M12 15a.75.75 0 0 1-.75-.75V4.56L8.03 8.03a.75.75 0 0 1-1.06-1.06l4.5-4.5a.75.75 0 0 1 1.06 0l4.5 4.5a.75.75 0 0 1-1.06 1.06L12.75 4.56v9.69a.75.75 0 0 1-.75.75Z"
                  />
                  <path
                    d="M1.5 12.75a.75.75 0 0 1 .75-.75h19.5a.75.75 0 0 1 0 1.5H2.25a.75.75 0 0 1-.75-.75Z"
                  />
                </svg>
              </button>
            </div>
          </div>

          <!-- Empty/Info States for Transcription -->
          <div
            v-if="
              !isLoading &&
              audioUrl &&
              (!transcript ||
                !transcript.segments ||
                transcript.segments.length === 0) &&
              !transcriptionError
            "
            class="mt-6 p-4 text-center text-gray-500 bg-gray-50 rounded-lg border border-gray-200"
          >
            <p>
              {{
                transcript &&
                transcript.text &&
                transcript.segments.length === 0
                  ? "Transcription available but no time segments found."
                  : "Audio loaded. Waiting for transcription or transcription is empty."
              }}
            </p>
          </div>

          <!-- Loading State -->
          <div
            v-if="isLoading"
            class="mt-6 p-6 text-center text-gray-700 bg-gray-50 rounded-lg border border-gray-200 flex flex-col items-center justify-center"
          >
            <svg
              class="animate-spin h-8 w-8 text-blue-600 mb-3"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            <p class="text-lg font-medium">{{ loadingMessage }}</p>
          </div>

          <!-- Error State for Transcription -->
          <div
            v-if="transcriptionError"
            class="mt-6 p-4 text-red-700 bg-red-100 border border-red-300 rounded-lg shadow-md"
          >
            <h4 class="font-bold mb-1">Transcription Error:</h4>
            <p class="text-sm">{{ transcriptionError }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Column: Explanation Panel (conditionally rendered) -->
    <div
      v-if="isExplanationPanelVisible"
      class="w-full md:w-1/3 md:sticky md:top-6 h-fit transition-all duration-300 ease-in-out"
    >
      <div class="bg-white shadow-xl rounded-lg p-6">
        <div class="flex flex-col gap-3 mb-4 border-b pb-3">
          <h3 class="text-xl font-semibold text-gray-800">Word Explanation</h3>
          <div class="flex items-center justify-between gap-2">
            <label
              for="userNativeLangSelect"
              class="text-sm text-gray-600 whitespace-nowrap"
              >Your Language:</label
            >
            <select
              id="userNativeLangSelect"
              v-model="userSelectedNativeLanguage"
              @change="handleUserNativeLanguageChange"
              title="Select your language for a secondary explanation"
              class="flex-grow text-sm rounded border-gray-300 focus:ring-blue-500 focus:border-blue-500 p-1.5 bg-gray-50 hover:bg-gray-100 text-gray-700"
            >
              <option value="English">English</option>
              <option value="Spanish">Español</option>
              <option value="French">Français</option>
              <option value="German">Deutsch</option>
              <option value="Chinese">中文</option>
              <option value="Japanese">日本語</option>
              <option value="Korean">한국어</option>
              {/* Add more as needed */}
            </select>
          </div>
        </div>

        <div v-if="selectedWord" class="min-h-[200px] space-y-3">
          <p class="mb-1">
            <span class="font-semibold text-gray-700">Word:</span>
            <span class="italic text-blue-600 text-lg ml-1">{{
              selectedWord
            }}</span>
          </p>

          <div v-if="isExplanationLoading" class="text-center py-4">
            <p class="text-gray-600">Loading explanations...</p>
            <svg
              class="animate-spin h-6 w-6 text-blue-500 mx-auto mt-2"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
          </div>
          <div v-else-if="explanationError" class="text-center py-4">
            <p class="text-red-500">{{ explanationError }}</p>
          </div>

          <div v-else class="space-y-4">
            <!-- Explanation in Audio Language -->
            <div class="p-3 bg-gray-50 rounded-md border border-gray-200">
              <h4 class="font-semibold text-sm text-gray-600 mb-1.5">
                Explanation (Audio: {{ audioSourceLanguage || "Unknown" }}):
              </h4>
              <p
                v-if="audioSourceExplanation"
                class="text-gray-800 leading-relaxed"
              >
                {{ audioSourceExplanation }}
              </p>
              <p v-else class="text-gray-500 text-sm italic">
                No explanation available for audio language.
              </p>
            </div>

            <!-- Toggle for User Native Language Explanation -->
            <button
              v-if="audioSourceExplanation || userNativeExplanation"
              @click="
                showUserNativeExplanationPanel = !showUserNativeExplanationPanel
              "
              class="w-full mt-2 px-3 py-2 text-sm font-medium text-blue-700 bg-blue-100 rounded-lg hover:bg-blue-200 active:bg-blue-300 focus:outline-none focus:ring-2 focus:ring-blue-400 flex items-center justify-center gap-1.5"
            >
              <svg
                v-if="showUserNativeExplanationPanel"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-4 h-4"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M19.5 12h-15m0 0l6.75-6.75M4.5 12l6.75 6.75"
                />
              </svg>
              <svg
                v-else
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-4 h-4"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M4.5 12h15m0 0l-6.75-6.75M19.5 12l-6.75 6.75"
                />
              </svg>
              <span
                >{{
                  showUserNativeExplanationPanel ? "Hide" : "Show"
                }}
                Explanation in {{ userSelectedNativeLanguage }}</span
              >
            </button>

            <!-- Explanation in User Native Language (conditionally rendered) -->
            <div
              v-if="showUserNativeExplanationPanel"
              class="p-3 bg-indigo-50 rounded-md border border-indigo-200 mt-2"
            >
              <h4 class="font-semibold text-sm text-indigo-700 mb-1.5">
                Explanation ({{ userSelectedNativeLanguage }}):
              </h4>
              <p
                v-if="userNativeExplanation"
                class="text-gray-800 leading-relaxed"
              >
                {{ userNativeExplanation }}
              </p>
              <p v-else class="text-gray-500 text-sm italic">
                No explanation available in {{ userSelectedNativeLanguage }}.
              </p>
            </div>
          </div>
        </div>
        <div
          v-else
          class="text-center py-4 text-gray-500 min-h-[200px] flex items-center justify-center"
        >
          <p>
            Double-click a word in the transcription to see its explanation
            here.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onUnmounted, nextTick, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import axios from "axios";
import { exportTranscription } from "../utils.js";

const API_URL = "";

const router = useRoute();

// ---- Worker Setup ----
const worker = ref(null);

// ---- App State & Core Refs ----
const inputMethod = ref("upload"); // 'upload', 'url', or 'rss'
const audioUrlInput = ref("");
const rssUrlInput = ref("");
const podcastTitle = ref("");
const podcastEpisodes = ref([]);
const selectedEpisodes = ref([]);
const rssError = ref(null);

const audioPlayer = ref(null);
const audioUrl = ref(null);
const currentAudioFile = ref(null);

const transcript = ref(null);
const currentTime = ref(0);
const isPlaying = ref(false);
const isLoading = ref(false);
const loadingMessage = ref("");
const transcriptionError = ref(null);
const audioError = ref(null);
const isHistoryView = ref(false); // New state for history view

const transcriptContainerRef = ref(null);
const segmentRefs = ref(new Map());
const showReturnToPlaybackButton = ref(false);

// New state for playback speed
const currentPlaybackRate = ref(1);

// State for export dropdown
const showExportOptions = ref(false);

// State for Word Explanation Feature
const isExplanationPanelVisible = ref(false);
const selectedWord = ref("");
const selectedContext = ref("");

const audioSourceLanguage = ref("English");
const userSelectedNativeLanguage = ref("English");

const audioSourceExplanation = ref("");
const userNativeExplanation = ref("");
const showUserNativeExplanationPanel = ref(false);

const isExplanationLoading = ref(false);
const explanationError = ref("");

// Audio Language Selection
const selectedAudioLanguage = ref("en"); // Default to English
const audioLanguages = ref([
  { name: "English", code: "en" },
  { name: "Spanish", code: "es" },
  { name: "French", code: "fr" },
  { name: "German", code: "de" },
  { name: "Japanese", code: "ja" },
  { name: "Chinese", code: "zh" },
  { name: "Auto-Detect", code: "" }, // Empty string for auto-detection
]);

// Translation State
const translatingSegmentId = ref(null); // To track the ID of the segment being translated
const hoveredSegmentId = ref(null); // To track segment hover for translate button visibility

// ---- Utilities ----
async function fetchWithProxyFallback(url) {
  try {
    // Attempt 1: Direct fetch
    console.log(
      `[Fetcher] Attempting direct fetch for: ${url.substring(0, 100)}...`
    );
    return await fetch(url);
  } catch (err) {
    if (err instanceof TypeError) {
      // Likely a CORS error or network failure
      console.warn(
        `[Fetcher] Direct fetch for ${url.substring(
          0,
          100
        )}... failed, retrying via proxy.`
      );

      // Read settings from localStorage to check for a proxy
      const llmSettings = JSON.parse(
        localStorage.getItem("llmSettings") || "{}"
      );
      const proxySetting = llmSettings.proxy;

      let proxyUrl = `${API_URL}/api/proxy?url=${encodeURIComponent(url)}`;
      if (proxySetting) {
        proxyUrl += `&proxy=${encodeURIComponent(proxySetting)}`;
      }

      try {
        const proxyResponse = await fetch(proxyUrl);
        if (!proxyResponse.ok) {
          let errorJson = {};
          try {
            errorJson = await proxyResponse.json();
          } catch (e) {
            // Ignore if response is not json
          }
          throw new Error(
            errorJson.error ||
              `Proxy request failed with status ${proxyResponse.status}`
          );
        }
        return proxyResponse;
      } catch (proxyErr) {
        throw new Error(`Failed to fetch via proxy: ${proxyErr.message}`);
      }
    } else {
      // It wasn't a CORS-like error, just a different network issue. Re-throw it.
      throw err;
    }
  }
}

// ---- Audio Handling ----

const handleTranscriptionCompletion = (result) => {
  transcript.value = {
    text: result.text,
    segments: result.chunks.map((chunk, index) => ({
      id: index,
      start: chunk.timestamp[0],
      end: chunk.timestamp[1],
      text: chunk.text.trim(),
    })),
    language: "English", // Whisper model used is English-only for now
  };
  audioSourceLanguage.value = "English";

  const historyItem = {
    id: Date.now(),
    timestamp: new Date().toISOString(),
    transcript: transcript.value.text,
    segments: transcript.value.segments,
    language: transcript.value.language,
    fileName: currentAudioFile.value.name,
    audioUrl: currentAudioFile.value.sourceUrl || null,
  };
  localStorage.setItem(
    `transcription-${historyItem.id}`,
    JSON.stringify(historyItem)
  );
};

const transcribeFile = async (file, sourceUrl = null, historyTitle = null) => {
  if (!worker.value) {
    transcriptionError.value =
      "Transcription worker is not ready. Please refresh the page.";
    return;
  }
  // 1. Reset state
  currentAudioFile.value = { name: historyTitle || file.name, sourceUrl };
  const objectUrl = URL.createObjectURL(file);
  audioUrl.value = objectUrl;

  transcript.value = null;
  transcriptionError.value = null;
  audioError.value = null;
  isLoading.value = true;
  loadingMessage.value = "Preparing audio for transcription...";
  isHistoryView.value = false;
  currentTime.value = 0;
  if (audioPlayer.value) audioPlayer.value.playbackRate = 1;

  try {
    // 2. Read audio data as a buffer
    const reader = new FileReader();
    reader.onload = async (e) => {
      try {
        const audioContext = new (window.AudioContext ||
          window.webkitAudioContext)();
        const decodedData = await audioContext.decodeAudioData(e.target.result);

        // 3. Resample to 16kHz, which is required by Whisper
        const targetSampleRate = 16000;
        const sourceSampleRate = decodedData.sampleRate;
        const pcmData = decodedData.getChannelData(0);

        let finalPcmData = pcmData;
        if (sourceSampleRate !== targetSampleRate) {
          // Resample the audio if necessary
          const ratio = targetSampleRate / sourceSampleRate;
          const newLength = Math.round(pcmData.length * ratio);
          const resampledData = new Float32Array(newLength);
          let offsetResult = 0;
          let offsetBuffer = 0;
          while (offsetResult < resampledData.length) {
            const nextOffsetBuffer = Math.round((offsetResult + 1) / ratio);
            let accum = 0,
              count = 0;
            for (
              let i = offsetBuffer;
              i < nextOffsetBuffer && i < pcmData.length;
              i++
            ) {
              accum += pcmData[i];
              count++;
            }
            resampledData[offsetResult] = accum / count;
            offsetResult++;
            offsetBuffer = nextOffsetBuffer;
          }
          finalPcmData = resampledData;
        }

        // 4. Send to worker
        loadingMessage.value = "Model loading (this may take a moment)...";
        worker.value.postMessage({ type: "transcribe", audio: finalPcmData });
      } catch (decodeErr) {
        console.error("Error decoding audio data:", decodeErr);
        transcriptionError.value = `Failed to decode audio data: ${decodeErr.message}`;
        isLoading.value = false;
      }
    };
    reader.onerror = (e) => {
      console.error("FileReader error:", e);
      transcriptionError.value = "Failed to read the audio file.";
      isLoading.value = false;
    };
    reader.readAsArrayBuffer(file);
  } catch (err) {
    console.error("Error processing audio file:", err);
    transcriptionError.value = `Failed to process audio file: ${err.message}`;
    isLoading.value = false;
  }
};

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) {
    transcriptionError.value = "No file selected.";
    return;
  }
  await transcribeFile(file);
};

const handleUrlTranscription = async () => {
  if (!audioUrlInput.value) {
    transcriptionError.value = "Please enter an audio URL.";
    return;
  }

  isLoading.value = true;
  loadingMessage.value = "Fetching audio from URL...";
  transcriptionError.value = null;
  audioError.value = null;

  try {
    const response = await fetchWithProxyFallback(audioUrlInput.value);
    if (!response.ok) {
      throw new Error(
        `Failed to fetch audio. Server responded with ${response.status}: ${response.statusText}`
      );
    }
    const blob = await response.blob();

    const urlParts = new URL(audioUrlInput.value).pathname.split("/");
    const fileName = urlParts.pop() || "audio_from_url.mp3";

    const audioFile = new File([blob], fileName, {
      type: blob.type || "audio/*",
    });

    await transcribeFile(audioFile, audioUrlInput.value);
  } catch (err) {
    console.error(
      "[AudioTranscriber] Failed to fetch or process audio from URL:",
      err
    );
    let errorMessage =
      "Failed to fetch audio from the provided URL. This could be due to a network error or CORS restrictions on the remote server.";
    if (err.name === "TypeError") {
      errorMessage =
        "Could not fetch audio from the URL. This is often caused by CORS policy. The audio server must allow requests from this website.";
    } else {
      errorMessage = `An error occurred: ${err.message}`;
    }
    transcriptionError.value = errorMessage;
    isLoading.value = false;
    loadingMessage.value = "";
  }
};

const handleFetchRssFeed = async () => {
  if (!rssUrlInput.value) {
    rssError.value = "Please enter an RSS feed URL.";
    return;
  }
  isLoading.value = true;
  loadingMessage.value = "Fetching and parsing RSS feed...";
  rssError.value = null;
  podcastEpisodes.value = [];
  selectedEpisodes.value = [];
  podcastTitle.value = "";

  try {
    const response = await fetchWithProxyFallback(rssUrlInput.value);
    if (!response.ok) {
      throw new Error(`Failed to fetch RSS feed. Status: ${response.status}`);
    }
    const xmlString = await response.text();

    const parser = new DOMParser();
    const doc = parser.parseFromString(xmlString, "application/xml");

    const errorNode = doc.querySelector("parsererror");
    if (errorNode) {
      throw new Error("Failed to parse the RSS feed. Please check the format.");
    }

    podcastTitle.value =
      doc.querySelector("channel > title")?.textContent || "Podcast";
    const items = doc.querySelectorAll("item");
    const episodes = [];
    items.forEach((item) => {
      const enclosure = item.querySelector('enclosure[url*="mp3"]');
      if (enclosure) {
        episodes.push({
          title: item.querySelector("title")?.textContent || "No Title",
          pubDate:
            item.querySelector("pubDate")?.textContent ||
            new Date().toISOString(),
          audioUrl: enclosure.getAttribute("url"),
          guid:
            item.querySelector("guid")?.textContent ||
            enclosure.getAttribute("url"),
        });
      }
    });
    podcastEpisodes.value = episodes;

    if (episodes.length === 0) {
      rssError.value = "No downloadable MP3 episodes found in this RSS feed.";
    }
  } catch (err) {
    console.error("[AudioTranscriber] Failed to fetch or parse RSS feed:", err);
    let errorMessage =
      "Failed to process RSS feed. This could be a network error or CORS issue.";
    if (err.name === "TypeError") {
      errorMessage =
        "Could not fetch the RSS feed due to CORS policy. The server must allow requests from this website.";
    } else {
      errorMessage = `An error occurred: ${err.message}`;
    }
    rssError.value = errorMessage;
  } finally {
    isLoading.value = false;
    loadingMessage.value = "";
  }
};

const handleBatchTranscription = async () => {
  if (selectedEpisodes.value.length === 0) {
    transcriptionError.value =
      "Please select at least one episode to transcribe.";
    return;
  }

  const episodesToTranscribe = [...selectedEpisodes.value];
  selectedEpisodes.value = []; // Clear selection

  for (let i = 0; i < episodesToTranscribe.length; i++) {
    const episode = episodesToTranscribe[i];
    loadingMessage.value = `Transcribing episode ${i + 1} of ${
      episodesToTranscribe.length
    }: ${episode.title}`;

    // Create a custom title for the history entry
    const historyTitle = `${podcastTitle.value} - ${episode.title} - ${new Date(
      episode.pubDate
    ).toLocaleDateString()}`;

    // This reuses the URL transcription logic.
    // We need to fetch the audio blob and create a file object from it.
    try {
      const response = await fetchWithProxyFallback(episode.audioUrl);
      if (!response.ok)
        throw new Error(`HTTP error! status: ${response.status}`);
      const blob = await response.blob();
      const audioFile = new File([blob], `${episode.guid}.mp3`, {
        type: "audio/mpeg",
      });

      // We pass the history title to be saved.
      await transcribeFile(audioFile, episode.audioUrl, historyTitle);
    } catch (err) {
      console.error(`Failed to transcribe ${episode.title}:`, err);
      // We could collect these errors and display them at the end.
      // For now, we'll just log it and continue.
      // Maybe set a global error message that the user can see.
      transcriptionError.value = `Failed to process episode "${episode.title}". Please check the console for details.`;
    }
  }

  loadingMessage.value = "Batch transcription complete!";
  // Maybe clear the message after a few seconds
  setTimeout(() => {
    if (loadingMessage.value === "Batch transcription complete!") {
      loadingMessage.value = "";
    }
  }, 5000);
};

const loadTranscriptionFromHistory = (id) => {
  const item = localStorage.getItem(`transcription-${id}`);
  if (item) {
    try {
      isHistoryView.value = true;
      const data = JSON.parse(item);
      transcript.value = {
        text: data.transcript,
        segments: data.segments,
        language: data.language,
      };
      audioSourceLanguage.value = data.language;
      currentAudioFile.value = { name: data.fileName }; // Mock file object for display

      if (data.audioUrl) {
        // If the audio source was a URL, we can stream it directly
        audioUrl.value = data.audioUrl;
        audioError.value = null; // Clear any previous errors
        console.log(
          `[AudioTranscriber] Loaded history for ID ${id}, streaming from URL: ${data.audioUrl}`
        );
      } else {
        // Otherwise, prompt user to upload the original file
        audioUrl.value = null;
        audioError.value =
          "This is a saved transcription. Please re-upload the original audio file to enable playback.";
      }

      console.log(
        `[AudioTranscriber] Loaded transcription ID ${id} from history.`
      );
    } catch (e) {
      console.error(
        `[AudioTranscriber] Failed to parse transcription from localStorage for ID: ${id}`,
        e
      );
      transcriptionError.value =
        "Failed to load the saved transcription. The data may be corrupt.";
    }
  } else {
    transcriptionError.value = `No transcription found in history for ID: ${id}`;
  }
};

const handlePlaybackFileLoad = (event) => {
  const file = event.target.files[0];
  if (!file) {
    return;
  }
  // We don't overwrite currentAudioFile.name, as that holds the original name from history
  audioUrl.value = URL.createObjectURL(file);
  audioError.value = null; // Clear the "please load a file" message
  console.log("[AudioTranscriber] Playback audio loaded for history item.");
};

onMounted(() => {
  // Initialize the worker and set up message handling
  worker.value = new Worker(new URL("../worker.js", import.meta.url), {
    type: "module",
  });

  worker.value.onmessage = (event) => {
    const { status, data } = event.data;
    switch (status) {
      case "progress":
        if (data.status === "downloading") {
          loadingMessage.value = `Downloading model: ${data.file} (${(
            data.progress || 0
          ).toFixed(2)}%)`;
        } else if (data.file) {
          // This covers other statuses like 'done' which still have a file associated.
          loadingMessage.value = `Loading model file: ${data.file}...`;
        } else if (data.status === "ready") {
          // This message is shown when the pipeline is fully initialized.
          loadingMessage.value = "Model initialized. Starting transcription...";
        }
        // Other progress messages without a file are ignored.
        console.log("loading file:", data.file)
        break;
      case "complete":
        isLoading.value = false;
        loadingMessage.value = "";
        handleTranscriptionCompletion(data);
        break;
      case "error":
        isLoading.value = false;
        loadingMessage.value = "";
        transcriptionError.value = `Transcription failed: ${data}`;
        break;
    }
  };

  // Add other on-mount logic
  const historyId = router.query.id;
  if (historyId) {
    loadTranscriptionFromHistory(historyId);
  }
  window.addEventListener("keydown", handleKeyDown);
});

onUnmounted(() => {
  // Terminate the worker to avoid memory leaks
  if (worker.value) {
    worker.value.terminate();
  }
  // Remove event listeners
  window.removeEventListener("keydown", handleKeyDown);
  if (audioUrl.value) {
    URL.revokeObjectURL(audioUrl.value);
    console.log(
      "[AudioTranscriber] Component unmounted, revoked audio URL:",
      audioUrl.value
    );
  }
});

const handleAudioMetadataLoaded = () => {
  console.log(
    "[AudioTranscriber] Audio metadata loaded. Duration:",
    audioPlayer.value?.duration
  );
  audioError.value = null;
  if (audioPlayer.value) {
    audioPlayer.value.playbackRate = currentPlaybackRate.value; // Apply current speed
  }
};

const handleAudioError = (event) => {
  console.error("[AudioTranscriber] Audio player error event:", event);
  let message = "Unknown audio error.";
  if (audioPlayer.value && audioPlayer.value.error) {
    const err = audioPlayer.value.error;
    message = `Error ${err.code}: ${err.message || "Failed to load audio."}`;
    if (err.code === MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED) {
      message = "Audio format not supported. Please try a different file type.";
    }
  }
  audioError.value = message;
};

const handleTimeUpdate = () => {
  if (audioPlayer.value) {
    currentTime.value = audioPlayer.value.currentTime;
  }
};

// ---- Playback Controls ----
const togglePlayback = () => {
  if (!audioPlayer.value) return;
  if (isPlaying.value) {
    audioPlayer.value.pause();
  } else {
    audioPlayer.value.play().catch((err) => {
      console.error("[AudioTranscriber] Error toggling playback:", err);
      audioError.value =
        "Could not start playback. Please ensure audio is loaded.";
    });
  }
};

const setPlaybackSpeed = () => {
  if (audioPlayer.value) {
    audioPlayer.value.playbackRate = parseFloat(currentPlaybackRate.value);
    console.log(
      "[AudioTranscriber] Playback speed set to:",
      currentPlaybackRate.value
    );
  }
};

const getCurrentSegmentIndex = () => {
  if (!transcript.value || !transcript.value.segments) return -1;
  return transcript.value.segments.findIndex((s) => isCurrentSegment(s));
};

const navigateToSegment = (direction) => {
  if (!transcript.value || !transcript.value.segments || !audioPlayer.value)
    return;

  const segments = transcript.value.segments;
  let currentIdx = getCurrentSegmentIndex();

  if (currentIdx === -1) {
    // If no segment is active, try to find based on current time
    currentIdx = segments.findIndex((s) => currentTime.value < s.end);
    if (currentIdx === -1 && segments.length > 0)
      currentIdx = segments.length - 1; // default to last if time is past all
  }

  let targetIdx = -1;
  if (direction === "previous") {
    targetIdx = currentIdx > 0 ? currentIdx - 1 : 0;
  } else if (direction === "next") {
    targetIdx =
      currentIdx < segments.length - 1 ? currentIdx + 1 : segments.length - 1;
  }

  if (targetIdx !== -1 && segments[targetIdx]) {
    seekToSegment(segments[targetIdx]);
  }
};

const canNavigateSegment = (direction) => {
  if (
    !transcript.value ||
    !transcript.value.segments ||
    transcript.value.segments.length === 0
  )
    return false;
  const segments = transcript.value.segments;
  let currentIdx = getCurrentSegmentIndex();

  // If no segment is currently active, base it on the first/last segment for button enablement logic
  if (currentIdx === -1) {
    if (direction === "previous") return false; // Can't go previous if no segment is active (or we are before the first)
    if (direction === "next") return segments.length > 1; // Can go next if there's more than one segment
  }

  if (direction === "previous") {
    return currentIdx > 0;
  } else if (direction === "next") {
    return currentIdx < segments.length - 1;
  }
  return false;
};

// ---- Transcription Handling ----

// ---- UI & Playback Sync ----

const formatTime = (seconds) => {
  if (isNaN(seconds) || seconds < 0) return "0:00";
  const minutes = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${minutes}:${secs.toString().padStart(2, "0")}`;
};

const isCurrentSegment = (segment) => {
  const epsilon = 0.1; // Increased epsilon for more robust active state during fast updates or short segments
  return (
    currentTime.value >= segment.start - epsilon &&
    currentTime.value <= segment.end + epsilon
  );
};

const seekToSegment = (segment) => {
  if (audioPlayer.value && segment && typeof segment.start === "number") {
    audioPlayer.value.currentTime = segment.start;
    if (!isPlaying.value) {
      audioPlayer.value.play().catch((err) => {
        console.error("[AudioTranscriber] Error playing on seek:", err);
        audioError.value =
          "Could not start playback. Please ensure audio is loaded and browser allows autoplay.";
      });
    }
  }
};

const assignSegmentRef = (el, segment) => {
  if (el && segment && typeof segment.start === "number") {
    segmentRefs.value.set(segment.start, el);
  }
};

const userHasScrolledAway = ref(false);
const isReturningToPlayback = ref(false); // Flag to ignore scroll events during programmatic scroll

const handleTranscriptScroll = () => {
  if (isReturningToPlayback.value || !transcriptContainerRef.value) return;
  const currentActiveSegment = transcript.value?.segments?.find((s) =>
    isCurrentSegment(s)
  );
  if (currentActiveSegment) {
    const segmentElement = segmentRefs.value.get(currentActiveSegment.start);
    if (segmentElement) {
      if (!isElementInViewport(segmentElement)) {
        userHasScrolledAway.value = true;
      }
      checkAndShowReturnToPlaybackButton();
    }
  }
};

const scrollToCurrentSegment = async () => {
  const currentActiveSegment = transcript.value?.segments?.find((s) =>
    isCurrentSegment(s)
  );
  if (currentActiveSegment) {
    const segmentElement = segmentRefs.value.get(currentActiveSegment.start);
    if (segmentElement) {
      // Re-enable auto-scrolling and ignore the upcoming scroll event
      userHasScrolledAway.value = false;
      isReturningToPlayback.value = true;

      await nextTick();
      segmentElement.scrollIntoView({ behavior: "smooth", block: "center" });

      showReturnToPlaybackButton.value = false;

      // Reset the flag after the scroll animation is likely finished
      setTimeout(() => {
        isReturningToPlayback.value = false;
      }, 1000); // 1 second should be enough for the smooth scroll
    }
  }
};

const checkAndShowReturnToPlaybackButton = () => {
  const currentActiveSegment = transcript.value?.segments?.find((s) =>
    isCurrentSegment(s)
  );
  if (currentActiveSegment) {
    const segmentElement = segmentRefs.value.get(currentActiveSegment.start);
    if (segmentElement && transcriptContainerRef.value) {
      showReturnToPlaybackButton.value = !isElementInViewport(segmentElement);
    } else {
      showReturnToPlaybackButton.value = false;
    }
  } else {
    showReturnToPlaybackButton.value = false;
  }
};

const isElementInViewport = (el) => {
  if (!el || !transcriptContainerRef.value) return false;
  const parentRect = transcriptContainerRef.value.getBoundingClientRect();
  const elRect = el.getBoundingClientRect();
  const tolerance = 10;

  return (
    elRect.top >= parentRect.top - tolerance &&
    elRect.bottom <= parentRect.bottom + tolerance
  );
};

// ---- Word Explanation Logic ----

const handleWordDoubleClick = async (event) => {
  const selection = window.getSelection();
  const word = selection
    .toString()
    .trim()
    .replace(/[.,;:?!]$/, "");

  if (!word) {
    // Optionally clear parts of the panel if desired, or just return
    // clearExplanationPanel(false); // Example: keep word but clear explanations
    return;
  }

  let context = "";
  const CXT_WORDS_RADIUS = 10;
  if (event.target && event.target.textContent) {
    const fullText = event.target.textContent;
    const wordIndex = fullText.toLowerCase().indexOf(word.toLowerCase()); // Case-insensitive search for index
    if (wordIndex !== -1) {
      const beforeText = fullText.substring(0, wordIndex).trim();
      const afterText = fullText.substring(wordIndex + word.length).trim();
      const wordsBefore = beforeText.split(/\s+/).slice(-CXT_WORDS_RADIUS);
      const wordsAfter = afterText.split(/\s+/).slice(0, CXT_WORDS_RADIUS);
      context = [...wordsBefore, word, ...wordsAfter].join(" ");
    } else {
      context = event.target.textContent; // Fallback
    }
  } else {
    context = word; // Fallback if no target text
  }

  selectedWord.value = word;
  selectedContext.value = context;

  console.log(
    `Word selected: "${selectedWord.value}", Context: "${selectedContext.value}", Audio Lang: ${audioSourceLanguage.value}, User Native Lang: ${userSelectedNativeLanguage.value}`
  );

  audioSourceExplanation.value = "";
  userNativeExplanation.value = "";
  explanationError.value = "";
  isExplanationLoading.value = true;
  isExplanationPanelVisible.value = true; // Ensure panel is visible
  // showUserNativeExplanationPanel.value = false; // Reset to hidden by default on new word selection

  try {
    const llmSettings = JSON.parse(localStorage.getItem("llmSettings") || "{}");
    const response = await axios.post(`${API_URL}/api/explain-word`, {
      word: selectedWord.value,
      context: selectedContext.value,
      audio_language: audioSourceLanguage.value || "English", // Fallback if somehow null
      user_native_language: userSelectedNativeLanguage.value,
      llm_settings: llmSettings,
    });

    if (response.data) {
      audioSourceExplanation.value =
        response.data.audio_language_explanation || "";
      userNativeExplanation.value =
        response.data.native_language_explanation || "";

      if (!audioSourceExplanation.value && !userNativeExplanation.value) {
        explanationError.value =
          "No explanations were returned from the server.";
      }
    } else {
      explanationError.value =
        "Received an empty response from the explanation API.";
      // throw new Error("Empty response from explanation API."); // Or handle as error
    }
  } catch (err) {
    console.error("Error fetching explanation:", err);
    if (err.response && err.response.data && err.response.data.error) {
      explanationError.value = `API Error: ${err.response.data.error}`;
    } else if (err.message) {
      explanationError.value = `Request Error: ${err.message}`;
    } else {
      explanationError.value =
        "An unknown error occurred while fetching the explanation.";
    }
    audioSourceExplanation.value = "";
    userNativeExplanation.value = "";
  } finally {
    isExplanationLoading.value = false;
  }
};

const clearExplanationPanel = (clearAll = true) => {
  if (clearAll) {
    selectedWord.value = "";
    selectedContext.value = "";
    isExplanationPanelVisible.value = false;
  }
  audioSourceExplanation.value = "";
  userNativeExplanation.value = "";
  explanationError.value = "";
  isExplanationLoading.value = false;
  showUserNativeExplanationPanel.value = false; // Always hide native panel on clear
};

const handleUserNativeLanguageChange = () => {
  if (selectedWord.value) {
    console.log(
      "User native language changed to:",
      userSelectedNativeLanguage.value,
      "; Re-fetching for word:",
      selectedWord.value
    );
    // Set loading state specifically for the native explanation part if possible,
    // or for the whole panel if simpler.
    // For now, re-triggering the full double-click logic to get both explanations updated.
    // This ensures consistency if the native language is now the same as audio language.
    handleWordDoubleClick({ target: { textContent: selectedContext.value } }); // Simulate event to re-fetch
  }
};

// Watcher for when the main panel visibility changes
watch(isExplanationPanelVisible, (newValue) => {
  if (!newValue) {
    clearExplanationPanel(true); // Clear everything when panel is hidden
  }
});

// No need for a separate watcher on userSelectedNativeLanguage if handleUserNativeLanguageChange covers it.

// ---- Keyboard Shortcuts ----
const handleKeyDown = (event) => {
  // Ignore shortcuts if the user is typing in an input, select, or textarea
  const targetTagName = event.target.tagName.toUpperCase();
  if (["INPUT", "TEXTAREA", "SELECT"].includes(targetTagName)) {
    return;
  }

  // Shortcuts should only work when an audio source is available.
  if (!audioUrl.value) {
    return;
  }

  switch (event.key) {
    case " ":
      event.preventDefault();
      togglePlayback();
      break;
    case "ArrowLeft":
      event.preventDefault();
      navigateToSegment("previous");
      break;
    case "ArrowRight":
      event.preventDefault();
      navigateToSegment("next");
      break;
  }
};

// ---- Lifecycle Hooks ----

// ---- Methods ----

const toggleExplanationPanel = () => {
  isExplanationPanelVisible.value = !isExplanationPanelVisible.value;
  if (!isExplanationPanelVisible.value) {
    // If hiding panel
    clearExplanationPanel(true);
  }
};

// ---- Translation ----
async function translateSegment(segment, index) {
  if (!segment || !segment.text) return;
  if (!userSelectedNativeLanguage.value) {
    console.error("User native language not set for translation.");
    if (
      transcript.value &&
      transcript.value.segments &&
      transcript.value.segments[index]
    ) {
      transcript.value.segments[index].translationError =
        "User native language not selected in panel.";
    }
    return;
  }

  // Use segment.id if available and unique, otherwise fallback to index for safety if needed
  const segmentId = segment.id !== undefined ? segment.id : index;
  translatingSegmentId.value = segmentId;

  // Ensure the segment object is extensible if it comes directly from a non-Proxy source
  // This might not be strictly necessary if transcript.value.segments are already reactive objects.
  const targetSegment = transcript.value.segments[index];
  if (!targetSegment) {
    console.error(`Segment at index ${index} not found for translation.`);
    translatingSegmentId.value = null;
    return;
  }

  targetSegment.translationError = null; // Clear previous error
  targetSegment.translation = null; // Clear previous translation
  // targetSegment.translationLoading = true; // Optional: for more granular loading state on the segment itself

  try {
    console.log(
      `[AudioTranscriber] Translating segment ID ${segmentId}: "${segment.text.substring(
        0,
        30
      )}..." to ${userSelectedNativeLanguage.value}`
    );
    // Read LLM settings from localStorage
    let llmSettings = {};
    try {
      llmSettings = JSON.parse(localStorage.getItem("llmSettings") || "{}");
    } catch (e) {
      llmSettings = {};
    }
    const response = await axios.post(`${API_URL}/api/translate-sentence`, {
      sentence: segment.text,
      target_language: userSelectedNativeLanguage.value,
      llm_settings: llmSettings,
    });

    if (response.data && response.data.translated_sentence) {
      targetSegment.translation = response.data.translated_sentence;
      console.log(
        `[AudioTranscriber] Translation success for segment ID ${segmentId}: "${targetSegment.translation.substring(
          0,
          30
        )}..."`
      );
    } else if (response.data && response.data.error) {
      console.error(
        `[AudioTranscriber] Translation API error for segment ID ${segmentId}:`,
        response.data.error
      );
      targetSegment.translationError = response.data.error;
    } else {
      console.error(
        `[AudioTranscriber] Unexpected translation response for segment ID ${segmentId}:`,
        response.data
      );
      targetSegment.translationError =
        "Unexpected response from translation server.";
    }
  } catch (err) {
    console.error(
      `[AudioTranscriber] Failed to translate segment ID ${segmentId}:`,
      err
    );
    targetSegment.translationError =
      err.response?.data?.error ||
      "Network error or server issue during translation.";
  } finally {
    translatingSegmentId.value = null;
    // targetSegment.translationLoading = false;
  }
}

// Watch for current segment change and scroll into view if needed
watch(currentTime, async () => {
  if (!transcript.value || !transcript.value.segments) return;
  const currentActiveSegment = transcript.value.segments.find((s) =>
    isCurrentSegment(s)
  );
  if (currentActiveSegment && !userHasScrolledAway.value) {
    const segmentElement = segmentRefs.value.get(currentActiveSegment.start);
    if (segmentElement) {
      await nextTick();
      segmentElement.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  }
});
</script>

<style scoped>
/* Custom scrollbar (already present) */
.scroll-smooth::-webkit-scrollbar {
  width: 8px;
}
.scroll-smooth::-webkit-scrollbar-track {
  background: #e5e7eb; /* Tailwind gray-200 */
  border-radius: 4px;
}
.scroll-smooth::-webkit-scrollbar-thumb {
  background: #9ca3af; /* Tailwind gray-400 */
  border-radius: 4px;
}
.scroll-smooth::-webkit-scrollbar-thumb:hover {
  background: #6b7280; /* Tailwind gray-500 */
}

/* Ensure enough space for fixed header/footer if any, or for general layout */
.container {
  padding-bottom: 2rem;
}

/* You might want to ensure the main flex container allows vertical scroll if content overflows on small screens */
.max-w-full {
  /* Ensures it doesn't exceed viewport width, good for responsive */
}
</style>

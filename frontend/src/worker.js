console.debug = console.log;
import {env, pipeline } from "@huggingface/transformers";

// env.remoteHost = 'https://hf-mirror.com';
env.allowRemoteModels = false;
env.allowLocalModels = true;
env.localModelPath = '/models/';
env.backends.onnx.wasm.wasmPaths = '/wasm/';

class TranscriptionPipeline {
  static task = "automatic-speech-recognition";
  static model = "Xenova/whisper-tiny"; // Using the 'base' model for a good balance of speed and accuracy
  // static model = "openai/whisper-base"; // Using the 'base' model for a good balance of speed and accuracy
  static instance = null;

  static async getInstance(progress_callback = null) {
    if (this.instance === null) {
      this.instance = await pipeline(this.task, this.model, {
        progress_callback,
        // device: 'webgpu',
        // dtype: {"encoder_model": "fp32","decoder_model_merged": "q4"}
      });
    }
    console.log(this.instance);
    return this.instance;
  }
}

self.addEventListener("message", async (event) => {
  const { type, audio } = event.data;
  if (type === "transcribe") {
    try {
      // 1. Convert audio data to a format that can be processed by the model
      const audioData = new Float32Array(audio);

      // 2. Get the singleton instance of the pipeline
      const transcriber = await TranscriptionPipeline.getInstance((progress) => {
        // Post progress messages back to the main thread
        self.postMessage({ status: "progress", data: progress });
      });

      // 3. Transcribe the audio
      const output = await transcriber(audioData, {
        chunk_length_s: 30, // Process in 30-second chunks
        stride_length_s: 5, // Overlap chunks by 5 seconds
        batch_size: 8,
        return_timestamps: true,
        callback_function: (callback) => {
            // This function is called for each chunk, but we only need the final output
            // We can still use it for more granular progress if needed in the future
            const chunk = callback[0];
            if(chunk.is_last){
                // We can use this to signal completion of a chunk if needed
            }
        },
      });

      // 4. Post the final transcription result
      self.postMessage({ status: "complete", data: output });
    } catch (e) {
      // Post an error message if something goes wrong
      self.postMessage({ status: "error", data: e.message });
    }
  }
}); 
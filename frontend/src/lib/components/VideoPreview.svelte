<script lang="ts">
  import { onMount } from 'svelte';
  import { editor, currentSegment } from '$lib/stores/editor';

  let videoElement: HTMLVideoElement;
  let videoSrc = '';

  $: if (videoElement) {
    if ($editor.isPlaying && videoElement.paused) {
      videoElement.play().catch(e => console.error("Playback failed:", e));
    } else if (!$editor.isPlaying && !videoElement.paused) {
      videoElement.pause();
    }
    videoElement.playbackRate = $editor.playbackRate;
    videoElement.volume = $editor.volume;
  }

  $: if (videoElement && Math.abs((videoElement.currentTime * 1000) - $editor.currentTimeMs) > 100) {
    if (!$editor.isPlaying) {
      videoElement.currentTime = $editor.currentTimeMs / 1000;
    }
  }

  function handleTimeUpdate() {
    if ($editor.isPlaying && videoElement) {
      editor.setCurrentTime(Math.floor(videoElement.currentTime * 1000));
    }
  }

  function handlePlay() {
    editor.setPlaying(true);
  }

  function handlePause() {
    editor.setPlaying(false);
  }

  function setLayer(layer: 'source' | 'l1' | 'l2' | 'dual') {
    editor.setSubtitleLayer(layer);
  }
</script>

<div class="relative w-full h-full bg-[#0a0a12] flex flex-col">
  <div class="flex-1 relative flex items-center justify-center bg-black overflow-hidden group">
    {#if videoSrc}
      <video
        bind:this={videoElement}
        src={videoSrc}
        class="w-full h-full object-contain"
        on:timeupdate={handleTimeUpdate}
        on:play={handlePlay}
        on:pause={handlePause}
        controls
      />
    {:else}
      <div class="text-gray-600 flex flex-col items-center">
        <span class="material-symbols-outlined text-4xl mb-2">movie</span>
        <span>Pratinjau Video</span>
      </div>
    {/if}

    {#if $currentSegment}
      <div class="absolute bottom-12 left-0 right-0 flex flex-col items-center pointer-events-none px-8">
        {#if $editor.subtitleLayer === 'source' || $editor.subtitleLayer === 'dual'}
          <div class="text-white text-xl md:text-2xl font-bold bg-black/60 px-4 py-1 rounded mb-1 text-center max-w-3xl drop-shadow-md">
            {$currentSegment.source}
          </div>
        {/if}
        {#if $editor.subtitleLayer === 'l1'}
          <div class="text-[#ffe04a] text-xl md:text-2xl font-bold bg-black/60 px-4 py-1 rounded text-center max-w-3xl drop-shadow-md">
            {$currentSegment.l1 || $currentSegment.source}
          </div>
        {/if}
        {#if $editor.subtitleLayer === 'l2' || $editor.subtitleLayer === 'dual'}
          <div class="text-[#00ffcc] text-xl md:text-2xl font-bold bg-black/60 px-4 py-1 rounded text-center max-w-3xl drop-shadow-md">
            {$currentSegment.l2 || $currentSegment.l1 || $currentSegment.source}
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <div class="p-3 bg-[#0f0f1a] border-t border-gray-800 flex items-center justify-center gap-4">
    <button 
      class="px-4 py-1 text-sm rounded transition-all {$editor.subtitleLayer === 'source' ? 'bg-[#ff2d78]/20 text-[#ff2d78] border border-[#ff2d78]' : 'text-gray-400 hover:text-white'}"
      on:click={() => setLayer('source')}>
      Sumber
    </button>
    <button 
      class="px-4 py-1 text-sm rounded transition-all {$editor.subtitleLayer === 'l1' ? 'bg-[#ffe04a]/20 text-[#ffe04a] border border-[#ffe04a]' : 'text-gray-400 hover:text-white'}"
      on:click={() => setLayer('l1')}>
      Terjemahan
    </button>
    <button 
      class="px-4 py-1 text-sm rounded transition-all {$editor.subtitleLayer === 'l2' ? 'bg-[#00ffcc]/20 text-[#00ffcc] border border-[#00ffcc]' : 'text-gray-400 hover:text-white'}"
      on:click={() => setLayer('l2')}>
      Hasil Akhir
    </button>
    <button 
      class="px-4 py-1 text-sm rounded transition-all {$editor.subtitleLayer === 'dual' ? 'bg-white/20 text-white border border-white' : 'text-gray-400 hover:text-white'}"
      on:click={() => setLayer('dual')}>
      Dual
    </button>
  </div>
</div>
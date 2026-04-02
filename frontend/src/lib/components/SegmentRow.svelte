<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Segment } from '$lib/wails';
  import { editor } from '$lib/stores/editor';

  export let segment: Segment;
  export let index: number;
  export let isSelected: boolean = false;

  const dispatch = createEventDispatcher<{
    select: string;
  }>();

  function formatTime(ms: number) {
    const totalSeconds = Math.floor(ms / 1000);
    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const seconds = totalSeconds % 60;
    const millis = ms % 1000;

    return `${hours.toString().padStart(2, '0')}:${minutes
      .toString()
      .padStart(2, '0')}:${seconds.toString().padStart(2, '0')}.${millis
      .toString()
      .padStart(3, '0')}`;
  }

  function getStatusColor(status: string) {
    switch (status) {
      case 'pass':
        return 'text-[#00ffcc] border-[#00ffcc]';
      case 'warn':
        return 'text-[#ffe04a] border-[#ffe04a]';
      case 'error':
        return 'text-[#ff2d78] border-[#ff2d78]';
      default:
        return 'text-gray-400 border-gray-600';
    }
  }

  function handleClick() {
    dispatch('select', segment.id);
  }
</script>

<div
  class="p-3 border-b border-gray-800 hover:bg-gray-800/50 cursor-pointer transition-colors duration-200 {isSelected ? 'bg-gray-800/80 neon-border-secondary' : ''}"
  on:click={handleClick}
  on:keydown={(e) => e.key === 'Enter' && handleClick()}
  role="button"
  tabindex="0"
>
  <div class="flex items-center justify-between mb-1">
    <div class="flex items-center gap-2 text-xs">
      <span class="text-gray-500 w-6">{index + 1}</span>
      <span class="text-gray-400 font-mono">
        {formatTime(segment.startMs)} - {formatTime(segment.endMs)}
      </span>
    </div>
    
    <div class="flex gap-1">
      <div 
        class="w-2 h-2 rounded-full border border-solid {getStatusColor(segment.qaStatus || 'pending')} 
        {segment.qaStatus === 'pass' ? 'bg-[#00ffcc]/20' : 
         segment.qaStatus === 'warn' ? 'bg-[#ffe04a]/20' : 
         segment.qaStatus === 'error' ? 'bg-[#ff2d78]/20' : 'bg-transparent'}"
        title="Status: {segment.qaStatus || 'Pending'}"
      ></div>
    </div>
  </div>

  <div class="text-sm text-gray-200 truncate pr-4">
    {segment.l2 || segment.l1 || segment.source}
  </div>
</div>
<script lang="ts">
  import { onMount } from 'svelte';
  import { editor, selectedSegment } from '$lib/stores/editor';
  import SegmentRow from './SegmentRow.svelte';

  let listContainer: HTMLDivElement;

  function handleSelect(event: CustomEvent<string>) {
    editor.selectSegment(event.detail);
    const selected = $editor.segments.find((s) => s.id === event.detail);
    if (selected) {
      editor.setCurrentTime(selected.startMs);
    }
  }

  // Auto-scroll to selected segment
  $: if (listContainer && $selectedSegment) {
    const selectedElement = listContainer.querySelector('.neon-border-secondary');
    if (selectedElement) {
      selectedElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }
</script>

<div class="flex flex-col h-full bg-[#0a0a12] border-l border-gray-800">
  <div class="p-4 border-b border-gray-800 bg-[#0f0f1a]">
    <h2 class="text-lg font-semibold text-[#00ffcc]">Daftar Segmen</h2>
  </div>

  <div 
    class="flex-1 overflow-y-auto"
    bind:this={listContainer}
  >
    {#each $editor.segments as segment, index (segment.id)}
      <SegmentRow 
        {segment} 
        {index} 
        isSelected={$selectedSegment?.id === segment.id}
        on:select={handleSelect}
      />
    {/each}
    {#if $editor.segments.length === 0}
      <div class="flex items-center justify-center h-full text-gray-500">
        Belum ada segmen
      </div>
    {/if}
  </div>
</div>
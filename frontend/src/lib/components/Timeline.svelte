<script lang="ts">
  import { onMount } from 'svelte';
  import { editor, visibleSegments, selectedSegment } from '$lib/stores/editor';
  import { UpdateSegment } from '$lib/wails';
  import type { Segment } from '$lib/wails';

  let containerRef: HTMLDivElement;
  let timelineWidth = 0;
  
  $: msPerPixel = 60000 / ($editor.zoom * timelineWidth || 1000);

  function getLeft(timeMs: number) {
    return Math.max(0, (timeMs - $editor.scrollPositionMs) / msPerPixel);
  }

  function getWidth(startMs: number, endMs: number) {
    return Math.max(2, (endMs - startMs) / msPerPixel);
  }

  function handleWheel(e: WheelEvent) {
    if (e.ctrlKey) {
      e.preventDefault();
      if (e.deltaY < 0) {
        editor.zoomIn();
      } else {
        editor.zoomOut();
      }
    } else {
      editor.setScrollPosition($editor.scrollPositionMs + e.deltaY * msPerPixel);
    }
  }

  let draggingSegment: Segment | null = null;
  let dragType: 'start' | 'end' | 'move' | null = null;
  let initialDragX = 0;
  let initialStartMs = 0;
  let initialEndMs = 0;

  function handlePointerDown(e: PointerEvent, seg: Segment, type: 'start' | 'end' | 'move') {
    e.stopPropagation();
    e.preventDefault();
    editor.selectSegment(seg.id);
    draggingSegment = seg;
    dragType = type;
    initialDragX = e.clientX;
    initialStartMs = seg.startMs;
    initialEndMs = seg.endMs;
    window.addEventListener('pointermove', handlePointerMove);
    window.addEventListener('pointerup', handlePointerUp);
  }

  function handlePointerMove(e: PointerEvent) {
    if (!draggingSegment || !dragType) return;
    
    const deltaX = e.clientX - initialDragX;
    const deltaMs = deltaX * msPerPixel;
    const snapEnabled = !e.altKey;

    let newStart = initialStartMs;
    let newEnd = initialEndMs;

    if (dragType === 'start') {
      newStart = Math.max(0, initialStartMs + deltaMs);
    } else if (dragType === 'end') {
      newEnd = initialEndMs + deltaMs;
    } else if (dragType === 'move') {
      newStart = Math.max(0, initialStartMs + deltaMs);
      newEnd = initialEndMs + deltaMs;
    }

    if (snapEnabled) {
      const snapThreshold = 100 * msPerPixel;
      for (const other of $editor.segments) {
        if (other.id === draggingSegment.id) continue;
        
        if (dragType === 'start' || dragType === 'move') {
          if (Math.abs(newStart - other.endMs) < snapThreshold) {
            const shift = other.endMs + 83 - newStart;
            newStart += shift;
            if (dragType === 'move') newEnd += shift;
          }
        }
        
        if (dragType === 'end' || dragType === 'move') {
          if (Math.abs(newEnd - other.startMs) < snapThreshold) {
            const shift = other.startMs - 83 - newEnd;
            newEnd += shift;
            if (dragType === 'move') newStart += shift;
          }
        }
      }
    }

    if (newEnd <= newStart + 100) return;

    editor.updateSegment(draggingSegment.id, {
      startMs: Math.round(newStart),
      endMs: Math.round(newEnd)
    });
  }

  async function handlePointerUp() {
    if (draggingSegment) {
      const current = $editor.segments.find(s => s.id === draggingSegment!.id);
      if (current) {
        await UpdateSegment(current);
      }
    }
    draggingSegment = null;
    dragType = null;
    window.removeEventListener('pointermove', handlePointerMove);
    window.removeEventListener('pointerup', handlePointerUp);
  }

  function handleTimelineClick(e: MouseEvent) {
    if (!containerRef) return;
    const rect = containerRef.getBoundingClientRect();
    const clickX = e.clientX - rect.left;
    const timeMs = $editor.scrollPositionMs + clickX * msPerPixel;
    editor.setCurrentTime(Math.round(timeMs));
  }

  function checkOverlap(seg: Segment): boolean {
    return $editor.segments.some(other => 
      other.id !== seg.id && 
      ((seg.startMs >= other.startMs && seg.startMs < other.endMs) ||
       (seg.endMs > other.startMs && seg.endMs <= other.endMs) ||
       (seg.startMs <= other.startMs && seg.endMs >= other.endMs))
    );
  }
</script>

<div class="flex flex-col h-full bg-[#0a0a12] border-b border-gray-800">
  <div class="px-4 py-2 border-b border-gray-800 flex items-center justify-between bg-[#0f0f1a]">
    <h2 class="text-[#ffe04a] font-semibold flex items-center gap-2">
      <span class="material-symbols-outlined">linear_scale</span> Timeline
    </h2>
    <div class="text-xs text-gray-500 flex gap-4">
      <span>Zoom: {$editor.zoom.toFixed(1)}x</span>
      <span>Tahan Alt untuk mematikan snap</span>
    </div>
  </div>

  <div 
    class="relative flex-1 overflow-hidden"
    bind:this={containerRef}
    bind:clientWidth={timelineWidth}
    on:wheel={handleWheel}
    on:click={handleTimelineClick}
    role="presentation"
  >
    <div class="absolute inset-0 pointer-events-none" style="background-image: linear-gradient(to right, #1f2937 1px, transparent 1px); background-size: {1000 / msPerPixel}px 100%;"></div>

    {#each $visibleSegments as seg (seg.id)}
      {@const isSelected = $selectedSegment?.id === seg.id}
      {@const isOverlapping = checkOverlap(seg)}
      
      <div
        class="absolute top-1/2 -translate-y-1/2 h-16 rounded cursor-pointer border group select-none transition-colors
          {isSelected ? 'bg-[#00ffcc]/20 border-[#00ffcc] z-10' : 
           isOverlapping ? 'bg-[#ff2d78]/20 border-[#ff2d78]' : 'bg-[#ffe04a]/10 border-[#ffe04a]/50'}"
        style="left: {getLeft(seg.startMs)}px; width: {getWidth(seg.startMs, seg.endMs)}px;"
        on:pointerdown={(e) => handlePointerDown(e, seg, 'move')}
        role="presentation"
      >
        <div class="px-2 py-1 truncate text-xs font-mono text-gray-300">
          {seg.l2 || seg.l1 || seg.source}
        </div>
        
        {#if isSelected}
          <div 
            class="absolute left-0 top-0 bottom-0 w-2 cursor-w-resize bg-[#00ffcc]/50 hover:bg-[#00ffcc]"
            on:pointerdown={(e) => handlePointerDown(e, seg, 'start')}
            role="presentation"
          ></div>
          <div 
            class="absolute right-0 top-0 bottom-0 w-2 cursor-e-resize bg-[#00ffcc]/50 hover:bg-[#00ffcc]"
            on:pointerdown={(e) => handlePointerDown(e, seg, 'end')}
            role="presentation"
          ></div>
        {/if}
      </div>
    {/each}

    <div 
      class="absolute top-0 bottom-0 w-[2px] bg-[#ff2d78] pointer-events-none z-20 shadow-[0_0_8px_#ff2d78]"
      style="left: {getLeft($editor.currentTimeMs)}px;"
    >
      <div class="absolute top-0 -translate-x-1/2 bg-[#ff2d78] text-white text-[10px] px-1 py-0.5 rounded font-mono">
        {($editor.currentTimeMs / 1000).toFixed(2)}s
      </div>
    </div>
  </div>
</div>
<script lang="ts">
  import { onMount } from 'svelte';
  import { editor, selectedSegment } from '$lib/stores/editor';
  import { history } from '$lib/stores/history';
  import { GetSegments, SplitSegment, MergeSegments, DeleteSegment } from '$lib/wails';
  import { shortcuts } from '$lib/utils/shortcuts';
  
  import VideoPreview from '$lib/components/VideoPreview.svelte';
  import SegmentList from '$lib/components/SegmentList.svelte';
  import Timeline from '$lib/components/Timeline.svelte';
  import InlineEditor from '$lib/components/InlineEditor.svelte';

  onMount(async () => {
    try {
      const segments = await GetSegments($editor.projectId || 'default');
      editor.setSegments(segments);
      history.init(segments, $selectedSegment?.id || null);
    } catch (e) {
      console.error("Failed to load segments:", e);
    }
  });

  function handleShortcut(event: any) {
    const { key, code, ctrl, shift, meta, alt, originalEvent } = event.detail;

    if (code === 'Space') {
      originalEvent.preventDefault();
      editor.togglePlay();
      return;
    }

    if (code === 'ArrowLeft') {
      editor.setCurrentTime(Math.max(0, $editor.currentTimeMs - 5000));
      return;
    }
    if (code === 'ArrowRight') {
      editor.setCurrentTime($editor.currentTimeMs + 5000);
      return;
    }

    if (code === 'ArrowUp' || code === 'ArrowDown') {
      originalEvent.preventDefault();
      if (!$selectedSegment) {
        if ($editor.segments.length > 0) editor.selectSegment($editor.segments[0].id);
        return;
      }
      
      const index = $editor.segments.findIndex(s => s.id === $selectedSegment?.id);
      if (index !== -1) {
        const nextIndex = code === 'ArrowUp' ? Math.max(0, index - 1) : Math.min($editor.segments.length - 1, index + 1);
        editor.selectSegment($editor.segments[nextIndex].id);
        editor.setCurrentTime($editor.segments[nextIndex].startMs);
      }
      return;
    }

    if (ctrl || meta) {
      if (key === 'z' || key === 'Z') {
        originalEvent.preventDefault();
        const entry = shift ? history.redo() : history.undo();
        if (entry) {
          editor.setSegments(entry.segments);
          editor.selectSegment(entry.selectedSegmentId);
        }
        return;
      }
      if (key === 'y' || key === 'Y') {
        originalEvent.preventDefault();
        const entry = history.redo();
        if (entry) {
          editor.setSegments(entry.segments);
          editor.selectSegment(entry.selectedSegmentId);
        }
        return;
      }
      if (key === 's' || key === 'S') {
        originalEvent.preventDefault();
        editor.markClean();
        return;
      }
    }

    if ($selectedSegment) {
      if (code === 'KeyS' && !ctrl && !meta) {
        SplitSegment($selectedSegment.id, $editor.currentTimeMs, 0);
        return;
      }
      if (code === 'KeyM' && !ctrl && !meta) {
        const idx = $editor.segments.findIndex(s => s.id === $selectedSegment?.id);
        if (idx !== -1 && idx < $editor.segments.length - 1) {
          MergeSegments($selectedSegment.id, $editor.segments[idx + 1].id);
        }
        return;
      }
      if (code === 'Delete' || code === 'Backspace') {
        DeleteSegment($selectedSegment.id);
        return;
      }
    }

    if (code === 'Digit1') editor.setSubtitleLayer('source');
    if (code === 'Digit2') editor.setSubtitleLayer('l1');
    if (code === 'Digit3') editor.setSubtitleLayer('l2');
    if (code === 'Digit4') editor.setSubtitleLayer('dual');
  }
</script>

<svelte:head>
  <title>SubFlow - Editor</title>
</svelte:head>

<div class="w-full h-screen bg-[#0a0a12] text-white flex flex-col overflow-hidden" use:shortcuts on:shortcut={handleShortcut}>
  <div class="flex-1 grid grid-cols-[1fr_350px] grid-rows-[1fr_200px_auto] overflow-hidden gap-[1px] bg-gray-800">
    <div class="col-start-1 row-start-1 overflow-hidden bg-[#0a0a12]">
      <VideoPreview />
    </div>
    
    <div class="col-start-2 row-start-1 overflow-hidden bg-[#0a0a12]">
      <SegmentList />
    </div>
    
    <div class="col-start-1 col-end-3 row-start-2 overflow-hidden bg-[#0a0a12]">
      <Timeline />
    </div>

    <div class="col-start-1 col-end-3 row-start-3 overflow-visible bg-[#0a0a12]">
      <InlineEditor />
    </div>
  </div>
</div>
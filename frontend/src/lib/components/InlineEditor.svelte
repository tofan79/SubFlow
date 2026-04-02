<script lang="ts">
  import { selectedSegment, editor } from '$lib/stores/editor';
  import { history } from '$lib/stores/history';
  import { UpdateSegment, SplitSegment, MergeSegments, DeleteSegment, RetryL1, RetryL2 } from '$lib/wails';

  async function handleUpdate(field: 'l1' | 'l2' | 'startMs' | 'endMs', value: any) {
    if (!$selectedSegment) return;
    
    history.push('Update Segment', $editor.segments, $selectedSegment.id);
    const changes = { [field]: value };
    editor.updateSegment($selectedSegment.id, changes);
    
    await UpdateSegment({ ...$selectedSegment, ...changes });
  }

  async function handleSplit() {
    if (!$selectedSegment) return;
    history.push('Split Segment', $editor.segments, $selectedSegment.id);
    await SplitSegment($selectedSegment.id, $editor.currentTimeMs, 0);
  }

  async function handleMerge() {
    if (!$selectedSegment) return;
    const idx = $editor.segments.findIndex(s => s.id === $selectedSegment?.id);
    if (idx !== -1 && idx < $editor.segments.length - 1) {
      history.push('Merge Segment', $editor.segments, $selectedSegment.id);
      await MergeSegments($selectedSegment.id, $editor.segments[idx + 1].id);
    }
  }

  async function handleRetryL1() {
    if (!$selectedSegment) return;
    history.push('Retry L1', $editor.segments, $selectedSegment.id);
    await RetryL1($selectedSegment.id);
  }

  async function handleRetryL2() {
    if (!$selectedSegment) return;
    history.push('Retry L2', $editor.segments, $selectedSegment.id);
    await RetryL2($selectedSegment.id);
  }

  async function handleDelete() {
    if (!$selectedSegment) return;
    history.push('Delete Segment', $editor.segments, null);
    await DeleteSegment($selectedSegment.id);
  }
</script>

<div class="flex flex-col h-full bg-[#0a0a12] p-4">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-xl font-bold text-white tracking-wide">Editor</h2>
    <div class="flex gap-2">
      <button class="btn-neon-secondary flex items-center justify-center p-2 rounded" on:click={handleSplit} title="Pisah" disabled={!$selectedSegment}>
        <span class="material-symbols-outlined text-sm">call_split</span>
      </button>
      <button class="btn-neon-secondary flex items-center justify-center p-2 rounded" on:click={handleMerge} title="Gabung" disabled={!$selectedSegment}>
        <span class="material-symbols-outlined text-sm">merge</span>
      </button>
      <button class="btn-neon-secondary flex items-center justify-center p-2 rounded" on:click={handleRetryL1} title="Coba Ulang L1" disabled={!$selectedSegment}>
        <span class="material-symbols-outlined text-sm">refresh</span>
      </button>
      <button class="btn-neon-secondary flex items-center justify-center p-2 rounded" on:click={handleRetryL2} title="Coba Ulang L2" disabled={!$selectedSegment}>
        <span class="material-symbols-outlined text-sm text-[#00ffcc]">refresh</span>
      </button>
      <button class="flex items-center justify-center p-2 rounded text-[#ff2d78] hover:bg-[#ff2d78]/10 border border-[#ff2d78]/50 transition-colors" on:click={handleDelete} title="Hapus" disabled={!$selectedSegment}>
        <span class="material-symbols-outlined text-sm">delete</span>
      </button>
    </div>
  </div>

  {#if $selectedSegment}
    <div class="flex flex-col gap-4 overflow-y-auto pr-2">
      <div class="flex gap-4 items-center">
        <div class="flex-1 flex gap-2">
          <label class="flex flex-col w-full">
            <span class="text-xs text-gray-500 mb-1">Mulai (ms)</span>
            <input 
              type="number" 
              class="input-neon w-full bg-[#141422]" 
              value={$selectedSegment.startMs}
              on:change={(e) => handleUpdate('startMs', parseInt(e.currentTarget.value))}
            />
          </label>
          <label class="flex flex-col w-full">
            <span class="text-xs text-gray-500 mb-1">Akhir (ms)</span>
            <input 
              type="number" 
              class="input-neon w-full bg-[#141422]" 
              value={$selectedSegment.endMs}
              on:change={(e) => handleUpdate('endMs', parseInt(e.currentTarget.value))}
            />
          </label>
        </div>
      </div>

      <label class="flex flex-col">
        <span class="text-sm text-gray-400 mb-1">Sumber</span>
        <textarea 
          class="input-neon bg-[#141422] resize-none h-20 text-gray-400" 
          readonly
          value={$selectedSegment.source}
        ></textarea>
      </label>

      <label class="flex flex-col">
        <span class="text-sm text-[#ffe04a] mb-1">Terjemahan</span>
        <textarea 
          class="input-neon bg-[#141422] resize-y h-20" 
          value={$selectedSegment.l1}
          on:change={(e) => handleUpdate('l1', e.currentTarget.value)}
        ></textarea>
      </label>

      <label class="flex flex-col">
        <span class="text-sm text-[#00ffcc] mb-1">Hasil Akhir</span>
        <textarea 
          class="input-neon bg-[#141422] resize-y h-20" 
          value={$selectedSegment.l2}
          on:change={(e) => handleUpdate('l2', e.currentTarget.value)}
        ></textarea>
      </label>
    </div>
  {:else}
    <div class="flex-1 flex items-center justify-center text-gray-500">
      Pilih segmen dari daftar atau timeline
    </div>
  {/if}
</div>
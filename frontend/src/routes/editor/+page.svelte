<script lang="ts">
  let currentPhase = 1;

  const ph1Segments = [
    { id: '#001', time: '00:00:12.400 -> 00:00:15.800', acc: 98, active: true, text: 'The neon skyline of Neo-Tokyo stretched endlessly before us, a concrete jungle pulsating with electric life.' },
    { id: '#002', time: '00:00:15.900 -> 00:00:18.200', acc: 92, text: 'Every flicker of light told a story of ambition and survival in the lower sectors.' },
    { id: '#003', time: '00:00:18.500 -> 00:00:22.000', acc: 74, warn: true, text: "We didn't come here to play by the rules; we came to rewrite the code of the city." },
    { id: '#004', time: '00:00:22.300 -> 00:00:25.100', acc: 99, text: 'Initiating synchronization with the main grid. Stand by for neural link.' }
  ];

  const ph2Rows = [
    { id: '#042', src: 'The city screams in neon colors while the clouds weep acid rain above the chrome towers.', l1: 'Kota itu berteriak dalam warna-warna neon sementara awan menangis hujan asam di atas menara krom.', polishHtml: '<span class="text-error line-through mr-1">Kota itu berteriak</span> <span class="neon-text-secondary font-bold mr-1">Metropolis membara</span> dengan pijar neon saat langit <span class="text-error line-through mr-1">menangis</span> <span class="neon-text-secondary font-bold mr-1">mencurahkan</span> tangis asam di sela gedung-gedung krom.' },
    { id: '#043', src: 'His cybernetic eye flickered as he entered the dark alleyway.', l1: 'Mata sibernetikanya berkedip saat dia memasuki lorong gelap.', polishHtml: '<span class="neon-text-secondary font-bold mr-1">Optik</span> sibernetikanya <span class="text-error line-through mr-1">berkedip</span> <span class="neon-text-secondary font-bold mr-1">berderak</span> saat ia merangsek masuk ke dalam <span class="text-error line-through mr-1">lorong</span> <span class="neon-text-secondary font-bold mr-1">gang sempit</span> yang pekat.' },
    { id: '#044', src: 'System override initiated. All data protocols are now compromised.', l1: 'Pengambilalihan sistem dimulai. Semua protokol data sekarang terkompromi.', polishHtml: '<span class="neon-text-secondary font-bold mr-1">Inisiasi bypass sistem.</span> Seluruh protokol data telah <span class="text-error line-through mr-1">sekarang terkompromi</span> <span class="neon-text-secondary font-bold mr-1">ditembus sepenuhnya.</span>' },
  ];

  const ph3Errors = [
    { time: '00:00:04:12 - 00:00:08:00', text: 'Selamat datang di masa depan Tokyo.', type: 'SAFE' },
    { time: '00:00:08:05 - 00:00:10:15', text: 'Distrik Shibuya sekarang dikelola sepenuhnya oleh AI pusat.', type: 'CHECK CPS', info: 'CPS: 22 (Batas: 20)', active: true },
    { time: '00:00:10:12 - 00:00:13:00', text: 'Jangan lupa untuk mendaftarkan chip identitas Anda.', type: 'OVERLAP' }
  ];
</script>

<div class="max-w-[1600px] mx-auto flex flex-col pb-12">
  <!-- Top Phase Navigation Switcher (Internal to Editor) -->
  <div class="flex justify-center mb-6">
    <div class="bg-surface-container border border-outline-variant/30 rounded-full p-1 flex gap-1 shadow-lg font-headline">
      <button on:click={() => currentPhase = 1} class="px-6 py-2 rounded-full text-xs font-bold tracking-widest transition-all uppercase {currentPhase === 1 ? 'bg-primary/20 text-primary border border-primary/50 shadow-[0_0_10px_rgba(255,45,120,0.3)]' : 'text-slate-400 hover:text-white'}">
        Phase 1: Upload & ASR
      </button>
      <button on:click={() => currentPhase = 2} class="px-6 py-2 rounded-full text-xs font-bold tracking-widest transition-all uppercase {currentPhase === 2 ? 'bg-secondary/20 text-secondary border border-secondary/50 shadow-[0_0_10px_rgba(0,255,204,0.3)]' : 'text-slate-400 hover:text-white'}">
        Phase 2: L1 & L2 Pipeline
      </button>
      <button on:click={() => currentPhase = 3} class="px-6 py-2 rounded-full text-xs font-bold tracking-widest transition-all uppercase {currentPhase === 3 ? 'bg-tertiary/20 text-tertiary border border-tertiary/50 shadow-[0_0_10px_rgba(255,224,74,0.3)]' : 'text-slate-400 hover:text-white'}">
        Phase 3: Export & QA
      </button>
    </div>
  </div>

  {#if currentPhase === 1}
    <!-- ========================= PHASE 1 ========================= -->
    <header class="flex justify-between items-end mb-6 shrink-0">
      <div>
        <h1 class="text-3xl font-headline font-extrabold text-on-background tracking-tight mb-2">
          Fase 1: ASR & Editor Segmen
        </h1>
        <p class="text-xs font-body text-slate-400">Automated Speech Recognition & Segment Refinement</p>
      </div>
      <div class="flex gap-4">
        <span class="px-4 py-1.5 border border-[#00ffcc] text-[#00ffcc] bg-[rgba(0,255,204,0.1)] rounded-full text-[9px] font-label font-bold tracking-[0.2em] uppercase">STATUS: PROCESSING</span>
        <span class="px-4 py-1.5 border border-tertiary text-tertiary bg-tertiary/10 rounded-full text-[9px] font-label font-bold tracking-[0.2em] uppercase">ENGINE: WHISPER</span>
      </div>
    </header>

    <div class="flex-1 grid grid-cols-1 lg:grid-cols-[1.5fr_1fr] gap-6">
      
      <!-- Left Column (Video & ASR Setup) -->
      <div class="flex flex-col gap-6 pr-2">
        <!-- Video Player Mock -->
        <div class="bg-surface-container rounded-xl overflow-hidden border border-outline-variant/30 flex flex-col relative group">
          <div class="aspect-video bg-black flex items-center justify-center relative overflow-hidden">
            <!-- Simulated Video Content -->
            <div class="absolute inset-0 bg-gradient-to-br from-slate-900 to-[#0a0f18] flex items-center justify-center">
              <div class="w-64 h-64 border-[4px] border-[#00ffcc]/20 rounded-full flex items-center justify-center animate-[spin_10s_linear_infinite]">
                 <div class="w-48 h-48 border-[2px] border-primary/20 rounded-full animate-[spin_5s_linear_infinite_reverse]"></div>
              </div>
            </div>
            <!-- Subtitle Overlay -->
            <div class="absolute bottom-16 bg-black/80 px-6 py-3 rounded border border-white/10 z-10">
              <p class="text-white font-headline text-lg tracking-wide">Welcome to the future of translation.</p>
            </div>
          </div>
          <!-- Player Controls -->
          <div class="h-20 bg-surface-container-low px-6 flex flex-col justify-center border-t border-outline-variant/30 flex-shrink-0 relative">
            <div class="absolute top-[-5px] left-8 right-8 h-10 flex gap-1 items-end z-20">
               <!-- Waveform mock -->
               {#each Array(15) as _, i}
                  <div class="w-full rounded-t-full bg-secondary/50 h-[30%]" style="height: {Math.random() * 40 + 20}%"></div>
               {/each}
               <div class="w-full rounded-t-full bg-primary h-[80%] drop-shadow-[0_0_8px_#ff2d78]"></div>
               {#each Array(10) as _, i}
                  <div class="w-full rounded-t-full bg-secondary/50 h-[30%]" style="height: {Math.random() * 40 + 20}%"></div>
               {/each}
            </div>
            <div class="relative z-30 flex justify-between items-center mt-3">
              <div class="flex items-center gap-4">
                <button class="text-slate-400 hover:text-white"><span class="material-symbols-outlined text-sm">skip_previous</span></button>
                <button class="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center border border-primary text-primary shadow-[0_0_10px_rgba(255,45,120,0.3)]"><span class="material-symbols-outlined text-sm">play_arrow</span></button>
                <button class="text-slate-400 hover:text-white"><span class="material-symbols-outlined text-sm">skip_next</span></button>
                <span class="text-[10px] font-mono text-[#00ffcc] ml-4 font-bold">00:12:45 / 00:45:00</span>
              </div>
              <div class="flex items-center gap-4">
                <button class="text-slate-400 hover:text-white"><span class="material-symbols-outlined text-sm">volume_up</span></button>
                <button class="text-slate-400 hover:text-white"><span class="material-symbols-outlined text-sm">fullscreen</span></button>
              </div>
            </div>
          </div>
        </div>

        <!-- ASR Engine Setup Card -->
        <div class="bg-surface-container-lowest border border-primary/20 rounded-xl p-6 shadow-xl relative overflow-hidden">
          <div class="absolute left-0 top-0 bottom-0 w-1 bg-gradient-to-b from-primary via-[#ff2d78]/50 to-transparent"></div>
          <h3 class="flex items-center gap-3 text-lg font-headline font-bold text-white mb-6">
            <span class="material-symbols-outlined text-primary">settings_voice</span> ASR Engine Setup
          </h3>
          <div class="flex gap-4 mb-6">
            <div class="flex-1">
              <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">TRANSCRIPTION MODEL</label>
              <select class="w-full bg-surface-container border border-outline-variant rounded p-3 text-sm font-label text-slate-200 outline-none focus:border-primary">
                <option>whisper-large-v3</option>
              </select>
            </div>
            <div class="flex-1">
              <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">TARGET LANGUAGE</label>
              <select class="w-full bg-surface-container border border-outline-variant rounded p-3 text-sm font-label text-slate-200 outline-none focus:border-primary">
                <option>English (Source)</option>
              </select>
            </div>
          </div>
          <button class="w-full py-4 bg-primary hover:bg-[#ff4088] text-white font-headline font-bold uppercase tracking-widest text-sm rounded shadow-[0_0_15px_rgba(255,45,120,0.5)] flex justify-center items-center gap-2 transition-transform active:scale-[0.98]">
            <span class="material-symbols-outlined">bolt</span> Jalankan Transkripsi
          </button>
        </div>
      </div>

      <!-- Right Column (Editor Segmen) -->
      <div class="bg-surface-container-low border border-outline-variant border-r-0 border-y-0 lg:border-l lg:border-y lg:border-r border-outline-variant/30 rounded-xl flex flex-col flex-grow">
        <!-- Panel Header -->
        <div class="p-5 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container relative">
          <div class="absolute inset-y-0 left-0 w-1 bg-secondary"></div>
          <h2 class="flex items-center gap-3 text-base font-headline font-bold text-white pl-3">
            <span class="material-symbols-outlined text-secondary">segment</span> Editor Segmen <span class="text-slate-400 font-normal text-sm block ml-1">(Source EN)</span>
          </h2>
          <div class="flex bg-surface-container-highest border border-outline-variant/50 rounded overflow-hidden">
            <button class="px-3 py-1.5 text-[9px] font-label font-bold uppercase tracking-widest text-slate-300 hover:bg-slate-700 transition">MERGE</button>
            <div class="w-[1px] bg-outline-variant/50"></div>
            <button class="px-3 py-1.5 text-[9px] font-label font-bold uppercase tracking-widest text-slate-300 hover:bg-slate-700 transition">SPLIT</button>
          </div>
        </div>

        <!-- Segments List -->
        <div class="p-5 space-y-4">
          {#each ph1Segments as seg, i}
          <div class="border rounded-lg p-5 transition-all {seg.active ? 'border-primary bg-primary/5 shadow-[0_0_15px_rgba(255,45,120,0.15)] relative' : 'border-outline-variant bg-surface-container hover:border-slate-500'} cursor-pointer">
            {#if seg.active}
               <div class="absolute left-[-2px] inset-y-4 w-1 bg-primary rounded-full shadow-[0_0_8px_#ff2d78]"></div>
            {/if}
            <div class="flex justify-between items-center mb-3">
              <div class="flex items-center gap-3">
                <span class="text-[10px] font-label font-bold text-slate-300">{seg.id}</span>
                <span class="text-[10px] font-mono text-slate-500">{seg.time}</span>
              </div>
              <div class="flex items-center gap-1">
                <span class="text-[10px] font-label font-bold {seg.warn ? 'text-tertiary' : 'text-[#00ffcc]'}">{seg.acc}%</span>
                <span class="material-symbols-outlined text-[12px] {seg.warn ? 'text-tertiary' : 'text-[#00ffcc]'}">{seg.warn ? 'warning' : 'check_circle'}</span>
              </div>
            </div>
            <p class="text-sm font-body leading-relaxed text-slate-200 selection:bg-primary/50">{seg.text}</p>
          </div>
          {/each}
        </div>

        <!-- Panel Footer Controls -->
        <div class="p-5 mt-auto border-t border-outline-variant/30 bg-surface-container/80 flex justify-between">
          <button class="px-5 py-2 border border-outline-variant text-[10px] uppercase font-bold tracking-widest rounded text-slate-400 hover:bg-surface-variant hover:text-white">PREVIOUS PAGE</button>
          <button on:click={() => currentPhase = 2} class="px-10 py-2 bg-[#00ffcc] hover:bg-[#33ffe6] text-[#001a1a] text-[10px] uppercase font-bold tracking-widest rounded drop-shadow-[0_0_10px_rgba(0,255,204,0.4)]">NEXT PAGE</button>
        </div>
      </div>
    </div>
  {/if}


  {#if currentPhase === 2}
    <!-- ========================= PHASE 2 ========================= -->
    <header class="flex justify-between items-end mb-8 shrink-0">
      <div>
        <h1 class="text-3xl font-headline font-extrabold text-on-background tracking-tight mb-2">
          Fase 2: <span class="text-secondary drop-shadow-[0_0_8px_rgba(0,255,204,0.6)]">L1 & L2 Pipeline</span>
        </h1>
        <p class="text-xs font-body text-slate-400">Deep Neural Translation & Stylistic Refinement Engine</p>
      </div>
      <div class="flex gap-4">
        <span class="px-5 py-2 border border-secondary/30 text-secondary bg-secondary/10 rounded-full text-[9px] font-label font-bold tracking-[0.2em] uppercase neon-glow-secondary">ACTIVE SESSION: ALPHA-01</span>
      </div>
    </header>

    <div class="flex flex-col gap-6">
      
      <!-- Top Form: Context Engine -->
      <div class="bg-surface-container-lowest border border-primary/20 border-l-2 border-l-primary rounded-xl p-6 shadow-xl shrink-0">
        <h3 class="flex items-center gap-3 text-lg font-headline font-bold text-white mb-6">
          <span class="material-symbols-outlined text-primary drop-shadow-[0_0_8px_rgba(255,45,120,0.5)]">settings_suggest</span> Context Engine Configuration
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
          <div>
            <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">KATEGORI KONTEN (L1)</label>
            <select class="w-full bg-background border border-outline-variant/50 rounded-md p-3 text-sm font-label text-slate-200 outline-none focus:border-primary">
              <option>Sci-Fi Narrative</option>
            </select>
          </div>
          <div>
            <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">GAYA BAHASA (TONE L2)</label>
            <select class="w-full bg-background border border-outline-variant/50 rounded-md p-3 text-sm font-label text-slate-200 outline-none focus:border-primary">
              <option>Noir & Gritty</option>
            </select>
          </div>
          <div>
            <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">INJEKSI KAMUS (GLOSSARY)</label>
            <select class="w-full bg-background border border-outline-variant/50 rounded-md p-3 text-sm font-label text-slate-200 outline-none focus:border-primary">
              <option>Cyberpunk-Slang-v2.db</option>
            </select>
          </div>
        </div>
        <div class="flex gap-6 items-end">
          <div class="flex-1">
            <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">SCENE CONTEXT</label>
            <textarea rows="2" class="w-full bg-background border border-outline-variant/50 rounded-md p-3 text-sm font-body text-slate-300 outline-none focus:border-primary placeholder:text-slate-600 resize-none font-mono text-[11px]" placeholder="Describe the scene context here... (e.g., A rainy night in Shinjuku, character is feeling tense)"></textarea>
          </div>
          <button class="px-8 py-3.5 h-fit border border-primary text-primary hover:bg-primary/10 font-headline font-bold uppercase tracking-widest text-xs rounded shadow-[0_0_15px_rgba(255,45,120,0.3)] flex justify-center items-center gap-2 transition-transform">
            <span class="material-symbols-outlined text-[20px]">psychiatry</span> JALANKAN PIPELINE AI
          </button>
        </div>
      </div>

      <!-- Bottom Table: Pipeline Output -->
      <div class="bg-surface-container-low border border-secondary/20 rounded-xl overflow-hidden shadow-2xl flex flex-col">
        <div class="p-5 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container/50">
          <h2 class="text-xs font-label font-bold text-secondary uppercase tracking-widest drop-shadow-[0_0_8px_rgba(0,255,204,0.3)]">Translation Pipeline Output</h2>
          <div class="flex items-center gap-4 text-[9px] font-label text-slate-400">
             <div class="flex items-center gap-1.5"><div class="w-2 h-2 rounded-full bg-error drop-shadow-[0_0_5px_#ff5252]"></div> L1 Deletions</div>
             <div class="flex items-center gap-1.5"><div class="w-2 h-2 rounded-full bg-secondary drop-shadow-[0_0_5px_#00ffcc]"></div> L2 Refinements</div>
          </div>
        </div>
        
        <div class="w-full max-w-full">
          <table class="w-full text-left border-collapse table-fixed">
            <thead class="sticky top-0 bg-surface-container/95 backdrop-blur z-20">
              <tr class="border-b border-outline-variant/30 text-[9px] font-label font-bold text-slate-500 tracking-[0.2em] uppercase">
                <th class="py-4 px-6 w-20">ID</th>
                <th class="py-4 px-6 w-1/4">TEKS SUMBER (EN)</th>
                <th class="py-4 px-6 w-1/4">L1: LITERAL TRANSLATION</th>
                <th class="py-4 px-6 w-1/3">L2: POLISHED REWRITE</th>
                <th class="py-4 px-6 text-right w-32">STATUS</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-outline-variant/20">
              {#each ph2Rows as row}
                <tr class="hover:bg-surface-container-high/30 transition-colors group">
                  <td class="py-6 px-6 align-top">
                    <span class="text-[10px] font-mono text-slate-500">{row.id}</span>
                  </td>
                  <td class="py-6 px-6 align-top">
                    <p class="text-xs font-body leading-relaxed text-slate-300 pr-4">{row.src}</p>
                  </td>
                  <td class="py-6 px-6 align-top border-l border-outline-variant/10">
                    <p class="text-xs font-body italic leading-relaxed text-slate-400 pr-4">{row.l1}</p>
                  </td>
                  <td class="py-6 px-6 align-top border-l border-outline-variant/10">
                    <p class="text-xs font-body leading-relaxed text-slate-100">{@html row.polishHtml}</p>
                  </td>
                  <td class="py-6 px-6 align-top text-right">
                     <span class="inline-flex py-1 px-3 mt-1 rounded-sm text-[8px] font-bold uppercase tracking-widest border border-secondary text-secondary bg-secondary/10 shadow-[0_0_8px_rgba(0,255,204,0.2)]">
                       REFINED
                     </span>
                     <div class="flex flex-col mt-4 opacity-0 group-hover:opacity-100 transition-opacity gap-2 items-end">
                       <button class="w-8 h-8 rounded border border-[#00ffcc] text-[#00ffcc] flex items-center justify-center hover:bg-[#00ffcc] hover:text-[#001a1a] shadow-[0_0_8px_rgba(0,255,204,0.3)] transition" title="Save Segmen"><span class="material-symbols-outlined text-[16px]">save</span></button>
                       <button class="w-8 h-8 rounded bg-primary text-white flex items-center justify-center shadow-[0_0_8px_#ff2d78] hover:bg-[#ff4088] transition" title="Koreksi Manual"><span class="material-symbols-outlined text-[16px]">publish</span></button>
                     </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        
        <!-- Table Footer -->
        <div class="p-4 border-t border-outline-variant/30 flex justify-between items-center bg-surface-container">
          <p class="text-[9px] font-mono text-slate-500 font-bold tracking-widest uppercase">SHOWING 3 OF 152 NODES</p>
          <div class="flex items-center gap-1.5">
            <button class="w-8 h-8 flex items-center justify-center rounded border border-outline-variant/50 hover:bg-surface-variant transition text-slate-400"><span class="material-symbols-outlined text-[14px]">chevron_left</span></button>
            <button class="w-8 h-8 flex items-center justify-center rounded border border-primary text-primary bg-primary/10 font-bold drop-shadow-[0_0_5px_currentColor]">1</button>
            <button class="w-8 h-8 flex items-center justify-center rounded border border-outline-variant/50 hover:bg-surface-variant transition text-slate-400">2</button>
            <button class="w-8 h-8 flex items-center justify-center rounded border border-outline-variant/50 hover:bg-surface-variant transition text-slate-400">3</button>
            <button class="w-8 h-8 flex items-center justify-center rounded border border-outline-variant/50 hover:bg-surface-variant transition text-slate-400"><span class="material-symbols-outlined text-[14px]">chevron_right</span></button>
          </div>
        </div>
      </div>
    </div>
  {/if}

  {#if currentPhase === 3}
    <!-- ========================= PHASE 3 ========================= -->
    <header class="flex justify-between items-end mb-6 shrink-0">
      <div>
        <h1 class="text-3xl font-headline font-extrabold text-on-background tracking-tight mb-2">
          QA & Final Polish
        </h1>
        <p class="text-xs font-body text-slate-400">Tinjauan akhir sebelum ekspor final.</p>
      </div>
      <div class="flex gap-4">
        <button class="px-6 py-2 border border-secondary text-secondary hover:bg-secondary/10 rounded font-headline font-bold text-xs tracking-widest uppercase transition-all shadow-[0_0_12px_rgba(0,255,204,0.2)]">PREVIEW</button>
        <button class="px-8 py-2 border border-primary text-white bg-primary hover:bg-[#ff4088] rounded font-headline font-bold text-xs tracking-widest uppercase transition-all shadow-[0_0_15px_rgba(255,45,120,0.5)]">EKSPOR FINAL</button>
      </div>
    </header>

    <div class="grid grid-cols-1 lg:grid-cols-[1.2fr_1fr] gap-6">
      
      <!-- Left Column (Visuals & Audio) -->
      <div class="flex flex-col gap-6 overflow-hidden">
        <!-- Video Container -->
        <div class="bg-surface-container rounded-xl overflow-hidden border border-outline-variant/30 flex flex-col relative aspect-[21/9]">
           <div class="absolute top-0 left-0 right-0 p-3 bg-gradient-to-b from-black/80 to-transparent flex justify-between z-20">
             <span class="text-[9px] font-mono text-slate-300 uppercase tracking-widest">VIDEO PREVIEW — 4K PRORES</span>
             <div class="flex gap-1.5"><div class="w-2 h-2 rounded-full bg-secondary shadow-[0_0_5px_#00ffcc]"></div><div class="w-2 h-2 rounded-full bg-secondary/30"></div></div>
           </div>
           
           <!-- Simulated Cyberpunk Video Frame -->
           <div class="absolute inset-0 bg-[#0a0f18] flex items-center justify-center overflow-hidden">
              <div class="w-full h-full relative" style="background: repeating-linear-gradient(90deg, rgba(255,255,255,0.03) 0px, rgba(255,255,255,0.03) 1px, transparent 1px, transparent 40px), repeating-linear-gradient(0deg, rgba(255,255,255,0.03) 0px, rgba(255,255,255,0.03) 1px, transparent 1px, transparent 40px);"></div>
              <!-- Mock Buildings -->
              <div class="absolute left-10 bottom-0 w-32 h-64 bg-slate-900 border-t border-r border-[#00ffcc]/30 shadow-[-20px_0_40px_rgba(0,255,204,0.1)] flex flex-col items-center py-4">
                 <div class="w-2 h-32 bg-primary/20 rounded shadow-[0_0_10px_#ff2d78]"></div>
              </div>
              <div class="absolute right-20 bottom-0 w-40 h-80 bg-slate-900 border-t border-l border-primary/30 flex flex-col justify-between py-6 items-center">
                 <span class="text-primary font-bold tracking-[1em] text-vertical neon-glow-primary text-xl rotate-90 scale-x-150">C R E O</span>
              </div>
           </div>
        </div>

        <!-- Audio Timeline Box -->
        <div class="bg-surface-container-low border border-outline-variant/30 rounded-xl flex flex-col p-5 overflow-hidden h-64">
           <div class="flex justify-between items-center mb-4">
             <div class="flex items-center gap-3">
                <span class="text-[10px] font-label font-bold text-secondary uppercase tracking-widest drop-shadow-[0_0_5px_rgba(0,255,204,0.3)]">TIMELINE AUDIO</span>
                <span class="material-symbols-outlined text-[14px] text-slate-500 cursor-pointer hover:text-white">zoom_in</span>
                <span class="material-symbols-outlined text-[14px] text-slate-500 cursor-pointer hover:text-white">zoom_out</span>
             </div>
             <div class="flex items-center gap-4 text-[9px] font-label text-slate-400">
               <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 bg-primary/60 border border-primary"></div> Vocal L</div>
               <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 bg-secondary/60 border border-secondary"></div> Vocal R</div>
             </div>
           </div>
           
           <!-- Waveform Renderer Mock -->
           <div class="flex-1 bg-surface-container-highest/20 rounded border border-outline-variant/30 relative flex items-center px-4 gap-1 overflow-hidden">
               <!-- Playhead -->
               <div class="absolute left-[30%] top-0 bottom-0 w-px bg-primary shadow-[0_0_8px_#ff2d78] z-30">
                  <div class="absolute top-0 left-[-4px] w-0 h-0 border-l-[4px] border-l-transparent border-r-[4px] border-r-transparent border-t-[6px] border-t-primary drop-shadow-[0_0_5px_#ff2d78]"></div>
               </div>
               
               <!-- Stereo Bars -->
               {#each Array(80) as _, i}
                  <div class="w-2 h-full flex flex-col justify-center gap-0.5 opacity-80" style="transform: scaleY({Math.max(0.1, Math.sin(i*0.2) * 0.8 + 0.2)})">
                     <div class="w-full bg-primary/70 {i>20 && i<30 ? 'bg-primary shadow-[0_0_8px_#ff2d78]' : ''} rounded-sm {i%3===0 ? 'h-32' : i%2===0 ? 'h-20' : 'h-10'} duration-300"></div>
                     <div class="w-full bg-secondary/70 {i>20 && i<30 ? 'bg-secondary drop-shadow-[0_0_5px_#00ffcc]' : ''} rounded-sm {i%3===0 ? 'h-24' : i%2===0 ? 'h-16' : 'h-8'} duration-300"></div>
                  </div>
               {/each}
           </div>
        </div>
      </div>

      <!-- Right Column (Error List & Panel) -->
      <div class="flex flex-col gap-6 overflow-hidden">
        
        <!-- Segment List -->
        <div class="bg-surface-container-low border border-outline-variant/30 rounded-xl flex flex-col">
          <div class="p-5 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container relative">
            <h2 class="text-base font-headline font-bold text-white">Daftar Segmen</h2>
            <div class="flex gap-2">
              <span class="px-2 py-1 bg-error/10 border border-error text-error text-[8px] font-bold uppercase rounded-sm">3 ERROR</span>
              <span class="px-2 py-1 bg-secondary/10 border border-secondary text-secondary text-[8px] font-bold uppercase rounded-sm glow-text-secondary">12 AMAN</span>
            </div>
          </div>
          
          <div class="divide-y divide-outline-variant/20 custom-scrollbar">
            {#each ph3Errors as err}
              <div class="p-4 {err.active ? 'bg-surface-variant/30 border-l-[3px] border-l-tertiary' : 'hover:bg-surface-container transition'}">
                <div class="flex justify-between items-center mb-2">
                   <span class="text-[9px] font-mono text-slate-400">{err.time}</span>
                   <div class="flex items-center gap-1.5">
                      {#if err.type === 'SAFE'}
                        <span class="material-symbols-outlined text-[12px] text-secondary">check_circle</span>
                        <span class="text-[9px] font-bold text-secondary">SAFE</span>
                      {:else if err.type === 'CHECK CPS'}
                        <span class="material-symbols-outlined text-[12px] text-tertiary">warning</span>
                        <span class="text-[9px] font-bold text-tertiary">CHECK CPS</span>
                      {:else if err.type === 'OVERLAP'}
                        <span class="material-symbols-outlined text-[12px] text-error">error</span>
                        <span class="text-[9px] font-bold text-error">OVERLAP</span>
                      {/if}
                   </div>
                </div>
                <p class="text-xs font-body text-slate-200">{err.text}</p>
                {#if err.info}
                   <p class="text-[9px] font-mono text-tertiary mt-2 font-bold italic">{err.info}</p>
                {/if}
              </div>
            {/each}
          </div>
        </div>

        <!-- QA Editor Panel -->
        <div class="bg-surface-container border border-outline-variant/30 border-t-primary/50 relative rounded-xl flex flex-col overflow-hidden shadow-2xl">
          <div class="p-5 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container-highest/20">
            <h2 class="flex items-center gap-3 text-base font-headline font-bold text-white">
              <span class="material-symbols-outlined text-primary text-xl pb-1">auto_fix_high</span> QA & Editor Panel
            </h2>
            <button class="px-4 py-1.5 border border-primary/50 text-primary hover:bg-primary/10 text-[9px] font-bold tracking-widest uppercase rounded shadow-[0_0_8px_rgba(255,45,120,0.2)]">AUTO-FIX ALL</button>
          </div>
          
          <div class="p-5 flex flex-col gap-4">
             <div>
               <label class="block text-[8px] font-label uppercase text-slate-400 mb-2">EDIT TEKS</label>
               <textarea rows="3" class="w-full bg-background border border-outline-variant/50 rounded-md p-3 text-sm font-body text-slate-200 outline-none focus:border-tertiary resize-none">Distrik Shibuya sekarang dikelola sepenuhnya oleh AI pusat.</textarea>
             </div>
             
             <div class="flex gap-4">
               <div class="flex-1">
                 <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">IN POINT</label>
                 <input type="text" value="00:00:08:05" class="w-full bg-background border border-outline-variant/50 rounded p-2 text-xs font-mono text-tertiary outline-none focus:border-tertiary" />
               </div>
               <div class="flex-1">
                 <label class="block text-[8px] font-label uppercase tracking-widest text-slate-400 mb-2">OUT POINT</label>
                 <input type="text" value="00:00:10:15" class="w-full bg-background border border-outline-variant/50 rounded p-2 text-xs font-mono text-tertiary outline-none focus:border-tertiary" />
               </div>
             </div>

             <!-- AI Suggestion Box -->
             <div class="mt-2 bg-error/5 border border-error/20 rounded-md p-4 relative overflow-hidden">
               <h4 class="text-[9px] font-bold text-error uppercase tracking-widest flex items-center gap-2 mb-2">
                 <span class="material-symbols-outlined text-[14px]">psychology</span> SARAN AI
               </h4>
               <p class="text-xs font-body text-slate-300 italic mb-4 pr-2">
                 "Teks ini terlalu panjang untuk durasi 2 detik. Kurangi kata atau tambahkan durasi 0.5 detik untuk mencapai target 18 CPS."
               </p>
               <div class="flex gap-3">
                 <button class="flex-1 py-2 bg-primary text-white text-[9px] font-bold tracking-widest uppercase rounded shadow-[0_0_8px_#ff2d78] hover:bg-[#ff4088] transition">Terapkan Singkatan</button>
                 <button class="flex-[0.5] py-2 border border-outline-variant/50 text-slate-300 hover:text-white hover:bg-surface-variant text-[9px] font-bold tracking-widest uppercase rounded transition">Ignore</button>
               </div>
             </div>
          </div>
        </div>

      </div>

    </div>
  {/if}

</div>

<style>
  /* Local Scoped Customizations for the complex Layout */
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background-color: rgba(255,45,120,0.2);
    border-radius: 20px;
  }
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>
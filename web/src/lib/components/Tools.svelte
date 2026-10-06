<script lang="ts">
  import { LEVELS, type Tool } from '../types';
  let { tools, skills }: { tools: Tool[]; skills: string[] } = $props();
</script>

<section id="tools" class="sec">
  <div class="head"><div class="eyebrow">Лист 02</div><h2 class="h2">Инструменты и навыки</h2></div>
  <div>
    <div class="grid">
      {#each tools as t, i}
        <div class="tool">
          <div class="meta"><span>{String(i + 1).padStart(2, '0')}</span><span class="lvl">{LEVELS[t.level]}</span></div>
          <div class="name">{t.name}</div>
          <div class="segs" role="img" aria-label={`${t.level} из 5`}>
            {#each [1, 2, 3, 4, 5] as k}<span class:on={k <= t.level}></span>{/each}
          </div>
        </div>
      {/each}
    </div>
    <div class="chips">{#each skills as s}<span>{s}</span>{/each}</div>
  </div>
</section>

<style>
  .sec { display: grid; grid-template-columns: clamp(280px, 20%, 420px) minmax(0, 1fr); border-bottom: 1px solid var(--ink); scroll-margin-top: 53px; }
  .head { padding: 40px 32px; border-right: 1px solid var(--ink); display: flex; flex-direction: column; gap: 12px; }
  .h2 { font-size: 32px; }
  .grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); }
  .tool { padding: 24px; border-right: 1px dashed var(--dash); border-bottom: 1px dashed var(--dash); display: flex; flex-direction: column; gap: 14px; }
  .meta { display: flex; justify-content: space-between; font-size: 12px; opacity: 0.85; }
  .lvl { text-transform: uppercase; }
  .name { font-family: var(--sans); font-size: 22px; font-weight: 600; }
  .segs { display: flex; gap: 4px; }
  .segs span { flex: 1; height: 6px; border: 1px solid var(--ink); }
  .segs span.on { height: 8px; border: 0; background: var(--accent); }
  .chips { padding: 24px; display: flex; flex-wrap: wrap; gap: 8px; }
  .chips span { font-size: 13px; border: 1px solid var(--ink); padding: 8px 12px; }

  @media (max-width: 900px) {
    .sec { display: block; }
    .head { padding: 28px 20px 20px; gap: 8px; border-right: 0; border-bottom: 1px solid var(--ink); }
    .h2 { font-size: 30px; }
    .grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
    .tool { padding: 16px; gap: 10px; }
    .name { font-size: 17px; }
    .segs { gap: 3px; }
    .segs span { height: 4px; }
    .segs span.on { height: 6px; }
    .chips { padding: 20px; gap: 6px; }
    .chips span { font-size: 12px; padding: 7px 10px; }
  }
</style>

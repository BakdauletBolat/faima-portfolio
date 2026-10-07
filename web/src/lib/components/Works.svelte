<script lang="ts">
  import { CATS, type Work } from '../types';
  let { works }: { works: Work[] } = $props();
  let filter = $state<string>('all');
  const shown = $derived(filter === 'all' ? works : works.filter((w) => w.cat === filter));
</script>

<section id="works">
  <div class="bar">
    <div class="title"><span class="eyebrow">Лист 03</span><h2 class="h2">Работы</h2></div>
    <div class="filters" role="group" aria-label="Фильтр работ">
      {#each CATS as c}
        <button class:on={filter === c.key} aria-pressed={filter === c.key} onclick={() => (filter = c.key)}>{c.label}</button>
      {/each}
    </div>
  </div>
  <div class="grid">
    {#each shown as w, i (w.id)}
      <article>
        <div class="cover">{#if w.cover}<img src={w.cover} alt={w.title} loading="lazy" />{/if}</div>
        <div class="body">
          <div class="meta">
            <span>АИ-{String(i + 1).padStart(2, '0')} · {w.type}</span>
            <span>{w.area} · {w.year}</span>
          </div>
          <h3>{w.title}</h3>
          <p>{w.scope}</p>
        </div>
        {#if w.pdf}
          <a class="dl" href={w.pdf} download={`${w.title}.pdf`}><span>↓ Скачать .pdf</span><span>{w.size}</span></a>
        {/if}
      </article>
    {/each}
  </div>
</section>

<style>
  #works { border-bottom: 1px solid var(--ink); scroll-margin-top: 53px; }
  .bar { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--ink); }
  .title { padding: 32px; display: flex; align-items: baseline; gap: 24px; }
  .filters { display: flex; align-self: stretch; }
  .filters button {
    border: 0; border-left: 1px solid var(--ink); cursor: pointer; font-size: 13px; text-transform: uppercase;
    letter-spacing: 0.06em; padding: 0 24px; background: transparent; color: var(--ink); font-family: var(--mono);
  }
  .filters button:hover { background: var(--hover); }
  .filters button.on { background: var(--ink); color: var(--paper); font-weight: 700; }
  .grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); }
  article { border-right: 1px solid var(--ink); border-bottom: 1px solid var(--ink); display: flex; flex-direction: column; }
  article:nth-child(3n) { border-right: 0; }
  .cover { border-bottom: 1px solid var(--ink); background: var(--hover); }
  .cover img { display: block; width: 100%; height: auto; object-fit: contain; }
  .body { padding: 20px 24px; display: flex; flex-direction: column; gap: 10px; flex: 1; }
  .meta { display: flex; justify-content: space-between; gap: 12px; font-size: 12px; letter-spacing: 0.06em; text-transform: uppercase; }
  h3 { margin: 0; font-family: var(--sans); font-size: 24px; font-weight: 700; letter-spacing: -0.02em; }
  p { margin: 0; font-size: 14px; line-height: 1.55; flex: 1; white-space: pre-wrap; }
  .dl {
    display: flex; justify-content: space-between; padding: 14px 24px; border-top: 1px solid var(--ink); color: var(--ink);
    font-size: 13px; text-transform: uppercase; letter-spacing: 0.06em; font-weight: 700;
  }
  .dl:hover { background: var(--ink); color: var(--paper); }

  @media (max-width: 900px) {
    #works { scroll-margin-top: 53px; }
    .bar { display: block; }
    .title { padding: 28px 20px 20px; flex-direction: column; gap: 8px; border-bottom: 1px solid var(--ink); }
    .filters { overflow-x: auto; }
    .filters button { flex-shrink: 0; border-left: 0; border-right: 1px solid var(--ink); font-size: 12px; padding: 0 16px; min-height: 48px; }
    .grid { display: block; }
    article { border-right: 0; }
    article:last-child { border-bottom: 0; }
    .body { padding: 16px 20px; gap: 8px; }
    .meta { letter-spacing: 0.04em; }
    h3 { font-size: 21px; }
    .dl { align-items: center; padding: 0 20px; min-height: 52px; }
  }
  @media (min-width: 1800px) {
    .grid { grid-template-columns: repeat(4, minmax(0, 1fr)); }
    article:nth-child(3n) { border-right: 1px solid var(--ink); }
    article:nth-child(4n) { border-right: 0; }
  }
  @media (min-width: 901px) and (max-width: 1200px) {
    .grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
    article:nth-child(3n) { border-right: 1px solid var(--ink); }
    article:nth-child(2n) { border-right: 0; }
  }
</style>

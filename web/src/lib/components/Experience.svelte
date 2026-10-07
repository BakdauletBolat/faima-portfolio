<script lang="ts">
  import type { Experience } from '../types';
  let { items }: { items: Experience[] } = $props();
</script>

<section id="exp" class="sec">
  <div class="head"><div class="eyebrow">Лист 01</div><h2 class="h2">Опыт работы</h2></div>
  <div class="rows">
    {#each items as e}
      <div class="row">
        <div class="period">{e.period}</div>
        <div class="who"><div class="role">{e.role}</div><div class="company">{e.company}</div></div>
        <div class="desc">{e.desc}</div>
      </div>
    {/each}
  </div>
</section>

<style>
  .sec { display: grid; grid-template-columns: clamp(280px, 20%, 420px) minmax(0, 1fr); border-bottom: 1px solid var(--ink); scroll-margin-top: 53px; }
  .head { padding: 40px 32px; border-right: 1px solid var(--ink); display: flex; flex-direction: column; gap: 12px; }
  .row { display: grid; grid-template-columns: 200px minmax(0, 1fr) minmax(0, 1.2fr); border-bottom: 1px dashed var(--dash); }
  .row:last-child { border-bottom: 0; }
  .row > div { padding: 28px 24px; }
  .period { border-right: 1px dashed var(--dash); font-size: 14px; }
  .who { display: flex; flex-direction: column; gap: 6px; }
  .role { font-family: var(--sans); font-size: 22px; font-weight: 700; letter-spacing: -0.02em; }
  .company { font-size: 14px; }
  .desc { font-size: 15px; line-height: 1.6; white-space: pre-wrap; }

  @media (max-width: 900px) {
    .sec { display: block; }
    .head { padding: 28px 20px 20px; gap: 8px; border-right: 0; border-bottom: 1px solid var(--ink); }
    .row { display: flex; flex-direction: column; gap: 8px; padding: 20px; }
    .row > div { padding: 0; }
    .period { border-right: 0; font-size: 13px; }
    .role { font-size: 19px; line-height: 1.2; }
    .company { font-size: 13px; }
    .desc { font-size: 14px; margin-top: 4px; }
  }
</style>

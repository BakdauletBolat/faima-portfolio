<script lang="ts">
  let { resume }: { resume: string } = $props();
  let open = $state(false);
  const links = [
    ['#exp', 'Л.01 Опыт'], ['#tools', 'Л.02 Инструменты'],
    ['#works', 'Л.03 Работы'], ['#contact', 'Л.04 Контакты'],
  ];
</script>

<header>
  <div class="logo">ФЖ — Портфолио</div>
  <nav class="desk" aria-label="Разделы">
    {#each links as [href, label]}<a {href}>{label}</a>{/each}
  </nav>
  {#if resume}<a class="resume" href={resume} download="Резюме.pdf">↓ Резюме.pdf</a>{/if}
  <button class="burger" aria-expanded={open} onclick={() => (open = !open)}>{open ? 'Закрыть' : 'Меню'}</button>
</header>
{#if open}
  <nav class="mob" aria-label="Разделы">
    {#each links as [href, label]}<a {href} onclick={() => (open = false)}>{label}</a>{/each}
  </nav>
{/if}
{#if resume}
  <a class="resume-bar" href={resume} download="Резюме.pdf"><span>↓ Скачать резюме</span><span>.pdf</span></a>
{/if}

<style>
  header {
    display: grid; grid-template-columns: auto 1fr auto; border-bottom: 1px solid var(--ink);
    font-size: 13px; letter-spacing: 0.06em; text-transform: uppercase;
  }
  .logo { padding: 16px 24px; border-right: 1px solid var(--ink); font-weight: 700; }
  .desk { display: flex; }
  .desk a { padding: 16px 24px; border-right: 1px solid var(--ink); color: var(--ink); }
  .desk a:hover { background: var(--ink); color: var(--paper); }
  .resume { padding: 16px 24px; background: var(--ink); color: var(--paper); font-weight: 700; }
  .resume:hover { background: var(--accent); color: var(--paper); }
  .burger, .mob, .resume-bar { display: none; }

  @media (max-width: 900px) {
    header { display: flex; justify-content: space-between; align-items: stretch; position: sticky; top: 0; z-index: 10; background: var(--paper); }
    .logo { padding: 0 20px; display: flex; align-items: center; border-right: 0; }
    .desk, .resume { display: none; }
    .burger {
      display: block; border: 0; border-left: 1px solid var(--ink); background: transparent; color: var(--ink);
      font-size: 13px; letter-spacing: 0.06em; text-transform: uppercase; font-weight: 700; padding: 0 20px; min-height: 52px; cursor: pointer;
    }
    .mob { display: flex; flex-direction: column; border-bottom: 1px solid var(--ink); background: var(--paper); font-size: 15px; text-transform: uppercase; letter-spacing: 0.06em; position: sticky; top: 53px; z-index: 10; }
    .mob a { padding: 16px 20px; border-bottom: 1px dashed var(--dash); color: var(--ink); }
    .mob a:last-child { border-bottom: 0; }
    .resume-bar {
      display: flex; justify-content: space-between; align-items: center; padding: 0 20px; min-height: 56px;
      background: var(--ink); color: var(--paper); font-size: 14px; text-transform: uppercase; letter-spacing: 0.06em; font-weight: 700;
    }
    .resume-bar span:last-child { font-weight: 400; }
  }
</style>

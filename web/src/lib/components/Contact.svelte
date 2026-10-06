<script lang="ts">
  import type { Contacts } from '../types';
  let { contacts, name }: { contacts: Contacts; name: string } = $props();
  const parts = $derived(name.split(' '));
  const short = $derived(parts.length > 1 ? `${parts[1]} ${parts[0][0]}.` : name);
  const tel = $derived(contacts.phone.replace(/[^\d+]/g, ''));
</script>

<section id="contact" class="sec">
  <div class="main">
    <div class="eyebrow">Лист 04 · Контакты</div>
    <a class="mail" href={`mailto:${contacts.email}`}>{contacts.email}</a>
    <div class="links">
      {#if contacts.telegram}<a href={contacts.telegram} target="_blank" rel="noopener">Telegram ↗</a>{/if}
      {#if contacts.phone}<a href={`tel:${tel}`}>{contacts.phone}</a>{/if}
    </div>
  </div>
  <div class="stamp">
    <div><div class="l">Разработала</div><div class="v">{short}</div></div>
    <div><div class="l">Стадия</div><div class="v">Портфолио</div></div>
    <div><div class="l">Лист</div><div class="v">04 / 04</div></div>
    <div><div class="l">Год</div><div class="v">{new Date().getFullYear()}</div></div>
  </div>
</section>

<style>
  .sec { display: grid; grid-template-columns: minmax(0, 1fr) 560px; scroll-margin-top: 53px; }
  .main { padding: 48px; border-right: 1px solid var(--ink); display: flex; flex-direction: column; gap: 24px; }
  .mail { font-family: var(--sans); font-size: 64px; font-weight: 700; letter-spacing: -0.04em; color: var(--ink); overflow-wrap: anywhere; }
  .mail:hover { text-decoration: underline; }
  .links { display: flex; gap: 32px; font-size: 15px; flex-wrap: wrap; }
  .links a { color: var(--ink); }
  .stamp { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); font-size: 12px; text-transform: uppercase; letter-spacing: 0.06em; align-content: start; }
  .stamp > div { padding: 14px 16px; border-bottom: 1px solid var(--ink); }
  .stamp > div:nth-child(odd) { border-right: 1px solid var(--ink); }
  .l { opacity: 0.8; }
  .v { font-size: 15px; margin-top: 4px; }

  @media (max-width: 900px) {
    .sec { display: block; }
    .main { padding: 28px 20px; gap: 16px; border-right: 0; border-bottom: 1px solid var(--ink); }
    .mail { font-size: 30px; letter-spacing: -0.03em; }
    .links { flex-direction: column; gap: 0; margin: 0 -20px -28px; font-size: 14px; }
    .links a { padding: 16px 20px; border-top: 1px dashed var(--dash); }
    .stamp > div { padding: 12px 16px; }
    .v { font-size: 13px; }
    .stamp > div:nth-last-child(-n + 2) { border-bottom: 0; }
  }
</style>

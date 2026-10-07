<script lang="ts">
  import type { Profile } from '../types';
  let { profile }: { profile: Profile } = $props();
  const [first, ...rest] = $derived(profile.name.split(' '));
</script>

<section class="hero">
  <div class="text">
    <div class="eyebrow">Раздел АИ · Интерьеры · М 1:1</div>
    <h1>{first}<br />{rest.join(' ')}</h1>
    <p>{profile.tagline}</p>
    {#if profile.available}<div class="badge">● Беру новые объекты</div>{/if}
  </div>
  <div class="side">
    <div class="photo">
      {#if profile.heroImage}<img src={profile.heroImage} alt="Фото объекта" />{/if}
    </div>
    <div class="stats">
      {#each profile.stats as s}
        <div><div class="l">{s.label}</div><div class="v">{s.value}</div></div>
      {/each}
    </div>
  </div>
</section>

<style>
  .hero { display: grid; grid-template-columns: minmax(0, 1fr) clamp(560px, 38%, 920px); border-bottom: 1px solid var(--ink); }
  .text { padding: 64px 48px; display: flex; flex-direction: column; justify-content: space-between; gap: 48px; border-right: 1px solid var(--ink); }
  h1 { margin: 0; font-family: var(--sans); font-size: clamp(88px, 6.2vw, 150px); line-height: 0.95; font-weight: 700; letter-spacing: -0.045em; overflow-wrap: anywhere; }
  p { margin: 0; font-size: 18px; line-height: 1.6; max-width: 640px; white-space: pre-wrap; }
  .badge { align-self: flex-start; border: 1px dashed var(--ink); padding: 10px 14px; font-size: 13px; letter-spacing: 0.08em; text-transform: uppercase; }
  .photo { aspect-ratio: 1; border-bottom: 1px solid var(--ink); background: var(--hover); }
  .photo img { width: 100%; height: 100%; object-fit: cover; }
  .stats { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); font-size: 12px; text-transform: uppercase; letter-spacing: 0.06em; }
  .stats > div { padding: 14px 16px; border-right: 1px solid var(--ink); }
  .stats > div:last-child { border-right: 0; }
  .l { opacity: 0.8; }
  .v { font-size: 22px; font-weight: 700; margin-top: 4px; }

  @media (min-width: 1800px) {
    .text { padding: 80px 72px; }
    p { font-size: 20px; }
  }

  @media (max-width: 900px) {
    .hero { display: flex; flex-direction: column; }
    .text { padding: 32px 20px 28px; gap: 20px; border-right: 0; border-bottom: 1px solid var(--ink); }
    h1 { font-size: 40px; line-height: 1; letter-spacing: -0.035em; }
    p { font-size: 15px; }
    .badge { padding: 8px 12px; font-size: 12px; }
    .stats { font-size: 12px; letter-spacing: 0.04em; }
    .stats > div { padding: 12px 14px; }
    .v { font-size: 17px; }
  }
</style>

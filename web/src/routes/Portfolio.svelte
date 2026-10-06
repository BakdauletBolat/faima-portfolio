<script lang="ts">
  import { api } from '../lib/api';
  import type { Content } from '../lib/types';
  import Header from '../lib/components/Header.svelte';
  import Hero from '../lib/components/Hero.svelte';
  import Experience from '../lib/components/Experience.svelte';
  import Tools from '../lib/components/Tools.svelte';
  import Works from '../lib/components/Works.svelte';
  import Contact from '../lib/components/Contact.svelte';

  let content = $state<Content | null>(null);
  let failed = $state(false);

  api.content().then((c) => {
    content = c;
    document.title = `${c.profile.name} — Портфолио`;
  }).catch(() => (failed = true));
</script>

<div class="sheet">
  {#if content}
    <Header resume={content.profile.resumePdf} />
    <main>
      <Hero profile={content.profile} />
      <Experience items={content.experience} />
      <Tools tools={content.tools} skills={content.skills} />
      <Works works={content.works} />
      <Contact contacts={content.contacts} name={content.profile.name} />
    </main>
  {:else}
    <p class="state">{failed ? 'Не удалось загрузить данные' : 'Загрузка…'}</p>
  {/if}
</div>

<style>
  .state { padding: 48px; margin: 0; text-transform: uppercase; letter-spacing: 0.08em; font-size: 13px; }
</style>

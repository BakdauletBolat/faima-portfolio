<script lang="ts">
  import Portfolio from './routes/Portfolio.svelte';
  import Admin from './routes/Admin.svelte';

  const isAdmin = location.pathname.replace(/\/+$/, '') === '/admin';
</script>

{#if isAdmin}
  <Admin />
{:else}
  <Portfolio />
{/if}

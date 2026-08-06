<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';

  // For a single-store juice shop, redirect straight to the menu
  onMount(() => {
    push('/menu');
  });
</script>

<div class="max-w-7xl mx-auto px-4 py-16 text-center">
  <span class="text-5xl animate-float inline-block">🍹</span>
  <p class="text-driftwood-mid mt-4">Loading the menu...</p>
</div>

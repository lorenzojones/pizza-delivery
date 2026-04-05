<script>
  import { onMount } from 'svelte';
  import { fetchRestaurant } from '../lib/api.js';
  import MenuItemCard from '../components/MenuItemCard.svelte';
  import CartSidebar from '../components/CartSidebar.svelte';

  export let params = {};

  let restaurant = null;
  let loading = true;
  let error = null;
  let activeCategory = 0;

  onMount(async () => {
    try {
      restaurant = await fetchRestaurant(params.id || 'enos-shack');
      loading = false;
    } catch (err) {
      error = err.message;
      loading = false;
    }
  });

  function scrollToCategory(index) {
    activeCategory = index;
    const el = document.getElementById(`category-${index}`);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }
</script>

{#if loading}
  <div class="max-w-7xl mx-auto px-4 py-6">
    <div class="animate-pulse">
      <div class="h-48 bg-sand rounded-2xl mb-6"></div>
      <div class="h-8 bg-sand rounded w-1/3 mb-3"></div>
      <div class="h-4 bg-sand rounded w-1/4 mb-6"></div>
      <div class="space-y-3">
        {#each Array(5) as _}
          <div class="h-20 bg-sand rounded-2xl"></div>
        {/each}
      </div>
    </div>
  </div>
{:else if error}
  <div class="max-w-7xl mx-auto px-4 py-16 text-center">
    <span class="text-6xl mb-4 block">🥥</span>
    <h2 class="text-xl font-bold text-driftwood mb-2">Couldn't load the menu</h2>
    <p class="text-driftwood-mid">{error}</p>
    <a href="#/" class="inline-block mt-4 text-ocean hover:text-ocean-dark font-bold no-underline">
      ← Back to home
    </a>
  </div>
{:else}
  <!-- Hero banner -->
  <div class="relative h-48 md:h-64 overflow-hidden">
    <img
      src={restaurant.image}
      alt={restaurant.name}
      class="w-full h-full object-cover"
    />
    <div class="absolute inset-0 bg-gradient-to-t from-driftwood/80 to-transparent"></div>
    <div class="absolute bottom-0 left-0 right-0 p-6 max-w-7xl mx-auto">
      <a href="#/" class="text-white/80 hover:text-white text-sm mb-2 inline-flex items-center gap-1 no-underline">
        ← Back
      </a>
      <h1 class="font-display text-3xl md:text-4xl text-white">{restaurant.name}</h1>
    </div>
  </div>

  <div class="max-w-7xl mx-auto px-4 py-6">
    <!-- Info bar -->
    <div class="bg-white rounded-2xl shadow-md p-4 mb-6 flex flex-wrap items-center gap-4 border border-sand">
      <div class="flex items-center gap-1">
        <span class="bg-ocean text-white text-sm font-bold px-2.5 py-1 rounded-full flex items-center gap-1">
          ⭐ {restaurant.rating}
        </span>
        <span class="text-xs text-driftwood-mid">({restaurant.reviewCount} reviews)</span>
      </div>
      <div class="w-px h-5 bg-sand-dark"></div>
      <span class="text-sm text-driftwood-mid">🕐 {restaurant.deliveryTime}</span>
      <div class="w-px h-5 bg-sand-dark"></div>
      <span class="text-sm text-palm font-bold">Pickup ready</span>
      {#if restaurant.tags && restaurant.tags.length > 0}
        <div class="flex gap-1.5 ml-auto">
          {#each restaurant.tags as tag}
            <span class="text-xs bg-ocean-light text-ocean font-bold px-2 py-0.5 rounded-full">{tag}</span>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Menu + Cart -->
    <div class="flex gap-6">
      <div class="flex-1 min-w-0">
        <!-- Category tabs -->
        <div class="bg-white rounded-2xl shadow-sm mb-4 sticky top-16 z-10 border border-sand">
          <div class="flex overflow-x-auto hide-scrollbar">
            {#each restaurant.menu as cat, i}
              <button
                on:click={() => scrollToCategory(i)}
                class="px-5 py-3 text-sm font-bold whitespace-nowrap border-b-2 transition
                  {activeCategory === i ? 'text-ocean border-ocean' : 'text-driftwood-mid border-transparent hover:text-driftwood'}"
              >
                {cat.category}
              </button>
            {/each}
          </div>
        </div>

        <!-- Menu sections -->
        {#each restaurant.menu as category, i}
          <div id="category-{i}" class="mb-6">
            <h2 class="text-lg font-surf text-driftwood mb-3 px-1">{category.category}</h2>
            <div class="space-y-3">
              {#each category.items as item}
                <MenuItemCard {item} {restaurant} />
              {/each}
            </div>
          </div>
        {/each}
      </div>

      <!-- Cart sidebar (desktop) -->
      <div class="hidden lg:block w-80 flex-shrink-0">
        <CartSidebar />
      </div>
    </div>
  </div>
{/if}

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
      restaurant = await fetchRestaurant(params.id);
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
      <div class="h-64 bg-gray-200 rounded-xl mb-6"></div>
      <div class="h-8 bg-gray-200 rounded w-1/3 mb-3"></div>
      <div class="h-4 bg-gray-200 rounded w-1/4 mb-6"></div>
      <div class="space-y-3">
        {#each Array(5) as _}
          <div class="h-20 bg-gray-200 rounded-lg"></div>
        {/each}
      </div>
    </div>
  </div>
{:else if error}
  <div class="max-w-7xl mx-auto px-4 py-16 text-center">
    <span class="text-6xl mb-4 block">😞</span>
    <h2 class="text-xl font-bold text-je-dark mb-2">Restaurant not found</h2>
    <p class="text-je-grey">{error}</p>
    <a href="#/restaurants" class="inline-block mt-4 text-je-orange hover:text-je-orange-dark font-semibold">
      ← Back to restaurants
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
    <div class="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent"></div>
    <div class="absolute bottom-0 left-0 right-0 p-6 max-w-7xl mx-auto">
      <a href="#/restaurants" class="text-white/80 hover:text-white text-sm mb-2 inline-flex items-center gap-1 no-underline">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
        </svg>
        Back
      </a>
      <h1 class="text-2xl md:text-3xl font-extrabold text-white">{restaurant.name}</h1>
    </div>
  </div>

  <div class="max-w-7xl mx-auto px-4 py-6">
    <!-- Restaurant info bar -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-6 flex flex-wrap items-center gap-4">
      <div class="flex items-center gap-1">
        <div class="flex items-center gap-1 bg-je-green text-white text-sm font-bold px-2.5 py-1 rounded">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
          </svg>
          {restaurant.rating}
        </div>
        <span class="text-xs text-je-grey">({restaurant.reviewCount} reviews)</span>
      </div>

      <div class="w-px h-5 bg-je-grey-border"></div>

      <div class="flex items-center gap-1 text-sm text-je-grey">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        {restaurant.deliveryTime}
      </div>

      <div class="w-px h-5 bg-je-grey-border"></div>

      <span class="text-sm text-je-grey">
        {restaurant.deliveryFee === 0 ? 'Free delivery' : `£${restaurant.deliveryFee.toFixed(2)} delivery`}
      </span>

      <div class="w-px h-5 bg-je-grey-border"></div>

      <span class="text-sm text-je-grey">Min order £{restaurant.minOrder.toFixed(2)}</span>

      {#if restaurant.tags && restaurant.tags.length > 0}
        <div class="flex gap-1.5 ml-auto">
          {#each restaurant.tags as tag}
            <span class="text-xs bg-je-orange-light text-je-orange font-medium px-2 py-0.5 rounded-full">{tag}</span>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Main content: menu + cart sidebar -->
    <div class="flex gap-6">
      <!-- Menu -->
      <div class="flex-1 min-w-0">
        <!-- Category tabs -->
        <div class="bg-white rounded-xl shadow-sm mb-4 sticky top-16 z-10">
          <div class="flex overflow-x-auto hide-scrollbar">
            {#each restaurant.menu as cat, i}
              <button
                on:click={() => scrollToCategory(i)}
                class="px-5 py-3 text-sm font-medium whitespace-nowrap border-b-2 transition
                  {activeCategory === i ? 'text-je-orange border-je-orange' : 'text-je-grey border-transparent hover:text-je-dark'}"
              >
                {cat.category}
              </button>
            {/each}
          </div>
        </div>

        <!-- Menu sections -->
        {#each restaurant.menu as category, i}
          <div id="category-{i}" class="mb-6">
            <h2 class="text-lg font-bold text-je-dark mb-3 px-1">{category.category}</h2>
            <div class="bg-white rounded-xl shadow-sm overflow-hidden">
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

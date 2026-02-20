<script>
  import { link } from 'svelte-spa-router';
  export let restaurant;
</script>

<a href="/restaurant/{restaurant.id}" use:link class="block bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-lg transition-all duration-200 group no-underline">
  <!-- Image -->
  <div class="relative h-40 overflow-hidden">
    <img
      src={restaurant.image}
      alt={restaurant.name}
      class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
      loading="lazy"
    />
    {#if restaurant.promoted}
      <div class="absolute top-3 left-3 bg-je-orange text-white text-xs font-bold px-2 py-1 rounded">
        Sponsored
      </div>
    {/if}
    {#if !restaurant.isOpen}
      <div class="absolute inset-0 bg-black/50 flex items-center justify-center">
        <span class="text-white font-bold text-lg">Currently Closed</span>
      </div>
    {/if}
  </div>

  <!-- Content -->
  <div class="p-4">
    <div class="flex items-start justify-between gap-2">
      <h3 class="font-bold text-je-dark text-base group-hover:text-je-orange transition-colors">{restaurant.name}</h3>
      <div class="flex items-center gap-1 bg-je-green text-white text-xs font-bold px-2 py-1 rounded flex-shrink-0">
        <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
        </svg>
        {restaurant.rating}
      </div>
    </div>

    <!-- Cuisine tags -->
    <p class="text-je-grey text-xs mt-1">{restaurant.cuisine.join(' · ')}</p>

    <!-- Meta info -->
    <div class="flex items-center gap-3 mt-3 text-xs text-je-grey">
      <span class="flex items-center gap-1">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        {restaurant.deliveryTime}
      </span>
      <span>·</span>
      {#if restaurant.deliveryFee === 0}
        <span class="text-je-green font-semibold">Free delivery</span>
      {:else}
        <span>£{restaurant.deliveryFee.toFixed(2)} delivery</span>
      {/if}
      <span>·</span>
      <span>Min £{restaurant.minOrder.toFixed(2)}</span>
    </div>

    <!-- Tags -->
    {#if restaurant.tags && restaurant.tags.length > 0}
      <div class="flex flex-wrap gap-1.5 mt-3">
        {#each restaurant.tags as tag}
          <span class="text-xs bg-je-orange-light text-je-orange font-medium px-2 py-0.5 rounded-full">{tag}</span>
        {/each}
      </div>
    {/if}
  </div>
</a>

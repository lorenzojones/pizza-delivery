<script>
  import { link } from 'svelte-spa-router';
  export let restaurant;
</script>

<a href="/menu" use:link class="block bg-white rounded-2xl overflow-hidden shadow-md hover:shadow-xl transition-all duration-200 group no-underline border border-sand">
  <div class="relative h-40 overflow-hidden">
    <img
      src={restaurant.image}
      alt={restaurant.name}
      class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
      loading="lazy"
    />
  </div>
  <div class="p-4">
    <div class="flex items-start justify-between gap-2">
      <h3 class="font-bold text-driftwood text-base group-hover:text-ocean transition-colors">{restaurant.name}</h3>
      <div class="flex items-center gap-1 bg-ocean text-white text-xs font-bold px-2 py-1 rounded-full flex-shrink-0">
        ⭐ {restaurant.rating}
      </div>
    </div>
    <p class="text-driftwood-mid text-xs mt-1">{restaurant.cuisine.join(' · ')}</p>
    <div class="flex items-center gap-3 mt-3 text-xs text-driftwood-mid">
      <span>🕐 {restaurant.deliveryTime}</span>
      <span>·</span>
      <span class="text-palm font-semibold">Pickup ready</span>
    </div>
  </div>
</a>

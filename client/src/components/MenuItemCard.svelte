<script>
  import { cart } from '../lib/cart.js';
  export let item;
  export let restaurant;

  let added = false;

  function addToCart() {
    cart.addItem(restaurant, item);
    added = true;
    setTimeout(() => added = false, 1000);
  }
</script>

<button
  on:click={addToCart}
  class="flex items-center justify-between p-4 bg-white rounded-2xl border border-sand hover:border-ocean/30 hover:shadow-lg transition-all duration-200 group w-full text-left"
>
  <div class="flex-1 pr-4">
    <div class="flex items-center gap-2">
      {#if item.emoji}
        <span class="text-xl">{item.emoji}</span>
      {/if}
      <h4 class="font-bold text-driftwood text-sm group-hover:text-ocean transition">{item.name}</h4>
      {#if item.popular}
        <span class="text-[10px] bg-sunset text-white font-bold px-1.5 py-0.5 rounded-full">POPULAR</span>
      {/if}
    </div>
    {#if item.description}
      <p class="text-driftwood-mid text-xs mt-1 line-clamp-2">{item.description}</p>
    {/if}
    <p class="text-ocean font-extrabold text-sm mt-2">${item.price} <span class="text-xs font-medium text-driftwood-mid">MXN</span></p>
  </div>

  <div class="flex items-center gap-3">
    {#if item.image}
      <img
        src={item.image}
        alt={item.name}
        class="w-20 h-20 rounded-xl object-cover flex-shrink-0 group-hover:scale-105 transition-transform duration-300"
        loading="lazy"
      />
    {/if}

    <div
      class="w-9 h-9 rounded-full flex items-center justify-center flex-shrink-0 transition-all duration-200
        {added ? 'bg-palm text-white scale-110' : 'bg-ocean text-white hover:bg-ocean-dark'}"
    >
      {#if added}
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
        </svg>
      {:else}
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
      {/if}
    </div>
  </div>
</button>

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

<div class="flex items-center justify-between p-4 bg-white border-b border-je-grey-light hover:bg-gray-50 transition group">
  <div class="flex-1 pr-4">
    <div class="flex items-center gap-2">
      <h4 class="font-semibold text-je-dark text-sm">{item.name}</h4>
      {#if item.popular}
        <span class="text-[10px] bg-je-orange text-white font-bold px-1.5 py-0.5 rounded">POPULAR</span>
      {/if}
    </div>
    {#if item.description}
      <p class="text-je-grey text-xs mt-1 line-clamp-2">{item.description}</p>
    {/if}
    <p class="text-je-dark font-bold text-sm mt-2">£{item.price.toFixed(2)}</p>
  </div>

  <div class="flex items-center gap-3">
    {#if item.image}
      <img
        src={item.image}
        alt={item.name}
        class="w-20 h-20 rounded-lg object-cover flex-shrink-0"
        loading="lazy"
      />
    {/if}

    <button
      on:click={addToCart}
      class="w-9 h-9 rounded-full flex items-center justify-center flex-shrink-0 transition-all duration-200
        {added ? 'bg-je-green text-white scale-110' : 'bg-je-orange text-white hover:bg-je-orange-dark'}"
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
    </button>
  </div>
</div>

<script>
  import { cart, cartItemCount, cartSubtotal, cartTotal } from '../lib/cart.js';
</script>

<div class="bg-white rounded-2xl shadow-lg overflow-hidden sticky top-20 border border-sand">
  <div class="p-4 bg-gradient-to-r from-ocean to-ocean-dark text-white">
    <h3 class="font-surf text-lg">Your Order 🤙</h3>
  </div>

  {#if $cartItemCount === 0}
    <div class="p-8 text-center">
      <span class="text-5xl block mb-3">🥥</span>
      <p class="text-driftwood-mid text-sm">Nothing here yet...</p>
      <p class="text-driftwood-mid text-xs mt-1">Add some juice to get started!</p>
    </div>
  {:else}
    <!-- Items -->
    <div class="max-h-80 overflow-y-auto">
      {#each $cart.items as item}
        <div class="flex items-center justify-between px-4 py-3 border-b border-sand-light">
          <div class="flex items-center gap-3 flex-1">
            <div class="flex items-center gap-1">
              <button
                on:click={() => cart.updateQuantity(item.id, item.quantity - 1)}
                class="w-6 h-6 rounded-full bg-sand hover:bg-sunset-light text-driftwood-mid hover:text-sunset flex items-center justify-center text-xs transition font-bold"
              >-</button>
              <span class="text-sm font-bold text-ocean w-5 text-center">{item.quantity}</span>
              <button
                on:click={() => cart.updateQuantity(item.id, item.quantity + 1)}
                class="w-6 h-6 rounded-full bg-sand hover:bg-palm-light text-driftwood-mid hover:text-palm flex items-center justify-center text-xs transition font-bold"
              >+</button>
            </div>
            <span class="text-sm text-driftwood truncate">{item.name}</span>
          </div>
          <span class="text-sm font-bold text-driftwood ml-2">${item.price * item.quantity}</span>
        </div>
      {/each}
    </div>

    <!-- Totals -->
    <div class="p-4 space-y-2 border-t border-sand">
      <div class="flex justify-between text-base font-bold pt-1">
        <span>Total</span>
        <span>${$cartSubtotal} MXN</span>
      </div>
    </div>

    <!-- Checkout button -->
    <div class="p-4 pt-0">
      <a
        href="#/checkout"
        class="block w-full text-center rounded-xl py-3 font-bold text-sm transition no-underline bg-sunset hover:bg-sunset-dark text-white"
      >
        Order Now 🏄 · ${$cartSubtotal} MXN
      </a>
    </div>
  {/if}
</div>

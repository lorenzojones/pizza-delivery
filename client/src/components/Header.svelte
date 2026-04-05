<script>
  import { link } from 'svelte-spa-router';
  import { cart, cartItemCount, cartSubtotal } from '../lib/cart.js';

  let showCartPreview = false;

  function toggleCart() {
    showCartPreview = !showCartPreview;
  }

  function formatPrice(price) {
    return `$${price} MXN`;
  }
</script>

<header class="bg-ocean sticky top-0 z-50 shadow-lg">
  <div class="max-w-7xl mx-auto px-4">
    <div class="flex items-center justify-between h-16">
      <!-- Logo -->
      <a href="/" use:link class="flex items-center gap-2 text-white no-underline">
        <span class="text-2xl">🍹</span>
        <span class="text-xl font-display tracking-wide">Eno's Shack</span>
      </a>

      <!-- Location (desktop) -->
      <div class="hidden md:flex items-center bg-white/20 backdrop-blur-sm rounded-full px-4 py-2 flex-1 max-w-sm mx-8">
        <span class="mr-2">📍</span>
        <span class="text-sm text-white/90 truncate">Sayulita, Nayarit, Mexico</span>
      </div>

      <!-- Right side -->
      <div class="flex items-center gap-4">
        <!-- Menu link -->
        <a href="/menu" use:link class="hidden md:flex items-center gap-1 text-white text-sm font-bold hover:text-sand transition no-underline">
          Menu
        </a>

        <!-- Cart button -->
        <div class="relative">
          <button
            on:click={toggleCart}
            class="flex items-center gap-2 bg-sunset hover:bg-sunset-dark text-white rounded-full px-4 py-2 text-sm font-bold transition shadow-sm"
          >
            <span class="text-lg">🛒</span>
            {#if $cartItemCount > 0}
              <span>{formatPrice($cartSubtotal)}</span>
              <span class="bg-white text-sunset w-5 h-5 rounded-full text-xs flex items-center justify-center font-black">{$cartItemCount}</span>
            {/if}
          </button>

          <!-- Cart dropdown -->
          {#if showCartPreview && $cartItemCount > 0}
            <div class="absolute right-0 top-12 w-80 bg-white rounded-2xl shadow-2xl border border-sand-dark z-50 overflow-hidden">
              <div class="p-4 border-b border-sand bg-sand-light">
                <h3 class="font-bold text-driftwood font-surf text-lg">Your Order</h3>
                <p class="text-xs text-driftwood-mid mt-0.5">from Eno's Shack</p>
              </div>
              <div class="max-h-64 overflow-y-auto">
                {#each $cart.items as item}
                  <div class="flex items-center justify-between px-4 py-3 border-b border-sand-light">
                    <div class="flex items-center gap-3">
                      <span class="bg-sunset-light text-sunset text-xs font-bold w-6 h-6 rounded-full flex items-center justify-center">{item.quantity}</span>
                      <span class="text-sm text-driftwood">{item.name}</span>
                    </div>
                    <span class="text-sm font-bold text-driftwood">${item.price * item.quantity}</span>
                  </div>
                {/each}
              </div>
              <div class="p-4 bg-sand-light">
                <div class="flex justify-between text-sm mb-3">
                  <span class="text-driftwood-mid">Total</span>
                  <span class="font-bold text-driftwood">{formatPrice($cartSubtotal)}</span>
                </div>
                <a
                  href="#/checkout"
                  on:click={() => showCartPreview = false}
                  class="block w-full bg-sunset hover:bg-sunset-dark text-white text-center rounded-xl py-3 font-bold text-sm transition no-underline"
                >
                  Checkout {formatPrice($cartSubtotal)}
                </a>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</header>

<!-- Click outside to close cart -->
{#if showCartPreview}
  <button class="fixed inset-0 z-40 bg-transparent cursor-default" on:click={() => showCartPreview = false} aria-label="Close cart"></button>
{/if}

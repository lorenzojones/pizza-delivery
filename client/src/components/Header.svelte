<script>
  import { link } from 'svelte-spa-router';
  import { cart, cartItemCount, cartSubtotal } from '../lib/cart.js';

  let showCartPreview = false;

  function toggleCart() {
    showCartPreview = !showCartPreview;
  }
</script>

<!-- Top bar -->
<header class="bg-je-orange sticky top-0 z-50 shadow-md">
  <div class="max-w-7xl mx-auto px-4">
    <div class="flex items-center justify-between h-16">
      <!-- Logo -->
      <a href="/" use:link class="flex items-center gap-2 text-white no-underline">
        <svg class="w-8 h-8" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
        </svg>
        <span class="text-xl font-extrabold tracking-tight">SliceNow</span>
      </a>

      <!-- Location bar (desktop) -->
      <div class="hidden md:flex items-center bg-white rounded-full px-4 py-2 flex-1 max-w-md mx-8">
        <svg class="w-5 h-5 text-je-orange mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
        <span class="text-sm text-je-grey truncate">London, United Kingdom</span>
      </div>

      <!-- Right side -->
      <div class="flex items-center gap-4">
        <!-- Account -->
        <button class="hidden md:flex items-center gap-1 text-white text-sm font-medium hover:text-je-orange-light transition">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
          </svg>
          Log in
        </button>

        <!-- Cart button -->
        <div class="relative">
          <button
            on:click={toggleCart}
            class="flex items-center gap-2 bg-je-green hover:bg-je-green-dark text-white rounded-full px-4 py-2 text-sm font-semibold transition shadow-sm"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"/>
            </svg>
            {#if $cartItemCount > 0}
              <span>£{$cartSubtotal.toFixed(2)}</span>
              <span class="bg-white text-je-green w-5 h-5 rounded-full text-xs flex items-center justify-center font-bold">{$cartItemCount}</span>
            {/if}
          </button>

          <!-- Cart dropdown -->
          {#if showCartPreview && $cartItemCount > 0}
            <div class="absolute right-0 top-12 w-80 bg-white rounded-xl shadow-2xl border border-je-grey-border z-50 overflow-hidden">
              <div class="p-4 border-b border-je-grey-border">
                <h3 class="font-bold text-je-dark">Your order</h3>
                <p class="text-xs text-je-grey mt-0.5">from {$cart.restaurantName}</p>
              </div>
              <div class="max-h-64 overflow-y-auto">
                {#each $cart.items as item}
                  <div class="flex items-center justify-between px-4 py-3 border-b border-je-grey-light">
                    <div class="flex items-center gap-3">
                      <span class="bg-je-orange-light text-je-orange text-xs font-bold w-6 h-6 rounded-full flex items-center justify-center">{item.quantity}</span>
                      <span class="text-sm text-je-dark">{item.name}</span>
                    </div>
                    <span class="text-sm font-semibold text-je-dark">£{(item.price * item.quantity).toFixed(2)}</span>
                  </div>
                {/each}
              </div>
              <div class="p-4 bg-je-grey-light">
                <div class="flex justify-between text-sm mb-1">
                  <span class="text-je-grey">Subtotal</span>
                  <span class="font-semibold">£{$cartSubtotal.toFixed(2)}</span>
                </div>
                <div class="flex justify-between text-sm mb-3">
                  <span class="text-je-grey">Delivery</span>
                  <span class="font-semibold">£{$cart.deliveryFee.toFixed(2)}</span>
                </div>
                <a
                  href="#/checkout"
                  on:click={() => showCartPreview = false}
                  class="block w-full bg-je-green hover:bg-je-green-dark text-white text-center rounded-lg py-3 font-semibold text-sm transition no-underline"
                >
                  Go to checkout · £{($cartSubtotal + $cart.deliveryFee).toFixed(2)}
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

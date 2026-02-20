<script>
  import { cart, cartItemCount, cartSubtotal, cartTotal } from '../lib/cart.js';

  $: meetsMinimum = $cartSubtotal >= $cart.minOrder;
  $: amountNeeded = ($cart.minOrder - $cartSubtotal).toFixed(2);
</script>

<div class="bg-white rounded-xl shadow-sm overflow-hidden sticky top-20">
  <div class="p-4 bg-je-orange text-white">
    <h3 class="font-bold text-lg">Your order</h3>
    {#if $cart.restaurantName}
      <p class="text-sm opacity-90 mt-0.5">from {$cart.restaurantName}</p>
    {/if}
  </div>

  {#if $cartItemCount === 0}
    <div class="p-8 text-center">
      <svg class="w-16 h-16 mx-auto text-je-grey-border mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z"/>
      </svg>
      <p class="text-je-grey text-sm">Your basket is empty</p>
      <p class="text-je-grey text-xs mt-1">Add items to get started</p>
    </div>
  {:else}
    <!-- Items -->
    <div class="max-h-80 overflow-y-auto">
      {#each $cart.items as item}
        <div class="flex items-center justify-between px-4 py-3 border-b border-je-grey-light">
          <div class="flex items-center gap-3 flex-1">
            <!-- Quantity controls -->
            <div class="flex items-center gap-1">
              <button
                on:click={() => cart.updateQuantity(item.id, item.quantity - 1)}
                class="w-6 h-6 rounded-full bg-je-grey-light hover:bg-je-grey-border text-je-grey flex items-center justify-center text-xs transition"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/>
                </svg>
              </button>
              <span class="text-sm font-bold text-je-dark w-5 text-center">{item.quantity}</span>
              <button
                on:click={() => cart.updateQuantity(item.id, item.quantity + 1)}
                class="w-6 h-6 rounded-full bg-je-orange-light hover:bg-je-orange text-je-orange hover:text-white flex items-center justify-center text-xs transition"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </button>
            </div>
            <span class="text-sm text-je-dark truncate">{item.name}</span>
          </div>
          <span class="text-sm font-semibold text-je-dark ml-2">£{(item.price * item.quantity).toFixed(2)}</span>
        </div>
      {/each}
    </div>

    <!-- Totals -->
    <div class="p-4 space-y-2 border-t border-je-grey-border">
      <div class="flex justify-between text-sm">
        <span class="text-je-grey">Subtotal</span>
        <span class="font-medium">£{$cartSubtotal.toFixed(2)}</span>
      </div>
      <div class="flex justify-between text-sm">
        <span class="text-je-grey">Delivery fee</span>
        <span class="font-medium">£{$cart.deliveryFee.toFixed(2)}</span>
      </div>
      <div class="flex justify-between text-base font-bold pt-2 border-t border-je-grey-light">
        <span>Total</span>
        <span>£{$cartTotal.toFixed(2)}</span>
      </div>
    </div>

    <!-- Checkout button -->
    <div class="p-4 pt-0">
      {#if !meetsMinimum}
        <p class="text-xs text-red-500 mb-2 text-center">Add £{amountNeeded} more to meet the minimum order</p>
      {/if}
      <a
        href="#/checkout"
        class="block w-full text-center rounded-lg py-3 font-bold text-sm transition no-underline
          {meetsMinimum ? 'bg-je-green hover:bg-je-green-dark text-white' : 'bg-gray-200 text-gray-400 pointer-events-none'}"
      >
        Go to checkout · £{$cartTotal.toFixed(2)}
      </a>
    </div>
  {/if}
</div>

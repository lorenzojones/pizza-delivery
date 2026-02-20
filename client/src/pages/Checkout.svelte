<script>
  import { push } from 'svelte-spa-router';
  import { cart, cartItemCount, cartSubtotal, cartTotal } from '../lib/cart.js';
  import { placeOrder } from '../lib/api.js';

  let customer = { name: '', email: '', phone: '' };
  let address = { line1: '', line2: '', city: 'London', postcode: '' };
  let paymentMethod = 'card';
  let submitting = false;
  let error = null;

  let cardNumber = '';
  let cardExpiry = '';
  let cardCvc = '';

  $: canSubmit = customer.name && customer.email && customer.phone && address.line1 && address.postcode && $cartItemCount > 0;

  async function handleSubmit() {
    if (!canSubmit) return;
    submitting = true;
    error = null;

    try {
      const order = await placeOrder({
        restaurantId: $cart.restaurantId,
        items: $cart.items.map(i => ({ id: i.id, name: i.name, price: i.price, quantity: i.quantity })),
        customer,
        address: `${address.line1}${address.line2 ? ', ' + address.line2 : ''}, ${address.city}, ${address.postcode}`,
        paymentMethod
      });

      cart.clear();
      push(`/order/${order.id}`);
    } catch (err) {
      error = err.message;
      submitting = false;
    }
  }
</script>

<div class="max-w-4xl mx-auto px-4 py-6">
  <a href="#/restaurants" class="text-je-grey hover:text-je-dark text-sm mb-4 inline-flex items-center gap-1 no-underline">
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
    </svg>
    Back to menu
  </a>

  <h1 class="text-2xl font-extrabold text-je-dark mb-6">Checkout</h1>

  {#if $cartItemCount === 0}
    <div class="bg-white rounded-xl shadow-sm p-12 text-center">
      <span class="text-6xl mb-4 block">🛒</span>
      <h2 class="text-xl font-bold text-je-dark mb-2">Your basket is empty</h2>
      <p class="text-je-grey mb-4">Add some delicious items before checking out</p>
      <a href="#/restaurants" class="inline-block bg-je-orange hover:bg-je-orange-dark text-white font-semibold px-6 py-3 rounded-lg transition no-underline">
        Browse restaurants
      </a>
    </div>
  {:else}
    <form on:submit|preventDefault={handleSubmit}>
      <div class="flex flex-col lg:flex-row gap-6">
        <!-- Left: forms -->
        <div class="flex-1 space-y-6">
          <!-- Delivery details -->
          <div class="bg-white rounded-xl shadow-sm p-6">
            <h2 class="font-bold text-je-dark mb-4 flex items-center gap-2">
              <span class="w-7 h-7 bg-je-orange text-white rounded-full flex items-center justify-center text-sm font-bold">1</span>
              Delivery details
            </h2>
            <div class="grid md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-je-grey mb-1">Address line 1 *</label>
                <input bind:value={address.line1} required
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="123 High Street" />
              </div>
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-je-grey mb-1">Address line 2</label>
                <input bind:value={address.line2}
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="Flat 4B" />
              </div>
              <div>
                <label class="block text-sm font-medium text-je-grey mb-1">City</label>
                <input bind:value={address.city}
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" />
              </div>
              <div>
                <label class="block text-sm font-medium text-je-grey mb-1">Postcode *</label>
                <input bind:value={address.postcode} required
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="EC1A 1BB" />
              </div>
            </div>
          </div>

          <!-- Personal details -->
          <div class="bg-white rounded-xl shadow-sm p-6">
            <h2 class="font-bold text-je-dark mb-4 flex items-center gap-2">
              <span class="w-7 h-7 bg-je-orange text-white rounded-full flex items-center justify-center text-sm font-bold">2</span>
              Personal details
            </h2>
            <div class="grid md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-je-grey mb-1">Full name *</label>
                <input bind:value={customer.name} required
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="John Smith" />
              </div>
              <div>
                <label class="block text-sm font-medium text-je-grey mb-1">Email *</label>
                <input bind:value={customer.email} type="email" required
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="john@email.com" />
              </div>
              <div>
                <label class="block text-sm font-medium text-je-grey mb-1">Phone *</label>
                <input bind:value={customer.phone} type="tel" required
                  class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="07123 456789" />
              </div>
            </div>
          </div>

          <!-- Payment -->
          <div class="bg-white rounded-xl shadow-sm p-6">
            <h2 class="font-bold text-je-dark mb-4 flex items-center gap-2">
              <span class="w-7 h-7 bg-je-orange text-white rounded-full flex items-center justify-center text-sm font-bold">3</span>
              Payment
            </h2>

            <!-- Payment method selector -->
            <div class="flex gap-3 mb-4">
              <button type="button" on:click={() => paymentMethod = 'card'}
                class="flex-1 p-3 rounded-lg border-2 text-sm font-medium transition
                  {paymentMethod === 'card' ? 'border-je-orange bg-je-orange-light text-je-orange' : 'border-je-grey-border text-je-grey hover:border-je-grey'}">
                💳 Card
              </button>
              <button type="button" on:click={() => paymentMethod = 'cash'}
                class="flex-1 p-3 rounded-lg border-2 text-sm font-medium transition
                  {paymentMethod === 'cash' ? 'border-je-orange bg-je-orange-light text-je-orange' : 'border-je-grey-border text-je-grey hover:border-je-grey'}">
                💵 Cash
              </button>
            </div>

            {#if paymentMethod === 'card'}
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-je-grey mb-1">Card number</label>
                  <input bind:value={cardNumber}
                    class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="1234 5678 9012 3456" />
                </div>
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-je-grey mb-1">Expiry</label>
                    <input bind:value={cardExpiry}
                      class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="MM/YY" />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-je-grey mb-1">CVC</label>
                    <input bind:value={cardCvc}
                      class="w-full px-4 py-3 rounded-lg border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm" placeholder="123" />
                  </div>
                </div>
              </div>
            {:else}
              <p class="text-sm text-je-grey bg-je-grey-light rounded-lg p-4">
                💡 Please have the exact amount ready. Our riders may not carry change.
              </p>
            {/if}
          </div>
        </div>

        <!-- Right: Order summary -->
        <div class="w-full lg:w-80 flex-shrink-0">
          <div class="bg-white rounded-xl shadow-sm overflow-hidden sticky top-20">
            <div class="p-4 bg-je-dark text-white">
              <h3 class="font-bold">Order summary</h3>
              <p class="text-sm text-gray-400 mt-0.5">from {$cart.restaurantName}</p>
            </div>

            <div class="max-h-60 overflow-y-auto">
              {#each $cart.items as item}
                <div class="flex justify-between px-4 py-3 border-b border-je-grey-light text-sm">
                  <div class="flex items-center gap-2">
                    <span class="font-bold text-je-orange">{item.quantity}x</span>
                    <span class="text-je-dark">{item.name}</span>
                  </div>
                  <span class="font-semibold">£{(item.price * item.quantity).toFixed(2)}</span>
                </div>
              {/each}
            </div>

            <div class="p-4 space-y-2 border-t border-je-grey-border">
              <div class="flex justify-between text-sm">
                <span class="text-je-grey">Subtotal</span>
                <span>£{$cartSubtotal.toFixed(2)}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-je-grey">Delivery fee</span>
                <span>£{$cart.deliveryFee.toFixed(2)}</span>
              </div>
              <div class="flex justify-between font-bold text-base pt-2 border-t border-je-grey-light">
                <span>Total</span>
                <span>£{$cartTotal.toFixed(2)}</span>
              </div>
            </div>

            <div class="p-4 pt-0">
              {#if error}
                <p class="text-red-500 text-xs mb-2 text-center">{error}</p>
              {/if}
              <button
                type="submit"
                disabled={!canSubmit || submitting}
                class="w-full rounded-lg py-3 font-bold text-sm transition
                  {canSubmit && !submitting
                    ? 'bg-je-green hover:bg-je-green-dark text-white cursor-pointer'
                    : 'bg-gray-200 text-gray-400 cursor-not-allowed'}"
              >
                {#if submitting}
                  <span class="inline-flex items-center gap-2">
                    <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Placing order...
                  </span>
                {:else}
                  Place order · £{$cartTotal.toFixed(2)}
                {/if}
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  {/if}
</div>

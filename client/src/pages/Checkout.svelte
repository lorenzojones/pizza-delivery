<script>
  import { push } from 'svelte-spa-router';
  import { cart, cartItemCount, cartSubtotal, cartTotal } from '../lib/cart.js';
  import { placeOrder } from '../lib/api.js';

  let customer = { name: '', email: '', phone: '' };
  let address = 'Pickup at Eno\'s Shack, Calle Delfines 12, Sayulita';
  let orderType = 'pickup';
  let paymentMethod = 'cash';
  let submitting = false;
  let error = null;
  let notes = '';

  $: canSubmit = customer.name && customer.phone && $cartItemCount > 0;

  async function handleSubmit() {
    if (!canSubmit) return;
    submitting = true;
    error = null;

    try {
      const order = await placeOrder({
        restaurantId: $cart.restaurantId,
        items: $cart.items.map(i => ({ id: i.id, name: i.name, price: i.price, quantity: i.quantity })),
        customer,
        address: orderType === 'pickup' ? 'Pickup at Eno\'s Shack' : address,
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
  <a href="#/menu" class="text-driftwood-mid hover:text-driftwood text-sm mb-4 inline-flex items-center gap-1 no-underline">
    ← Back to menu
  </a>

  <h1 class="text-2xl font-surf text-ocean mb-6">Checkout 🤙</h1>

  {#if $cartItemCount === 0}
    <div class="bg-white rounded-2xl shadow-md p-12 text-center border border-sand">
      <span class="text-6xl mb-4 block">🥥</span>
      <h2 class="text-xl font-bold text-driftwood mb-2">Your order is empty</h2>
      <p class="text-driftwood-mid mb-4">Add some juice to get started!</p>
      <a href="#/menu" class="inline-block bg-ocean hover:bg-ocean-dark text-white font-bold px-6 py-3 rounded-xl transition no-underline">
        See the Menu
      </a>
    </div>
  {:else}
    <form on:submit|preventDefault={handleSubmit}>
      <div class="flex flex-col lg:flex-row gap-6">
        <!-- Left: forms -->
        <div class="flex-1 space-y-6">
          <!-- Order type -->
          <div class="bg-white rounded-2xl shadow-md p-6 border border-sand">
            <h2 class="font-bold text-driftwood mb-4 flex items-center gap-2">
              <span class="w-7 h-7 bg-ocean text-white rounded-full flex items-center justify-center text-sm font-bold">1</span>
              Order Type
            </h2>
            <div class="flex gap-3">
              <button type="button" on:click={() => orderType = 'pickup'}
                class="flex-1 p-4 rounded-xl border-2 text-sm font-bold transition text-center
                  {orderType === 'pickup' ? 'border-ocean bg-ocean-light text-ocean' : 'border-sand text-driftwood-mid hover:border-sand-dark'}">
                🏖️ Pickup at Shack
              </button>
              <button type="button" on:click={() => orderType = 'beach'}
                class="flex-1 p-4 rounded-xl border-2 text-sm font-bold transition text-center
                  {orderType === 'beach' ? 'border-ocean bg-ocean-light text-ocean' : 'border-sand text-driftwood-mid hover:border-sand-dark'}">
                🏄 Beach Delivery
              </button>
            </div>
            {#if orderType === 'beach'}
              <div class="mt-4">
                <label class="block text-sm font-medium text-driftwood-mid mb-1">Where on the beach?</label>
                <input bind:value={address}
                  class="w-full px-4 py-3 rounded-xl border border-sand focus:border-ocean focus:ring-1 focus:ring-ocean outline-none text-sm"
                  placeholder="Near the main break, blue umbrella" />
              </div>
            {/if}
          </div>

          <!-- Your details -->
          <div class="bg-white rounded-2xl shadow-md p-6 border border-sand">
            <h2 class="font-bold text-driftwood mb-4 flex items-center gap-2">
              <span class="w-7 h-7 bg-ocean text-white rounded-full flex items-center justify-center text-sm font-bold">2</span>
              Your Details
            </h2>
            <div class="grid md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-driftwood-mid mb-1">Name *</label>
                <input bind:value={customer.name} required
                  class="w-full px-4 py-3 rounded-xl border border-sand focus:border-ocean focus:ring-1 focus:ring-ocean outline-none text-sm" placeholder="Your name" />
              </div>
              <div>
                <label class="block text-sm font-medium text-driftwood-mid mb-1">Phone / WhatsApp *</label>
                <input bind:value={customer.phone} type="tel" required
                  class="w-full px-4 py-3 rounded-xl border border-sand focus:border-ocean focus:ring-1 focus:ring-ocean outline-none text-sm" placeholder="+52 322 xxx xxxx" />
              </div>
              <div>
                <label class="block text-sm font-medium text-driftwood-mid mb-1">Email</label>
                <input bind:value={customer.email} type="email"
                  class="w-full px-4 py-3 rounded-xl border border-sand focus:border-ocean focus:ring-1 focus:ring-ocean outline-none text-sm" placeholder="you@email.com" />
              </div>
            </div>
          </div>

          <!-- Payment -->
          <div class="bg-white rounded-2xl shadow-md p-6 border border-sand">
            <h2 class="font-bold text-driftwood mb-4 flex items-center gap-2">
              <span class="w-7 h-7 bg-ocean text-white rounded-full flex items-center justify-center text-sm font-bold">3</span>
              Payment
            </h2>
            <div class="flex gap-3">
              <button type="button" on:click={() => paymentMethod = 'cash'}
                class="flex-1 p-3 rounded-xl border-2 text-sm font-bold transition
                  {paymentMethod === 'cash' ? 'border-ocean bg-ocean-light text-ocean' : 'border-sand text-driftwood-mid hover:border-sand-dark'}">
                💵 Cash (MXN)
              </button>
              <button type="button" on:click={() => paymentMethod = 'card'}
                class="flex-1 p-3 rounded-xl border-2 text-sm font-bold transition
                  {paymentMethod === 'card' ? 'border-ocean bg-ocean-light text-ocean' : 'border-sand text-driftwood-mid hover:border-sand-dark'}">
                💳 Card
              </button>
            </div>
            {#if paymentMethod === 'cash'}
              <p class="text-sm text-driftwood-mid bg-sand-light rounded-xl p-4 mt-3">
                🤙 Pay in pesos when you pick up. We also accept USD at today's rate.
              </p>
            {:else}
              <p class="text-sm text-driftwood-mid bg-sand-light rounded-xl p-4 mt-3">
                💳 Pay by card at the shack when your order is ready.
              </p>
            {/if}
          </div>

          <!-- Special notes -->
          <div class="bg-white rounded-2xl shadow-md p-6 border border-sand">
            <label class="block font-bold text-driftwood mb-2">Special Requests</label>
            <textarea bind:value={notes} rows="2"
              class="w-full px-4 py-3 rounded-xl border border-sand focus:border-ocean focus:ring-1 focus:ring-ocean outline-none text-sm resize-none"
              placeholder="Extra ice? No sugar? Let us know..."></textarea>
          </div>
        </div>

        <!-- Right: Order summary -->
        <div class="w-full lg:w-80 flex-shrink-0">
          <div class="bg-white rounded-2xl shadow-md overflow-hidden sticky top-20 border border-sand">
            <div class="p-4 bg-gradient-to-r from-ocean to-ocean-dark text-white">
              <h3 class="font-surf text-lg">Order Summary 🍹</h3>
            </div>

            <div class="max-h-60 overflow-y-auto">
              {#each $cart.items as item}
                <div class="flex justify-between px-4 py-3 border-b border-sand-light text-sm">
                  <div class="flex items-center gap-2">
                    <span class="font-bold text-ocean">{item.quantity}x</span>
                    <span class="text-driftwood">{item.name}</span>
                  </div>
                  <span class="font-bold">${item.price * item.quantity}</span>
                </div>
              {/each}
            </div>

            <div class="p-4 space-y-2 border-t border-sand">
              <div class="flex justify-between font-bold text-base pt-1">
                <span>Total</span>
                <span>${$cartSubtotal} MXN</span>
              </div>
            </div>

            <div class="p-4 pt-0">
              {#if error}
                <p class="text-coral text-xs mb-2 text-center">{error}</p>
              {/if}
              <button
                type="submit"
                disabled={!canSubmit || submitting}
                class="w-full rounded-xl py-3 font-bold text-sm transition
                  {canSubmit && !submitting
                    ? 'bg-sunset hover:bg-sunset-dark text-white cursor-pointer'
                    : 'bg-sand text-driftwood-mid cursor-not-allowed'}"
              >
                {#if submitting}
                  <span class="inline-flex items-center gap-2">
                    <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Blending your order...
                  </span>
                {:else}
                  Place Order · ${$cartSubtotal} MXN 🏄
                {/if}
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  {/if}
</div>

<script>
  import { onMount, onDestroy } from 'svelte';
  import { fetchOrder } from '../lib/api.js';

  export let params = {};

  let order = null;
  let loading = true;
  let error = null;
  let pollInterval;

  const statusSteps = ['confirmed', 'preparing', 'ready', 'delivering', 'delivered'];
  const statusIcons = {
    confirmed: '✅',
    preparing: '👨‍🍳',
    ready: '📦',
    delivering: '🛵',
    delivered: '🎉'
  };
  const statusLabels = {
    confirmed: 'Confirmed',
    preparing: 'Preparing',
    ready: 'Ready',
    delivering: 'On its way',
    delivered: 'Delivered'
  };

  $: currentStepIndex = order ? statusSteps.indexOf(order.status) : 0;
  $: progressPercent = order ? ((currentStepIndex) / (statusSteps.length - 1)) * 100 : 0;

  onMount(async () => {
    await loadOrder();
    // Poll for updates every 5 seconds
    pollInterval = setInterval(loadOrder, 5000);
  });

  onDestroy(() => {
    if (pollInterval) clearInterval(pollInterval);
  });

  async function loadOrder() {
    try {
      order = await fetchOrder(params.id);
      if (order.status === 'delivered' && pollInterval) {
        clearInterval(pollInterval);
      }
      if (loading) loading = false;
    } catch (err) {
      error = err.message;
      loading = false;
      if (pollInterval) clearInterval(pollInterval);
    }
  }
</script>

<div class="max-w-3xl mx-auto px-4 py-6">
  {#if loading}
    <div class="animate-pulse space-y-6">
      <div class="h-8 bg-gray-200 rounded w-1/3"></div>
      <div class="h-48 bg-gray-200 rounded-xl"></div>
      <div class="h-32 bg-gray-200 rounded-xl"></div>
    </div>
  {:else if error}
    <div class="bg-white rounded-xl shadow-sm p-12 text-center">
      <span class="text-6xl mb-4 block">❓</span>
      <h2 class="text-xl font-bold text-je-dark mb-2">Order not found</h2>
      <p class="text-je-grey">{error}</p>
      <a href="#/" class="inline-block mt-4 text-je-orange hover:text-je-orange-dark font-semibold no-underline">
        Go to homepage
      </a>
    </div>
  {:else}
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center gap-3 mb-1">
        <h1 class="text-2xl font-extrabold text-je-dark">Order #{order.id}</h1>
        <span class="text-2xl">{statusIcons[order.status]}</span>
      </div>
      <p class="text-je-grey text-sm">from {order.restaurantName}</p>
    </div>

    <!-- Progress tracker -->
    <div class="bg-white rounded-xl shadow-sm p-6 mb-6">
      <div class="flex items-center justify-between mb-2">
        <h2 class="font-bold text-je-dark">Order status</h2>
        <span class="text-sm font-semibold text-je-orange">{statusLabels[order.status]}</span>
      </div>

      <!-- Progress bar -->
      <div class="relative mb-8 mt-6">
        <div class="h-2 bg-je-grey-light rounded-full">
          <div
            class="h-2 bg-je-green rounded-full transition-all duration-1000"
            style="width: {progressPercent}%"
          ></div>
        </div>

        <!-- Step markers -->
        <div class="flex justify-between absolute -top-3 w-full">
          {#each statusSteps as step, i}
            <div class="flex flex-col items-center">
              <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm border-2 transition-all
                {i <= currentStepIndex
                  ? 'bg-je-green border-je-green text-white'
                  : 'bg-white border-je-grey-border text-je-grey'}">
                {#if i < currentStepIndex}
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
                  </svg>
                {:else if i === currentStepIndex}
                  <span class="text-xs">{statusIcons[step]}</span>
                {:else}
                  <span class="text-xs">{i + 1}</span>
                {/if}
              </div>
              <span class="text-[10px] text-je-grey mt-2 whitespace-nowrap {i <= currentStepIndex ? 'font-semibold text-je-dark' : ''}">{statusLabels[step]}</span>
            </div>
          {/each}
        </div>
      </div>

      <!-- Status message -->
      {#if order.statusHistory && order.statusHistory.length > 0}
        <div class="bg-je-grey-light rounded-lg p-4 mt-4">
          <p class="text-sm font-medium text-je-dark">
            {order.statusHistory[order.statusHistory.length - 1].message}
          </p>
          <p class="text-xs text-je-grey mt-1">
            {new Date(order.statusHistory[order.statusHistory.length - 1].time).toLocaleTimeString()}
          </p>
        </div>
      {/if}

      {#if order.status !== 'delivered'}
        <p class="text-sm text-je-grey mt-3 flex items-center gap-2">
          <svg class="w-4 h-4 animate-pulse text-je-orange" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          Estimated delivery: {order.estimatedDelivery}
        </p>
      {/if}
    </div>

    <!-- Order details -->
    <div class="bg-white rounded-xl shadow-sm overflow-hidden mb-6">
      <div class="p-4 border-b border-je-grey-light">
        <h2 class="font-bold text-je-dark">Order details</h2>
      </div>

      {#each order.items as item}
        <div class="flex justify-between px-4 py-3 border-b border-je-grey-light text-sm">
          <div class="flex items-center gap-2">
            <span class="font-bold text-je-orange">{item.quantity}x</span>
            <span>{item.name}</span>
          </div>
          <span class="font-semibold">£{(item.price * item.quantity).toFixed(2)}</span>
        </div>
      {/each}

      <div class="p-4 space-y-2">
        <div class="flex justify-between text-sm">
          <span class="text-je-grey">Subtotal</span>
          <span>£{order.subtotal.toFixed(2)}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-je-grey">Delivery fee</span>
          <span>£{order.deliveryFee.toFixed(2)}</span>
        </div>
        <div class="flex justify-between font-bold text-base pt-2 border-t border-je-grey-light">
          <span>Total</span>
          <span>£{order.total.toFixed(2)}</span>
        </div>
      </div>
    </div>

    <!-- Delivery address -->
    <div class="bg-white rounded-xl shadow-sm p-4 mb-6">
      <h2 class="font-bold text-je-dark mb-2">Delivering to</h2>
      <p class="text-sm text-je-grey">{order.address}</p>
      <p class="text-sm text-je-grey mt-1">{order.customer.name} · {order.customer.phone}</p>
    </div>

    <!-- Activity log -->
    <div class="bg-white rounded-xl shadow-sm p-4">
      <h2 class="font-bold text-je-dark mb-4">Activity</h2>
      <div class="space-y-3">
        {#each [...order.statusHistory].reverse() as event}
          <div class="flex items-start gap-3">
            <div class="w-2 h-2 rounded-full bg-je-green mt-1.5 flex-shrink-0"></div>
            <div>
              <p class="text-sm text-je-dark">{event.message}</p>
              <p class="text-xs text-je-grey">{new Date(event.time).toLocaleTimeString()}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- New order button -->
    <div class="text-center mt-8">
      <a href="#/restaurants" class="inline-block bg-je-orange hover:bg-je-orange-dark text-white font-bold px-8 py-3 rounded-lg transition no-underline">
        Order again
      </a>
    </div>
  {/if}
</div>

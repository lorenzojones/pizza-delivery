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
    preparing: '🔪',
    ready: '🍹',
    delivering: '🏃',
    delivered: '🤙'
  };
  const statusLabels = {
    confirmed: 'Confirmed',
    preparing: 'Blending',
    ready: 'Ready!',
    delivering: 'On its way',
    delivered: 'Enjoy!'
  };

  $: currentStepIndex = order ? statusSteps.indexOf(order.status) : 0;
  $: progressPercent = order ? ((currentStepIndex) / (statusSteps.length - 1)) * 100 : 0;

  onMount(async () => {
    await loadOrder();
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
      <div class="h-8 bg-sand rounded w-1/3"></div>
      <div class="h-48 bg-sand rounded-2xl"></div>
      <div class="h-32 bg-sand rounded-2xl"></div>
    </div>
  {:else if error}
    <div class="bg-white rounded-2xl shadow-md p-12 text-center border border-sand">
      <span class="text-6xl mb-4 block">🥥</span>
      <h2 class="text-xl font-bold text-driftwood mb-2">Order not found</h2>
      <p class="text-driftwood-mid">{error}</p>
      <a href="#/" class="inline-block mt-4 text-ocean hover:text-ocean-dark font-bold no-underline">
        Back to Eno's
      </a>
    </div>
  {:else}
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center gap-3 mb-1">
        <h1 class="text-2xl font-surf text-ocean">Order #{order.id}</h1>
        <span class="text-2xl">{statusIcons[order.status]}</span>
      </div>
      <p class="text-driftwood-mid text-sm">from Eno's Shack</p>
    </div>

    <!-- Progress tracker -->
    <div class="bg-white rounded-2xl shadow-md p-6 mb-6 border border-sand">
      <div class="flex items-center justify-between mb-2">
        <h2 class="font-bold text-driftwood">Order Status</h2>
        <span class="text-sm font-bold text-ocean">{statusLabels[order.status]}</span>
      </div>

      <!-- Progress bar -->
      <div class="relative mb-8 mt-6">
        <div class="h-2 bg-sand rounded-full">
          <div
            class="h-2 bg-ocean rounded-full transition-all duration-1000"
            style="width: {progressPercent}%"
          ></div>
        </div>

        <!-- Step markers -->
        <div class="flex justify-between absolute -top-3 w-full">
          {#each statusSteps as step, i}
            <div class="flex flex-col items-center">
              <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm border-2 transition-all
                {i <= currentStepIndex
                  ? 'bg-ocean border-ocean text-white'
                  : 'bg-white border-sand-dark text-driftwood-mid'}">
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
              <span class="text-[10px] text-driftwood-mid mt-2 whitespace-nowrap {i <= currentStepIndex ? 'font-bold text-driftwood' : ''}">{statusLabels[step]}</span>
            </div>
          {/each}
        </div>
      </div>

      <!-- Status message -->
      {#if order.statusHistory && order.statusHistory.length > 0}
        <div class="bg-ocean-light rounded-xl p-4 mt-4">
          <p class="text-sm font-bold text-driftwood">
            {order.statusHistory[order.statusHistory.length - 1].message}
          </p>
          <p class="text-xs text-driftwood-mid mt-1">
            {new Date(order.statusHistory[order.statusHistory.length - 1].time).toLocaleTimeString()}
          </p>
        </div>
      {/if}

      {#if order.status !== 'delivered'}
        <p class="text-sm text-driftwood-mid mt-3 flex items-center gap-2">
          <span class="animate-pulse text-ocean">🍹</span>
          Estimated: {order.estimatedDelivery}
        </p>
      {/if}
    </div>

    <!-- Order details -->
    <div class="bg-white rounded-2xl shadow-md overflow-hidden mb-6 border border-sand">
      <div class="p-4 border-b border-sand">
        <h2 class="font-bold text-driftwood">Your Juices</h2>
      </div>

      {#each order.items as item}
        <div class="flex justify-between px-4 py-3 border-b border-sand-light text-sm">
          <div class="flex items-center gap-2">
            <span class="font-bold text-ocean">{item.quantity}x</span>
            <span>{item.name}</span>
          </div>
          <span class="font-bold">${item.price * item.quantity} MXN</span>
        </div>
      {/each}

      <div class="p-4 space-y-2">
        <div class="flex justify-between font-bold text-base pt-1">
          <span>Total</span>
          <span>${order.total} MXN</span>
        </div>
      </div>
    </div>

    <!-- Pickup info -->
    <div class="bg-white rounded-2xl shadow-md p-4 mb-6 border border-sand">
      <h2 class="font-bold text-driftwood mb-2">📍 Pickup at</h2>
      <p class="text-sm text-driftwood-mid">{order.address || "Eno's Shack, Calle Delfines 12, Sayulita"}</p>
      <p class="text-sm text-driftwood-mid mt-1">{order.customer.name} · {order.customer.phone}</p>
    </div>

    <!-- Activity log -->
    <div class="bg-white rounded-2xl shadow-md p-4 border border-sand">
      <h2 class="font-bold text-driftwood mb-4">Activity</h2>
      <div class="space-y-3">
        {#each [...order.statusHistory].reverse() as event}
          <div class="flex items-start gap-3">
            <div class="w-2 h-2 rounded-full bg-ocean mt-1.5 flex-shrink-0"></div>
            <div>
              <p class="text-sm text-driftwood">{event.message}</p>
              <p class="text-xs text-driftwood-mid">{new Date(event.time).toLocaleTimeString()}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- Order again -->
    <div class="text-center mt-8">
      <a href="#/menu" class="inline-block bg-sunset hover:bg-sunset-dark text-white font-bold px-8 py-3 rounded-xl transition no-underline">
        Order Again 🍹
      </a>
    </div>
  {/if}
</div>

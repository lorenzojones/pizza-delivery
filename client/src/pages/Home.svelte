<script>
  import { push } from 'svelte-spa-router';
</script>

<!-- Hero Section -->
<div class="relative overflow-hidden wave-divider" style="background: linear-gradient(135deg, #2B8A9E 0%, #1E6F80 50%, #F4845F 100%);">
  <div class="max-w-7xl mx-auto px-4 py-20 md:py-32 relative z-10">
    <div class="max-w-2xl">
      <p class="text-white/80 text-sm font-bold uppercase tracking-widest mb-3">Sayulita, Mexico</p>
      <h1 class="font-display text-5xl md:text-7xl text-white leading-tight mb-2">
        Eno's Shack
      </h1>
      <p class="font-surf text-xl md:text-2xl text-mango mb-4">
        Fresh Juice & Good Vibes 🤙
      </p>
      <p class="text-white/90 text-lg mb-8 max-w-lg">
        Hand-pressed fruit juices, smoothie bowls & tropical blends made fresh daily from local Nayarit fruit.
      </p>

      <button
        on:click={() => push('/menu')}
        class="bg-sunset hover:bg-sunset-dark text-white font-extrabold px-8 py-4 rounded-2xl shadow-lg transition text-lg inline-flex items-center gap-2"
      >
        See the Menu 🍹
      </button>
    </div>
  </div>

  <!-- Decorative elements -->
  <div class="absolute top-10 right-10 text-8xl opacity-20 animate-sway hidden md:block">🌴</div>
  <div class="absolute bottom-20 right-32 text-6xl opacity-15 animate-float hidden md:block">🏄</div>
</div>

<!-- Juice categories -->
<div class="max-w-7xl mx-auto px-4 py-12">
  <h2 class="text-2xl font-bold text-driftwood mb-2">What's your vibe?</h2>
  <p class="text-driftwood-mid mb-6">Choose your wave, we'll blend the rest.</p>
  <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
    {#each categories as cat}
      <button
        on:click={() => push('/menu')}
        class="flex flex-col items-center gap-3 p-5 bg-white rounded-2xl border border-sand hover:border-ocean/30 hover:shadow-lg transition group"
      >
        <span class="text-4xl group-hover:animate-float">{cat.emoji}</span>
        <span class="text-sm font-bold text-driftwood group-hover:text-ocean transition">{cat.name}</span>
      </button>
    {/each}
  </div>
</div>

<!-- About / Story section -->
<div class="bg-white py-12">
  <div class="max-w-7xl mx-auto px-4">
    <div class="grid md:grid-cols-2 gap-12 items-center">
      <div>
        <h2 class="font-surf text-3xl text-ocean mb-4">Born on the Beach</h2>
        <p class="text-driftwood-mid leading-relaxed mb-4">
          Eno's Shack started as a tiny stand on the Sayulita beach break. Every morning, we hit the waves at dawn, then head to the market for the freshest Nayarit fruit — mangos from La Cruz, papayas from San Pancho, pineapples straight from the campo.
        </p>
        <p class="text-driftwood-mid leading-relaxed mb-4">
          No artificial anything. No frozen concentrates. Just real fruit, real vibes, and that barefoot Sayulita feeling in every sip.
        </p>
        <div class="flex gap-4 mt-6">
          <div class="text-center">
            <span class="text-3xl font-extrabold text-ocean">100%</span>
            <p class="text-xs text-driftwood-mid mt-1">Natural</p>
          </div>
          <div class="text-center">
            <span class="text-3xl font-extrabold text-sunset">Local</span>
            <p class="text-xs text-driftwood-mid mt-1">Nayarit Fruit</p>
          </div>
          <div class="text-center">
            <span class="text-3xl font-extrabold text-palm">Daily</span>
            <p class="text-xs text-driftwood-mid mt-1">Fresh Pressed</p>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-gradient-to-br from-ocean-light to-ocean/10 rounded-2xl p-6 flex flex-col items-center justify-center text-center">
          <span class="text-5xl mb-2">🌊</span>
          <span class="font-surf text-ocean text-sm">Surf & Sip</span>
        </div>
        <div class="bg-gradient-to-br from-sunset-light to-sunset/10 rounded-2xl p-6 flex flex-col items-center justify-center text-center">
          <span class="text-5xl mb-2">🥭</span>
          <span class="font-surf text-sunset text-sm">Tropical Blends</span>
        </div>
        <div class="bg-gradient-to-br from-palm-light to-palm/10 rounded-2xl p-6 flex flex-col items-center justify-center text-center">
          <span class="text-5xl mb-2">🥬</span>
          <span class="font-surf text-palm text-sm">Green Cleanses</span>
        </div>
        <div class="bg-gradient-to-br from-papaya to-mango/20 rounded-2xl p-6 flex flex-col items-center justify-center text-center">
          <span class="text-5xl mb-2">🫐</span>
          <span class="font-surf text-driftwood text-sm">Açaí Bowls</span>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- How it works -->
<div class="max-w-7xl mx-auto px-4 py-12">
  <h2 class="text-2xl font-bold text-driftwood mb-8 text-center">How It Works</h2>
  <div class="grid md:grid-cols-3 gap-8">
    {#each steps as step, i}
      <div class="text-center">
        <div class="w-16 h-16 bg-ocean-light rounded-2xl flex items-center justify-center mx-auto mb-4">
          <span class="text-3xl">{step.emoji}</span>
        </div>
        <h3 class="font-bold text-driftwood mb-2">{step.title}</h3>
        <p class="text-driftwood-mid text-sm">{step.desc}</p>
      </div>
    {/each}
  </div>
</div>

<!-- Popular items teaser -->
<div class="bg-gradient-to-r from-ocean to-ocean-dark py-12">
  <div class="max-w-7xl mx-auto px-4">
    <div class="flex flex-col md:flex-row items-center gap-8">
      <div class="flex-1 text-white">
        <h2 class="font-surf text-3xl mb-3">Today's Best Sellers</h2>
        <p class="text-white/80 mb-6">Fresh pressed this morning. Get yours before they're gone.</p>
        <div class="flex flex-wrap gap-3">
          {#each bestSellers as item}
            <div class="bg-white/10 backdrop-blur-sm rounded-xl px-4 py-3 flex items-center gap-2">
              <span class="text-xl">{item.emoji}</span>
              <div>
                <p class="font-bold text-white text-sm">{item.name}</p>
                <p class="text-white/70 text-xs">${item.price} MXN</p>
              </div>
            </div>
          {/each}
        </div>
      </div>
      <button
        on:click={() => push('/menu')}
        class="bg-sunset hover:bg-sunset-dark text-white font-bold px-8 py-4 rounded-2xl shadow-lg transition text-lg whitespace-nowrap"
      >
        Order Now 🏄
      </button>
    </div>
  </div>
</div>

<!-- Location & Hours -->
<div class="max-w-7xl mx-auto px-4 py-12">
  <div class="bg-white rounded-2xl shadow-lg border border-sand overflow-hidden">
    <div class="grid md:grid-cols-2">
      <div class="p-8">
        <h2 class="font-surf text-2xl text-ocean mb-4">Find Us</h2>
        <div class="space-y-4">
          <div class="flex items-start gap-3">
            <span class="text-xl">📍</span>
            <div>
              <p class="font-bold text-driftwood">Calle Delfines 12</p>
              <p class="text-sm text-driftwood-mid">Sayulita, Nayarit, Mexico 63734</p>
              <p class="text-sm text-driftwood-mid">Just steps from the main beach break</p>
            </div>
          </div>
          <div class="flex items-start gap-3">
            <span class="text-xl">🕐</span>
            <div>
              <p class="font-bold text-driftwood">Open Daily</p>
              <p class="text-sm text-driftwood-mid">7:00 AM — 5:00 PM</p>
              <p class="text-sm text-palm font-semibold">Open now</p>
            </div>
          </div>
          <div class="flex items-start gap-3">
            <span class="text-xl">📱</span>
            <div>
              <p class="font-bold text-driftwood">Contact</p>
              <p class="text-sm text-driftwood-mid">Walk in, or order ahead for pickup!</p>
            </div>
          </div>
        </div>
      </div>
      <div class="bg-gradient-to-br from-ocean-light to-sand p-8 flex items-center justify-center">
        <div class="text-center">
          <span class="text-7xl block mb-4 animate-sway">🏖️</span>
          <p class="font-display text-2xl text-ocean">See you on the beach!</p>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Footer -->
<footer class="bg-driftwood text-white py-12">
  <div class="max-w-7xl mx-auto px-4">
    <div class="grid md:grid-cols-3 gap-8">
      <div>
        <h3 class="font-display text-2xl mb-3">Eno's Shack</h3>
        <p class="text-gray-400 text-sm">Fresh fruit juices & smoothie bowls in Sayulita, Mexico. Made with love, local fruit, and good vibes.</p>
      </div>
      <div>
        <h4 class="font-bold mb-3 text-sm">Quick Links</h4>
        <ul class="space-y-2 text-sm text-gray-400">
          <li><a href="#/menu" class="hover:text-mango transition no-underline text-gray-400">Full Menu</a></li>
          <li><a href="#/" class="hover:text-mango transition no-underline text-gray-400">Our Story</a></li>
          <li><a href="#/" class="hover:text-mango transition no-underline text-gray-400">Catering</a></li>
        </ul>
      </div>
      <div>
        <h4 class="font-bold mb-3 text-sm">Sayulita Life</h4>
        <ul class="space-y-2 text-sm text-gray-400">
          <li>Surf Reports</li>
          <li>Local Events</li>
          <li>Community Board</li>
        </ul>
      </div>
    </div>
    <div class="border-t border-gray-700 mt-8 pt-8 text-center text-gray-500 text-xs">
      © 2026 Eno's Shack, Sayulita. All rights reserved. 🌴🤙🏄
    </div>
  </div>
</footer>

<script context="module">
  const categories = [
    { name: 'Signature Juices', emoji: '🍹' },
    { name: 'Smoothie Bowls', emoji: '🥣' },
    { name: 'Smoothies', emoji: '🥤' },
    { name: 'Fresh Pressed', emoji: '🍊' },
    { name: 'Extras', emoji: '✨' },
  ];

  const steps = [
    { title: 'Pick your juice', desc: 'Browse our menu of fresh-pressed juices, smoothie bowls & tropical blends.', emoji: '🍹' },
    { title: 'We blend it fresh', desc: 'Every drink is made to order with fruit from the morning market.', emoji: '🔪' },
    { title: 'Grab & enjoy', desc: 'Pick up at the shack and enjoy with your toes in the sand.', emoji: '🏖️' },
  ];

  const bestSellers = [
    { name: 'Sayulita Sunrise', emoji: '🌅', price: 85 },
    { name: 'Pipeline Punch', emoji: '🏄', price: 90 },
    { name: 'Acai Shack Bowl', emoji: '🫐', price: 120 },
    { name: 'Ginger Shot', emoji: '🔥', price: 45 },
  ];
</script>

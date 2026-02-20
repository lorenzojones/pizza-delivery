<script>
  import { push } from 'svelte-spa-router';

  let postcode = '';

  function handleSearch() {
    push('/restaurants' + (postcode ? `?q=${encodeURIComponent(postcode)}` : ''));
  }
</script>

<!-- Hero Section -->
<div class="bg-je-orange">
  <div class="max-w-7xl mx-auto px-4 py-16 md:py-24">
    <div class="max-w-2xl">
      <h1 class="text-4xl md:text-6xl font-extrabold text-white leading-tight">
        Your favourite pizza,<br/>delivered.
      </h1>
      <p class="text-white/90 text-lg mt-4 mb-8">
        Order from the best local pizzerias with fast delivery to your door.
      </p>

      <!-- Search bar -->
      <form on:submit|preventDefault={handleSearch} class="flex gap-2">
        <div class="flex-1 relative">
          <svg class="w-5 h-5 text-je-grey absolute left-4 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          <input
            type="text"
            bind:value={postcode}
            placeholder="Enter your postcode"
            class="w-full pl-12 pr-4 py-4 rounded-xl text-base border-0 outline-none focus:ring-2 focus:ring-je-green shadow-lg"
          />
        </div>
        <button
          type="submit"
          class="bg-je-green hover:bg-je-green-dark text-white font-bold px-8 py-4 rounded-xl shadow-lg transition text-base"
        >
          Find food
        </button>
      </form>
    </div>
  </div>
</div>

<!-- Cuisine categories -->
<div class="max-w-7xl mx-auto px-4 py-12">
  <h2 class="text-2xl font-bold text-je-dark mb-6">Popular cuisines</h2>
  <div class="grid grid-cols-3 md:grid-cols-6 gap-4">
    {#each cuisines as cuisine}
      <button
        on:click={() => push(`/restaurants?cuisine=${cuisine.name}`)}
        class="flex flex-col items-center gap-3 p-4 bg-white rounded-xl hover:shadow-md transition group"
      >
        <span class="text-4xl">{cuisine.emoji}</span>
        <span class="text-sm font-medium text-je-dark group-hover:text-je-orange transition">{cuisine.name}</span>
      </button>
    {/each}
  </div>
</div>

<!-- How it works -->
<div class="bg-white py-12">
  <div class="max-w-7xl mx-auto px-4">
    <h2 class="text-2xl font-bold text-je-dark mb-8 text-center">How it works</h2>
    <div class="grid md:grid-cols-3 gap-8">
      {#each steps as step, i}
        <div class="text-center">
          <div class="w-16 h-16 bg-je-orange-light rounded-2xl flex items-center justify-center mx-auto mb-4">
            <span class="text-je-orange text-2xl font-extrabold">{i + 1}</span>
          </div>
          <h3 class="font-bold text-je-dark mb-2">{step.title}</h3>
          <p class="text-je-grey text-sm">{step.desc}</p>
        </div>
      {/each}
    </div>
  </div>
</div>

<!-- App download CTA -->
<div class="max-w-7xl mx-auto px-4 py-12">
  <div class="bg-je-dark rounded-2xl p-8 md:p-12 flex flex-col md:flex-row items-center gap-8">
    <div class="flex-1">
      <h2 class="text-2xl md:text-3xl font-bold text-white mb-3">Get the SliceNow app</h2>
      <p class="text-gray-400 mb-6">Order even faster with our mobile app. Track your delivery in real time.</p>
      <div class="flex gap-3">
        <button class="bg-white text-je-dark font-semibold px-6 py-3 rounded-lg text-sm hover:bg-gray-100 transition">
          App Store
        </button>
        <button class="bg-white text-je-dark font-semibold px-6 py-3 rounded-lg text-sm hover:bg-gray-100 transition">
          Google Play
        </button>
      </div>
    </div>
    <div class="text-8xl">📱🍕</div>
  </div>
</div>

<!-- Footer -->
<footer class="bg-je-dark text-white py-12 mt-8">
  <div class="max-w-7xl mx-auto px-4">
    <div class="grid md:grid-cols-4 gap-8">
      <div>
        <h3 class="font-bold text-lg mb-4">SliceNow</h3>
        <p class="text-gray-400 text-sm">The best pizza delivery in your area. Fresh, fast, delicious.</p>
      </div>
      <div>
        <h4 class="font-semibold mb-3 text-sm">Discover</h4>
        <ul class="space-y-2 text-sm text-gray-400">
          <li><a href="#/restaurants" class="hover:text-white transition no-underline text-gray-400">Restaurants</a></li>
          <li><a href="#/" class="hover:text-white transition no-underline text-gray-400">Top cuisines</a></li>
          <li><a href="#/" class="hover:text-white transition no-underline text-gray-400">Offers</a></li>
        </ul>
      </div>
      <div>
        <h4 class="font-semibold mb-3 text-sm">Company</h4>
        <ul class="space-y-2 text-sm text-gray-400">
          <li>About us</li>
          <li>Careers</li>
          <li>Blog</li>
        </ul>
      </div>
      <div>
        <h4 class="font-semibold mb-3 text-sm">Help</h4>
        <ul class="space-y-2 text-sm text-gray-400">
          <li>Contact</li>
          <li>FAQs</li>
          <li>Terms & Conditions</li>
        </ul>
      </div>
    </div>
    <div class="border-t border-gray-700 mt-8 pt-8 text-center text-gray-500 text-xs">
      © 2026 SliceNow. All rights reserved. This is a demo application.
    </div>
  </div>
</footer>

<script context="module">
  const cuisines = [
    { name: 'Pizza', emoji: '🍕' },
    { name: 'Italian', emoji: '🇮🇹' },
    { name: 'American', emoji: '🇺🇸' },
    { name: 'Kebab', emoji: '🥙' },
    { name: 'Gourmet', emoji: '👨‍🍳' },
    { name: 'Desserts', emoji: '🍰' },
  ];

  const steps = [
    { title: 'Choose a restaurant', desc: 'Browse local pizza restaurants near you and discover new favourites.' },
    { title: 'Pick your food', desc: 'Build your perfect order from the menu. Add sides, drinks, and desserts.' },
    { title: 'Fast delivery', desc: 'Sit back and relax. Track your order in real time until it arrives.' },
  ];
</script>

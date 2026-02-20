<script>
  import { onMount } from 'svelte';
  import { querystring } from 'svelte-spa-router';
  import { fetchRestaurants } from '../lib/api.js';
  import RestaurantCard from '../components/RestaurantCard.svelte';

  let restaurants = [];
  let loading = true;
  let searchQuery = '';
  let selectedCuisine = '';
  let selectedSort = '';

  const cuisineFilters = ['All', 'Pizza', 'Italian', 'American', 'Gourmet', 'Kebab'];
  const sortOptions = [
    { value: '', label: 'Recommended' },
    { value: 'rating', label: 'Top rated' },
    { value: 'deliveryTime', label: 'Fastest delivery' },
    { value: 'deliveryFee', label: 'Lowest delivery fee' },
  ];

  onMount(() => {
    const params = new URLSearchParams($querystring);
    searchQuery = params.get('q') || '';
    selectedCuisine = params.get('cuisine') || '';
    loadRestaurants();
  });

  async function loadRestaurants() {
    loading = true;
    try {
      const params = {};
      if (searchQuery) params.q = searchQuery;
      if (selectedCuisine && selectedCuisine !== 'All') params.cuisine = selectedCuisine;
      if (selectedSort) params.sort = selectedSort;
      restaurants = await fetchRestaurants(params);
    } catch (err) {
      console.error(err);
    }
    loading = false;
  }

  function handleSearch(e) {
    e.preventDefault();
    loadRestaurants();
  }

  function selectCuisine(c) {
    selectedCuisine = c === 'All' ? '' : c;
    loadRestaurants();
  }

  function changeSort(e) {
    selectedSort = e.target.value;
    loadRestaurants();
  }
</script>

<div class="max-w-7xl mx-auto px-4 py-6">
  <!-- Search & Filter bar -->
  <div class="bg-white rounded-xl shadow-sm p-4 mb-6">
    <form on:submit={handleSearch} class="flex gap-3 mb-4">
      <div class="flex-1 relative">
        <svg class="w-5 h-5 text-je-grey absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Search for restaurants or dishes..."
          class="w-full pl-10 pr-4 py-3 rounded-lg bg-je-grey-light border border-je-grey-border focus:border-je-orange focus:ring-1 focus:ring-je-orange outline-none text-sm"
        />
      </div>
      <button type="submit" class="bg-je-orange hover:bg-je-orange-dark text-white font-semibold px-6 py-3 rounded-lg text-sm transition">
        Search
      </button>
    </form>

    <!-- Cuisine filter pills -->
    <div class="flex gap-2 overflow-x-auto hide-scrollbar pb-1">
      {#each cuisineFilters as cuisine}
        <button
          on:click={() => selectCuisine(cuisine)}
          class="px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition
            {(selectedCuisine === cuisine || (!selectedCuisine && cuisine === 'All'))
              ? 'bg-je-orange text-white'
              : 'bg-je-grey-light text-je-grey hover:bg-je-grey-border'}"
        >
          {cuisine}
        </button>
      {/each}
    </div>
  </div>

  <!-- Sort & count bar -->
  <div class="flex items-center justify-between mb-4">
    <p class="text-sm text-je-grey">
      {#if loading}
        Searching...
      {:else}
        <span class="font-semibold text-je-dark">{restaurants.length}</span> restaurants available
      {/if}
    </p>
    <select on:change={changeSort} value={selectedSort} class="text-sm border border-je-grey-border rounded-lg px-3 py-2 bg-white outline-none focus:border-je-orange">
      {#each sortOptions as opt}
        <option value={opt.value}>{opt.label}</option>
      {/each}
    </select>
  </div>

  <!-- Restaurant grid -->
  {#if loading}
    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-5">
      {#each Array(6) as _}
        <div class="bg-white rounded-xl overflow-hidden shadow-sm animate-pulse">
          <div class="h-40 bg-gray-200"></div>
          <div class="p-4 space-y-3">
            <div class="h-5 bg-gray-200 rounded w-3/4"></div>
            <div class="h-3 bg-gray-200 rounded w-1/2"></div>
            <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if restaurants.length === 0}
    <div class="text-center py-16">
      <span class="text-6xl mb-4 block">🔍</span>
      <h2 class="text-xl font-bold text-je-dark mb-2">No restaurants found</h2>
      <p class="text-je-grey">Try a different search or browse all restaurants</p>
      <button on:click={() => { searchQuery = ''; selectedCuisine = ''; loadRestaurants(); }} class="mt-4 text-je-orange hover:text-je-orange-dark font-semibold text-sm">
        Clear filters
      </button>
    </div>
  {:else}
    <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-5">
      {#each restaurants as restaurant}
        <RestaurantCard {restaurant} />
      {/each}
    </div>
  {/if}
</div>

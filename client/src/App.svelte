<script>
  import Router from 'svelte-spa-router';
  import Header from './components/Header.svelte';
  import Home from './pages/Home.svelte';
  import RestaurantList from './pages/RestaurantList.svelte';
  import RestaurantDetail from './pages/RestaurantDetail.svelte';
  import Checkout from './pages/Checkout.svelte';
  import OrderTracking from './pages/OrderTracking.svelte';

  const routes = {
    '/': Home,
    '/menu': RestaurantDetail,
    '/restaurants': RestaurantList,
    '/restaurant/:id': RestaurantDetail,
    '/checkout': Checkout,
    '/order/:id': OrderTracking,
  };
</script>

<div class="min-h-screen bg-coconut">
  <Header />
  <Router {routes} />
</div>

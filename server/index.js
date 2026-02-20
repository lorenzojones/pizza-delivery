const express = require('express');
const cors = require('cors');
const path = require('path');
const { v4: uuidv4 } = require('uuid');
const restaurants = require('./data/restaurants.json');

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());

// Serve static frontend files
app.use(express.static(path.join(__dirname, '..', 'client', 'public')));

// In-memory orders store
const orders = new Map();

// --- API Routes ---

// Get all restaurants (supports search & filter)
app.get('/api/restaurants', (req, res) => {
  const { q, cuisine, sort, minRating } = req.query;
  let results = [...restaurants];

  // Text search
  if (q) {
    const query = q.toLowerCase();
    results = results.filter(r =>
      r.name.toLowerCase().includes(query) ||
      r.cuisine.some(c => c.toLowerCase().includes(query)) ||
      r.menu.some(cat => cat.items.some(item => item.name.toLowerCase().includes(query)))
    );
  }

  // Filter by cuisine
  if (cuisine) {
    results = results.filter(r =>
      r.cuisine.some(c => c.toLowerCase() === cuisine.toLowerCase())
    );
  }

  // Filter by minimum rating
  if (minRating) {
    results = results.filter(r => r.rating >= parseFloat(minRating));
  }

  // Sort
  if (sort === 'rating') {
    results.sort((a, b) => b.rating - a.rating);
  } else if (sort === 'deliveryFee') {
    results.sort((a, b) => a.deliveryFee - b.deliveryFee);
  } else if (sort === 'deliveryTime') {
    results.sort((a, b) => {
      const aTime = parseInt(a.deliveryTime);
      const bTime = parseInt(b.deliveryTime);
      return aTime - bTime;
    });
  } else {
    // Default: promoted first, then by rating
    results.sort((a, b) => {
      if (a.promoted && !b.promoted) return -1;
      if (!a.promoted && b.promoted) return 1;
      return b.rating - a.rating;
    });
  }

  // Return summary data (no menu)
  const summaries = results.map(({ menu, ...rest }) => rest);
  res.json(summaries);
});

// Get single restaurant with full menu
app.get('/api/restaurants/:id', (req, res) => {
  const restaurant = restaurants.find(r => r.id === req.params.id);
  if (!restaurant) {
    return res.status(404).json({ error: 'Restaurant not found' });
  }
  res.json(restaurant);
});

// Place an order
app.post('/api/orders', (req, res) => {
  const { restaurantId, items, customer, address, paymentMethod } = req.body;

  if (!restaurantId || !items || !items.length || !customer || !address) {
    return res.status(400).json({ error: 'Missing required fields' });
  }

  const restaurant = restaurants.find(r => r.id === restaurantId);
  if (!restaurant) {
    return res.status(404).json({ error: 'Restaurant not found' });
  }

  // Calculate totals
  const subtotal = items.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  const deliveryFee = restaurant.deliveryFee;
  const total = subtotal + deliveryFee;

  if (subtotal < restaurant.minOrder) {
    return res.status(400).json({
      error: `Minimum order is £${restaurant.minOrder.toFixed(2)}`
    });
  }

  const order = {
    id: uuidv4().slice(0, 8).toUpperCase(),
    restaurantId,
    restaurantName: restaurant.name,
    items,
    customer,
    address,
    paymentMethod: paymentMethod || 'card',
    subtotal,
    deliveryFee,
    total,
    status: 'confirmed',
    statusHistory: [
      { status: 'confirmed', time: new Date().toISOString(), message: 'Order confirmed' }
    ],
    estimatedDelivery: restaurant.deliveryTime,
    createdAt: new Date().toISOString()
  };

  orders.set(order.id, order);

  // Simulate order progress
  simulateOrderProgress(order.id);

  res.status(201).json(order);
});

// Get order status
app.get('/api/orders/:id', (req, res) => {
  const order = orders.get(req.params.id.toUpperCase());
  if (!order) {
    return res.status(404).json({ error: 'Order not found' });
  }
  res.json(order);
});

// Simulate order lifecycle
function simulateOrderProgress(orderId) {
  const stages = [
    { delay: 15000, status: 'preparing', message: 'Restaurant is preparing your order' },
    { delay: 45000, status: 'ready', message: 'Your order is ready for pickup' },
    { delay: 60000, status: 'delivering', message: 'Your rider is on the way' },
    { delay: 90000, status: 'delivered', message: 'Your order has been delivered. Enjoy!' }
  ];

  stages.forEach(({ delay, status, message }) => {
    setTimeout(() => {
      const order = orders.get(orderId);
      if (order) {
        order.status = status;
        order.statusHistory.push({
          status,
          time: new Date().toISOString(),
          message
        });
      }
    }, delay);
  });
}

// SPA fallback - serve index.html for all non-API routes
app.get('/{*splat}', (req, res) => {
  res.sendFile(path.join(__dirname, '..', 'client', 'public', 'index.html'));
});

app.listen(PORT, () => {
  console.log(`Pizza Delivery API running on http://localhost:${PORT}`);
});

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

  if (q) {
    const query = q.toLowerCase();
    results = results.filter(r =>
      r.name.toLowerCase().includes(query) ||
      r.cuisine.some(c => c.toLowerCase().includes(query)) ||
      r.menu.some(cat => cat.items.some(item => item.name.toLowerCase().includes(query)))
    );
  }

  if (cuisine) {
    results = results.filter(r =>
      r.cuisine.some(c => c.toLowerCase() === cuisine.toLowerCase())
    );
  }

  if (minRating) {
    results = results.filter(r => r.rating >= parseFloat(minRating));
  }

  if (sort === 'rating') {
    results.sort((a, b) => b.rating - a.rating);
  } else if (sort === 'price') {
    results.sort((a, b) => {
      const aAvg = a.menu.flatMap(c => c.items).reduce((s, i) => s + i.price, 0) / a.menu.flatMap(c => c.items).length;
      const bAvg = b.menu.flatMap(c => c.items).reduce((s, i) => s + i.price, 0) / b.menu.flatMap(c => c.items).length;
      return aAvg - bAvg;
    });
  }

  const summaries = results.map(({ menu, ...rest }) => rest);
  res.json(summaries);
});

// Get single restaurant with full menu
app.get('/api/restaurants/:id', (req, res) => {
  const restaurant = restaurants.find(r => r.id === req.params.id);
  if (!restaurant) {
    return res.status(404).json({ error: 'Not found' });
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

  const subtotal = items.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  const total = subtotal;

  const order = {
    id: uuidv4().slice(0, 8).toUpperCase(),
    restaurantId,
    restaurantName: restaurant.name,
    items,
    customer,
    address,
    paymentMethod: paymentMethod || 'cash',
    subtotal,
    deliveryFee: 0,
    total,
    status: 'confirmed',
    statusHistory: [
      { status: 'confirmed', time: new Date().toISOString(), message: 'Order confirmed! We\'re on it 🤙' }
    ],
    estimatedDelivery: '10-15 min',
    createdAt: new Date().toISOString()
  };

  orders.set(order.id, order);
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
    { delay: 10000, status: 'preparing', message: 'Blending your juice fresh! 🍹' },
    { delay: 25000, status: 'ready', message: 'Your order is ready for pickup! 🤙' },
    { delay: 40000, status: 'delivering', message: 'Heading your way on the beach! 🏃' },
    { delay: 55000, status: 'delivered', message: 'Enjoy your juice! See you next time 🌊' }
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

// SPA fallback
app.get('/{*splat}', (req, res) => {
  res.sendFile(path.join(__dirname, '..', 'client', 'public', 'index.html'));
});

app.listen(PORT, () => {
  console.log(`Eno's Shack Juice Bar running on http://localhost:${PORT}`);
});

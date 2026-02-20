const API_BASE = '/api';

export async function fetchRestaurants(params = {}) {
  const query = new URLSearchParams();
  if (params.q) query.set('q', params.q);
  if (params.cuisine) query.set('cuisine', params.cuisine);
  if (params.sort) query.set('sort', params.sort);
  if (params.minRating) query.set('minRating', params.minRating);

  const qs = query.toString();
  const res = await fetch(`${API_BASE}/restaurants${qs ? '?' + qs : ''}`);
  if (!res.ok) throw new Error('Failed to fetch restaurants');
  return res.json();
}

export async function fetchRestaurant(id) {
  const res = await fetch(`${API_BASE}/restaurants/${id}`);
  if (!res.ok) throw new Error('Restaurant not found');
  return res.json();
}

export async function placeOrder(orderData) {
  const res = await fetch(`${API_BASE}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(orderData)
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Failed to place order');
  return data;
}

export async function fetchOrder(id) {
  const res = await fetch(`${API_BASE}/orders/${id}`);
  if (!res.ok) throw new Error('Order not found');
  return res.json();
}

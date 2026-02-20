import { writable, derived } from 'svelte/store';

function createCart() {
  const { subscribe, set, update } = writable({
    restaurantId: null,
    restaurantName: '',
    deliveryFee: 0,
    minOrder: 0,
    items: []
  });

  return {
    subscribe,
    addItem: (restaurant, item) => {
      update(cart => {
        // If switching restaurants, clear the cart
        if (cart.restaurantId && cart.restaurantId !== restaurant.id) {
          return {
            restaurantId: restaurant.id,
            restaurantName: restaurant.name,
            deliveryFee: restaurant.deliveryFee,
            minOrder: restaurant.minOrder,
            items: [{ ...item, quantity: 1 }]
          };
        }

        // Same restaurant - check if item exists
        const existing = cart.items.find(i => i.id === item.id);
        if (existing) {
          existing.quantity += 1;
          return { ...cart, items: [...cart.items] };
        }

        return {
          ...cart,
          restaurantId: restaurant.id,
          restaurantName: restaurant.name,
          deliveryFee: restaurant.deliveryFee,
          minOrder: restaurant.minOrder,
          items: [...cart.items, { ...item, quantity: 1 }]
        };
      });
    },
    removeItem: (itemId) => {
      update(cart => {
        const items = cart.items.filter(i => i.id !== itemId);
        if (items.length === 0) {
          return { restaurantId: null, restaurantName: '', deliveryFee: 0, minOrder: 0, items: [] };
        }
        return { ...cart, items };
      });
    },
    updateQuantity: (itemId, quantity) => {
      update(cart => {
        if (quantity <= 0) {
          const items = cart.items.filter(i => i.id !== itemId);
          if (items.length === 0) {
            return { restaurantId: null, restaurantName: '', deliveryFee: 0, minOrder: 0, items: [] };
          }
          return { ...cart, items };
        }
        const items = cart.items.map(i =>
          i.id === itemId ? { ...i, quantity } : i
        );
        return { ...cart, items };
      });
    },
    clear: () => {
      set({ restaurantId: null, restaurantName: '', deliveryFee: 0, minOrder: 0, items: [] });
    }
  };
}

export const cart = createCart();

export const cartItemCount = derived(cart, $cart =>
  $cart.items.reduce((sum, item) => sum + item.quantity, 0)
);

export const cartSubtotal = derived(cart, $cart =>
  $cart.items.reduce((sum, item) => sum + (item.price * item.quantity), 0)
);

export const cartTotal = derived([cart, cartSubtotal], ([$cart, $subtotal]) =>
  $subtotal + $cart.deliveryFee
);

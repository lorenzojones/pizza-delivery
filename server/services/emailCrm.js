const contacts = new Map();
const emailLog = [];

function addOrUpdateContact(email, data) {
  const existing = contacts.get(email);
  if (existing) {
    Object.assign(existing, data, { updatedAt: new Date().toISOString() });
    return existing;
  }

  const contact = {
    email,
    name: data.name || '',
    phone: data.phone || '',
    marketingOptIn: data.marketingOptIn ?? false,
    orderCount: 0,
    totalSpent: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };
  contacts.set(email, contact);
  return contact;
}

function getContact(email) {
  return contacts.get(email) || null;
}

function getAllContacts() {
  return Array.from(contacts.values());
}

function recordOrder(email, orderTotal) {
  const contact = contacts.get(email);
  if (contact) {
    contact.orderCount += 1;
    contact.totalSpent += orderTotal;
    contact.updatedAt = new Date().toISOString();
  }
}

function updatePreferences(email, preferences) {
  const contact = contacts.get(email);
  if (!contact) return null;
  if (preferences.marketingOptIn !== undefined) {
    contact.marketingOptIn = preferences.marketingOptIn;
  }
  contact.updatedAt = new Date().toISOString();
  return contact;
}

function sendEmail(to, subject, body) {
  const entry = {
    id: emailLog.length + 1,
    to,
    subject,
    body,
    sentAt: new Date().toISOString(),
  };
  emailLog.push(entry);
  console.log(`[Email CRM] Sent to ${to}: "${subject}"`);
  return entry;
}

function sendOrderConfirmation(order) {
  const itemLines = order.items
    .map(i => `  ${i.quantity}x ${i.name} - £${(i.price * i.quantity).toFixed(2)}`)
    .join('\n');

  const body = [
    `Hi ${order.customer.name},`,
    '',
    `Thanks for your order from ${order.restaurantName}!`,
    '',
    `Order #${order.id}`,
    '---',
    itemLines,
    '---',
    `Subtotal: £${order.subtotal.toFixed(2)}`,
    `Delivery: £${order.deliveryFee.toFixed(2)}`,
    `Total: £${order.total.toFixed(2)}`,
    '',
    `Delivering to: ${order.address}`,
    `Estimated delivery: ${order.estimatedDelivery}`,
    '',
    'Track your order at: /order/' + order.id,
  ].join('\n');

  return sendEmail(order.customer.email, `Order confirmed - #${order.id}`, body);
}

function sendWelcomeEmail(contact) {
  const body = [
    `Hi ${contact.name},`,
    '',
    'Welcome to SliceNow!',
    '',
    "You've signed up for our newsletter. We'll keep you posted on exclusive deals, new restaurants, and tasty offers.",
    '',
    'Hungry? Browse our restaurants and place your first order today!',
  ].join('\n');

  return sendEmail(contact.email, 'Welcome to SliceNow!', body);
}

function getEmailLog(email) {
  if (email) {
    return emailLog.filter(e => e.to === email);
  }
  return [...emailLog];
}

module.exports = {
  addOrUpdateContact,
  getContact,
  getAllContacts,
  recordOrder,
  updatePreferences,
  sendEmail,
  sendOrderConfirmation,
  sendWelcomeEmail,
  getEmailLog,
};

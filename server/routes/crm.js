const express = require('express');
const crm = require('../services/emailCrm');

const router = express.Router();

router.post('/subscribe', (req, res) => {
  const { email, name } = req.body;
  if (!email) {
    return res.status(400).json({ error: 'Email is required' });
  }

  const existing = crm.getContact(email);
  const contact = crm.addOrUpdateContact(email, {
    name: name || '',
    marketingOptIn: true,
  });

  if (!existing) {
    crm.sendWelcomeEmail(contact);
  }

  res.status(existing ? 200 : 201).json({ subscribed: true, contact });
});

router.post('/unsubscribe', (req, res) => {
  const { email } = req.body;
  if (!email) {
    return res.status(400).json({ error: 'Email is required' });
  }

  const contact = crm.updatePreferences(email, { marketingOptIn: false });
  if (!contact) {
    return res.status(404).json({ error: 'Contact not found' });
  }

  res.json({ unsubscribed: true, contact });
});

router.get('/contacts', (_req, res) => {
  res.json(crm.getAllContacts());
});

router.get('/contacts/:email', (req, res) => {
  const contact = crm.getContact(req.params.email);
  if (!contact) {
    return res.status(404).json({ error: 'Contact not found' });
  }
  res.json(contact);
});

router.put('/contacts/:email/preferences', (req, res) => {
  const { marketingOptIn } = req.body;
  const contact = crm.updatePreferences(req.params.email, { marketingOptIn });
  if (!contact) {
    return res.status(404).json({ error: 'Contact not found' });
  }
  res.json(contact);
});

router.get('/emails', (req, res) => {
  const { email } = req.query;
  res.json(crm.getEmailLog(email));
});

module.exports = router;

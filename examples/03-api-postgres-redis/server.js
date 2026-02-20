const express = require('express');
const { Client } = require('pg');
const redis = require('redis');

const app = express();
const port = process.env.PORT || 4000;
const dbUrl = process.env.DATABASE_URL || 'postgresql://app:secret@postgres:5432/appdb?sslmode=disable';
const redisUrl = process.env.REDIS_URL || 'redis://redis:6379';

app.use(express.json());

const pg = new Client({ connectionString: dbUrl });
const redisClient = redis.createClient({ url: redisUrl });

async function init() {
  await pg.connect();
  await redisClient.connect();
  await pg.query(`
    CREATE TABLE IF NOT EXISTS items (
      id SERIAL PRIMARY KEY,
      name TEXT NOT NULL,
      created_at TIMESTAMPTZ DEFAULT NOW()
    )
  `);
}

init().catch((e) => console.error('Init error:', e.message));

app.get('/health', (req, res) => res.json({ status: 'ok' }));

app.get('/items', async (req, res) => {
  try {
    const r = await pg.query('SELECT id, name, created_at FROM items ORDER BY id');
    res.json(r.rows);
  } catch (e) {
    res.status(500).json({ error: e.message });
  }
});

app.post('/items', async (req, res) => {
  const name = req.body?.name || 'item';
  try {
    const r = await pg.query('INSERT INTO items (name) VALUES ($1) RETURNING id, name, created_at', [name]);
    const cacheKey = 'items:count';
    await redisClient.incr(cacheKey);
    res.status(201).json(r.rows[0]);
  } catch (e) {
    res.status(500).json({ error: e.message });
  }
});

app.get('/stats', async (req, res) => {
  try {
    const count = await redisClient.get('items:count');
    const r = await pg.query('SELECT COUNT(*) FROM items');
    res.json({
      redis_count: count ? parseInt(count, 10) : 0,
      postgres_count: parseInt(r.rows[0].count, 10),
    });
  } catch (e) {
    res.status(500).json({ error: e.message });
  }
});

app.listen(port, () => console.log('API listening on port', port));

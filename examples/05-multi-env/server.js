const express = require('express');
const redis = require('redis');

const app = express();
const port = process.env.PORT || 3000;
const redisUrl = process.env.REDIS_URL || 'redis://redis:6379';

const client = redis.createClient({ url: redisUrl });
client.connect().catch(() => {});

const KEY = 'counter';

app.get('/', async (req, res) => {
  try {
    const n = await client.incr(KEY);
    const count = await client.get(KEY);
    res.type('html');
    res.send(`
      <!DOCTYPE html>
      <html><head><title>Counter</title></head>
      <body>
        <h1>Counter (Redis)</h1>
        <p>Count: <strong>${count}</strong></p>
        <p>Refresh the page to increment.</p>
      </body></html>
    `);
  } catch (e) {
    res.status(500).send('Redis error: ' + e.message);
  }
});

app.get('/health', (req, res) => res.send('ok'));

app.listen(port, () => console.log('Counter listening on port', port));

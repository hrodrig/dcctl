import os
from flask import Flask
from redis import Redis

app = Flask(__name__)
redis_url = os.environ.get("REDIS_URL", "redis://redis:6379/0")
redis = Redis.from_url(redis_url)
port = int(os.environ.get("FLASK_PORT", "8000"))

@app.route('/')
def hello():
    redis.incr('hits')
    counter = redis.get('hits') or b'0'
    return f"This webpage has been viewed {counter.decode()} time(s)"

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=port, debug=True)

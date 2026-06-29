# Spark Python SDK

Python client library for the Spark platform API.

## Installation

```bash
pip install spark-sdk
```

## Usage

```python
from spark import Spark

client = Spark(access_token="your_token")

# Get current user
user = client.me()

# List live streams
streams = client.streams.list(is_live=True)

# Search
results = client.search.query("gaming")
```
"""

from .client import Spark

__all__ = ["Spark"]

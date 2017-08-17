import base64
import datetime
import hashlib
import hmac
import json

import requests
encoded_hashed_body = base64.b64encode(hashlib.sha256("").digest())
print encoded_hashed_body

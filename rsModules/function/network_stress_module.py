import time
import requests

def send_large_post_request(network_mode, net_url, net_port):
    # Create large data
    if network_mode=="preprocess":
        large_data = 'x' * 10**8  # 10MB data
    else:
        large_data = 'x' * 10**4
    url = f"{net_url}:{net_port}/post"
    response = requests.post(url, data=large_data)
    print(f"Sent {len(large_data)} bytes to {url}, received {len(response.content)} bytes")
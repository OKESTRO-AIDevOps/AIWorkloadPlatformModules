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

def send_large_get_request(net_url, net_port):
    url = f"{net_url}:{net_port}/get"
    response = requests.get(url)
    print(f"Received {len(response.content)} bytes from {url}")

def network_stress_func(duration, net_url:str='http://localhost', net_port:str='5001', network_mode:str="preprocess"):
    print(f"Network I/O Test mode : {network_mode}")
    print(f"Generating network I/O load. (Duration: {duration} s)")
    end_time = time.time() + duration
    while time.time() < end_time:
        send_large_post_request(network_mode, net_url, net_port)
        send_large_get_request(net_url, net_port)
        time.sleep(1)%
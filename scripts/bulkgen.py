import random
import time
import struct
import sys
import os
import argparse

parser = argparse.ArgumentParser(description="Generate binary observation data for bulk upload")
parser.add_argument("node_id", help="Node ID")
parser.add_argument("password", help="Node password")
parser.add_argument("quantity_ids", type=int, nargs="+", help="Quantity IDs to generate")
parser.add_argument("-n", "--num", type=int, default=1000, help="Number of observations (default: 1000)")
parser.add_argument("-o", "--output", default=None, help="Output file (default: bulk_<node_id>.bin)")
parser.add_argument("--upload", action="store_true", help="Upload to bulk service after generating")
args = parser.parse_args()

longitude = round(random.uniform(-180.0, 180.0), 6)
latitude = round(random.uniform(-90.0, 90.0), 6)
script_dir = os.path.dirname(os.path.abspath(__file__))
project_dir = os.path.dirname(script_dir)
outfile = args.output or os.path.join(project_dir, "output", f"bulk_{args.node_id}.bin")


def make_observation():
    quantity_id = random.choice(args.quantity_ids)
    value = round(random.uniform(-50.0, 50.0), 2)
    timestamp = time.time() - random.uniform(0, 86400)
    return struct.pack('!Bfdff', quantity_id, value, timestamp, longitude, latitude)


def generate():
    data = b''
    for _ in range(args.num):
        data += make_observation()
    os.makedirs(os.path.dirname(outfile), exist_ok=True)
    with open(outfile, 'wb') as f:
        f.write(data)
    print(f"Wrote {args.num} observations ({len(data)} bytes) to {outfile}")
    return outfile


def upload(filepath):
    import urllib.request

    url = "http://localhost:8081/upload"
    boundary = "----BulkTransfer"
    body = (
        f"--{boundary}\r\n"
        f'Content-Disposition: form-data; name="node_id"\r\n\r\n'
        f"{args.node_id}\r\n"
        f"--{boundary}\r\n"
        f'Content-Disposition: form-data; name="file"; filename="{os.path.basename(filepath)}"\r\n'
        f"Content-Type: application/octet-stream\r\n\r\n"
    ).encode()
    with open(filepath, "rb") as f:
        body += f.read()
    body += f"\r\n--{boundary}--\r\n".encode()

    req = urllib.request.Request(url, data=body)
    req.add_header("Content-Type", f"multipart/form-data; boundary={boundary}")
    req.add_header("Authorization", f"Bearer {args.password}")

    resp = urllib.request.urlopen(req)
    print(resp.read().decode())


if __name__ == "__main__":
    fpath = generate()
    if args.upload:
        upload(fpath)

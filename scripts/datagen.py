import random
import time
import struct
import sys

from paho.mqtt import client as mqtt_client
import os
from os.path import join, dirname
from dotenv import load_dotenv

dotenv_path = join(dirname(dirname(__file__)), '.env')
load_dotenv(dotenv_path)

broker = "localhost"
port = 1883
topic = 'o/' + sys.argv[1]

client_id = f'publish-{random.randint(0, 1000)}'
username = sys.argv[1]
password = sys.argv[2]
max_sensor = int(sys.argv[3])


def make_observation():
    quantity_id = random.randint(1, max_sensor)
    value = round(random.uniform(-50.0, 50.0), 2)
    timestamp = time.time()
    longitude = round(random.uniform(-180.0, 180.0), 6)
    latitude = round(random.uniform(-90.0, 90.0), 6)
    return struct.pack('!Bffff', quantity_id, value, timestamp, longitude, latitude)


def connect_mqtt():
    def on_connect(client, userdata, flags, reason_code, properties):
        if reason_code == 0:
            print("Connected to MQTT Broker!")
        else:
            print(f"Failed to connect, return code {reason_code}")

    client = mqtt_client.Client(
        client_id=client_id,
        callback_api_version=mqtt_client.CallbackAPIVersion.VERSION2,
    )
    client.username_pw_set(username, password)
    client.on_connect = on_connect
    client.connect(broker, port)
    return client


def publish(client):
    msg_count = 1
    while True:
        time.sleep(1)
        payload = make_observation()
        result = client.publish(topic, payload)
        status = result.rc
        if status == 0:
            print(f"[{msg_count}] Sent {len(payload)} bytes to `{topic}`")
        else:
            print(f"Failed to send message to topic {topic}")
        msg_count += 1


def run():
    client = connect_mqtt()
    client.loop_start()
    try:
        publish(client)
    except KeyboardInterrupt:
        print("\nStopping...")
    finally:
        client.loop_stop()
        client.disconnect()


if __name__ == '__main__':
    run()

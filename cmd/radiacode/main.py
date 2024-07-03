import os
import time

import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS
from radiacode import RealTimeData, RadiaCode

org = "reactor"
bucket = "dlog"


def main():
    rc_conn = RadiaCode()

    client = influxdb_client.InfluxDBClient(
        url=os.environ.get("OPENREACTOR_INFLUXDB_URL"),
        token=os.environ.get("OPENREACTOR_INFLUXDB_TOKEN"),
        org=org
    )
    write_api = client.write_api(write_options=SYNCHRONOUS)

    while True:
        data_buf = rc_conn.data_buf()

        last = None
        for v in data_buf:
            if isinstance(v, RealTimeData):
                if last is None or last.dt < v.dt:
                    last = v

        if last is None:
            continue

        p = influxdb_client.Point("radiacode").field("cps", last.count_rate)
        write_api.write(bucket=bucket, org=org, record=p)

        time.sleep(1)


if __name__ == '__main__':
    main()

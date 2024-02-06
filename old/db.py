import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS

org = "reactor"
bucket = "test"

client = influxdb_client.InfluxDBClient(
    url="https://influxdb-reactor.westland.as34553.net",
    token="NBoL8jI4cPbziQyLSqvyzHcH3rVMex3359-OnWIATHIMScwP8eJDR2ysD7zW9ijWyRJ9i8C9cXwN6e4rJr3LOQ==",
    org=org
)

write_api = client.write_api(write_options=SYNCHRONOUS)


def write(point):
    write_api.write(bucket=bucket, org=org, record=point)

# for value in range(5):
#     point = (
#         Point("measurement1")
#         .tag("tagname1", "tagvalue1")
#         .field("field1", value)
#     )
#     print("Writing point")
#     write_api.write(bucket=bucket, org=org, record=point)
#     time.sleep(0.25)

#
# query_api = client.query_api()
#
# query = f"""from(bucket: "{bucket}")
#  |> range(start: -10m)
#  |> filter(fn: (r) => r._measurement == "measurement1")"""
# tables = query_api.query(query, org=org)
#
# for table in tables:
#     for record in table.records:
#         print(record.get_start())

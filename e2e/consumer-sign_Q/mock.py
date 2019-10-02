# Core-stack | Worker Nonce :: Consumer-sign_Q/mock
# Simulates a consumer connected to worker nonce's output (Sign Q perspective)
# For testing purposes
# @author Guillaume Lethuillier <guillaume.lethuillier@consensys.net>

from kafka import KafkaConsumer
from time import sleep
import trace_pb2
from importlib import import_module


KAFKA_URL = 'localhost:9092'
WORKER_NONCE_OUT_TOPIC = 'nonce-out'

#####################
# Preliminary check #
#####################

NUMBER_OF_RETRIES = 10

assert KAFKA_URL, 'Kafka address is not set'

print("Kafka address: {0}".format(KAFKA_URL))

retries = 0
while True:

    if retries >= NUMBER_OF_RETRIES:
        print()
        sys.exit('ERROR: unable to connect to Kafka')

    print('Checking Kafka connection...')

    try:
        consumer = KafkaConsumer(WORKER_NONCE_OUT_TOPIC,
                                 bootstrap_servers=[KAFKA_URL])
        break
    except:
        retries += 1

    sleep(6)

print('Connected to Kafka')

############
# CONSUMER #
############

def deserialize(message, typ):
    module_, class_ = typ.rsplit('.', 1)
    class_ = getattr(import_module(module_), class_)
    rv = class_()
    rv.ParseFromString(message)
    return rv


print('Ready!')
for message in consumer:

    print('offset:', message[2])
    
    try:
        print(deserialize(message[6], 'trace_pb2.Trace'))
    except:
        print('raw value: ', message[6])

    print()

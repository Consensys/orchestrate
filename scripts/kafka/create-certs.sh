BASE_DIR=$(git rev-parse --show-toplevel)
SECRET_DIR=${BASE_DIR}/config/kafka/secrets


CA_DIR=${SECRET_DIR}/ca
CA_PASSWORD=secret
CA_KEY_FILE=${CA_DIR}/ca.key
CA_CERT_FILE=${CA_DIR}/ca.crt

rm -rf ${SECRET_DIR}
mkdir -p ${CA_DIR}

# 1 - Create a Certificate authority
openssl \
  req \
  -new \
  -days 365 \
  -x509 \
  -subj "/CN=Kafka-Security-CA" \
  -keyout ${CA_KEY_FILE} \
  -out ${CA_CERT_FILE} \
  -passin pass:${CA_PASSWORD} \
  -passout pass:${CA_PASSWORD}


for i in kafka client
do
echo "------------------------------- $i -------------------------------"

CERT_DIR=${SECRET_DIR}/${i}
KAFKA_SSL_KEYSTORE_FILE=${CERT_DIR}/${i}.keystore.jks
KAFKA_SSL_KEYSTORE_FILE_P12=${CERT_DIR}/${i}.keystore.p12
KAFKA_SSL_KEYSTORE_FILE_PEM=${CERT_DIR}/${i}.keystore.pem
KAFKA_SSL_TRUSTSTORE_FILE=${CERT_DIR}/${i}.truststore.jks
KAFKA_SSL_CERTIFICATE_REQUEST_FILE=${CERT_DIR}/${i}.cert-req.csr
KAFKA_SSL_SIGNED_CERTIFICATE_FILE=${CERT_DIR}/${i}.cert-signed.crt
KAFKA_SSL_KEY_CREDENTIALS_FILE=${CERT_DIR}/${i}.sslkey-creds
KAFKA_SSL_KEY_CREDENTIALS=${i}_secret

mkdir -p ${CERT_DIR}


# # 2.1 - Create keystore
keytool \
  -genkey \
  -keyalg RSA \
  -keystore ${KAFKA_SSL_KEYSTORE_FILE} \
  -validity 365 \
  -storepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -keypass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -alias $i \
  -dname "CN=${i}" 

# 2.2 - Create a certificate request file to be signed by the CA
keytool \
  -keystore ${KAFKA_SSL_KEYSTORE_FILE} \
  -alias $i \
  -certreq -file ${KAFKA_SSL_CERTIFICATE_REQUEST_FILE} \
  -storepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -keypass ${KAFKA_SSL_KEY_CREDENTIALS}

# 2.3 - Sign certificate request from CA
openssl \
  x509 \
  -req \
  -CA ${CA_CERT_FILE} \
  -CAkey ${CA_KEY_FILE} \
  -in ${KAFKA_SSL_CERTIFICATE_REQUEST_FILE}\
  -out ${KAFKA_SSL_SIGNED_CERTIFICATE_FILE} \
  -days 365 \
  -CAcreateserial \
  -passin pass:${CA_PASSWORD}

# 2.4 - Sign and import the CA cert into the keystore
keytool \
  -keystore ${KAFKA_SSL_KEYSTORE_FILE} \
  -alias CARoot \
  -import -file ${CA_CERT_FILE} \
  -storepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -keypass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -noprompt

# 2.5 - Sign and import the host certificate into the keystore
keytool \
  -keystore ${KAFKA_SSL_KEYSTORE_FILE} \
  -alias ${i} \
  -import -file ${KAFKA_SSL_SIGNED_CERTIFICATE_FILE} \
  -storepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -keypass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -noprompt

# 2.6 - Create truststore and import the CA cert
keytool \
  -keystore ${KAFKA_SSL_TRUSTSTORE_FILE} \
  -alias CARoot \
  -import -file ${CA_CERT_FILE} \
  -storepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -keypass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -noprompt

# 3 - Save credentials
echo ${KAFKA_SSL_KEY_CREDENTIALS} > ${KAFKA_SSL_KEY_CREDENTIALS_FILE}

# 4 - Get PEM files from truststore and keystore
keytool \
  -importkeystore \
  -srckeystore ${KAFKA_SSL_TRUSTSTORE_FILE} \
  -destkeystore ${KAFKA_SSL_KEYSTORE_FILE_P12} \
  -srcstoretype JKS \
  -deststoretype PKCS12 \
  -srcstorepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -deststorepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -noprompt

openssl \
  pkcs12 \
  -in ${KAFKA_SSL_KEYSTORE_FILE_P12} \
  -nokeys \
  -passin \
  pass:${KAFKA_SSL_KEY_CREDENTIALS} \
  | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > ${CERT_DIR}/ca.cer.pem


keytool \
  -importkeystore \
  -srckeystore ${KAFKA_SSL_KEYSTORE_FILE} \
  -destkeystore ${KAFKA_SSL_KEYSTORE_FILE_P12} \
  -srcstoretype JKS \
  -deststoretype PKCS12 \
  -srcstorepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -deststorepass ${KAFKA_SSL_KEY_CREDENTIALS} \
  -noprompt

openssl \
  pkcs12 \
  -in ${KAFKA_SSL_KEYSTORE_FILE_P12} \
  -nodes -nocerts \
  -passin \
  pass:${KAFKA_SSL_KEY_CREDENTIALS} \
  | sed -ne '/-BEGIN PRIVATE KEY-/,/-END PRIVATE KEY-/p' > ${CERT_DIR}/client.key.pem

done
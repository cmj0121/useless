#! /usr/bin/env sh
#
# generate strongswan configuration via extra configuration
#    - DOMAIN
#    - DEVICE
#    - VLAN
#    - LAN
#    - DNS
#
# ref: https://wiki.archlinux.org/index.php/StrongSwan

IPSEC_DIR=/etc/ipsec.d
# environment and default value
UNIQUE=${UNIQUE:=replace}
DOMAIN=${DOMAIN:=vpn.example.com}
DEVICE=${DEVICE:=eth0}
VLAN=${VLAN:=10.10.0.0/24}
LAN=${LAN:=192.168.0.0/24}
DNS=${DNS:=8.8.8.8}
USR=${USR:=guest}
P12_PASSWORD=${P12_PASSWORD:=p12_secret}


function run() {
	iptables -t NAT -A POSTROUTING -s "${VLAN}" -o "${DEVICE}" -j MASQUERADE
	exec ipsec start --nofork "$@"
}

if [ -e ${IPSEC_DIR}/ipsec.conf ]; then
    echo "already configured"
	run
    exit 0
fi

mkdir -p ${IPSEC_DIR}/{private,certs,cacerts}

cat > "${IPSEC_DIR}/ipsec.conf" << _EOF_
config setup
    uniqueids    = ${UNIQUE}    # allow multiple device or other


conn %default
    keyexchange = ikev2    # method of key exchange
    dpdaction   = clear    # the connection is closed with no further actions taken
    dpddelay    = 300s
    rekey       = no    # do NOT renegotiate a connection if it is about to expire

    left        = %any
    leftca      = ca.cert.pem
    leftcert    = server.cert.pem
    leftsubnet  = 0.0.0.0/0

    right           = %any
    rightdns        = ${DNS}
    rightsourceip   = ${VLAN}
    rightsubnets    = ${LAN}
_EOF_

for usr in ${USR}; do
    cat >> "${IPSEC_DIR}/ipsec.conf" << _EOF_


conn IPSec-IKEv2
    keyexchange  = ikev2
    ike          = aes256-sha256-modp1024,3des-sha1-modp1024,aes256-sha1-modp1024!
    esp          = aes256-sha256,3des-sha1,aes256-sha1!
    auto         = add

    leftid          = "${DOMAIN}"
    leftsendcert    = always
    leftauth        = pubkey
    rightauth       = pubkey
    rightid         = "${usr}@${DOMAIN}"
    rightcert       = ${usr}.client.cert.pem
_EOF_
done

cat > "${IPSEC_DIR}/ipsec.secrets" << _EOF_
: RSA server.pem
_EOF_

# generate CA
cd ${IPSEC_DIR}
ipsec pki --gen --type rsa --size 4096 --outform pem > private/key.pem
chmod 600 private/key.pem
ipsec pki --self \
          --ca --lifetime 3650 --outform pem \
          --in private/key.pem --type rsa \
          --dn "C=CH, O=strongSwan, CN=strongSwan Root CA" \
          > cacerts/cert.pem

## ---- server ----
if [ ! -e private/private/server.pem ]; then
    ipsec pki --gen --type rsa --size 2048 --outform pem > private/server.pem
    chmod 600 private/server.pem
    ipsec pki --pub --in private/server.pem --type rsa | \
        ipsec pki --issue --lifetime 730 --outform pem \
                  --cacert cacerts/cert.pem \
                  --cakey private/key.pem \
                  --dn "C=CH, O=strongSwan, CN=vpn.example.com" \
                  --san ${DOMAIN} \
                  --flag serverAuth \
                  --flag ikeIntermediate \
                  > certs/server.cert.pem
fi

## ---- client ----
for usr in ${USR}; do
    if [ ! -e private/${usr}.client.pem ]; then
        ipsec pki --gen --type rsa --size 2048 --outform pem > private/${usr}.client.pem
        chmod 600 private/${usr}.client.pem
        ipsec pki --pub --in private/${usr}.client.pem --type rsa | \
            ipsec pki --issue --lifetime 730 --outform pem \
                      --cacert cacerts/cert.pem \
                      --cakey private/key.pem \
                      --dn "C=CH, O=strongSwan, CN=vpn.example.com" \
                      --san ${usr}@${DOMAIN} \
                      > certs/${usr}.client.cert.pem

        openssl pkcs12 -export \
               -inkey ${IPSEC_DIR}/private/${usr}.client.pem \
               -in ${IPSEC_DIR}/certs/${usr}.client.cert.pem \
               -name "${usr}@${DOMAIN}" \
               -certfile ${IPSEC_DIR}/cacerts/cert.pem \
               -caname "strongSwan Root CA" \
               -out ${IPSEC_DIR}/${usr}.client.cert.p12 \
               -passout pass:${P12_PASSWORD}

        # for MAC
        UUID1=$(uuidgen)
        UUID2=$(uuidgen)
        UUID3=$(uuidgen)
        UUID4=$(uuidgen)
        UUID5=$(uuidgen)

        cat > ${usr}.client.mobileconfig <<_EOF_
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
 <key>PayloadContent</key>
 <array>
  <dict>
   <key>Password</key>
   <string>${P12_PASSWORD}</string>
   <key>PayloadCertificateFileName</key>
   <string>client.cert.p12</string>
   <key>PayloadContent</key>
   <data>
$(base64 ${IPSEC_DIR}/${usr}.client.cert.p12)
   </data>
   <key>PayloadDescription</key>
   <string>添加 PKCS#12 格式的证书</string>
   <key>PayloadDisplayName</key>
   <string>client.cert.p12</string>
   <key>PayloadIdentifier</key>
   <string>com.apple.security.pkcs12.${UUID1}</string>
   <key>PayloadType</key>
   <string>com.apple.security.pkcs12</string>
   <key>PayloadUUID</key>
   <string>${UUID1}</string>
   <key>PayloadVersion</key>
   <integer>1</integer>
  </dict>
  <dict>
   <key>PayloadCertificateFileName</key>
   <string>ca.cer</string>
   <key>PayloadContent</key>
   <data>
$(base64 ${IPSEC_DIR}/cacerts/cert.pem)
   </data>
   <key>PayloadDescription</key>
   <string>添加 CA 根证书</string>
   <key>PayloadDisplayName</key>
   <string>strongSwan Root CA</string>
   <key>PayloadIdentifier</key>
   <string>com.apple.security.root.${UUID2}</string>
   <key>PayloadType</key>
   <string>com.apple.security.root</string>
   <key>PayloadUUID</key>
   <string>${UUID2}</string>
   <key>PayloadVersion</key>
   <integer>1</integer>
  </dict>
  <dict>
   <key>IKEv2</key>
   <dict>
    <key>AuthenticationMethod</key>
    <string>Certificate</string>
    <key>ChildSecurityAssociationParameters</key>
    <dict>
     <key>DiffieHellmanGroup</key>
     <integer>2</integer>
     <key>EncryptionAlgorithm</key>
     <string>3DES</string>
     <key>IntegrityAlgorithm</key>
     <string>SHA1-96</string>
     <key>LifeTimeInMinutes</key>
     <integer>1440</integer>
    </dict>
    <key>DeadPeerDetectionRate</key>
    <string>Medium</string>
    <key>DisableMOBIKE</key>
    <integer>0</integer>
    <key>DisableRedirect</key>
    <integer>0</integer>
    <key>EnableCertificateRevocationCheck</key>
    <integer>0</integer>
    <key>EnablePFS</key>
    <integer>0</integer>
    <key>IKESecurityAssociationParameters</key>
    <dict>
     <key>DiffieHellmanGroup</key>
     <integer>2</integer>
     <key>EncryptionAlgorithm</key>
     <string>3DES</string>
     <key>IntegrityAlgorithm</key>
     <string>SHA1-96</string>
     <key>LifeTimeInMinutes</key>
     <integer>1440</integer>
    </dict>
    <key>LocalIdentifier</key>
    <string>client@${DOMAIN}</string>
    <key>PayloadCertificateUUID</key>
    <string>${UUID1}</string>
    <key>RemoteAddress</key>
    <string>${DOMAIN}</string>
    <key>RemoteIdentifier</key>
    <string>${DOMAIN}</string>
    <key>UseConfigurationAttributeInternalIPSubnet</key>
    <integer>0</integer>
   </dict>
   <key>IPv4</key>
   <dict>
    <key>OverridePrimary</key>
    <integer>1</integer>
   </dict>
   <key>PayloadDescription</key>
   <string>Configures VPN settings</string>
   <key>PayloadDisplayName</key>
   <string>VPN</string>
   <key>PayloadIdentifier</key>
   <string>com.apple.vpn.managed.${UUID3}</string>
   <key>PayloadType</key>
   <string>com.apple.vpn.managed</string>
   <key>PayloadUUID</key>
   <string>${UUID3}</string>
   <key>PayloadVersion</key>
   <real>1</real>
   <key>Proxies</key>
   <dict>
    <key>HTTPEnable</key>
    <integer>0</integer>
    <key>HTTPSEnable</key>
    <integer>0</integer>
   </dict>
   <key>UserDefinedName</key>
   <string>VPN (IKEv2)</string>
   <key>VPNType</key>
   <string>IKEv2</string>
  </dict>
  <dict>
   <key>PayloadCertificateFileName</key>
   <string>server.cer</string>
   <key>PayloadContent</key>
   <data>
$(base64 ${IPSEC_DIR}/certs/server.cert.pem)
   </data>
   <key>PayloadDescription</key>
   <string>添加 PKCS#1 格式的證書</string>
   <key>PayloadDisplayName</key>
   <string>${DOMAIN}</string>
   <key>PayloadIdentifier</key>
   <string>com.apple.security.pkcs1.${UUID4}</string>
   <key>PayloadType</key>
   <string>com.apple.security.pkcs1</string>
   <key>PayloadUUID</key>
   <string>${UUID4}</string>
   <key>PayloadVersion</key>
   <integer>1</integer>
  </dict>
 </array>
 <key>PayloadDisplayName</key>
 <string>VPN</string>
 <key>PayloadIdentifier</key>
 <string>com.github.vimagick.strongswan</string>
 <key>PayloadRemovalDisallowed</key>
 <false/>
 <key>PayloadType</key>
 <string>Configuration</string>
 <key>PayloadUUID</key>
 <string>${UUID5}</string>
 <key>PayloadVersion</key>
 <integer>1</integer>
</dict>
</plist>
_EOF_
    fi
done

run

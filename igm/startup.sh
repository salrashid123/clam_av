#! /bin/bash    

echo "deb http://http.debian.net/debian/ stretch main contrib non-free" > /etc/apt/sources.list
echo "deb http://http.debian.net/debian/ stretch-updates main contrib non-free" >> /etc/apt/sources.list
echo "deb http://security.debian.org/ stretch/updates main contrib non-free" >> /etc/apt/sources.list
apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y -qq \
    clamav-daemon \
    clamav-freshclam \
    libclamunrar7 \
    wget

wget -O /var/lib/clamav/main.cvd http://database.clamav.net/main.cvd
wget -O /var/lib/clamav/daily.cvd http://database.clamav.net/daily.cvd
wget -O /var/lib/clamav/bytecode.cvd http://database.clamav.net/bytecode.cvd
chown clamav:clamav /var/lib/clamav/*.cvd

sed -i 's/^Foreground .*$/Foreground true/g' /etc/clamav/clamd.conf 
sed -i 's/MaxScanSize 100M/MaxScanSize 500M/g' /etc/clamav/clamd.conf 
sed -i 's/MaxFileSize 25M/MaxFileSize 500M/g' /etc/clamav/clamd.conf
sed -i 's/StreamMaxLength 25M/StreamMaxLength 500M/g' /etc/clamav/clamd.conf 
echo "TCPSocket 3310" >> /etc/clamav/clamd.conf 

service clamav-daemon restart
service clamav-freshclam restart
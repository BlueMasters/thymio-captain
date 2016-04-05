#!/bin/bash


# Install package to build aseba
apt-get update
apt-get install libboost-dev libqt4-dev qt4-dev-tools libqwt5-qt4-dev libudev-dev libxml2-dev cmake g++ git make -y

# Build aseba
./compile_aseba.sh

# Copy asebahttp to /usr/bin
cp aseba/build-aseba/switches/http/asebahttp /usr/bin

# Install python3 and pip
apt-get install python3 -y
apt-get install python3-pip -y
pip3 install httplib2
python3 thymiolib/setup.py install

# Create a script who lauch asebahttp with good parameters
echo "#!/bin/bash" > /usr/local/bin/thymioHttp
echo "asebahttp -p 3000 -v -a /home/pi/thymio_events.aesl ser:name=Thymio-II" > /usr/local/bin/thymioHttp
chmod +x /usr/local/bin/thymioHttp

# Create infinite loop script who check if Thymio is connected and start thymioHttp
echo "#!/bin/bash" > /usr/local/bin/cronThymioHttp
echo "Interval="1" # polling interval in seconds" >> /usr/local/bin/cronThymioHttp
echo "while true; do" >> /usr/local/bin/cronThymioHttp
echo "  if lsusb | grep 0617; then" >> /usr/local/bin/cronThymioHttp
echo "    /usr/local/bin/thymioHttp" >> /usr/local/bin/cronThymioHttp
echo "  fi" >> /usr/local/bin/cronThymioHttp
echo "  sleep \"\$Interval\"" >> /usr/local/bin/cronThymioHttp
echo "done" >> /usr/local/bin/cronThymioHttp
chmod +x /usr/local/bin/cronThymioHttp

# Add script in crontab to lauch it automatically on boot
crontab -l > tmp_crontab
if [ $? -ne 0 ]; then
  echo "" > tmp_crontab
fi
grep -Fq "@reboot /usr/local/bin/cronThymioHttp" tmp_crontab || echo "@reboot /usr/local/bin/cronThymioHttp" >> tmp_crontab
crontab tmp_crontab
rm tmp_crontab 

# Install Apache2 server and flask
apt-get install apache2 -y
pip3 install Flask
pip3 install flask-restful
apt-get install libapache2-mod-wsgi-py3 -y
cp -r flask_dev /var/www/
cp flask_dev.conf /etc/apache2/sites-available/
a2enmod wsgi
a2ensite flask_dev.conf
a2dissite 000-default.conf
service apache2 restart


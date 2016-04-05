#!/bin/sh
# create build tree
mkdir -p aseba/build-dashel aseba/build-enki aseba/build-aseba
cd aseba
# fetch and compile dashel
git clone git://github.com/aseba-community/dashel.git
cd build-dashel
cmake ../dashel -DCMAKE_BUILD_TYPE=RelWithDebInfo -DBUILD_SHARED_LIBS=ON
make
cd ..
# fetch and compile enki
git clone git://github.com/enki-community/enki.git
cd build-enki
cmake ../enki -DCMAKE_BUILD_TYPE=RelWithDebInfo
make
cd ..
# fetch and compile aseba, telling it where to find dashel and enki
git clone --recursive git://github.com/aseba-community/aseba.git
cd aseba && git checkout release-1.5.x && cd ..
cd build-aseba
export dashel_DIR=../build-dashel
cmake ../aseba -DCMAKE_BUILD_TYPE=RelWithDebInfo -Ddashel_DIR=../build-dashel -DDASHEL_INCLUDE_DIR=../dashel -DDASHEL_LIBRARY=../build-dashel/libdashel.so -DENKI_INCLUDE_DIR=../enki -DENKI_LIBRARY=../build-enki/enki/libenki.a -DENKI_VIEWER_LIBRARY=../build-enki/viewer/libenkiviewer.a
make
cd ../..

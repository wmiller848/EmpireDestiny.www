#
FROM fedora:22

RUN dnf install -y glibc-devel glibc-headers glibc-static
RUN dnf install -y gcc gcc-c++ cpp pcre-devel
RUN dnf install -y openssl-devel zlib-devel
RUN dnf install -y cmake make autoconf pkgconfig
RUN dnf install -y git mercurial tar xz wget
RUN dnf install -y bzip2 freetype* fontconfig

WORKDIR /opt/vendor
RUN wget --no-cookies --no-check-certificate --header "Cookie: gpw_e24=https%3A%2F%2Fwww.oracle.com%2F; oraclelicense=accept-securebackup-cookie" "https://download.oracle.com/otn-pub/java/jdk/8u45-b14/jdk-8u45-linux-x64.tar.gz"

RUN tar -xf jdk-8u45-linux-x64.tar.gz
RUN rm jdk-8u45-linux-x64.tar.gz

ENV JAVA_HOME=/opt/vendor/jdk1.8.0_45
ENV JRE_HOME=/opt/vendor/jdk1.8.0_45/jre
ENV PATH=$PATH:/opt/vendor/jdk-8u45-linux-x64/bin:/opt/vendor/jdk1.8.0_45/jre/bin

RUN wget https://nodejs.org/dist/v5.7.0/node-v5.7.0-linux-x64.tar.xz
RUN tar -xf node-v5.7.0-linux-x64.tar.xz
RUN rm node-v5.7.0-linux-x64.tar.xz
ENV PATH=$PATH:/opt/vendor/node-v5.7.0-linux-x64/bin

WORKDIR /opt

RUN npm install --global gulp-cli

ADD ./package.json /opt/empire.www/package.json
WORKDIR /opt/empire.www
RUN npm install

ADD ./ /opt/empire.www
RUN gulp
VOLUME /opt/empire.www

CMD gulp watch

FROM centos:6.4

RUN rpm -Uvh http://download.fedoraproject.org/pub/epel/6/i386/epel-release-6-8.noarch.rpm
RUN yum install -y npm

ADD . /src
RUN cd /src; npm install --unsafe-perm

EXPOSE 3000
CMD ["node", "/src/app.js", "-c", "/src/config.js"]

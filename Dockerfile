FROM centos
ADD main /
ENV ip "127.0.0.1"
ENV port "8888"
EXPOSE 8888:8888
CMD ./main -startType server
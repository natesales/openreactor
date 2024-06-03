FROM debian:bullseye

ARG SVC
COPY svc-$SVC /svc

EXPOSE 80

CMD ["/svc"]

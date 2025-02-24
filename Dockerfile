FROM debian:stable-slim
# COPY source destination
COPY chirpy /bin/goserver

CMD ["/bin/goserver"]





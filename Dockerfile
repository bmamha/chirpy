FROM debian:stable-slim
# COPY source destination
COPY chirpy /bin/goserver
COPY index.html index.html
COPY assets/ assets/
ENV PORT=8080
CMD ["/bin/goserver"]





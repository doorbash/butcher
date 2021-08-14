FROM scratch
ADD butcher /butcher
ADD config.json /config.json
EXPOSE 53
CMD [ "/butcher", "-c", "config.json", "0.0.0.0:53" ]
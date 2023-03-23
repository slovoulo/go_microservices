
#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

#copy fromthe first broker image to this one
COPY listenerService /app

#Build the first docker image, create a  much smaller docker image then copy the executable from first to second smaller image
CMD [ "/app/listenerService" ]

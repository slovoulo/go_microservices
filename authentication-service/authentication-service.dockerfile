#base go image
FROM golang:1.19.5-alpine as builder

#Make a folder called authapp in the docker image
RUN mkdir /authapp

#Copy everything from current folder (.) to the (authapp) folder
COPY . /authapp

#MAke authapp the new working directory
WORKDIR /authapp

#Build the go code
#CGO_ENABLED=0 is an environment variable
#go build is the build -0 command
#brokerService is the name of the app once built
#./cmd is the directory containing the project files we are building from
RUN CGO_ENABLED=0 go build -o authService ./cmd/api

#Run the chmod command to ensure the app is executable
RUN chmod +x /authapp/authService

#build a tiny docker image
FROM alpine:latest

RUN mkdir /authapp

#copy fromthe first broker image to this one
COPY --from=builder /authapp/authService /authapp

#Build the first docker image, create a  much smaller docker image then copy the executable from first to second smaller image
CMD [ "/authapp/authService" ]
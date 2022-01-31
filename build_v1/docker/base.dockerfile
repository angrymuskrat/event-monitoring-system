FROM ubuntu:20.04

# tor for downloading images
RUN apt-get update
RUN apt-get install tor -y

# postgres client
#RUN apt-get install -y lsb-release
#RUN apt-get install gnupg -y
#RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
#RUN apt-get install wget -y
#RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
#RUN apt-get update
#  RUN apt-get -y install postgresql-13

# npm for front
RUN apt-get update
RUN apt-get install nodejs -y
RUN apt-get install npm -y

RUN apt-get install screen nano -y
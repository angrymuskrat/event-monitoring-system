FROM base_event_service:v1

COPY start_tor.sh start_tor.sh

#EXPOSE 5432
EXPOSE 17115
EXPOSE 17112

RUN mkdir /monitoring
RUN mkdir -p /data/images


COPY backend /monitoring/backend
COPY front /monitoring/front
COPY storage /monitoring/storage
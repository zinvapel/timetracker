FROM golang:1.15

ENV TZ=Europe/Moscow
RUN ln -fs /usr/share/zoneinfo/Europe/Moscow /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

COPY . /app/.

WORKDIR /app

RUN go build -o /bin/timetracker .

ENTRYPOINT ["timetracker"]
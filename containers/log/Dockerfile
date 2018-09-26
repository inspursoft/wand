FROM library/ubuntu:14.04

RUN rm /etc/rsyslog.d/* && rm /etc/rsyslog.conf

ADD containers/log/rsyslog.conf /etc/rsyslog.conf

# rotate logs weekly
# notes: file name cannot contain dot, or the script will not run
ADD containers/log/rotate.sh /etc/cron.weekly/rotate

# rsyslog configuration file for docker
ADD containers/log/rsyslog_docker.conf /etc/rsyslog.d/

VOLUME /var/log/board/

EXPOSE 514

CMD cron && rsyslogd -n

FROM jenkins/jenkins:2.127

ENV JAVA_OPTS "-Djenkins.install.runSetupWizard=false \
               -Dhudson.model.User.allowNonExistentUser=true"
USER root

COPY containers/jenkins/sources.list /tmp/sources.list
COPY containers/jenkins/jobs /tmp/jobs

ADD containers/jenkins/site-packages/python-jenkins-0.4.15.tar.gz /usr/share/jenkins
ADD containers/jenkins/site-packages/setuptools-38.5.1.tar.gz /usr/share/jenkins
ADD containers/jenkins/site-packages/multi_key_dict-2.0.3.tar.gz /usr/share/jenkins

COPY containers/jenkins/pythonenv.sh /usr/share/jenkins/pythonenv.sh
COPY containers/jenkins/plugins.txt /usr/share/jenkins/plugins.txt
COPY containers/jenkins/addNode.py /usr/share/jenkins/addNode.py
COPY containers/jenkins/jenkins.sh /usr/local/bin/jenkins.sh
COPY containers/jenkins/init.sh /usr/share/jenkins/init.sh

RUN  sed -i 's/^root\:x\:0\:/root\:x\:0\:root\,jenkins/' /etc/group \
      && cat /tmp/sources.list > /etc/apt/sources.list \    
      && apt-get update \
      && apt-get install -y --allow-unauthenticated sudo libltdl7 curl \
      && rm -rf /var/lib/apt/lists/* \
      && echo "jenkins ALL=NOPASSWD: ALL" >> /etc/sudoers \
      && chmod u+x /usr/share/jenkins/pythonenv.sh \
      && chmod u+x /usr/local/bin/jenkins.sh \
      && chmod u+x /usr/share/jenkins/init.sh \
      && /usr/local/bin/plugins.sh /usr/share/jenkins/plugins.txt
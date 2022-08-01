from golang:1.18.4-bullseye

RUN apt update
RUN apt install -y libpcap0.8 libpcap0.8-dev libpcap-dev vim net-tools
RUN apt install -y tcpdump netcat openssh-server lsof
RUN apt install -y supervisor
RUN echo "PermitRootLogin yes" > /etc/ssh/sshd_config
RUN echo "root:pass" | chpasswd
RUN ssh-keygen -t rsa -f /root/.ssh/id_rsa -N ""
RUN cat /root/.ssh/id_rsa.pub > /root/.ssh/authorized_keys

# Incidentally creates /run/sshd
# RUN service ssh start
# Run the service in the foreground when starting the container
# CMD ["/usr/sbin/sshd", "-D"]

# ssh-keyscan ubuntu1 > /root/.ssh/known_hosts
# -o StrictHostKeyChecking=no

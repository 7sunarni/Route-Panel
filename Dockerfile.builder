FROM debian:11.9 

RUN echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free" > /etc/apt/sources.list \
    && echo "deb-src http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye main contrib non-free" >> /etc/apt/sources.list \
    && echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free" >> /etc/apt/sources.list \
    && echo "deb-src http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-updates main contrib non-free" >> /etc/apt/sources.list \
    && echo "deb http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free" >> /etc/apt/sources.list \
    && echo "deb-src http://mirrors.tuna.tsinghua.edu.cn/debian/ bullseye-backports main contrib non-free" >> /etc/apt/sources.list \
    && echo "deb https://security.debian.org/debian-security bullseye-security main contrib non-free" >> /etc/apt/sources.list \
    && echo "deb-src https://security.debian.org/debian-security bullseye-security main contrib non-free" >> /etc/apt/sources.list \
    && apt update \
    && apt install -y pkg-config libwebkit2gtk-4.0-dev libgtk-3-dev libucl1 upx-ucl libx11-dev build-essential \
    && rm -rf /var/lib/apt/lists/* \
    && curl https://nodejs.org/dist/v20.11.1/node-v20.11.1-linux-x64.tar.xz -L -o node-v20.11.1-linux-x64.tar.xz \
    && curl https://go.dev/dl/go1.21.7.linux-amd64.tar.gz -L -o go1.21.7.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go1.21.7.linux-amd64.tar.gz \
    && tar -xf node-v20.11.1-linux-x64.tar.xz --directory=/opt/ \
    && rm -f node-v20.11.1-linux-x64.tar.xz go1.21.7.linux-amd64.tar.gz

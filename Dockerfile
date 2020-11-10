from golang:1.12

RUN apt-get update && apt-get install --no-install-recommends -y \
    bc \
    gcc-multilib \
    libssl-dev  \
    llvm-dev \
    libjemalloc-dev \
    libnuma-dev \
    python-sphinx \
    libelf-dev \
    selinux-policy-dev \
    libunbound-dev \
    autoconf \
    automake \
    libtool

ADD .travis /src/travis
WORKDIR /src/travis

ENV OVN_SRCDIR=/src/
RUN sh ./test_prepare.sh

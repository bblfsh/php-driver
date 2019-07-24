# This file can be used directly with Docker.
#
# Prerequisites:
#   go mod vendor
#   bblfsh-sdk release
#
# However, the preferred way is:
#   go run ./build.go driver:tag
#
# This will regenerate all necessary files before building the driver.

#==============================
# Stage 1: Native Driver Build
#==============================
FROM php:7-alpine3.6 as native

# add dependency files
ADD https://getcomposer.org/installer /tmp/composer-setup.php


# install build dependencies
RUN php /tmp/composer-setup.php --install-dir=/bin/


ADD native /native
WORKDIR /native

# build native driver
RUN composer.phar install


#================================
# Stage 1.1: Native Driver Tests
#================================
FROM native as native_test

# run native driver tests
RUN ./vendor/bin/phpunit tests/


#=================================
# Stage 2: Go Driver Server Build
#=================================
FROM golang:1.10-alpine as driver

ENV DRIVER_REPO=github.com/bblfsh/php-driver
ENV DRIVER_REPO_PATH=/go/src/$DRIVER_REPO

ADD go.* $DRIVER_REPO_PATH/
ADD vendor $DRIVER_REPO_PATH/vendor
ADD driver $DRIVER_REPO_PATH/driver

WORKDIR $DRIVER_REPO_PATH/

ENV GO111MODULE=on GOFLAGS=-mod=vendor

# workaround for https://github.com/golang/go/issues/28065
ENV CGO_ENABLED=0

# build server binary
RUN go build -o /tmp/driver ./driver/main.go
# build tests
RUN go test -c -o /tmp/fixtures.test ./driver/fixtures/

#=======================
# Stage 3: Driver Build
#=======================
FROM php:7-alpine3.6



LABEL maintainer="source{d}" \
      bblfsh.language="php"

WORKDIR /opt/driver

# copy static files from driver source directory
ADD ./native/ast ./bin/native
ADD ./native/src ./bin/src


# copy build artifacts for native driver
COPY --from=native /native/vendor ./bin/vendor


# copy driver server binary
COPY --from=driver /tmp/driver ./bin/

# copy tests binary
COPY --from=driver /tmp/fixtures.test ./bin/
# move stuff to make tests work
RUN ln -s /opt/driver ../build
VOLUME /opt/fixtures

# copy driver manifest and static files
ADD .manifest.release.toml ./etc/manifest.toml

ENTRYPOINT ["/opt/driver/bin/driver"]
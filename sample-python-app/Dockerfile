# Use the UBI minimal base image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL name="sample-python-app" \
      description="Print Current UTC time with Platform details" \
      target-file="Dockerfile"

# Set environment variables
ENV PYTHON_VERSION=3.9 \
    PATH=$HOME/.local/bin/:$PATH \
    PYTHONUNBUFFERED=1 \
    PYTHONIOENCODING=UTF-8 \
    LC_ALL=en_US.UTF-8 \
    LANG=en_US.UTF-8

# Install Python 3.9
RUN microdnf update -y && \
    microdnf install -y python39 && \
    microdnf -y clean all && rm -rf /var/cache/yum


# Install pipenv and setuptools
RUN pip3.9 install --no-cache-dir --upgrade pip && \
    pip3.9 install --no-cache-dir --upgrade pipenv

# Copy files to workdir
WORKDIR /opt/app-root/
COPY src /opt/app-root/src
COPY ./Pipfile* /opt/app-root/

# install python dependency using pipenv
RUN PIPENV_VENV_IN_PROJECT=1 PIPENV_IGNORE_VIRTUALENVS=1 pipenv install --system --deploy --ignore-pipfile

# drop root privilege
USER 1001

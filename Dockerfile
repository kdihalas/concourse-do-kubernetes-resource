FROM python:3.7-alpine3.7

RUN apk update \
    && apk add --no-cache bash

# Copy assets
COPY resources/ /opt/resource/
ADD requirements.txt .
RUN pip install -r requirements.txt

CMD ["/bin/bash"]

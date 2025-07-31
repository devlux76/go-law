FROM mcr.microsoft.com/playwright/python:latest

WORKDIR /app

COPY fetch.py .

RUN pip install --no-cache-dir playwright requests pytest && \
    playwright install --with-deps chromium

VOLUME ["/srv/bulk-data"]

ENTRYPOINT ["python", "fetch.py"]
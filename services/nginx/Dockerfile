FROM nginx:1.29.0

COPY ./nginx.conf /etc/nginx/nginx.conf

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s \
    CMD ["curl", "--fail", "--silent", "http://localhost:6666/health"]

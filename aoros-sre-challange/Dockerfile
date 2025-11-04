FROM python:3-slim

WORKDIR /app

# Create a simple HTML file
RUN echo "Hello World" > index.html

# Expose port 8080
EXPOSE 8080

# Run a simple HTTP server on port 8080
CMD ["python3", "-m", "http.server", "8080"]
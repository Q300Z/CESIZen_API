FROM golang

# Add alias in bashrc
RUN echo "alias prisma='go run github.com/steebchen/prisma-client-go'" >> /root/.bashrc

# Install air
RUN go install github.com/air-verse/air@latest

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

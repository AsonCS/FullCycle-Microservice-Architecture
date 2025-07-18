# Use a base image with Java installed
FROM openjdk:26-jdk-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the JAR file into the container
COPY balances-0.0.1.jar app.jar

# Expose the port your application listens on (if applicable)
EXPOSE 3003

# Define the command to run your application
CMD ["java", "-jar", "app.jar"]
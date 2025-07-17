plugins {
    kotlin("jvm") version "2.1.10"
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
    maven {
        url = uri("http://packages.confluent.io/maven/")
        isAllowInsecureProtocol = true
    }
}

dependencies {
    // https://mvnrepository.com/artifact/org.apache.kafka/kafka-clients
    // implementation("org.apache.kafka:kafka-clients:4.0.0")

    implementation(kotlin("stdlib"))
    implementation(kotlin("reflect"))

    implementation("org.apache.kafka:kafka-clients:3.9.1")
    // implementation("org.apache.kafka:kafka-streams:2.1.0")
    implementation("org.apache.kafka:connect-runtime:2.3.1")
    // implementation("io.confluent:kafka-json-serializer:5.0.1")
    implementation("org.slf4j:slf4j-api:1.7.6")
    implementation("org.slf4j:slf4j-log4j12:1.7.6")
    // implementation("com.fasterxml.jackson.core:jackson-databind:[2.8.11.1,)")
    // implementation ("com.fasterxml.jackson.module:jackson-module-kotlin:[2.8.11.1,)")
    implementation("com.google.code.gson:gson:2.8.9")

    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.10.2")

    testImplementation(kotlin("test"))
}

tasks.test {
    useJUnitPlatform()
}
plugins {
    alias(libs.plugins.kotlin.jvm)
    alias(libs.plugins.ktor)
    alias(libs.plugins.kotlin.plugin.serialization)
}

group = "br.com.wallet"
version = "0.0.1"

application {
    mainClass = "br.com.wallet.MainKt"
}

tasks.jar.configure {
    manifest {
        attributes(mapOf("Main-Class" to "br.com.wallet.MainKt"))
    }
    configurations["compileClasspath"].forEach { file: File ->
        from(zipTree(file.absoluteFile))
    }
    duplicatesStrategy = DuplicatesStrategy.INCLUDE
}

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
    // implementation("com.google.code.gson:gson:2.8.9")

    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.10.2")

    implementation("mysql:mysql-connector-java:8.0.15")

    implementation(libs.ktor.server.core)
    implementation(libs.ktor.serialization.kotlinx.json)
    implementation(libs.ktor.server.content.negotiation)
    implementation(libs.postgresql)
    implementation(libs.h2)
    implementation(libs.exposed.core)
    implementation(libs.exposed.jdbc)
    implementation(libs.ktor.server.host.common)
    implementation(libs.ktor.server.status.pages)
    implementation(libs.ktor.server.netty)
    implementation(libs.logback.classic)
    implementation(libs.ktor.server.config.yaml)
    implementation(libs.ktor.server.config.yaml.jvm)
    testImplementation(libs.ktor.server.test.host)
    testImplementation(libs.kotlin.test.junit)

    testImplementation(kotlin("test"))
}

tasks.test {
    useJUnitPlatform()
}
package br.com.wallet.database

import io.ktor.server.application.*
import org.jetbrains.exposed.sql.Database
import java.sql.Connection
import java.sql.DriverManager

private const val driver = "com.mysql.cj.jdbc.Driver"
//private const val driver = "com.mysql.jdbc.Driver"
private val Application.password
    get() = environment.config
        .property("mysql.password")
        .getString()
private val Application.url
    get() = environment.config
        .property("mysql.url")
        .getString()
private val Application.user
    get() = environment.config
        .property("mysql.user")
        .getString()

fun Application.getConnection(): Connection {
    Class.forName(driver)
    return DriverManager.getConnection(url, user, password)
}

fun Application.getDatabase(): Database {
    Class.forName(driver)
    return Database.connect(
        //driver = driver,
        password = password,
        url = url,
        user = user
    )
}

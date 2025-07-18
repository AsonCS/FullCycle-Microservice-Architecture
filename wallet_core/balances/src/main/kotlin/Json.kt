package br.com.wallet

import kotlinx.serialization.json.Json

val json = Json {
    explicitNulls = true
    prettyPrint = true
}

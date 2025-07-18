package br.com.wallet.kafka

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class BalanceUpdated(
    @SerialName("name")
    val name: String,
    @SerialName("payload")
    val payload: Payload
) {
    @Serializable
    data class Payload(
        @SerialName("origin_account_balance")
        val originAccountBalance: Float,
        @SerialName("origin_account_id")
        val originAccountId: String,
        @SerialName("destination_account_balance")
        val destinationAccountBalance: Float,
        @SerialName("destination_account_id")
        val destinationAccountId: String
    )
}

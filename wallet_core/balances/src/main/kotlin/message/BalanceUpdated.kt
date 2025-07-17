package org.example.message

import com.google.gson.annotations.SerializedName

data class BalanceUpdated(
    @SerializedName("name")
    val name: String,
    @SerializedName("payload")
    val payload: Payload
) : Message() {
    data class Payload(
        @SerializedName("origin_account_balance")
        val originAccountBalance: Float,
        @SerializedName("origin_account_id")
        val originAccountId: String,
        @SerializedName("destine_account_balance")
        val destineAccountBalance: Float,
        @SerializedName("destine_account_id")
        val destineAccountId: String
    )
}

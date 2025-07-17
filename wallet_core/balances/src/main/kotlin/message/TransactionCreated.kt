package org.example.message

import com.google.gson.annotations.SerializedName

data class TransactionCreated(
    @SerializedName("name")
    val name: String,
    @SerializedName("payload")
    val payload: Payload
) : Message() {
    data class Payload(
        @SerializedName("amount")
        val amount: String,
        @SerializedName("origin_account_id")
        val originAccountId: String,
        @SerializedName("destine_account_id")
        val destineAccountId: String,
        @SerializedName("transaction_id")
        val transactionId: String,
    )
}

package br.com.wallet.database

import br.com.wallet.entity.Account
import br.com.wallet.gateway.AccountGateway
import kotlinx.coroutines.Dispatchers.IO
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.jetbrains.exposed.sql.transactions.transaction

class AccountService(
    database: Database
) : AccountGateway {
    object Accounts : Table("accounts") {
        val id = varchar(
            name = "id",
            length = 255
        ).uniqueIndex()
        val balance = float(
            name = "balance"
        )
        val updatedAt = varchar(
            name = "updated_at",
            length = 255
        )
    }

    init {
        transaction(database) {
            //SchemaUtils.create(Accounts)
        }
    }

    override suspend fun delete(
        accountId: String
    ) {
        dbQuery {
            Accounts
                .deleteWhere { Accounts.id eq accountId }
        }
    }

    override suspend fun findAll(): List<Account> = dbQuery {
        Accounts.selectAll()
            .map {
                Account(
                    id = it[Accounts.id],
                    balance = it[Accounts.balance],
                    updatedAt = it[Accounts.updatedAt]
                )
            }
    }

    override suspend fun findById(
        id: String
    ): Account? = dbQuery {
        selectAllWhereId(id)
    }

    override suspend fun upsert(
        vararg accounts: Account
    ) = dbQuery {
        /*
        MySQL doesn't support specifying conflict keys in UPSERT clause, dialect: MySQL
        accounts.forEach { account ->
            Accounts.upsert(
                Accounts.balance,
                where = { Accounts.id eq account.id },
            ) {
                it[balance] = account.balance
                it[id] = account.id
            }
        }
        */
        accounts.forEach { account ->
            val record = selectAllWhereId(account.id)
            if (record != null) {
                Accounts
                    .update(
                        where = { Accounts.id eq account.id }
                    ) {
                        it[balance] = account.balance
                        it[Accounts.id] = account.id
                    }
            } else {
                Accounts
                    .insert {
                        it[balance] = account.balance
                        it[Accounts.id] = account.id
                    }
            }
        }
    }

    private fun selectAllWhereId(
        id: String
    ) = Accounts
        .selectAll()
        /*.select(
            Accounts.balance
        )*/
        .where(
            Accounts.id eq id
        )
        .map {
            Account(
                id = it[Accounts.id],
                balance = it[Accounts.balance],
                updatedAt = it[Accounts.updatedAt]
            )
        }.singleOrNull()

    private suspend fun <T> dbQuery(block: suspend () -> T): T =
        newSuspendedTransaction(IO) { block() }
}

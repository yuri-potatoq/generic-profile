using Dapper;
using Microsoft.Data.Sqlite;

namespace furry_profile;

public class TransactionManager: IDisposable, IAsyncDisposable {
	private SqliteTransaction _tx { get; set; }
	private SqliteConnection _db { get; set; }

	public TransactionManager(SqliteConnection db) {
		_tx = db.BeginTransaction();
		_db = db;
	}

	public Task<int> RunAsync(string sql) {
		return _db.ExecuteAsync(sql, transaction: _tx);
	}
	
	public Task<T> QueryFirstAsync<T>(string sql) {
		return _db.QueryFirstAsync<T>(sql, transaction: _tx);
	}
	
	public void Dispose() {
		_tx.Dispose();
	}

	public async ValueTask DisposeAsync() {
		await _tx.DisposeAsync();
	}
}

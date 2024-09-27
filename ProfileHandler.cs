using System.Transactions;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;

namespace furry_profile;


public class ProfileHandler {
	private readonly ProfileService _profileService;
	private readonly ILogger<ProfileHandler> _logger;

	public ProfileHandler(ProfileService profileService , ILogger<ProfileHandler> logger) {
		_profileService = profileService;
		_logger = logger;
	}
	
	public async Task<IResult> GetEnrollmentAsync(int enrollmentId) {
		return Results.Ok(await _profileService.GetFullEnrolmentAsync(enrollmentId));
	}

	public async Task<IResult> PartialUpdateAsync([FromBody] PartialProfileRequest partial) {
	    /*
	     
	     //TODO: use transactions to all updates made
	     
	     var transactionOptions = new TransactionOptions {
		    IsolationLevel = IsolationLevel.ReadCommitted,
		    Timeout = TimeSpan.FromMinutes(1)
	    };
	    using TransactionScope tran = new TransactionScope(TransactionScopeOption.Required, transactionOptions);
	    // conn.EnlistTransaction(Transaction.Current);
	    */

	    int enrollmentId;
	    if (partial.Id is not null) {
		    var exists = await _profileService.CheckEnrolmentAsync(partial.Id.Value);
		    if (!exists) throw new Exception("enrollment ID does not exists!");
		    enrollmentId = partial.Id.Value;
	    } else {
		    enrollmentId = await _profileService.NewEnrollmentAsync();
	    }
	    
	    // partial checks
	    await RunPartialHandler(enrollmentId, partial.ChildInfos?.ToCommand(), _profileService.UpdateChildProfileAsync);
	    await RunPartialHandler(enrollmentId, partial.Parent?.ToCommand(), _profileService.UpdateChildParentsAsync);
	    await RunPartialHandler(enrollmentId, partial.Address?.ToCommand(), _profileService.UpdateAddressAsync);
	    await RunPartialHandler(enrollmentId, partial.ToModalitiesCommand(), _profileService.UpdateModalitiesAsync);
	    if (partial.EnrollmentShift is not null) {
		    await _profileService.UpdateShiftAsync(enrollmentId,
			    Enum.Parse<EnrollmentShift>(partial.EnrollmentShift));
	    }
	    
        return Results.Ok(await _profileService.GetFullEnrolmentAsync(enrollmentId));
    }

    async Task RunPartialHandler<T>(int enrollmentId, T? entity, PartialHandler<T> handler) {
	    if (entity is not null) {
		    await handler(enrollmentId, entity!);
	    }
    }
}



delegate Task PartialHandler<T>(int enrollmentId, T entity);


static class NopPartialHandler<T> {
	public static async Task Handle(int enrollmentId, T entity) {
		return;
	}
}

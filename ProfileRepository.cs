using System.ComponentModel;
using System.Text.Json.Serialization;
using System.Transactions;
using Dapper;
using Microsoft.Data.Sqlite;

namespace furry_profile;


[JsonConverter(typeof(JsonStringEnumConverter))]
public enum Gender {
	None,
	Male,
	Female
}

[JsonConverter(typeof(JsonStringEnumConverter))]
public enum EnrollmentShift {
	None,
	Morning,
	Afternoon,
}


[JsonConverter(typeof(JsonStringEnumConverter))]
public enum Modalities {
	Football,
	Basketball,
	Swimming,
	Yoga,
	Volleyball,
}

public class Address {
	public string ZipCode {get; set; } 
	public string Street {get; set; }
	public string City {get; set; }
	public string State {get; set; }
	public int HouseNumber {get; set; }
}

public class ChildProfile {
	public string fullName { get; set; }
	public string birthdate { get; set; } 
	public Gender gender { get; set; }
	public string medicalInfo { get; set; }
}

public class ChildParent {
	public string fullName { get; set; }
	public string email { get; set; }
	public string phoneNumber { get; set; }
}

public class EnrollmentState {
	public int Id { get; set; }
	public ChildParent childParent { get; set; }
	public ChildProfile childProfile { get; set; }
	public Address address { get; set; }
	public List<Modalities> modalities { get; set; }
	public EnrollmentShift enrollmentShift { get; set; }
	public bool terms { get; set; }
}

public interface IProfileRepository {
	public Task UpdateChildProfileAsync(int enrollmentId, ChildProfile childProfile);
	public Task UpdateChildParentsAsync(int enrollmentId, ChildParent childParent);
	public Task<List<ChildProfile>> GetChildProfileAsync(int enrollmentId);
	public Task<List<ChildParent>> GetChildParentsAsync(int enrollmentId);
	public Task<int> NewEnrollmentAsync();
	public Task<List<Address>> GetAddressAsync(int enrollmentId);
	public Task UpdateAddressAsync(int enrollmentId, Address address);
	public Task UpdateModalitiesAsync(int enrollmentId, List<Modalities> modalities);
	public Task<List<Modalities>> GetModalitiesAsync(int enrollmentId);

	public Task UpdateShiftAsync(int enrollmentId, EnrollmentShift shift);
	public Task<List<EnrollmentShift>> GetShiftAsync(int enrollmentId);
	public Task<bool> CheckEnrollmentAsync(int enrollmentId);
	public Task UpdateTermAsync(int enrollmentId, bool term);
	public Task<bool> GetTermAsync(int enrollmentId);
	public Task SubmitEnrollmentAsync(int enrollmentId);
} 

public class ProfileRepository : IProfileRepository {
	private readonly SqliteConnection _db;
	private readonly ILogger<ProfileRepository> _logger;

	public ProfileRepository(SqliteConnection db, ILogger<ProfileRepository> logger) {
		_db = db;
		_logger = logger;
	}

	public async Task<int> NewEnrollmentAsync() {
		var r = await _db.QueryFirstAsync<int>(@"INSERT INTO enrollments DEFAULT VALUES RETURNING ID;");
		_logger.LogDebug("[ProfileRepository][NewEnrollmentAsync] result: {r}", r);
		return r;
	}
	
	public async Task<bool> CheckEnrollmentAsync(int enrollmentId) {
		var r = await _db.QueryFirstAsync<int>(
			@"SELECT COUNT(*) FROM enrollments WHERE ID = @enrollmentId;", new { enrollmentId });
		_logger.LogDebug("[ProfileRepository][CheckEnrollmentAsync] result: {r}", r);
		return r > 0;
	}
	
	public async Task SubmitEnrollmentAsync(int enrollmentId) {
		var r = await _db.ExecuteAsync(
			@"UPDATE enrollments SET registered = true WHERE ID = @enrollmentId;", new { enrollmentId });
		_logger.LogDebug("[ProfileRepository][SubmitEnrollmentAsync] result: {r}", r);
	}
	
	public async Task UpdateChildParentsAsync(int enrollmentId, ChildParent childParent) {
		var r = await _db.ExecuteAsync(@"
			INSERT INTO parents(
				email
				,phone_number
				,full_name
			) VALUES (@email, @phoneNumber, @fullName);

			INSERT INTO child_parents (enrollment_fk, parents_fk)
			VALUES (@enrollmentId, last_insert_rowid());
			", new {
			childParent.fullName,
			childParent.email,
			childParent.phoneNumber,
			enrollmentId,
		});
		
		_logger.LogDebug("[ProfileRepository][UpdateChildParentsAsync] result: {r}", r);
	}
	
	public async Task UpdateChildProfileAsync(int enrollmentId, ChildProfile childProfile) {
		var r = await _db.ExecuteAsync(@"			
			INSERT INTO child_profile (
				full_name
				,birthdate
				,gender
				,medical_info
			) VALUES (@fullName, @birthdate, @gender, @medicalInfo);

			INSERT INTO enrollments_child (enrollment_fk, child_profile_fk) 
				VALUES (@enrollmentId, last_insert_rowid());
		", new {
			childProfile.fullName,
			childProfile.birthdate,
			gender = childProfile.gender.ToString(),
			childProfile.medicalInfo,
			enrollmentId,
		});
		
		_logger.LogDebug("[ProfileRepository][UpdateChildProfileAsync] result: {r}", r);
	}
	
	public async Task UpdateAddressAsync(int enrollmentId, Address address) {
		var r = await _db.ExecuteAsync(@"			
			INSERT INTO addresses(
				zipcode
				,street
				,house_number
				,state
				,city
			) VALUES (@zipcode, @street, @houseNumber, @state, @city);

			INSERT INTO enrollments_addresses (enrollment_fk, addresses_fk) 
				VALUES (@enrollmentId, last_insert_rowid());
		", new {
			zipcode = address.ZipCode,
			street = address.Street,
			houseNumber = address.HouseNumber,
			state = address.State,
			city = address.City,
			enrollmentId,
		});
		
		_logger.LogDebug("[ProfileRepository][UpdateAddressAsync] result: {r}", r);
	}
	
	public async Task<List<Address>> GetAddressAsync(int enrollmentId) {
		var r = await _db.QueryMultipleAsync(@"			
			SELECT 
				a.city
				,a.house_number as 'houseNumber'
				,a.state
				,a.street
				,a.zipcode
			FROM addresses a
			INNER JOIN enrollments_addresses ea 
				ON ea.addresses_fk = a.ID 
			WHERE ea.enrollment_fk =  @enrollmentId
			ORDER BY a.ID DESC
		", new { enrollmentId });

		_logger.LogDebug("[ProfileRepository][GetChildProfileAsync] result: {r}", r);
		return r.Read<Address>().ToList();
	}
	
	public async Task<List<ChildProfile>> GetChildProfileAsync(int enrollmentId) {
		var r = await _db.QueryMultipleAsync(@"			
			SELECT 
				full_name as 'fullName'
				,birthdate
				,gender
				,medical_info as 'medicalInfo'
			FROM child_profile cp
			JOIN enrollments_child ec
				ON ec.child_profile_fk = cp.ID
			WHERE ec.enrollment_fk = @enrollmentId
			ORDER BY cp.ID DESC
		", new { enrollmentId });

		_logger.LogDebug("[ProfileRepository][GetChildProfileAsync] result: {r}", r);
		return r.Read<ChildProfile>().ToList();
	}
	
	public async Task<List<ChildParent>> GetChildParentsAsync(int enrollmentId) {
		var r = await _db.QueryMultipleAsync(@"
			SELECT
				p.email
				,p.full_name as 'fullName'
				,p.phone_number as 'phoneNumber'
			FROM parents p
			JOIN child_parents ec 
				ON ec.parents_fk = p.ID
			WHERE ec.enrollment_fk = @enrollmentId
			ORDER BY p.ID DESC;
		", new { enrollmentId });

		_logger.LogDebug("[ProfileRepository][GetChildParentsAsync] result: {r}", r);
		return r.Read<ChildParent>().ToList();
	}
	
	public async Task UpdateModalitiesAsync(int enrollmentId, List<Modalities> modalities) { 
		var queryParams = 
			from m in modalities
				select new { name=m.ToString(), enrollmentId };

		await _db.ExecuteAsync(@"DELETE FROM student_modalities where enrollment_fk = @enrollmentId;"
			, new { enrollmentId });
		await _db.ExecuteAsync(@"
				INSERT INTO student_modalities(enrollment_fk, modalities_fk) 
					VALUES (@enrollmentId, (SELECT m.ID FROM modalities m where m.name = @name));
			", queryParams);
		
		_logger.LogDebug("[ProfileRepository][UpdateAddressAsync]");
	}

	public async Task<List<Modalities>> GetModalitiesAsync(int enrollmentId) {
		var r = await _db.QueryMultipleAsync(@"			
			SELECT m.name FROM modalities m
			LEFT JOIN student_modalities sm 
			    ON m.ID = sm.modalities_fk
			WHERE sm.enrollment_fk = @enrollmentId
		", new { enrollmentId });

		var modalities = 
			from name in r.Read<string>() 
			select Enum.Parse<Modalities>(name);
		
		_logger.LogDebug("[ProfileRepository][GetModalitiesAsync] result: {modalities}", modalities);
		return modalities.ToList();
	}
	
	
	public async Task UpdateShiftAsync(int enrollmentId, EnrollmentShift shift) {
		var r = await _db.ExecuteAsync(@"
				DELETE FROM enrollments_shifts where enrollment_fk = @enrollmentId;

				INSERT INTO enrollments_shifts(enrollment_fk, enrollments_shift_fk) 
					VALUES (@enrollmentId, (SELECT ID FROM enrollments_shift where name = @shift));
			", new {
			shift = shift.ToString(),
			enrollmentId,
		});
		
		_logger.LogDebug("[ProfileRepository][UpdateShiftsAsync] result: {r}", r);
	}
	
	public async Task<List<EnrollmentShift>> GetShiftAsync(int enrollmentId) {
		var r = await _db.QueryMultipleAsync(@"			
			SELECT es.name FROM enrollments_shift es
			LEFT JOIN enrollments_shifts ess
			    ON es.ID = ess.enrollments_shift_fk
			WHERE ess.enrollment_fk = @enrollmentId;
		", new { enrollmentId });
		
		var shifts = 
			from name in r.Read<string>() 
			select Enum.Parse<EnrollmentShift>(name);
		
		_logger.LogDebug("[ProfileRepository][GetShiftAsync] result: {r}", r);
		return shifts.ToList();
	}
	
	public async Task UpdateTermAsync(int enrollmentId, bool term) {
		var r = await _db.ExecuteAsync(@"
				INSERT INTO enrollments_terms(enrollment_fk, terms_agreement) 
					VALUES (@enrollmentId, @term);
			", new {
			term,
			enrollmentId,
		});
		
		_logger.LogDebug("[ProfileRepository][UpdateTermAsync] result: {r}", r);
	}
	
	public async Task<bool> GetTermAsync(int enrollmentId) {
		var r = await _db.QuerySingleAsync<int>(@"			
			SELECT count(*) FROM enrollments_terms 
				WHERE terms_agreement = true AND enrollment_fk = @enrollmentId;
		", new { enrollmentId });
		
		_logger.LogDebug("[ProfileRepository][GetTermAsync] result: {r}", r);
		return r > 0;
	}
}

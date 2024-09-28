namespace furry_profile;




public class ProfileService {
	private readonly IProfileRepository _repo;
	private readonly ILogger<ProfileService> _logger;

	public ProfileService(IProfileRepository repo, ILogger<ProfileService> logger) {
		_repo = repo;
		_logger = logger;
	}

	public async Task<int> NewEnrollmentAsync() {
		return await _repo.NewEnrollmentAsync();
	}

	public async Task UpdateChildProfileAsync(int enrollmentId, ChildProfile childProfile) {
		_logger.LogDebug("[ProfileService][UpdateChildProfileAsync] for enrollment: {}", enrollmentId);

		// TODO: check update conditions;
		await _repo.UpdateChildProfileAsync(enrollmentId, childProfile);
	}

	public async Task UpdateChildParentsAsync(int enrollmentId, ChildParent childParent) {
		_logger.LogDebug("[ProfileService][UpdateChildProfileAsync] for enrollment: {}", enrollmentId);

		// TODO: check update conditions;
		await _repo.UpdateChildParentsAsync(enrollmentId, childParent);
	}

	public async Task<EnrollmentState> GetFullEnrolmentAsync(int enrollmentId) {
		if (!await _repo.CheckEnrollmentAsync(enrollmentId)) {
			throw new Exception("Fail to GetFullEnrolmentAsync, enrollment not exists!");
		}

		var childProfile = await _repo.GetChildProfileAsync(enrollmentId);
		var childParent = await _repo.GetChildParentsAsync(enrollmentId);
		var address = await _repo.GetAddressAsync(enrollmentId);
		var shift = await _repo.GetShiftAsync(enrollmentId);
		var modalities = await _repo.GetModalitiesAsync(enrollmentId);
		var terms = await _repo.GetTermAsync(enrollmentId);

		return new EnrollmentState {
			Id = enrollmentId,
			childProfile = childProfile.FirstOrDefault(new ChildProfile { gender = Gender.None }),
			childParent = childParent.FirstOrDefault(new ChildParent()),
			address = address.FirstOrDefault(new Address()),
			modalities = modalities,
			enrollmentShift = shift.FirstOrDefault(EnrollmentShift.None),
			terms = terms,
		};
	}

	public async Task UpdateAddressAsync(int enrollmentId, Address address) {
		_logger.LogDebug("[ProfileService][UpdateAddressAsync] for enrollment: {}", enrollmentId);
		await _repo.UpdateAddressAsync(enrollmentId, address);
	}

	public async Task UpdateModalitiesAsync(int enrollmentId, List<Modalities> modalities) {
		_logger.LogDebug("[ProfileService][UpdateModalitiesAsync] for enrollment: {}", enrollmentId);
		await _repo.UpdateModalitiesAsync(enrollmentId, modalities);
	}

	public async Task UpdateShiftAsync(int enrollmentId, EnrollmentShift shift) {
		_logger.LogDebug("[ProfileService][UpdateShiftAsync] for enrollment: {}", enrollmentId);
		await _repo.UpdateShiftAsync(enrollmentId, shift);
	}

	public async Task UpdateTermsAsync(int enrollmentId, bool terms) {
		_logger.LogDebug("[ProfileService][UpdateTermsAsync] for enrollment: {}", enrollmentId);
		await _repo.UpdateTermAsync(enrollmentId, terms);
	}

	public async Task<bool> CheckEnrolmentAsync(int enrollmentId) {
		_logger.LogDebug("[ProfileService][CheckEnrolmentAsync] for enrollment: {}", enrollmentId);
		return await _repo.CheckEnrollmentAsync(enrollmentId);
	}

	public async Task SubmitEnrollmentAsync(int enrollmentId) {
		_logger.LogDebug("[ProfileService][SubmitEnrollment] for enrollment: {}", enrollmentId);
		
		if (!await _repo.CheckEnrollmentAsync(enrollmentId)) {
			throw new InvalidOperationException("Fail to GetFullEnrolmentAsync, enrollment not exists!");
		}
		
		var childProfile = await _repo.GetChildProfileAsync(enrollmentId);
		assertMinRelation(childProfile.Count, 1, "childInfos");
		
		var childParent = await _repo.GetChildParentsAsync(enrollmentId);
		assertMinRelation(childParent.Count, 1, "childParent");
		
		var address = await _repo.GetAddressAsync(enrollmentId);
		assertMinRelation(address.Count, 1, "address");
		
		var shift = await _repo.GetShiftAsync(enrollmentId);
		assertMinRelation(shift.Count, 1, "enrollmentShift");
		
		var modalities = await _repo.GetModalitiesAsync(enrollmentId);
		assertMinRelation(modalities.Count, 1, "modalities");
		
		var terms = await _repo.GetTermAsync(enrollmentId);
		if (!terms) throw new InvalidOperationException("terms agreement not accepted");

		await _repo.SubmitEnrollmentAsync(enrollmentId);
	}

	private void assertMinRelation(int total, int min, string relation) {
		if (total < min) {
			throw new InvalidOperationException($"[{relation}] does not satisfy minimum requirements!");
		}
	}
}


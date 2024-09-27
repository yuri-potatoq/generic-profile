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
	    _logger.LogDebug("[ProfileService][GetFullEnrolmentAsync] for enrollment: {}", enrollmentId);
	    return await _repo.GetFullEnrollment(enrollmentId);
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
    
    public async Task BulkUpdateEnrollmentStateAsync(int enrollmentId) {
	    _logger.LogDebug("[ProfileService][BatchUpdateEnrollmentAsync] for enrollment: {}", enrollmentId);
    }
}


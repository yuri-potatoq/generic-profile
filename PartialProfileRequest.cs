using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using FluentValidation;

namespace furry_profile;

public class SubmitProfileRequest {
	public int Id { get; set; }
}

public class PartialProfileRequest {
	public int? Id { get; set; }
	
	[JsonPropertyName("childInfos")]
	public ChildInfo? ChildInfos { get; set; }
	
	[JsonPropertyName("address")]
	public AddressInfo? Address { get; set; }
	
	[JsonPropertyName("parent")]
	public ParentInfo? Parent { get; set; }
	
	[JsonPropertyName("modalities")]
	public List<string>? Modalities { get; set; }
	
	[JsonPropertyName("enrollmentShift")]
	public string? EnrollmentShift { get; set; }
	
	[JsonPropertyName("terms")]
	public bool? Terms { get; set; }
	
	private static int calculateAge(DateOnly birthDate) {
		var dt = birthDate.ToDateTime(TimeOnly.MinValue);
		DateTime today = DateTime.Today;
		int age = today.Year - dt.Year;
		if (dt.Date > today.AddYears(-age)) age--;
		return age;
	}

	public record ChildInfo(
		string FullName,
		DateOnly Birthdate,
		string Gender,
		string MedicalRecords
	) {
		public class Validator : AbstractValidator<ChildInfo> {
			public Validator() {
				RuleFor(x => x.FullName).Length(3, 100);
				RuleFor(x => x.Birthdate)
					.Must(birthdate =>
						calculateAge(birthdate) > 1 && calculateAge(birthdate) <= 14)
					.WithMessage("Child should have between 1 and 14 years");
				RuleFor(x => x.Gender)
					.IsEnumName(typeof(Gender), caseSensitive: false);
				RuleFor(x => x.MedicalRecords).Length(3, 500);
			}
		}

		public ChildProfile ToCommand() {
			return new ChildProfile {
				fullName = FullName,
				birthdate = Birthdate.ToString("O"),
				medicalInfo = MedicalRecords,
				gender = Gender switch {
					"Female" => furry_profile.Gender.Female,
					"Male" => furry_profile.Gender.Male,
					_ => throw new ArgumentOutOfRangeException($"Not expected direction value: {Gender}")
				},
			};
		}
	}

	public record AddressInfo(
		string ZipCode,
		string Street,
		string City,
		string State,
		int Number
	) {
		public class Validator : AbstractValidator<AddressInfo> {
			public Validator() {
				RuleFor(x => x.ZipCode).Length(8, 8);
				RuleFor(x => x.Street).Length(3, 100);
				RuleFor(x => x.City).Length(3, 100);
				RuleFor(x => x.State).Length(3, 100);
				RuleFor(x => x.Number).NotEmpty();
			}
		}
		
		public Address ToCommand() {
			return new Address {
				City = City,
				State = State,
				Street = Street,
				ZipCode = ZipCode,
				HouseNumber = Number
			};
		}
	}

	public record ParentInfo(
		string Email,
		string PhoneNumber,
		string FullName
	) {
		public class Validator : AbstractValidator<ParentInfo> {
			public Validator() {
				RuleFor(x => x.Email).EmailAddress();
				RuleFor(x => x.PhoneNumber).Length(8, 14);
				RuleFor(x => x.FullName).Length(3, 100);
			}
		}

		public ChildParent ToCommand() {
			return new ChildParent {
				fullName = FullName,
				email = Email,
				phoneNumber = PhoneNumber
			};
		}
	}

	public List<Modalities> ToModalitiesCommand() {
		return Modalities?.Select(m => Enum.Parse<Modalities>(m)).ToList();
	}
	
	public class Validator : AbstractValidator<PartialProfileRequest> {
		public Validator() {
			RuleFor(x => x.Id)
				.GreaterThan(0);
			RuleFor(x => x.ChildInfos)
				.SetValidator(new ChildInfo.Validator())
				.When(partial => partial.ChildInfos is not null);
			RuleFor(x => x.Address)
				.SetValidator(new AddressInfo.Validator())
				.When(partial => partial.Address is not null);
			RuleFor(x => x.Parent)
				.SetValidator(new ParentInfo.Validator())
				.When(partial => partial.Parent is not null);
			RuleFor(x => x.Modalities)
				.ForEach(m => 
					m.IsEnumName(typeof(Modalities), caseSensitive: true))
				.When(partial => partial.Modalities is not null);
			RuleFor(x => x.EnrollmentShift)
				.IsEnumName(typeof(EnrollmentShift), caseSensitive: true)
				.When(partial => partial.EnrollmentShift is not null);
		}
	}
}


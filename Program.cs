using System.Net;
using Dapper;
using Microsoft.AspNetCore.Diagnostics;
using FluentValidation;
using Microsoft.Data.Sqlite;

namespace  furry_profile;


enum ErrorCode : int {
	Unknown = 1,
}

record ErrorMessage(
	ErrorCode code,
	string message
);

internal class Program {
    private static async Task<int> Main(string[] args) {
        var builder = WebApplication.CreateBuilder(args);
        
        var connectionString = 
	        builder.Configuration.GetConnectionString("FurryProfile") ?? "Data Source=furry_profile.db;Cache=Shared";
	
        builder.Services.AddSingleton<SqliteConnection>(_ => new SqliteConnection(connectionString));
        builder.Services.AddScoped<IProfileRepository, ProfileRepository>();
        builder.Services.AddScoped<ProfileService>();
        builder.Services.AddScoped<ProfileHandler>();
        
        builder.Services.AddValidatorsFromAssemblyContaining<Program>();
        builder.Services.AddEndpointsApiExplorer();
        
        builder.Services.AddSwaggerGen();
        builder.Services.AddControllers();
        builder.Services.AddProblemDetails();

        var app = builder.Build();
	
        app.UseExceptionHandler(exceptionHandlerApp 
	        => exceptionHandlerApp.Run(async context => {
		        var exception = context.Features.Get<IExceptionHandlerFeature>();
		        var err = exception?.Error;
		        await Results.Problem(err?.Message, statusCode: 500).ExecuteAsync(context);
	        }));
        
        app.UseSwagger();
        app.UseSwaggerUI();
        app.UseHttpsRedirection();
        
        using (var serviceScope = app.Services.CreateScope()) {
	        var services = serviceScope.ServiceProvider;
	        var profileHandler = services.GetRequiredService<ProfileHandler>();

	        var profileGroup = app.MapGroup("/profile");
        
	        profileGroup.MapPatch("/", profileHandler.PartialUpdateAsync)
		        .AddEndpointFilter<ModelValidationFilter<PartialProfileRequest>>();
	        profileGroup.MapGet("/{enrollmentId}", profileHandler.GetEnrollmentAsync);
        }
		
        await app.RunAsync("http://[::]:5000");
        return 0;
    }

    public async Task EnsureMigrations(IServiceProvider services) {
	    using var db = services.CreateScope().ServiceProvider.GetRequiredService<SqliteConnection>();
	    // TODO: build up migrations
	    await db.ExecuteAsync("");
    }
}

public class ModelValidationFilter<T> : IEndpointFilter {
	public async ValueTask<object?> InvokeAsync(EndpointFilterInvocationContext context, EndpointFilterDelegate next) {
		T? argToValidate = context.GetArgument<T>(0);
		IValidator<T>? validator = context.HttpContext.RequestServices.GetService<IValidator<T>>();

		if (validator is not null) {
			var validationResult = await validator.ValidateAsync(argToValidate!);
			if (!validationResult.IsValid)
			{
				return Results.ValidationProblem(validationResult.ToDictionary(),
					statusCode: (int)HttpStatusCode.UnprocessableEntity);
			}
		}
		return await next.Invoke(context);
	}
}

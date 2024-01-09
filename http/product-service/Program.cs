using System.Reflection;
using FluentValidation;
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;

var builder = WebApplication.CreateBuilder(args);

var dbHost = Environment.GetEnvironmentVariable("POSTGRES_HOST");
var dbPort = Environment.GetEnvironmentVariable("POSTGRES_PORT");
var dbName = Environment.GetEnvironmentVariable("POSTGRES_DATABASE");
var dbUsername = Environment.GetEnvironmentVariable("POSTGRES_USERNAME");
var dbPassword = Environment.GetEnvironmentVariable("POSTGRES_PASSWORD");

string connectionString = $"Host={dbHost};Port={dbPort};Database={dbName};Username={dbUsername};Password={dbPassword}";

builder.Services.AddDbContext<ProductContext>(options => options.UseNpgsql(connectionString));

builder.Services.AddValidatorsFromAssembly(Assembly.GetExecutingAssembly());

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

app.MapSwagger();
app.UseSwaggerUI();

using (var scope = app.Services.CreateScope())
{
    var context = scope.ServiceProvider.GetRequiredService<ProductContext>();
    await context.Database.EnsureCreatedAsync();
}

app.UseStatusCodePages(async statusCodeContext
    => await Results.Problem(statusCode: statusCodeContext.HttpContext.Response.StatusCode)
                 .ExecuteAsync(statusCodeContext.HttpContext));

app.MapGet("/api/products", (ProductContext db) => db.Products.ToListAsync());

app.MapGet("/api/products/{id}", async (ProductContext db, Guid id) =>
    await db.Products.FindAsync(id)
        is Product product
            ? Results.Ok(product)
            : Results.NotFound());

app.MapPost("/api/products", async (ProductContext db, IValidator<Product> validator, Product product) =>
{
    var result = await validator.ValidateAsync(product);
    if (!result.IsValid)
    {
        return Results.ValidationProblem(result.ToDictionary());
    }

    db.Products.Add(product);
    await db.SaveChangesAsync();

    return Results.Created($"/api/products/{product.Id}", product);
});

app.MapPut("/api/products/{id}", async (ProductContext db, IValidator<Product> validator, Guid id, Product request) =>
{
    var product = await db.Products.FindAsync(id);

    if (product is null) return Results.NotFound();

    var result = await validator.ValidateAsync(request);
    if (!result.IsValid)
    {
        return Results.ValidationProblem(result.ToDictionary());
    }

    product.Name = request.Name;
    product.Description = request.Description;
    product.Price = request.Price;

    await db.SaveChangesAsync();

    return Results.NoContent();
});

app.MapDelete("/api/products/{id}", async (ProductContext db, Guid id) =>
{
    if (await db.Products.FindAsync(id) is Product product)
    {
        db.Products.Remove(product);
        await db.SaveChangesAsync();
        return Results.NoContent();
    }

    return Results.NotFound();
});


app.Run();

using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDbContext<UserDb>(opt => opt.UseInMemoryDatabase("Users"));
builder.Services.AddDatabaseDeveloperPageExceptionFilter();
var app = builder.Build();

app.MapGet("/users", async (UserDb db) =>
    await db.Users.Include(x => x.Addresses).ToListAsync());

app.MapGet("/users/{id}", async (int id, UserDb db) =>
    await db.Users.Include(x => x.Addresses).FirstOrDefaultAsync(x => x.Id == id)
        is User user
            ? Results.Ok(user)
            : Results.NotFound());

app.MapPost("/users", async (User user, UserDb db) =>
{
    await db.Users.AddAsync(user);
    await db.SaveChangesAsync();

    return Results.Created($"/users/{user.Id}", user);
});

app.MapPut("/users/{id}", async (int id, User req, UserDb db) =>
{
    var user = await db.Users.FindAsync(id);

    if (user is null) return Results.NotFound();

    user.Name = req.Name;
    user.Addresses = req.Addresses;

    await db.SaveChangesAsync();

    return Results.Ok(user);
});

// Addresses
app.MapGet("/users/{userid}/addresses/{id}", async (int userId, int id, UserDb db) =>
    await db.Addresses.Include(x => x.User)
        .Where(x => x.Id == id && x.User.Id == userId).FirstOrDefaultAsync()
            is Address address
            ? Results.Ok(address)
            : Results.NotFound());

app.MapPost("/users/{id}/addresses", async (int id, Address address, UserDb db) =>
{
    var user = await db.Users.FindAsync(id);

    if (user is null) return Results.NotFound();

    user.Addresses.Add(address);
    await db.SaveChangesAsync();

    return Results.Created($"/users/{user.Id}", user);
});

app.Run();

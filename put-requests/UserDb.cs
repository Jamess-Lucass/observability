using Microsoft.EntityFrameworkCore;

class UserDb : DbContext
{
    public UserDb(DbContextOptions<UserDb> options) : base(options) { }

    public DbSet<User> Users => Set<User>();
    public DbSet<Address> Addresses => Set<Address>();
}
using Microsoft.EntityFrameworkCore;

public class OrderContext : DbContext
{
    public OrderContext(DbContextOptions<OrderContext> options) : base(options) { }

    public DbSet<Order> Orders { get; set; }
}

public record Order
{
    public Guid Id { get; set; }
    public Guid ProductId { get; set; }
    public uint Quantity { get; set; }

    [Precision(18, 2)]
    public decimal Price { get; set; }
}
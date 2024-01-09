using System.ComponentModel.DataAnnotations.Schema;
using Microsoft.EntityFrameworkCore;

public class ProductContext : DbContext
{
    public ProductContext(DbContextOptions<ProductContext> options) : base(options) { }

    public DbSet<Product> Products { get; set; }
}

public record Product
{
    public Guid Id { get; set; }

    [Column(TypeName = "varchar(128)")]
    public string Name { get; set; } = string.Empty;

    [Column(TypeName = "varchar(1024)")]
    public string Description { get; set; } = string.Empty;

    [Precision(18, 2)]
    public decimal Price { get; set; }
}
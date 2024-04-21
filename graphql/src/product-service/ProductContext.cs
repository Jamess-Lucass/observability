using System.ComponentModel.DataAnnotations.Schema;
using HotChocolate.ApolloFederation.Resolvers;
using HotChocolate.ApolloFederation.Types;
using Microsoft.EntityFrameworkCore;

public class ProductContext : DbContext
{
    public ProductContext(DbContextOptions<ProductContext> options) : base(options) { }

    public DbSet<Product> Products { get; set; }
}

public record Product : IResponse
{
    [Key]
    public Guid Id { get; set; }

    [Column(TypeName = "varchar(128)")]
    public string Name { get; set; } = string.Empty;

    [Column(TypeName = "varchar(1024)")]
    public string Description { get; set; } = string.Empty;

    [Precision(18, 2)]
    public decimal Price { get; set; }

    [GraphQLIgnore]
    public bool IsDeleted { get; set; }

    // Resolvers
    [ReferenceResolver]
    public static async Task<Product?> ResolveReference(
        Guid id,
        [Service] ProductContext db
    )
    {
        return await db.Products.FindAsync(id);
    }
}

[Shareable]
public record Error
{
    public string Message { get; set; } = string.Empty;
    public string Path { get; set; } = string.Empty;
}

[Shareable]
public record ErrorPayload : IResponse
{
    public IEnumerable<Error>? Errors { get; set; }
}

[UnionType("Response")]
public interface IResponse
{
}
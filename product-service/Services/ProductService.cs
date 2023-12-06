using Grpc.Core;
using Microsoft.EntityFrameworkCore;

namespace ProductServiceGrpc.Services;

public class ProductService : Product.ProductBase
{
    private readonly ILogger<ProductService> _logger;
    private readonly ProductContext _db;

    public ProductService(ILogger<ProductService> logger, ProductContext db)
    {
        _logger = logger;
        _db = db;
    }

    public override async Task<GetAllProductsResponse> GetAllProducts(GetAllProductsRequest request, ServerCallContext context)
    {
        var products = await _db.Products.Select(x => new GetProductResponse
        {
            Id = x.Id.ToString(),
            Name = x.Name,
            Description = x.Description
        }).ToListAsync();

        var response = new GetAllProductsResponse();
        response.Value.AddRange(products);

        return response;
    }

    public override async Task<GetProductResponse> GetProduct(GetProductRequest request, ServerCallContext context)
    {
        if (!Guid.TryParse(request.Id, out var id))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid Id"));
        }

        var product = await _db.Products.FindAsync(id);
        if (product is null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Not found"));
        }

        return new GetProductResponse
        {
            Id = product.Id.ToString(),
            Name = product.Name,
            Description = product.Description
        };
    }
}
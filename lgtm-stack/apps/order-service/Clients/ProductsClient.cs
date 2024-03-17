public class ProductsClient : IProductsClient
{
    private readonly HttpClient _httpClient;

    public ProductsClient(HttpClient httpClient)
    {
        _httpClient = httpClient;
    }

    public async Task<ProductResponse?> GetAsync(Guid id)
    {
        var response = await _httpClient.GetAsync($"/api/products/{id}");

        if (!response.IsSuccessStatusCode)
        {
            return null;
        }

        return await response.Content.ReadFromJsonAsync<ProductResponse>();
    }
}

public interface IProductsClient
{
    Task<ProductResponse?> GetAsync(Guid id);
}

public record ProductResponse
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public decimal Price { get; set; }
}
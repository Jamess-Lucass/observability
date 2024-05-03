var convertor = new Dictionary<string, decimal>()
{
    {"USD", 1.25327m}
};

var product = new Product(5.00m, "GBP");

var USDPricing = convertor["USD"] * product.Price;

var order = new Order(USDPricing, "USD");

Console.WriteLine(order.Price);


class Order
{
    public Order(decimal price, string currency)
    {
        Price = price;
        Currency = currency;
    }

    public decimal Price { get; set; }
    public string Currency { get; set; } = string.Empty;
}

class Product
{
    public Product(decimal price, string currency)
    {
        Price = price;
        Currency = currency;
    }

    public decimal Price { get; set; }
    public string Currency { get; set; } = string.Empty;
}
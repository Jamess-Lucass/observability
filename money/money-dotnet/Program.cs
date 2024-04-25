var order = new Order()
{
    Items = new List<OrderItem>()
    {
        new OrderItem(0.333m),
        new OrderItem(0.333m),
        new OrderItem(0.333m),
    }
};

var totalPrice = order.Items.Sum(x => x.Price);

Console.WriteLine(totalPrice);


class Order
{
    public IEnumerable<OrderItem> Items { get; set; } = new List<OrderItem>();
}

class OrderItem
{
    public OrderItem(decimal price)
    {
        Price = price;
    }

    public decimal Price { get; set; }
}
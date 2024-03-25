public class User
{
    public int Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public ICollection<Address> Addresses { get; set; } = new List<Address>();
}

public class Address
{
    public int Id { get; set; }
    public string AddressLine1 { get; set; } = string.Empty;
    public string PostCode { get; set; } = string.Empty;

    public User User { get; } = new();
}
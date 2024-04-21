using FluentValidation;

public class CreateProductInput
{
    public CreateProductInput(
            string name,
            string description,
            decimal price
        )
    {
        Name = name;
        Description = description;
        Price = price;
    }

    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public decimal Price { get; set; }
}

public class ProductValidator : AbstractValidator<CreateProductInput>
{
    public ProductValidator()
    {
        RuleFor(x => x.Name).Length(2, 128);
        RuleFor(x => x.Description).Length(2, 1024);
        RuleFor(x => x.Price).GreaterThan(0).PrecisionScale(18, 2, false);
    }
}
